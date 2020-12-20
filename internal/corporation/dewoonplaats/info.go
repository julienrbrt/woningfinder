package dewoonplaats

import (
	"net/url"

	"github.com/woningfinder/woningfinder/internal/corporation"
)

// Info defines the static information about this Housing  Corporation
var Info = corporation.Corporation{
	APIEndpoint: &url.URL{Scheme: "https", Host: "www.dewoonplaats.nl", Path: "/wh_services"},
	Name:        "De Woonplaats",
	URL:         "https://dewoonplaats.nl",
	Cities: []corporation.City{
		{Name: "Enschede"},
		{Name: "Zwolle"},
		{Name: "Aatlen"},
		{Name: "Dinxperlo"},
		{Name: "Winterswijk"},
		{Name: "Neede"},
		{Name: "Wehl"},
	},
	SelectionMethod: []corporation.SelectionMethod{
		{
			Method: corporation.SelectionRandom,
		},
		{
			Method: corporation.SelectionFirstComeFirstServed,
		},
	},
}
