package tasgrid

import (
	"encoding/csv"
	"fmt"
	"io"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"unicode"

	utm "github.com/im7mortal/UTM"
)

// GridPoint holds all the necessary information pertaining to a grid point, calculated from the
// name of the map provided and the three-figure easting and northing
type GridPoint struct {
	MapName      string
	fullEasting  float64
	fullNorthing float64
	decimalLat   float64
	decimalLong  float64
	latDegs      float64
	latMins      float64
	latSecs      float64
	longDegs     float64
	longMins     float64
	longSecs     float64
}

// NewGridPoint creates a GridPoint from the supplied map name and 3-digit
// easting and northing, and calculates the full easting and northing, as well as
// the latitude and longitude of the record in decimal degrees and degrees, minutes and
// seconds
func NewGridPoint(name, textEasting, textNorthing string, mg MapGrid) (GridPoint, error) {
	var mapName string            // Name of map in database, converted to uppercase
	var gp GridPoint              // Main grid point to be returned
	var stringFullEasting string  // Full UTM easting as a string
	var stringFullNorthing string // Full UTM northing as a string
	var numEasting int            // integer version of 3-digit grid reference easting from database
	var numNorthing int           // integer version of 3-digit grid reference northing from database
	var firstEasting string       // Westernmost (lowest) easting in map sheet, as a string
	var firstNorthing string      // Southernmost (lowest) northing in map sheet, as a string
	var lastEasting string        // Easternmost (highest) easting in map sheet, as a string
	var lastNorthing string       // Northernomst (highest) northing in map sheet, as a string
	var mapRangeW float64         // Full UTM version (numeric) of westernmost easting in map sheet
	var mapRangeS float64         // Full UTM version (numeric) of southernmost easting in map sheet
	var mapRangeE float64         // Full UTM version (numeric) of easternmost easting in map sheet
	var mapRangeN float64         // Full UTM version (numeric) of northernmost easting in map sheet

	// Only proceed if the information consists of a three-letter map name, and three-figure
	// easting and northings, assuming it is information from a TASMAP 1:100,000-series map
	if len(name) != 3 ||
		len(textEasting) != 3 ||
		len(textNorthing) != 3 {
		return GridPoint{}, nil
	}

	// If the map contains digits, it is a NZ map, return empty gridpoint but no error
	for _, c := range name {
		if unicode.IsDigit(c) {
			return GridPoint{}, nil
		}
	}

	mapName = strings.ToUpper(name)
	gp = GridPoint{MapName: mapName}

	// Convert string easting and northing to integers and return an error if problems arise
	numEasting, err := strconv.Atoi(textEasting)
	if err != nil {
		return GridPoint{}, fmt.Errorf("ERROR parsing grid: I can't convert easting %v to an integer", textEasting)
	}
	numNorthing, err = strconv.Atoi(textNorthing)
	if err != nil {
		return GridPoint{}, fmt.Errorf("ERROR parsing grid: I can't convert northing %v to an integer", textNorthing)
	}

	firstEasting = mg[mapName].eastingStart
	firstNorthing = mg[mapName].northingStart
	lastEasting = mg[mapName].eastingEnd
	lastNorthing = mg[mapName].northingEnd

	// Calculate the range of acceptable eastings and northings from that map sheet
	mapRangeW, err = strconv.ParseFloat(firstEasting+"000", 64)
	if err != nil {
		return GridPoint{}, fmt.Errorf("ERROR converting westernmost easting %v to a float", textNorthing)
	}
	mapRangeE, err = strconv.ParseFloat(lastEasting+"000", 64)
	if err != nil {
		return GridPoint{}, fmt.Errorf("ERROR converting easternmost easting %v to a float", textNorthing)
	}
	mapRangeS, err = strconv.ParseFloat(firstNorthing+"000", 64)
	if err != nil {
		return GridPoint{}, fmt.Errorf("ERROR converting southernmost easting %v to a float", textNorthing)
	}
	mapRangeN, err = strconv.ParseFloat(lastNorthing+"000", 64)
	if err != nil {
		return GridPoint{}, fmt.Errorf("ERROR converting northernmost easting %v to a float", textNorthing)
	}

	// If we don't have figures in the required fields, the map name may have been wrong - ignore and return an error
	if len(firstEasting)+len(firstNorthing) == 0 {
		return GridPoint{}, fmt.Errorf("ERROR parsing grid: I'm having trouble geting values for map %v", name)
	}

	// Convert the all but the last digit of the starting easting and northing lines to integers
	numFirstEasting, err := strconv.Atoi(firstEasting[:1])
	numFirstNorthing, err := strconv.Atoi(firstNorthing[:2])

	// Extract the last two figures of the easting and northing starting lines to later determine how to calculate
	// the complete easting and northing (if it carries over 99). Return errors if needed
	eastingVariable, err := strconv.Atoi(firstEasting[len(firstEasting)-2:])
	if err != nil {
		return GridPoint{}, fmt.Errorf("ERROR parsing grid: I can't extract a number from %v", firstEasting)
	}
	northingVariable, err := strconv.Atoi(firstNorthing[len(firstNorthing)-2:])
	if err != nil {
		return GridPoint{}, fmt.Errorf("ERROR parsing grid: I can't extract a number from %v", firstNorthing)
	}

	// If the easting is greater than the first line easting on the map, append the first figure from the easting
	// starting line and add two zeros to get a complete easting. However if the easting is a smaller number,
	// we need to carry one because we wrap over 100.
	if numEasting > eastingVariable*10 {
		stringFullEasting = firstEasting[:1] + textEasting + "00"
	} else {
		newFirstEasting := strconv.Itoa(numFirstEasting + 1)
		stringFullEasting = newFirstEasting + textEasting + "00"
	}
	gp.fullEasting, err = strconv.ParseFloat(stringFullEasting, 64)
	if err != nil {
		return GridPoint{}, fmt.Errorf("ERROR parsing grid: I can't extract a number from %v", stringFullEasting)
	}

	// Ditto for the northing, but using the first two figures from the starting line
	if numNorthing > northingVariable*10 {
		stringFullNorthing = firstNorthing[:2] + textNorthing + "00"
	} else {
		newFirstNorthing := strconv.Itoa(numFirstNorthing + 1)
		stringFullNorthing = newFirstNorthing + textNorthing + "00"
	}
	gp.fullNorthing, err = strconv.ParseFloat(stringFullNorthing, 64)
	if err != nil {
		return GridPoint{}, fmt.Errorf("ERROR parsing grid: I can't extract a number from %v", stringFullNorthing)
	}

	// Check if easting and northing fall within the range of expected values for their map sheet
	if gp.fullEasting < mapRangeW ||
		gp.fullEasting > mapRangeE {
		return GridPoint{}, fmt.Errorf("ERROR parsing grid: easting %v is out of the expected range for map %v", textEasting, gp.MapName)
	}
	if gp.fullNorthing < mapRangeS ||
		gp.fullNorthing > mapRangeN {
		return GridPoint{}, fmt.Errorf("ERROR parsing grid: northing %v is out of the expected range for map %v", textNorthing, gp.MapName)
	}

	// Use utm library to calculate the decimal latitude and longitude. Treat King Island specimens in
	// zone 54 as if they were zone 55, using the zone 55 numbers in the TASMAP maps (the error is small
	// enough to be safely ignored)
	gp.decimalLat, gp.decimalLong, err = utm.ToLatLon(gp.fullEasting, gp.fullNorthing, 55, "G")
	if err != nil {
		return GridPoint{}, fmt.Errorf("Map name: %v, 3f easting: %v, 3f northing: %v", mapName, textEasting, textNorthing)
	}

	// Calculate the DMS lat and long
	gp.latDegs, gp.latMins, gp.latSecs = ddToDMS(gp.decimalLat)
	gp.longDegs, gp.longMins, gp.longSecs = ddToDMS(gp.decimalLong)

	return gp, nil
}

