package woningnet

import (
	"net/url"

	"github.com/woningfinder/woningfinder/internal/corporation"
	"github.com/woningfinder/woningfinder/internal/corporation/city"
)

var HengeloBorneInfo = corporation.Corporation{
	APIEndpoint: &url.URL{Scheme: "https", Host: "www.woningnethengeloborne.nl", Path: "/webapi"},
	Name:        "WoningNet Hengelo-Borne",
	URL:         "https://www.woningnethengeloborne.nl",
	Cities: []city.City{
		city.Hengelo,
		city.Borne,
		city.Hertme,
		city.Zenderen,
	},
	SelectionMethod: []corporation.SelectionMethod{
		corporation.SelectionFirstComeFirstServed,
		corporation.SelectionRandom,
		corporation.SelectionRegistrationDate,
	},
}
