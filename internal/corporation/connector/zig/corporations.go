package zig

import (
	"net/url"

	"github.com/julienrbrt/woningfinder/internal/corporation"
	"github.com/julienrbrt/woningfinder/internal/corporation/city"
)

var RoomspotInfo = corporation.Corporation{
	APIEndpoint: &url.URL{Scheme: "https", Host: "www.roomspot.nl"},
	Name:        "Roomspot",
	URL:         "https://www.roomspot.nl",
	Cities: []city.City{
		city.Enschede,
		city.Hengelo,
	},
	SelectionMethod: []corporation.SelectionMethod{
		corporation.SelectionRandom,
		corporation.SelectionRegistrationDate,
	},
}

var DeWoningZoekerInfo = corporation.Corporation{
	APIEndpoint: &url.URL{Scheme: "https", Host: "www.dewoningzoeker.nl"},
	Name:        "DeWoningZoeker",
	URL:         "https://www.dewoningzoeker.nl",
	Cities: []city.City{
		city.Zwolle,
		city.IJsselmuiden,
		city.Kampen,
		city.Wanneperveen,
		city.Vollenhove,
		city.Wijthmen,
		city.Zalk,
		city.Willemsoord,
		city.Giethoorn,
		city.Blokzijl,
		city.Hasselt,
		city.Zwartsluis,
		city.Genemuiden,
	},
	SelectionMethod: []corporation.SelectionMethod{
		corporation.SelectionRandom,
		corporation.SelectionRegistrationDate,
	},
}

var WoonnetHaaglanden = corporation.Corporation{
	APIEndpoint: &url.URL{Scheme: "https", Host: "www.woonnet-haaglanden.nl"},
	Name:        "Woonnet Haaglanden",
	URL:         "https://www.woonnet-haaglanden.nl",
	Cities: []city.City{
		city.Delft,
		city.DenHaag,
		city.Voorburg,
		city.Rijswijk,
		city.Zoetermeer,
		city.Maasland,
		city.Naaldwijk,
		city.Monster,
		city.Wassenaar,
		city.DeLier,
		city.Wateringen,
		city.Leidschendam,
		city.Gravenzande,
		city.Pijnacker,
		city.Poeldijk,
		city.DenHoorn,
		city.Schipluiden,
		city.Honselersdijk,
		city.Kwintsheul,
		city.Nootdorp,
		city.Delfgauw,
	},
	SelectionMethod: []corporation.SelectionMethod{
		corporation.SelectionRandom,
		corporation.SelectionRegistrationDate,
		corporation.SelectionFirstComeFirstServed,
	},
}
