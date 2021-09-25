package city

// https://nl.wikipedia.org/wiki/Wijken_en_buurten_in_Enschede
// https://allecijfers.nl/gemeente-overzicht/enschede/
var Enschede = City{
	Name: "Enschede",
	SuggestedDistrict: map[string][]string{
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
}

// https://nl.wikipedia.org/wiki/Wijken_en_buurten_in_Borne
var Borne = City{
	Name: "Borne",
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
	SuggestedDistrict: map[string][]string{
		"Kern": {"Haaksbergen Kern-1",
			"Haaksbergen Kern-2",
			"Haaksbergen Kern-3",
			"Haaksbergen Kern-4",
		},
	},
}

// https://nl.wikipedia.org/wiki/Wijken_en_buurten_in_Hengelo
var Hengelo = City{
	Name: "Hengelo OV",
}

// https://nl.wikipedia.org/wiki/Wijken_en_buurten_in_Losser
var Losser = City{
	Name: "Losser",
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
}
