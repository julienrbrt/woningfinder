package onshuis

import (
	"net/url"

	"github.com/woningfinder/woningfinder/internal/domain/entity"
)

// Info defines the static information about this Housing  Corporation
var Info = entity.Corporation{
	APIEndpoint: &url.URL{Scheme: "https", Host: "mijn.onshuis.com", Path: "/apps/com.itris.klantportaal"},
	Name:        "OnsHuis",
	URL:         "https://mijn.onshuis.com",
	Cities: []entity.City{
		{Name: "Enschede"},
		{Name: "Hengelo"},
	},
	SelectionMethod: []entity.SelectionMethod{
		{
			Method: entity.SelectionRandom,
		},
	},
}
