package city

// https://nl.wikipedia.org/wiki/Wijken_en_buurten_in_Utrecht
var Utrecht = City{
	Name:      "Utrecht",
	Latitude:  52.095942,
	Longitude: 5.110414,
	SuggestedDistrict: map[string][]string{
		"West": {
			"Welgelegen",
			"Den Hommel",
			"Oog in Al",
			"Halve Maan",
			"Lombok-Oost",
			"Leidseweg en omgeving",
			"Lombok-West",
			"Laan van Nieuw Guinea-Spinozaplantsoen",
			"Nieuw Engeland",
			"Schepenbuurt",
			"Lage weide",
		},
		"Noordwest": {
			"Pijlsweerd-Zuid",
			"Pijlsweerd-Noord",
			"Ondiep",
			"Tweede Daalsebuurt",
			"Egelantierstraat-Mariëndaalstraat",
			"Geuzenwijk",
			"De Driehoek",
			"Julianapark",
			"Elinkwijk",
			"Prins Bernhardplein",
			"Schaakbuurt",
			"Queeckhovenplein",
			"Zuilen-Noord",
			"Loevenhoutsedijk",
		},
		"Overvecht": {
			"Taag- en Rubicondreef",
			"Zamenhofdreef",
			"Wolga- en Donaudreef",
			"Neckardreef",
			"Amazone- en Nicaraguadreef",
			"Zambesidreef",
			"Tigris- en Bostondreef",
			"Bedrijvengebied",
			"Poldergebied Overvecht",
		},
		"Noordoost": {
			"Vogelenbuurt",
			"Lauwerecht",
			"Staatsliedenbuurt",
			"Tuinwijk-West",
			"Tuinwijk-Oost",
			"Tuindorp en Van Lieflandlaan-West",
			"Tuindorp-Oost",
			"Huizingalaan K. Doormanlaan",
			"Zeeheldenbuurt",
			"Wittevrouwen",
			"Voordorp",
		},
		"Oost": {
			"Buiten Wittevrouwen",
			"Tolsteegsingel",
			"Sterrenwijk",
			"Watervogelbuurt",
			"Lodewijk Napoleonplantsoen",
			"Rubenslaan",
			"Schildersbuurt",
			"Abstede",
			"Oudwijk",
			"Wilhelminapark",
			"De Uithof",
			"Rijnsweerd",
			"Galgenwaard en Kromhoutkazerne",
			"Maarschalkerweerd en Mereveld",
		},
		"Binnenstad": {
			"Wijk C",
			"Breedstraat en Plompetorengracht",
			"Lange Elisabethstraat Mariaplaats",
			"Neude Janskerkhof en Domplein",
			"Nobelstraat",
			"Springweg",
			"Lange Nieuwstraat",
			"Nieuwegracht-Oost",
			"Bleekstraat",
			"Hooch Boulandt Moreelsepark",
			"Hoog-Catharijne CS en Leidseveer",
		},
		"Zuid": {
			"Lunetten-Noord",
			"Lunetten-Zuid",
			"Bokkenbuurt",
			"Tolsteeg en Rotsoord",
			"Oud Hoograven",
			"Nieuw Hoograven",
		},
		"Zuidwest": {
			"Dichterswijk",
			"Rivierenwijk",
			"Bedrijvengebied Kanaleneiland",
			"Transwijk-Zuid",
			"Transwijk-Noord",
			"Kanaleneiland-Zuid",
			"Kanaleneiland-Noord",
		},
		"Leidsche Rijn": {
			"Papendorp",
			"Bedrijventerrein De Wetering",
			"Terwijde",
			"Parkwijk, 't Zand",
			"Hogeweide",
			"Langerak",
			"Strijkviertel",
		},
	},
}

// https://nl.wikipedia.org/wiki/Wijken_en_buurten_in_Zeist
var Zeist = City{
	Name:      "Zeist",
	Latitude:  52.08944,
	Longitude: 5.22785,
	SuggestedDistrict: map[string][]string{
		"Centrum": {
			"Carre",
			"Centrumschil-Zuid",
			"Centrumschil-Noord",
			"Lyceumkwartier",
			"Het Slot",
		},
		"Zeist-Noord": {
			"Patijnpark",
			"Dijnselburg en Verzetswijk",
			"Staatsliedenkwartier",
			"Mooi Zeist, Sanatorium en omgeving",
			"Vollenhove",
			"Utrechtseweg",
		},
		"Zeist-West": {
			"Griffensteijn en Kersbergen",
			"Nijenheim",
			"Crosestein",
			"Vogelwijk",
			"Brugakker en De Clomp",
			"Couwenhoven",
			"Blikkenburg",
			"Wiedegebied",
		},
		"Zeist-Oost": {
			"Hoge Dennen",
			"Kerckebosch",
			"Driebergseweg",
			"Station NS",
			"Zeister Bos",
			"Austerlitz",
		},
	},
}

var Bilthoven = City{
	Name:      "Bilthoven",
	Latitude:  52.129141,
	Longitude: 5.20454,
}

var Bunnik = City{
	Name:      "Bunnik",
	Latitude:  52.066344,
	Longitude: 5.194111,
}

var Nieuwegein = City{
	Name:      "Nieuwegein",
	Latitude:  52.036715,
	Longitude: 5.083612,
}

var Maarssen = City{
	Name:      "Maarssen",
	Latitude:  52.138485,
	Longitude: 5.040082,
}

var WijkBijDuurstede = City{
	Name:      "Wijk bij Duurstede",
	Latitude:  51.97516,
	Longitude: 5.33743,
}

