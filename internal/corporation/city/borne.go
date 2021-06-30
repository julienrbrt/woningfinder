package city

import "github.com/woningfinder/woningfinder/internal/corporation"

// https://nl.wikipedia.org/wiki/Wijken_en_buurten_in_Borne
var Borne = corporation.City{
	Name: "Borne",
	District: []string{
		"Bornsche Maten",
		"Centrum",
		"'t Wensink",
		"Dikkerslaan-Molenkampsweg",
		"Lettersveld",
		"Tichelkamp",
		"Stroom-Esch",
	},
}

var Hertme = corporation.City{
	Name: "Hertme",
}
