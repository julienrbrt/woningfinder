package city

import "github.com/woningfinder/woningfinder/internal/corporation"

// https://nl.wikipedia.org/wiki/Wijken_en_buurten_in_Hengelo
var Hengelo = corporation.City{
	Name: "Hengelo OV",
	District: []string{
		"Binnenstad",
		"Hengelose Es",
		"Noord",
		"Hasseler Es",
		"Groot Driene",
		"Berflo Es",
		"Wilderinkshoek",
		"Woolde",
		"Slangenbeek",
	},
}
