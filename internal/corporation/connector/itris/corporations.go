package itris

import (
	"net/url"

	"github.com/woningfinder/woningfinder/internal/city"
	"github.com/woningfinder/woningfinder/internal/corporation"
)

var OnsHuisInfo = corporation.Corporation{
	APIEndpoint: &url.URL{Scheme: "https", Host: "mijn.onshuis.com"},
	Name:        "OnsHuis",
	URL:         "https://www.onshuis.com",
	Cities: []city.City{
		city.Enschede,
		city.Hengelo,
	},
	SelectionMethod: []corporation.SelectionMethod{
		corporation.SelectionRandom,
	},
}

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
