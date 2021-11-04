package zig

import (
	"net/url"

	"github.com/woningfinder/woningfinder/internal/corporation"
	"github.com/woningfinder/woningfinder/internal/corporation/city"
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
	},
	SelectionMethod: []corporation.SelectionMethod{
		corporation.SelectionRandom,
		corporation.SelectionRegistrationDate,
	},
}
