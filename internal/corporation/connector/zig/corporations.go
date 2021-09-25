package zig

import (
	"net/url"

	"github.com/woningfinder/woningfinder/internal/city"
	"github.com/woningfinder/woningfinder/internal/corporation"
)

var RoomspotInfo = corporation.Corporation{
	APIEndpoint: &url.URL{Scheme: "https", Host: "www.roomspot.nl"},
	Name:        "Roomspot",
	URL:         "https://www.roomspot.nl",
	Cities: []city.City{
		city.Enschede,
		city.Hengelo,
	},
	SelectionMethod: []corporation.SelectionMethod{
		corporation.SelectionRandom,
		corporation.SelectionRegistrationDate,
	},
}
