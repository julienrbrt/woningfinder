package woningnet

import (
	"net/url"

	"github.com/woningfinder/woningfinder/internal/city"
	"github.com/woningfinder/woningfinder/internal/corporation"
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

var UtrechtInfo = corporation.Corporation{
	APIEndpoint: &url.URL{Scheme: "https", Host: "www.woningnetregioutrecht.nl", Path: "/webapi"},
	Name:        "WoningNet city",
	URL:         "https://www.woningnetregioutrecht.nl",
	Cities: []city.City{
		city.Utrecht,
		city.Zeist,
		city.Bilthoven,
		city.Bunnik,
		city.Nieuwegein,
		city.Maarssen,
		city.WijkBijDuurstede,
		city.DenDoler,
		city.Maartensdijk,
		city.Baambrugge,
		city.Wilnis,
		city.Woerden,
		city.Vianen,
		city.DeMeern,
		city.Papekop,
		city.Breukelen,
		city.DeBilt,
		city.IJsselstein,
		city.Vleuten,
		city.DriebergenRijsenburg,
		city.Mijdrecht,
		city.Linschoten,
		city.Odijk,
		city.Doorn,
	},
	SelectionMethod: []corporation.SelectionMethod{
		corporation.SelectionFirstComeFirstServed,
		corporation.SelectionRandom,
		corporation.SelectionRegistrationDate,
	},
}
