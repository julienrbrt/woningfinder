package dewoonplaats

import (
	"net/url"

	"github.com/woningfinder/woningfinder/internal/corporation"
	"github.com/woningfinder/woningfinder/internal/corporation/city"
)

var Info = corporation.Corporation{
	APIEndpoint: &url.URL{Scheme: "https", Host: "www.dewoonplaats.nl", Path: "/wh_services"},
	Name:        "De Woonplaats",
	URL:         "https://www.dewoonplaats.nl",
	Cities: []corporation.City{
		city.Enschede,
		city.Zwolle,
		city.Dinxperlo,
		city.Winterswijk,
		city.Neede,
		city.Wehl,
		city.Aalten,
		city.Groenlo,
		city.Bussum,
		city.Bredevoort,
		city.Ulft,
		city.Hengelo,
		city.Almelo,
	},
	SelectionMethod: []corporation.SelectionMethod{
		corporation.SelectionRandom,
		corporation.SelectionFirstComeFirstServed,
	},
}
