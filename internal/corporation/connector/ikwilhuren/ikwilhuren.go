package ikwilhuren

import (
	"net/url"

	"github.com/julienrbrt/woningfinder/internal/corporation"
	"github.com/julienrbrt/woningfinder/internal/corporation/city"
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
		city.Wijchen,
		city.Veenendaal,
		city.Valkenswaard,
		city.Barendrecht,
		city.Oisterwijk,
		city.AlphenAanDenRijn,
		city.Lelystad,
		city.Zeewolde,
		city.Schoonebeek,
		city.BerkelEnschot,
		city.Veldhoven,
		city.Vught,
		city.Montfoort,
		city.Oosterhout,
		city.Bleiswijk,
		city.HoekVanHolland,
		city.Emmen,
		city.Nuenen,
		city.Pijnacker,
		city.Almere,
		city.Zwijndrecht,
		city.Maassluis,
		city.Geldermalsen,
		city.Dordrecht,
		city.Weesp,
		city.Rosmalen,
		city.Goirle,
		city.Nijmegen,
		city.Dongen,
		city.Rhoon,
		city.Losser,
		city.Gorredijk,
		city.Haren,
		city.Leusden,
		city.Westervoort,
		city.Bergschenhoek,
		city.Veghel,
		city.Beilen,
		city.Ridderkerk,
		city.Schijndel,
		city.Raalte,
		city.Brielle,
		city.Meerkerk,
		city.BergenOpZoom,
		city.Wassenaar,
		city.BerkelEnRodenrijs,
		city.Holten,
		city.Rijssen,
		city.Waddinxveen,
		city.Didam,
		city.Duiven,
		city.Leeuwarden,
		city.KrimpenAanDenIJssel,
		city.Lichtenvoorde,
		city.Elburg,
		city.Vroomshoop,
		city.Assen,
		city.Wormerveer,
		city.Huizen,
		city.DeMeern,
		city.KrimpenAanDenIJssel,
		city.Hilvarenbeek,
		city.Burgum,
		city.Rijen,
		city.Odijk,
		city.Uithoorn,
		city.Leerdam,
		city.Culemborg,
		city.Lent,
		city.Hellevoetsluis,
		city.Colmschate,
		city.Leerdam,
		city.CapelleAanDenIJssel,
		city.Helmond,
		city.DenHam,
	},
	SelectionMethod: []corporation.SelectionMethod{
		corporation.SelectionFirstComeFirstServed,
		corporation.SelectionRandom,
	},
}
