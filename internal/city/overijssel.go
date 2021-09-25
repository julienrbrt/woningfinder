package city

// https://nl.wikipedia.org/wiki/Wijken_en_buurten_in_Enschede
// https://allecijfers.nl/gemeente-overzicht/enschede/
var Enschede = City{
	Name: "Enschede",
	District: map[string][]string{
		"Centrum": {
			"City",
			"Lasonder, Zeggelt",
			"De Laares",
			"De Bothoven",
			"Hogeland",
			"Getfert",
			"Veldkamp",
			"Stadsweide",
			"Boddenkamp",
		},
		"Hogeland - Velve": {
			"Velve-Lindenhof",
			"Wooldrik",
			"Hogeland",
			"Varvik-Diekman",
			"Sleutelkamp",
			"'t Weldink",
			"De Leuriks",
		},
		"Boswinkel - Stadsveld": {
			"Cromhoffsbleek",
			"Boswinkel",
			"Pathmos",
			"Stevenfenne",
			"Stadsveld",
			"Elferink-Heuwkamp",
			"'t Zwering",
			"Ruwenbos",
		},
		"Twekkelerveld": {
			"Tubantia-Toekomst",
			"Twekkelerveld",
		},
		"Enschede-Noord": {
			"Walhof-Roessingh",
			"Bolhaar",
			"Roombeek-Roomveldje",
			"Mekkelholt",
			"Deppenbroek",
			"Voortman-Amelink",
			"Drienerveld-U.T.",
		},
		"Ribbelt - Stokhorst": {
			"Schreurserve",
			"Ribbelt",
			"Stokhorst",
		},
		"Enschede-Zuid": {
			"Wesselerbrink",
			"Stroinkslanden",
			"Helmerhoek",
		},
		"Glanerbrug": {
			"Glanerveld",
			"Bentveld-Bultserve",
			"Schipholt-Glanermaten",
			"Eekmaat",
			"Oikos",
			"Dolphia",
		}},
}

// https://nl.wikipedia.org/wiki/Wijken_en_buurten_in_Almelo
var Almelo = City{
	Name: "Almelo",
	District: map[string][]string{
		"Binnenstad":           nil,
		"De Riet":              nil,
		"Noorderkwartier":      nil,
		"Sluitersveld":         nil,
		"Aalderinkshoek":       nil,
		"Nieuwstraat-Kwartier": nil,
		"Ossenkoppelerhoek":    nil,
		"Hofkamp":              nil,
		"Schelfhorst":          nil,
	},
}

// https://nl.wikipedia.org/wiki/Wijken_en_buurten_in_Borne
var Borne = City{
	Name: "Borne",
	District: map[string][]string{
		"Bornsche Maten":            nil,
		"Centrum":                   nil,
		"'t Wensink":                nil,
		"Dikkerslaan-Molenkampsweg": nil,
		"Lettersveld":               nil,
		"Tichelkamp":                nil,
		"Stroom-Esch":               nil,
	},
}

var Hertme = City{
	Name: "Hertme",
}

var Zenderen = City{
	Name: "Zenderen",
}

// https://nl.wikipedia.org/wiki/Wijken_en_buurten_in_Haaksbergen
var Haaksbergen = City{
	Name: "Haaksbergen",
	District: map[string][]string{
		"Kern": {"Haaksbergen Kern-1",
			"Haaksbergen Kern-2",
			"Haaksbergen Kern-3",
			"Haaksbergen Kern-4",
		},
		"Veldmaat 1":   nil,
		"Veldmaat 2":   nil,
		"Leemdijk":     nil,
		"Zienesch":     nil,
		"de Pas":       nil,
		"de Els":       nil,
		"Wolferink":    nil,
		"Hassinkbrink": nil,
	},
}

// https://nl.wikipedia.org/wiki/Wijken_en_buurten_in_Hengelo
var Hengelo = City{
	Name: "Hengelo OV",
	District: map[string][]string{
		"Binnenstad":     nil,
		"Hengelose Es":   nil,
		"Noord":          nil,
		"Hasseler Es":    nil,
		"Groot Driene":   nil,
		"Berflo Es":      nil,
		"Wilderinkshoek": nil,
		"Woolde":         nil,
		"Slangenbeek":    nil,
	},
}

// https://nl.wikipedia.org/wiki/Wijken_en_buurten_in_Losser
var Losser = City{
	Name: "Losser",
	District: map[string][]string{
		"Losser-Oost": nil,
		"Losser-West": nil,
	},
}

var Overdinkel = City{
	Name: "Overdinkel",
}

var DeLutte = City{
	Name: "De Lutte",
}

// https://nl.wikipedia.org/wiki/Wijken_en_buurten_in_Zwolle
var Zwolle = City{
	Name: "Zwolle",
	District: map[string][]string{
		"Binnenstad":            nil,
		"Diezerpoort":           nil,
		"Wipstrik":              nil,
		"Assendorp":             nil,
		"Kamperpoort-Veerallee": nil,
		"Poort van Zwolle":      nil,
		"Westenholte":           nil,
		"Stadshagen":            nil,
		"Holtenbroek":           nil,
		"Aalanden":              nil,
		"Vechtlanden":           nil,
		"Berkum":                nil,
		"Marsweteringlanden":    nil,
		"Schelle":               nil,
		"Ittersum":              nil,
		"Soestweteringlanden":   nil,
	},
}
