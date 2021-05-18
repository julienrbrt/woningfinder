package city

import "github.com/woningfinder/woningfinder/internal/corporation"

// https://nl.wikipedia.org/wiki/Wijken_en_buurten_in_Bussum
var Bussum = corporation.City{
	Name: "Bussum",
	District: []string{
		"Centrum",
		"Brediuskwartier",
		"Eng",
		"Spiegel",
	},
}
