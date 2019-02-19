package tasgrid

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
)

// TasMap holds map-unique information
type TasMap struct {
	zone          int
	alpha         string
	eastingStart  int
	northingStart int
}

// MapGrid holds a dictionary of TASMAP three-letter acronyms
// containing the map's necessary data to calculate the full grid reference
type MapGrid map[string]TasMap

// NewTasMapGrid returns a TasMap object
func NewTasMapGrid() *MapGrid {
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
		zone, err := strconv.Atoi(tasMap[1])
		alpha := tasMap[2]
		east, err := strconv.Atoi(tasMap[3])
		north, err := strconv.Atoi(tasMap[4])

		mapList[name] = TasMap{zone: zone, alpha: alpha, eastingStart: east, northingStart: north}

	}

	return &mapList
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
