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
		city.Langbroek,
		city.Lopik,
		city.Kockengen,
		city.Polsbroek,
		city.Hagestein,
		city.Leersum,
		city.Abcoude,
		city.Maarn,
		city.Leerdam,
		city.Kamerik,
		city.Zegveld,
		city.Cothen,
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
		city.Aalsmeer,
		city.Diemen,
		city.Zaandam,
		city.Hoofddorp,
		city.Krommenie,
		city.Kudelstaart,
		city.Landsmeer,
		city.Marken,
		city.NieuwVennep,
		city.Oostzaan,
		city.OuderkerkAanDeAmstel,
		city.Purmerend,
		city.Uithoorn,
		city.Vijfhuizen,
		city.Wormer,
		city.Zwanenburg,
		city.Badhoevedorp,
		city.Zaandijk,
		city.Zuidoostbeemster,
		city.DeKwakel,
		city.Lisserbroek,
		city.Purmerland,
		city.KoogaandeZaan,
		city.Assendelft,
		city.Wormerveer,
		city.Jisp,
		city.Rijsenhout,
		city.Middenbeemster,
		city.Monnickendam,
		city.Beinsdorp,
		city.Westzaan,
		city.Cruquius,
	},
	SelectionMethod: []corporation.SelectionMethod{
		corporation.SelectionFirstComeFirstServed,
		corporation.SelectionRandom,
		corporation.SelectionRegistrationDate,
	},
}

var AlmereInfo = corporation.Corporation{
	APIEndpoint: &url.URL{Scheme: "https", Host: "www.woningnetalmere.nl", Path: "/webapi"},
	Name:        "WoningNet Almere",
	URL:         "https://www.woningnetalmere.nl",
	Cities: []city.City{
		city.Almere,
	},
	SelectionMethod: []corporation.SelectionMethod{
		corporation.SelectionRandom,
		corporation.SelectionRegistrationDate,
	},
}

var WoonkeusInfo = corporation.Corporation{
	APIEndpoint: &url.URL{Scheme: "https", Host: "www.woonkeus.nl", Path: "/webapi"},
	Name:        "WoningNet Woonkeus Drechtsteden",
	URL:         "https://www.woonkeus.nl",
	Cities: []city.City{
		city.Alblasserdam,
		city.Dordrecht,
		city.Papendrecht,
		city.Sliedrecht,
		city.Zwijndrecht,
	},
	SelectionMethod: []corporation.SelectionMethod{
		corporation.SelectionRandom,
		corporation.SelectionRegistrationDate,
	},
}

var EemvalleiInfo = corporation.Corporation{
	APIEndpoint: &url.URL{Scheme: "https", Host: "www.woningneteemvallei.nl", Path: "/webapi"},
	Name:        "WoningNet Eemvallei",
	URL:         "https://www.woningneteemvallei.nl",
	Cities: []city.City{
		city.Achterveld,
		city.Amersfoort,
		city.Nijkerk,
		city.Nijkerkerveen,
		city.Soest,
		city.Soesterberg,
	},
	SelectionMethod: []corporation.SelectionMethod{
		corporation.SelectionRandom,
		corporation.SelectionRegistrationDate,
		corporation.SelectionFirstComeFirstServed,
	},
}

var WoonserviceInfo = corporation.Corporation{
	APIEndpoint: &url.URL{Scheme: "https", Host: "www.mijnwoonservice.nl", Path: "/webapi"},
	Name:        "WoningNet Woonservice",
	URL:         "https://www.mijnwoonservice.nl",
	Cities: []city.City{
		city.Beverwijk,
		city.Haarlem,
		city.Heemstede,
	},
	SelectionMethod: []corporation.SelectionMethod{
		corporation.SelectionRegistrationDate,
		corporation.SelectionFirstComeFirstServed,
	},
}

var MercatusInfo = corporation.Corporation{
	APIEndpoint: &url.URL{Scheme: "https", Host: "www.mercatuswoningaanbod.nl", Path: "/webapi"},
	Name:        "WoningNet Mercatus",
	URL:         "https://www.mercatuswoningaanbod.nl",
	Cities: []city.City{
		city.Bant,
		city.Emmeloord,
		city.Nagele,
	},
	SelectionMethod: []corporation.SelectionMethod{
		corporation.SelectionRegistrationDate,
	},
}

