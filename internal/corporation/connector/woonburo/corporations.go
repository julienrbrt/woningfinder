package woonburo

import (
	"net/url"
	"time"

	"github.com/woningfinder/woningfinder/internal/corporation"
	"github.com/woningfinder/woningfinder/internal/corporation/city"
	"github.com/woningfinder/woningfinder/internal/corporation/scheduler"
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
	SelectionTime: []time.Time{
		scheduler.CreateSelectionTime(19, 0),
	},
}
