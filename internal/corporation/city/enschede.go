package city

// https://nl.wikipedia.org/wiki/Wijken_en_buurten_in_Enschede
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
