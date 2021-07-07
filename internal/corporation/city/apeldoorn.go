package city

import "github.com/woningfinder/woningfinder/internal/corporation"

// https://nl.wikipedia.org/wiki/Wijken_en_buurten_in_Apeldoorn
var Apeldoorn = corporation.City{
	Name: "Apeldoorn",
	District: []string{
		"Centrum",
		"West",
		"Zuid",
		"Oost",
		"Noord",
	},
}
