package itris

import (
	"net/url"

	"github.com/woningfinder/woningfinder/internal/corporation"
	"github.com/woningfinder/woningfinder/internal/corporation/city"
)

var OnsHuisInfo = corporation.Corporation{
	APIEndpoint: &url.URL{Scheme: "https", Host: "mijn.onshuis.com"},
	Name:        "OnsHuis",
	URL:         "https://www.onshuis.com",
	Cities: []corporation.City{
		city.Enschede,
		city.Hengelo,
	},
	SelectionMethod: []corporation.SelectionMethod{
		corporation.SelectionRandom,
	},
}
