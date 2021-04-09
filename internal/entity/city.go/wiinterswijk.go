package city

import "github.com/woningfinder/woningfinder/internal/entity"

// https://nl.wikipedia.org/wiki/Wijken_en_buurten_in_Winterswijk
var Winterswijk = entity.City{
	Name: "Winterswijk",
	District: []entity.CityDistrict{
		{Name: "Centrale deel"},
		{Name: "Noordoost"},
		{Name: "Noordwest"},
		{Name: "Zuidoost"},
		{Name: "Zuidwest"},
	},
}
