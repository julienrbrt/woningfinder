package city

import "github.com/woningfinder/woningfinder/internal/corporation"

// https://nl.wikipedia.org/wiki/Wijken_en_buurten_in_Haaksbergen
var Haaksbergen = corporation.City{
	Name: "Haaksbergen",
	District: []corporation.CityDistrict{
		{Name: "Haaksbergen Kern-1"},
		{Name: "Haaksbergen Kern-2"},
		{Name: "Haaksbergen Kern-3"},
		{Name: "Haaksbergen Kern-4"},
		{Name: "Veldmaat 1"},
		{Name: "Veldmaat 2"},
		{Name: "Leemdijk"},
		{Name: "Zienesch"},
		{Name: "de Pas"},
		{Name: "de Els"},
		{Name: "Wolferink"},
		{Name: "Hassinkbrink"},
	},
}