var DenDoler = City{
	Name:      "Den Dolder",
	Latitude:  51.633405,
	Longitude: 4.36658,
}

var Maartensdijk = City{
	Name:      "Maartensdijk",
	Latitude:  52.15842,
	Longitude: 5.17294,
}

var Baambrugge = City{
	Name:      "Baambrugge",
	Latitude:  52.247595,
	Longitude: 4.990842,
}

var Wilnis = City{
	Name:      "Wilnis",
	Latitude:  52.195378,
	Longitude: 4.899067,
}

var Woerden = City{
	Name:      "Woerden",
	Latitude:  52.086757,
	Longitude: 4.891356,
}

var Vianen = City{
	Name:      "Vianen",
	Latitude:  51.989039,
	Longitude: 5.09372,
}

var DeMeern = City{
	Name:      "De Meern",
	Latitude:  52.080582,
	Longitude: 5.03658,
}

var Papekop = City{
	Name:      "Papekop",
	Latitude:  52.044843,
	Longitude: 4.853489,
}

var Breukelen = City{
	Name:      "Breukelen",
	Latitude:  52.172152,
	Longitude: 5.002549,
}

var DeBilt = City{
	Name:      "De Bilt",
	Latitude:  52.109701,
	Longitude: 5.178706,
}

var DriebergenRijsenburg = City{
	Name:      "Driebergen-Rijsenburg",
	Latitude:  52.051922,
	Longitude: 5.28267,
}

var IJsselstein = City{
	Name:      "IJsselstein",
	Latitude:  52.024638,
	Longitude: 5.045879,
}

var Vleuten = City{
	Name:      "Vleuten",
	Latitude:  52.106979,
	Longitude: 5.011092,
}

var Mijdrecht = City{
	Name:      "Mijdrecht",
	Latitude:  52.206018,
	Longitude: 4.864817,
}

var Linschoten = City{
	Name:      "Linschoten",
	Latitude:  52.06433,
	Longitude: 4.91286,
}

var Odijk = City{
	Name:      "Odijk",
	Latitude:  52.04947,
	Longitude: 5.23439,
}

var Doorn = City{
	Name:      "Doorn",
	Latitude:  52.033299,
	Longitude: 5.345864,
}

var Oudewater = City{
	Name:      "Oudewater",
	Latitude:  52.024527,
	Longitude: 4.866287,
}

var Vreeland = City{
	Name:      "Vreeland",
	Latitude:  52.209893,
	Longitude: 5.021718,
}

var Houten = City{
	Name:      "Houten",
	Latitude:  52.035909,
	Longitude: 5.167245,
}

var Vinkeveen = City{
	Name:      "Vinkeveen",
	Latitude:  52.214764,
	Longitude: 4.932163,
}

var Harmelen = City{
	Name:      "Harmelen",
	Latitude:  52.091534,
	Longitude: 4.963882,
}

var Langbroek = City{
	Name:      "Langbroek",
	Latitude:  52.01114,
	Longitude: 5.32597,
}

var Lopik = City{
	Name:      "Lopik",
	Latitude:  51.97528,
	Longitude: 4.94915,
}

var Kockengen = City{
	Name:      "Kockengen",
	Latitude:  52.151667,
	Longitude: 4.952222,
}

var Polsbroek = City{
	Name:      "Polsbroek",
	Latitude:  51.978056,
	Longitude: 4.851944,
}

var Hagestein = City{
	Name:      "Hagestein",
	Latitude:  51.979722,
	Longitude: 5.120556,
}

var Leersum = City{
	Name:      "Leersum",
	Latitude:  52.010658,
	Longitude: 5.432013,
}

var Abcoude = City{
	Name:      "Abcoude",
	Latitude:  52.271248,
	Longitude: 4.974134,
}

var Maarn = City{
	Name:      "Maarn",
	Latitude:  52.064015,
	Longitude: 5.370910,
}

var Leerdam = City{
	Name:      "Leerdam",
	Latitude:  51.891664,
	Longitude: 5.093719,
}

var Kamerik = City{
	Name:      "Kamerik",
	Latitude:  52.111790,
	Longitude: 4.894307,
}

var Zegveld = City{
	Name:      "Zegveld",
	Latitude:  52.115317,
	Longitude: 4.837293,
}

var Cothen = City{
	Name:      "Cothen",
	Latitude:  51.995363,
	Longitude: 5.310461,
}

var Soest = City{
	Name:      "Soest",
	Latitude:  52.175035,
	Longitude: 5.291844,
}

var Soesterberg = City{
	Name:      "Soesterberg",
	Latitude:  52.119543,
	Longitude: 5.304264,
}

var Amersfoort = City{
	Name:      "Amersfoort",
	Latitude:  52.155581,
	Longitude: 5.389371,
}

var Leerbroek = City{
	Name:      "Leerbroek",
	Latitude:  51.916492,
	Longitude: 5.038154,
}

var Achterveld = City{
	Name:      "Achterveld",
	Latitude:  52.136229,
	Longitude: 5.497583,
}

var Elst = City{
	Name:      "Elst",
	Latitude:  51.985794,
	Longitude: 5.498410,
}

var Rhenen = City{
	Name:      "Rhenen",
	Latitude:  51.965837,
	Longitude: 5.571089,
}

var Veenendaal = City{
	Name:      "Veenendaal",
	Latitude:  52.026185,
	Longitude: 5.556038,
}

var Meerkerk = City{
	Name:      "Meerkerk",
	Latitude:  51.918506,
	Longitude: 4.993374,
}
