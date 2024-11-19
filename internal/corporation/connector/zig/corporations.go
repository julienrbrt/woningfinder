package zig

import (
	"net/url"

	"github.com/julienrbrt/woningfinder/internal/corporation"
	"github.com/julienrbrt/woningfinder/internal/corporation/city"
)

var WoningHurenInfo = corporation.Corporation{
	APIEndpoint: &url.URL{Scheme: "https", Host: "www.woninghuren.nl"},
	Name:        "WoningHuren",
	URL:         "https://www.woninghuren.nl",
	Cities: []city.City{
		city.Aadorp,
		city.Ootmarsum,
		city.Vriezenveen,
		city.Bornerbroek,
		city.Denekamp,
		city.Hengelo,
		city.Borne,
		city.Hertme,
		city.Zenderen,
		city.Overdinkel,
		city.Haaksbergen,
		city.Losser,
		city.DeLutte,
		city.Enschede,
		city.Zwolle,
		city.Dinxperlo,
		city.Winterswijk,
		city.Neede,
		city.Wehl,
		city.Aalten,
		city.Groenlo,
		city.Bussum,
		city.Bredevoort,
		city.Ulft,
		city.Almelo,
		city.Apeldoorn,
		city.Heerenberg,
		city.Oldenzaal,
		city.Enter,
		city.Vroomshoop,
		city.Rijssen,
		city.Nijverdal,
		city.Hellendoorn,
		city.Goor,
		city.Wierden,
	},
	SelectionMethod: []corporation.SelectionMethod{
		corporation.SelectionRandom,
		corporation.SelectionFirstComeFirstServed,
	},
}

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
		city.Scheerwolde,
		city.Windesheim,
		city.Nieuwleusen,
		city.Ommen,
		city.Heeten,
		city.Wijhe,
		city.Dedemsvaart,
		city.Hardenberg,
		city.Slagharen,
		city.Raalte,
		city.Staphorst,
		city.Dalfsen,
	},
	SelectionMethod: []corporation.SelectionMethod{
		corporation.SelectionRandom,
		corporation.SelectionRegistrationDate,
	},
}

var MercatusInfo = corporation.Corporation{
	APIEndpoint: &url.URL{Scheme: "https", Host: "woningaanbod.mercatus.nl"},
	Name:        "Mercatus",
	URL:         "https://woningaanbod.mercatus.nl",
	Cities: []city.City{
		city.Bant,
		city.Emmeloord,
		city.Nagele,
		city.Kraggenburg,
		city.Luttelgeest,
		city.Ens,
		city.Tollebeek,
	},
	SelectionMethod: []corporation.SelectionMethod{
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
		city.Maasdijk,
	},
	SelectionMethod: []corporation.SelectionMethod{
		corporation.SelectionRandom,
		corporation.SelectionRegistrationDate,
		corporation.SelectionFirstComeFirstServed,
	},
}
