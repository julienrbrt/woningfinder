package city

import "github.com/woningfinder/woningfinder/internal/corporation"

// https://nl.wikipedia.org/wiki/Wijken_en_buurten_in_Winterswijk
var Winterswijk = corporation.City{
	Name: "Winterswijk",
	District: []string{
		"Centrale deel",
		"Noordoost",
		"Noordwest",
		"Zuidoost",
		"Zuidwest",
	},
}
