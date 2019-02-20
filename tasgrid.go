package tasgrid

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"

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
	mapName      string
	easting3d    string
	northing3d   string
	fullEasting  float64
	fullNorthing float64
	decimalLat   float64
	decimalLong  float64
}

// NewGridPoint creates a GridPoint from the supplied map name and 3-digit
// easting and northing
func NewGridPoint(name, eas, nor string, mg MapGrid) (gp GridPoint) {
	gp = GridPoint{mapName: name, easting3d: eas, northing3d: nor}
	var strEasting string
	var strNorthing string

	easting, err := strconv.Atoi(eas)
	checkError(err)
	northing, err := strconv.Atoi(nor)
	checkError(err)

	ess := mg[name].eastingStart
	nss := mg[name].northingStart

	easStart, err := strconv.Atoi(ess[:1])
	norStart, err := strconv.Atoi(nss[:2])

	easEnd, err := strconv.Atoi(ess[len(ess)-2:])
	checkError(err)
	norEnd, err := strconv.Atoi(nss[len(nss)-2:])
	checkError(err)

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

	fmt.Println(eas, easEnd*10, easStart, nor, norEnd*10, norStart, gp.fullEasting, gp.fullNorthing)

	return gp
}

// GetDecimalLat returns the latitude in decimal degrees of the GridPoint
func (gp GridPoint) GetDecimalLat() (dLat string) {
	dLat = strconv.FormatFloat(gp.decimalLat, 'f', 6, 64)
	return dLat
}

// GetDecimalLong returns the longitude in decimal degrees of the GridPoint
func (gp GridPoint) GetDecimalLong() (dLong string) {
	dLong = strconv.FormatFloat(gp.decimalLong, 'f', 6, 64)
	return dLong
}

// MapGrid holds a dictionary of TASMAP three-letter acronyms
// containing the map's necessary data to calculate the full grid reference
type MapGrid map[string]TasMap

// NewTasMapGrid returns a TasMap object
func NewTasMapGrid() MapGrid {
	mapList := MapGrid{}

	mapFile, err := os.Open("../mypacks/tasgrid/mapinfo.csv")
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

	return mapList
}

func checkError(err error) {
	if err != nil {
		// fmt.Printf("Error type: %T\n", err)
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
