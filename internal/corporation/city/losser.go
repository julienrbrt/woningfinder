package city

import "github.com/woningfinder/woningfinder/internal/corporation"

// https://nl.wikipedia.org/wiki/Wijken_en_buurten_in_Losser
var Losser = corporation.City{
	Name: "Losser",
	District: []string{
		"Losser-Oost",
		"Losser-West",
	},
}

var Overdinkel = corporation.City{
	Name: "Overdinkel",
}
