package woonburo

import (
	"net/url"

	"github.com/julienrbrt/woningfinder/internal/corporation"
	"github.com/julienrbrt/woningfinder/internal/corporation/city"
	"github.com/julienrbrt/woningfinder/internal/corporation/scheduler"
)

var AlmeloInfo = corporation.Corporation{
	APIEndpoint: &url.URL{Scheme: "https", Host: "www.woonburoalmelo", Path: "/wsAlmelo/WoningenModule/Service.asmx"},
	Name:        "Woonburo Almelo",
	URL:         "https://www.woonburoalmelo.nl",
	Cities: []city.City{
		city.Almelo,
		city.Hengelo,
		city.Aadorp,
		city.Bornerbroek,
		city.Denekamp,
		city.Haaksbergen,
		city.Wierden,
	},
	SelectionMethod: []corporation.SelectionMethod{
		corporation.SelectionRegistrationDate,
	},
	SelectionTime: scheduler.CreateSelectionTime(19, 0),
}
