package ikwilhuren

import (
	"net/url"

	"github.com/woningfinder/woningfinder/internal/corporation"
	"github.com/woningfinder/woningfinder/internal/corporation/city"
)

var Info = corporation.Corporation{
	APIEndpoint: &url.URL{Scheme: "https", Host: "mijn.ikwilhuren.nu"},
	Name:        "MVGM Wonen",
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
		city.Wierden,
		city.Hertogenbosch,
		city.Maastricht,
		city.Heerenveen,
		city.Deventer,
		city.Terborg,
		city.Haaften,
		city.Loosdrecht,
		city.Deurne,
		city.Soest,
		city.OudBeijerland,
		city.Uden,
		city.Doetinchem,
		city.HendrikIdoAmbacht,
		city.CapelleAanDenIJssel,
		city.Kortenhoef,
		city.Hoogeveen,
		city.Voorhout,
		city.Voorburg,
		city.Amersfoort,
		city.Dieren,
		city.Ede,
		city.Breda,
		city.Arnhem,
		city.Hoogvliet,
		city.Oss,
		city.Heemstede,
		city.Echt,
		city.Delft,
		city.Terneuzen,
		city.Roosendaal,
		city.Harderwijk,
		city.Hardenberg,
		city.Zoetermeer,
		city.Winschoten,
		city.Hilversum,
		city.Groesbeek,
		city.Naaldwijk,
		city.Heerlen,
		city.Sliedrecht,
		city.Leerbroek,
		city.Groningen,
		city.Kaatsheuvel,
		city.Prinsenbeek,
		city.Eindhoven,
		city.Papendrecht,
		city.Oegstgeest,
		city.Goes,
		city.DenHaag,
		city.Roermond,
		city.Leiden,
		city.Tilburg,
		city.Ermelo,
		city.Vlissingen,
		city.Gorinchem,
		city.Waalwijk,
		city.Rotterdam,
		city.Oldenzaal,
		city.Schiedam,
		city.Alkmaar,
		city.Blaricum,
		city.Best,
		city.Heerhugowaard,
		city.Borculo,
		city.Hulst,
		city.Voorst,
		city.Haarlem,
	},
	SelectionMethod: []corporation.SelectionMethod{
		corporation.SelectionFirstComeFirstServed,
	},
}
