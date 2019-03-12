# tasgrid

Golang library to calculate full UTM grid reference for Tasmanian 6-figure grid references

Many Tasmanian agencies and research organisations have traditionally geocoded biological observations using a map name and 6-figure grid reference. The Tasmanian Herbarium uses a three-letter code to designate a 1:100,000-series map (for example, HUO for Huon), and the following 6 figures to designate the point location within a 50 m radius.

The full UTM grid reference can be calculated from the map code and 6-figure grid reference, and latitude and longitude derived from this:

* Most of Tasmania falls within UTM Zone 55 G, simplifying calculations.
* The missing coordinates from the grid reference easting and northing can be calculated using the starting points written in each map

## Base data

The file `mapinfo.csv` contains a matrix of map names, zone, zone alpha, as well as the lowest easting and northing coordinates for each map.

## Usage

`go get github.com/kurankat/tasgrid`

Create a `MapGrid` object to hold the necessary data for each map by calling `NewTasMapGrid()`. `MapGrid` is a map of `TasMap` objects, with the map's name as the index.

`var fullGrid = *tasgrid.NewTasMapGrid()`

You can then call `NewGridPoint` to generate a grid point with all its spatial information, and call its methods.

`gridPoint, err = tasgrid.NewGridPoint(mapname, easting, northing, fullGrid)`

To retrieve a point's full-length easting or northing (as zone 55 G) as a string:

`gridPoint.GetFullEasting()`

`gridPoint.GetFullNorthing()`

To retrieve a point's decimal latitude and longitude as strings:

`gridPoint.GetDecimalLat()`

`gridPoint.GetDecimalLong()`

To retrieve a point's latitude and longitude seconds as strings:

`gridPoint.GetLatSeconds()`

`gridPoint.GetLongSeconds()`
