package city

import "github.com/woningfinder/woningfinder/internal/domain/entity"

var Enschede = entity.City{
	Name: "Enschede",
	District: []entity.CityDistrict{
		{Name: "Roombeek"},
	},
}
