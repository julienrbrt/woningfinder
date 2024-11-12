package itris

import (
	"net/url"

	"github.com/julienrbrt/woningfinder/internal/corporation"
	"github.com/julienrbrt/woningfinder/internal/corporation/city"
)

var MijandeWonenInfo = corporation.Corporation{
	APIEndpoint: &url.URL{Scheme: "https", Host: "mijn.mijande.nl"},
	Name:        "Mijande Wonen",
	URL:         "https://www.mijande.nl",
	Cities: []city.City{
		city.Denekamp,
		city.Ootmarsum,
		city.Vriezenveen,
	},
	SelectionMethod: []corporation.SelectionMethod{
		corporation.SelectionRegistrationDate,
	},
}
