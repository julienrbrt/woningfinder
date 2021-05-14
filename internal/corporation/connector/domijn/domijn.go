package domijn

import (
	"net/url"

	"github.com/woningfinder/woningfinder/internal/corporation"
	"github.com/woningfinder/woningfinder/internal/corporation/city"
)

var Info = corporation.Corporation{
	APIEndpoint: &url.URL{Scheme: "https", Host: "www.domijn.nl"},
	Name:        "Domijn",
	URL:         "https://www.domijn.nl",
	Cities: []corporation.City{
		city.Enschede,
		city.Haaksbergen,
		city.Losser,
		city.Overdinkel,
	},
	SelectionMethod: []corporation.SelectionMethod{
		corporation.SelectionRandom,
	},
}
