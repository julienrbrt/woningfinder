package onshuis

import (
	"net/url"

	"github.com/woningfinder/woningfinder/internal/corporation"
)

// Info defines the static information about this Housing  Corporation
var Info = corporation.Corporation{
	APIEndpoint: &url.URL{Scheme: "https", Host: "mijn.onshuis.com", Path: "/apps/com.itris.klantportaal"},
	Name:        "OnsHuis",
	URL:         "https://mijn.onshuis.com",
	Cities: []corporation.City{
		{Name: "Enschede"},
		{Name: "Hengelo"},
	},
	SelectionMethod: []corporation.SelectionMethod{
		{
			Method: corporation.SelectionRandom,
		},
		{
			Method: corporation.SelectionRandom,
		},
	},
}
