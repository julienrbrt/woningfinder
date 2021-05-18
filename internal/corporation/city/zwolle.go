package city

import "github.com/woningfinder/woningfinder/internal/corporation"

// https://nl.wikipedia.org/wiki/Wijken_en_buurten_in_Zwolle
var Zwolle = corporation.City{
	Name: "Zwolle",
	District: []string{
		"Binnenstad",
		"Diezerpoort",
		"Wipstrik",
		"Assendorp",
		"Kamperpoort-Veerallee",
		"Poort van Zwolle",
		"Westenholte",
		"Stadshagen",
		"Holtenbroek",
		"Aalanden",
		"Vechtlanden",
		"Berkum",
		"Marsweteringlanden",
		"Schelle",
		"Ittersum",
		"Soestweteringlanden",
	},
}
