package dewoonplaats

import (
	"net/url"

	"github.com/woningfinder/woningfinder/internal/corporation/scheduler"

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
		{Name: "Aalten"},
		{Name: "Groenlo"},
	},
	SelectionMethod: []entity.SelectionMethod{
		{
			Method: entity.SelectionRandom,
		},
		{
			Method: entity.SelectionFirstComeFirstServed,
		},
	},
	SelectionTime: scheduler.CreateSelectionTime(18, 0, 0),
}
