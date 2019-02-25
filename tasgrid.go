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

	utm "github.com/kurankat/UTM"
)

// TasMap holds map-unique information
type TasMap struct {
	zone          int
	alpha         string
	eastingStart  string
	northingStart string
}

// NewTasMap generates the necessary information for each map
func NewTasMap(mapInfo []string) TasMap {
	zone, err := strconv.Atoi(mapInfo[1])
	checkError(err)
	alpha := mapInfo[2]
	east := mapInfo[3]
	north := mapInfo[4]

	return TasMap{zone: zone, alpha: alpha, eastingStart: east, northingStart: north}
}

// GridPoint holds all the necessary information to calculate the full
// grid reference of a point from its 6-figure map grid reference
type GridPoint struct {
	MapName      string
	easting3d    string
	northing3d   string
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
// easting and northing
func NewGridPoint(name, eas, nor string, mg MapGrid) (GridPoint, error) {
	if len(name) != 3 || len(eas) != 3 || len(nor) != 3 {
		return GridPoint{}, nil
	}
	mapName := strings.ToUpper(name)
	gp := GridPoint{MapName: mapName, easting3d: eas, northing3d: nor}
	var strEasting string
	var strNorthing string

	easting, err := strconv.Atoi(eas)
	if err != nil {
		return GridPoint{}, fmt.Errorf("ERROR parsing grid: I can't convert easting %v to an integer", eas)
	}
	northing, err := strconv.Atoi(nor)
	if err != nil {
		return GridPoint{}, fmt.Errorf("ERROR parsing grid: I can't convert northing %v to an integer", nor)
	}

	ess := mg[mapName].eastingStart
	nss := mg[mapName].northingStart

	if len(ess)+len(nss) == 0 {
		return GridPoint{}, fmt.Errorf("ERROR parsing grid: I'm having trouble geting values for map %v", name)
	}

	easStart, err := strconv.Atoi(ess[:1])
	norStart, err := strconv.Atoi(nss[:2])

	easEnd, err := strconv.Atoi(ess[len(ess)-2:])
	if err != nil {
		return GridPoint{}, fmt.Errorf("ERROR parsing grid: I can't extract a number from %v", ess)
	}
	norEnd, err := strconv.Atoi(nss[len(nss)-2:])
	if err != nil {
		return GridPoint{}, fmt.Errorf("ERROR parsing grid: I can't extract a number from %v", nss)
	}

	if easting > easEnd*10 {
		strEasting = ess[:1] + eas + "00"
	} else {
		newEss := strconv.Itoa(easStart + 1)
		strEasting = newEss + eas + "00"
	}
	gp.fullEasting, err = strconv.ParseFloat(strEasting, 64)

	if northing > norEnd*10 {
		strNorthing = nss[:2] + nor + "00"
	} else {
		newNss := strconv.Itoa(norStart + 1)
		strNorthing = newNss + nor + "00"
	}
	gp.fullNorthing, err = strconv.ParseFloat(strNorthing, 64)

	gp.decimalLat, gp.decimalLong, err = utm.ToLatLon(gp.fullEasting, gp.fullNorthing, 55, "G")
	if err != nil {
		fmt.Printf("Map name: %v, 3f easting: %v, 3f northing: %v", mapName, eas, nor)
		panic(err)
	}

	gp.latDegs, gp.latMins, gp.latSecs = ddToDMS(gp.decimalLat)
	gp.longDegs, gp.longMins, gp.longSecs = ddToDMS(gp.decimalLong)

	return gp, nil
}

// GetFullEasting returns the full-length easting of the grid point
func (gp GridPoint) GetFullEasting() (easting string) {
	easting = strconv.FormatFloat(gp.fullEasting, 'f', 0, 64)
	return easting
}

// GetFullNorthing returns the full-length northing of the grid point
func (gp GridPoint) GetFullNorthing() (northing string) {
	northing = strconv.FormatFloat(gp.fullNorthing, 'f', 0, 64)
	return
}

// GetDecimalLat returns the latitude in decimal degrees of the grid point
func (gp GridPoint) GetDecimalLat() (dLat string) {
	dLat = strconv.FormatFloat(gp.decimalLat, 'f', 6, 64)
	return
}

// GetDecimalLong returns the longitude in decimal degrees of the grid point
func (gp GridPoint) GetDecimalLong() (dLong string) {
	dLong = strconv.FormatFloat(gp.decimalLong, 'f', 6, 64)
	return
}

// GetLatSeconds returns the seconds reading of the grid point's latitude
func (gp GridPoint) GetLatSeconds() (secs string) {
	if gp.latSecs < 0 {
		gp.latSecs = -gp.latSecs
	}
	secs = strconv.FormatFloat(gp.latSecs, 'f', 1, 64)
	return
}

// GetLongSeconds returns the seconds reading of the grid point's latitude
func (gp GridPoint) GetLongSeconds() (secs string) {
	if gp.longSecs < 0 {
		gp.longSecs = -gp.longSecs
	}
	secs = strconv.FormatFloat(gp.longSecs, 'f', 1, 64)
	return
}

// MapGrid holds a dictionary of TASMAP three-letter acronyms
// containing the map's necessary data to calculate the full grid reference
type MapGrid map[string]TasMap

// NewTasMapGrid returns a TasMap object
func NewTasMapGrid() *MapGrid {
	mapList := MapGrid{}

	gopath := os.Getenv("GOPATH")
	mapFile, err := os.Open(filepath.Join(gopath, "src/mypacks/tasgrid/mapinfo.csv"))
	checkError(err)
	defer mapFile.Close() // Defer closing until the program is done

	mapReader := csv.NewReader(mapFile)

	for {
		// Read line into memory
		tasMap, err := mapReader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			checkError(err)
		}

		name := tasMap[0]
		mapList[name] = NewTasMap(tasMap)
	}

	return &mapList
}

func ddToDMS(dd float64) (degs, mins, secs float64) {

	degs = math.Trunc(dd)                    // Degrees as float is the truncated decimal degrees
	minDiff := math.Abs(dd) - math.Abs(degs) // What remains right of the decimal point is the decimal mins
	dMins := minDiff * 60.0                  // Float minutes is decimal minutes * 60
	mins = math.Trunc(dMins)
	secDiff := dMins - mins
	secs = secDiff * 60.0

	return
}

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