// GetFullEasting returns the full-length easting of the grid point as a string
func (gp GridPoint) GetFullEasting() (easting string) {
	easting = strconv.FormatFloat(gp.fullEasting, 'f', 0, 64)
	return easting
}

// GetFullNorthing returns the full-length northing of the grid point as a string
func (gp GridPoint) GetFullNorthing() (northing string) {
	northing = strconv.FormatFloat(gp.fullNorthing, 'f', 0, 64)
	return
}

// GetDecimalLat returns the latitude in decimal degrees of the grid point as a string
func (gp GridPoint) GetDecimalLat() (dLat string) {
	dLat = strconv.FormatFloat(gp.decimalLat, 'f', 6, 64)
	return
}

// GetDecimalLong returns the longitude in decimal degrees of the grid point as a string
func (gp GridPoint) GetDecimalLong() (dLong string) {
	dLong = strconv.FormatFloat(gp.decimalLong, 'f', 6, 64)
	return
}

// GetLatSeconds returns the seconds reading of the grid point's latitude as a string
func (gp GridPoint) GetLatSeconds() (secs string) {
	if gp.latSecs < 0 {
		gp.latSecs = -gp.latSecs
	}
	secs = strconv.FormatFloat(gp.latSecs, 'f', 1, 64)
	return
}

