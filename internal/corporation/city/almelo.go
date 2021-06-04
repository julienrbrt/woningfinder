package city

import "github.com/woningfinder/woningfinder/internal/corporation"

// https://nl.wikipedia.org/wiki/Wijken_en_buurten_in_Almelo
var Almelo = corporation.City{
	Name: "Almelo",
	District: []string{
		"Binnenstad",
		"De Riet",
		"Noorderkwartier",
		"Sluitersveld",
		"Aalderinkshoek",
		"Nieuwstraat-Kwartier",
		"Ossenkoppelerhoek",
		"Hofkamp",
		"Schelfhorst",
	},
}
