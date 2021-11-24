package domijn

import (
	"net/url"

	"github.com/julienrbrt/woningfinder/internal/corporation"
	"github.com/julienrbrt/woningfinder/internal/corporation/city"
)

var Info = corporation.Corporation{
	APIEndpoint: &url.URL{Scheme: "https", Host: "www.domijn.nl"},
	Name:        "Domijn",
	URL:         "https://www.domijn.nl",
	Cities: []city.City{
		city.Enschede,
		city.Haaksbergen,
		city.Losser,
		city.Overdinkel,
		city.DeLutte,
	},
	SelectionMethod: []corporation.SelectionMethod{
		corporation.SelectionRandom,
	},
}
