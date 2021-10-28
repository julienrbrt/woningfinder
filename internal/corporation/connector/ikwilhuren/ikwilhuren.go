package ikwilhuren

import (
	"net/url"

	"github.com/woningfinder/woningfinder/internal/corporation"
	"github.com/woningfinder/woningfinder/internal/corporation/city"
)

var Info = corporation.Corporation{
	APIEndpoint: &url.URL{Scheme: "https", Host: "ikwilhuren.nu"},
	Name:        "ikwilhuren.nu",
	URL:         "https://ikwilhuren.nu",
	Cities: []city.City{
		city.Almelo,
		city.Amsterdam,
		city.Apeldoorn,
		city.Amstelveen,
		city.Borne,
		city.Hengelo,
		city.Utrecht,
		city.Enschede,
		city.Zwolle,
		city.Zeist,
		city.Zaandam,
		city.Woerden,
		city.WijkBijDuurstede,
		city.Vleuten,
		city.Vinkeveen,
		city.Purmerend,
		city.Oudewater,
		city.OuderkerkAanDeAmstel,
		city.Nieuwegein,
		city.Maarn,
		city.IJsselstein,
		city.Houten,
		city.Hoofddorp,
		city.DriebergenRijsenburg,
		city.Dinxperlo,
		city.Diemen,
		city.DeBilt,
		city.Bussum,
		city.Breukelen,
		city.Bilthoven,
		city.Badhoevedorp,
		city.Aalsmeer,
	},
	SelectionMethod: []corporation.SelectionMethod{
		corporation.SelectionFirstComeFirstServed,
	},
}
