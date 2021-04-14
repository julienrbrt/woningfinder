package city

import "github.com/woningfinder/woningfinder/internal/corporation"

// https://nl.wikipedia.org/wiki/Wijken_en_buurten_in_Winterswijk
var Winterswijk = corporation.City{
	Name: "Winterswijk",
	District: []corporation.CityDistrict{
		{Name: "Centrale deel"},
		{Name: "Noordoost"},
		{Name: "Noordwest"},
		{Name: "Zuidoost"},
		{Name: "Zuidwest"},
	},
}