// GetLongSeconds returns the seconds reading of the grid point's longitude as a string
func (gp GridPoint) GetLongSeconds() (secs string) {
	if gp.longSecs < 0 {
		gp.longSecs = -gp.longSecs
	}
	secs = strconv.FormatFloat(gp.longSecs, 'f', 1, 64)
	return
}

// GetDistance takes in a latitude and longitude and calculates the distance of that point to
// the GridPoint
func (gp GridPoint) GetDistance(lat, long string) (distance float64, err error) {
	// Parse floats from string lat and long
	fLat, err := strconv.ParseFloat(lat, 64)
	if err != nil {
		return distance, fmt.Errorf("I can't parse a latitude from %v", lat)
	}
	fLong, err := strconv.ParseFloat(long, 64)
	if err != nil {
		return distance, fmt.Errorf("I can't parse a longitude from %v", lat)
	}

	// Use UTM converter to derive full easting and northing
	easting, northing, _, _, err := utm.FromLatLon(fLat, fLong, false)
	if err != nil {
		return distance, fmt.Errorf("UTM converter has trouble with lat %v & long %v", lat, long)
	}

	eastingDistance := math.Abs(gp.fullEasting - easting)
	northingDistance := math.Abs(gp.fullNorthing - northing)

	distance = math.Sqrt((eastingDistance * eastingDistance) + (northingDistance * northingDistance))

	return
}

// TasMap holds map-unique information: the UTM zone, alphanumeric code, as well as the lowest
// easting and northing 1000m lines
type TasMap struct {
	zone                       int
	alpha                      string
	eastingStart, eastingEnd   string
	northingStart, northingEnd string
}

// newTasMap assigns the information provided in the argument (a slice of map information)
// to the correct fields
func newTasMap(mapInfo []string) TasMap {
	zone, err := strconv.Atoi(mapInfo[1])
	checkError(err)
	alpha := mapInfo[2]
	west := mapInfo[3]
	east := mapInfo[4]
	south := mapInfo[5]
	north := mapInfo[6]

	return TasMap{zone: zone, alpha: alpha, eastingStart: west, eastingEnd: east, northingStart: south, northingEnd: north}
}

// MapGrid holds a dictionary of TASMAP three-letter acronyms
// containing the map's necessary data to calculate the full grid reference
type MapGrid map[string]TasMap

// NewTasMapGrid returns a TasMap object
func NewTasMapGrid() *MapGrid {
	mapList := MapGrid{}

	// Locate the accessory data file "mapinfo.csv"
	gopath := os.Getenv("GOPATH")
	mapFile, err := os.Open(filepath.Join(gopath, "src/github.com/kurankat/tasgrid/mapinfo.csv"))
	checkError(err)
	defer mapFile.Close() // Defer closing until the program is done

	// Read it as a CSV file
	mapReader := csv.NewReader(mapFile)

	// Read each line into memory and use the data to create a tasMap object for each line
	for {
		tasMap, err := mapReader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			checkError(err)
		}

		name := tasMap[0]
		mapList[name] = newTasMap(tasMap)
	}

	return &mapList
}

// Convert decimal degrees to degrees, minutes, seconds
func ddToDMS(dd float64) (degs, mins, secs float64) {

	degs = math.Trunc(dd)                    // Degrees as float is the truncated decimal degrees
	minDiff := math.Abs(dd) - math.Abs(degs) // What remains right of the decimal point is the decimal mins
	dMins := minDiff * 60.0                  // Float minutes is decimal minutes * 60
	mins = math.Trunc(dMins)
	secDiff := dMins - mins
	secs = secDiff * 60.0

	return
}

// Check and handle errors (simplified)
func checkError(err error) {
	if err != nil {
		switch err.(type) {
		// If the error is a path error (such as file not being found) then print
		// individualised error message
		case *os.PathError:
			fmt.Println("\nI'm having trouble accessing a file. The system says:")
			fmt.Printf("\t* %v\n\n", err)
			os.Exit(1)
		default:
			// If we don't have a specific way of handling the error, print the error
			// to screen and exit
			panic(err)
		}
	}
}