var MiddenHollandInfo = corporation.Corporation{
	APIEndpoint: &url.URL{Scheme: "https", Host: "www.woningnetregiomiddenholland.nl", Path: "/webapi"},
	Name:        "WoningNet Midden-Holland",
	URL:         "https://www.woningnetregiomiddenholland.nl",
	Cities: []city.City{
		city.Bodegraven,
		city.Gouda,
		city.Waddinxveen,
		city.Zevenhuizen,
	},
	SelectionMethod: []corporation.SelectionMethod{
		corporation.SelectionRegistrationDate,
	},
}

var GroningenInfo = corporation.Corporation{
	APIEndpoint: &url.URL{Scheme: "https", Host: "www.woningnetgroningen.nl", Path: "/webapi"},
	Name:        "WoningNet Groningen",
	URL:         "https://www.woningnetgroningen.nl",
	Cities: []city.City{
		city.Groningen,
	},
	SelectionMethod: []corporation.SelectionMethod{
		corporation.SelectionRegistrationDate,
		corporation.SelectionRandom,
	},
}

var BovenGroningenInfo = corporation.Corporation{
	APIEndpoint: &url.URL{Scheme: "https", Host: "www.woningnetbovengroningen.nl", Path: "/webapi"},
	Name:        "WoningNet Boven Groningen",
	URL:         "https://www.woningnetbovengroningen.nl",
	Cities: []city.City{
		city.Eenrum,
		city.Ezinge,
		city.Kloosterburen,
		city.Winsum,
	},
	SelectionMethod: []corporation.SelectionMethod{
		corporation.SelectionRegistrationDate,
	},
}

var GooiVechtstreekInfo = corporation.Corporation{
	APIEndpoint: &url.URL{Scheme: "https", Host: "www.woningnetgooienvechtstreek.nl", Path: "/webapi"},
	Name:        "WoningNet Gooi en Vechtstreek",
	URL:         "https://www.woningnetgooienvechtstreek.nl",
	Cities: []city.City{
		city.Graveland,
		city.Ankeveen,
		city.Blaricum,
		city.Bussum,
		city.Hilversum,
		city.Huizen,
		city.Kortenhoef,
		city.Muiden,
	},
	SelectionMethod: []corporation.SelectionMethod{
		corporation.SelectionRegistrationDate,
	},
}

var HuiswaartsInfo = corporation.Corporation{
	APIEndpoint: &url.URL{Scheme: "https", Host: "www.huiswaarts.nu", Path: "/webapi"},
	Name:        "WoningNet Huiswaarts",
	URL:         "https://www.huiswaarts.nu",
	Cities: []city.City{
		city.Barneveld,
		city.Bennekom,
		city.Ede,
		city.Elst,
		city.Lunteren,
		city.Rhenen,
		city.Scherpenzeel,
		city.Terschuur,
		city.Veenendaal,
		city.Wageningen,
		city.Wekerom,
	},
	SelectionMethod: []corporation.SelectionMethod{
		corporation.SelectionRandom,
		corporation.SelectionRegistrationDate,
		corporation.SelectionFirstComeFirstServed,
	},
}

var WoongaardInfo = corporation.Corporation{
	APIEndpoint: &url.URL{Scheme: "https", Host: "www.woongaard.com", Path: "/webapi"},
	Name:        "WoningNet Woongaard",
	URL:         "https://www.woongaard.com",
	Cities: []city.City{
		city.Arkel,
		city.BenedenLeeuwen,
		city.Culemborg,
		city.Est,
		city.Geldermalsen,
		city.Gorinchem,
		city.HardinxveldGiessendam,
		city.Hedel,
		city.Leerbroek,
		city.Leerdam,
		city.Maurik,
		city.Meerkerk,
		city.Opheusden,
		city.Rossum,
		city.Tiel,
		city.Veen,
		city.WijkEnAalburg,
		city.Zaltbommel,
	},
	SelectionMethod: []corporation.SelectionMethod{
		corporation.SelectionRandom,
		corporation.SelectionRegistrationDate,
		corporation.SelectionFirstComeFirstServed,
	},
}
