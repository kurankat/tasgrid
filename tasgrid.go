package tasgrid

var mapNames = [...]string{
	"ART",
	"BAB",
	"BRE",
	"CIR",
	"CON",
	"DAV",
	"DEN",
	"DER",
	"ESK",
	"FOR",
	"FOS",
	"FRA",
	"FRE",
	"GEO",
	"GOO",
	"HEL",
	"HUO",
	"KIL",
	"KIN",
	"LAD",
	"LAK",
	"LIT",
	"MAR",
	"MEA",
	"MER",
	"NIN",
	"NIV",
	"OLD",
	"OLG",
	"PAT",
	"PAU",
	"PIE",
	"POR",
	"PRO",
	"SAN",
	"SEC",
	"SHA",
	"SOP",
	"SOR",
	"SPE",
	"STO",
	"SWA",
	"SWC",
	"TAB",
	"TAM",
	"THR",
	"TYE",
	"WED",
	"WEL"}

// TasMap holds map-unique information
type TasMap struct {
	zone          int
	alpha         string
	eastingStart  int
	northingStart int
}

var mapRow = map[string]int{}
