package dewoonplaats

import (
	"net/url"
	"time"

	"github.com/woningfinder/woningfinder/internal/domain/entity"
)

// Info defines the static information about this Housing  Corporation
var Info = entity.Corporation{
	APIEndpoint: &url.URL{Scheme: "https", Host: "www.dewoonplaats.nl", Path: "/wh_services"},
	Name:        "De Woonplaats",
	URL:         "https://dewoonplaats.nl",
	Cities: []entity.City{
		{Name: "Enschede"},
		{Name: "Zwolle"},
		{Name: "Aatlen"},
		{Name: "Dinxperlo"},
		{Name: "Winterswijk"},
		{Name: "Neede"},
		{Name: "Wehl"},
	},
	SelectionMethod: []entity.SelectionMethod{
		{
			Method: entity.SelectionRandom,
		},
		{
			Method: entity.SelectionFirstComeFirstServed,
		},
	},
	SelectionTime: time.Date(2021, time.January, 1, 18, 0, 0, 0, time.UTC),
}
