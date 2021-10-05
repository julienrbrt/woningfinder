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

var UtrechtInfo = corporation.Corporation{
	APIEndpoint: &url.URL{Scheme: "https", Host: "www.woningnetregioutrecht.nl", Path: "/webapi"},
	Name:        "WoningNet Utrecht",
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
		city.Oudewater,
		city.Vreeland,
		city.Houten,
		city.Vinkeveen,
		city.Harmelen,
	},
	SelectionMethod: []corporation.SelectionMethod{
		corporation.SelectionFirstComeFirstServed,
		corporation.SelectionRandom,
		corporation.SelectionRegistrationDate,
	},
}

var AmsterdamInfo = corporation.Corporation{
	APIEndpoint: &url.URL{Scheme: "https", Host: "www.woningnetregioamsterdam.nl", Path: "/webapi"},
	Name:        "WoningNet Stadsregio Amsterdam",
	URL:         "https://www.woningnetregioamsterdam.nl",
	Cities: []city.City{
		city.Amsterdam,
		city.Amstelveen,
		city.Zaandam,
	},
	SelectionMethod: []corporation.SelectionMethod{
		corporation.SelectionFirstComeFirstServed,
		corporation.SelectionRandom,
		corporation.SelectionRegistrationDate,
	},
}

var RooftrackInfo = corporation.Corporation{
	APIEndpoint: &url.URL{Scheme: "https", Host: "www.rooftrack.nl", Path: "/webapi"},
	Name:        "Rooftrack - Vrijesectorhuurwoningen",
	URL:         "https://www.rooftrack.nl",
	Cities: []city.City{
		city.Utrecht,
	},
	SelectionMethod: []corporation.SelectionMethod{
		corporation.SelectionFirstComeFirstServed,
	},
}
