package tasgrid

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
)

// MapGrid holds a dictionary of TASMAP three-letter acronyms
// containing the map's necessary data to calculate the full grid reference
type MapGrid map[string]TasMap

// TasMap holds map-unique information
type TasMap struct {
	zone          int
	alpha         string
	eastingStart  string
	northingStart string
}

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
	fullEasting  string
	fullNorthing string
}

// NewGridPoint creates a GridPoint from the supplied map name and 3-digit
// easting and northing
func NewGridPoint(name, eas, nor string, mg MapGrid) (gp GridPoint) {
	gp = GridPoint{mapName: name, easting3d: eas, northing3d: nor}

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
		gp.fullEasting = ess[:1] + eas + "00"
	} else {
		newEss := strconv.Itoa(easStart + 1)
		gp.fullEasting = newEss + eas + "00"
	}

	if northing > norEnd*10 {
		gp.fullNorthing = nss[:2] + nor + "00"
	} else {
		newNss := strconv.Itoa(norStart + 1)
		gp.fullNorthing = newNss + nor + "00"
	}

	fmt.Println(eas, easEnd*10, easStart, nor, norEnd*10, norStart, gp.fullEasting, gp.fullNorthing)

	return gp
}

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
