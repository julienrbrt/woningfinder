package city

// https://nl.wikipedia.org/wiki/Wijken_en_buurten_in_Utrecht
var Utrecht = City{
	Name:        "Utrecht",
	Coordinates: []float64{5.110414, 52.095942},
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
	Name:        "Zeist",
	Coordinates: []float64{5.22785, 52.08944},
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
	Name:        "Bilthoven",
	Coordinates: []float64{5.20454, 52.129141},
}

var Bunnik = City{
	Name:        "Bunnik",
	Coordinates: []float64{5.194111, 52.066344},
}

var Nieuwegein = City{
	Name:        "Nieuwegein",
	Coordinates: []float64{5.083612, 52.036715},
}

var Maarssen = City{
	Name:        "Maarssen",
	Coordinates: []float64{5.040082, 52.138485},
}

var WijkBijDuurstede = City{
	Name:        "Wijk bij Duurstede",
	Coordinates: []float64{5.33743, 51.97516},
}

var DenDoler = City{
	Name:        "Den Dolder",
	Coordinates: []float64{4.36658, 51.633405},
}

var Maartensdijk = City{
	Name:        "Maartensdijk",
	Coordinates: []float64{5.17294, 52.15842},
}

var Baambrugge = City{
	Name:        "Baambrugge",
	Coordinates: []float64{4.990842, 52.247595},
}

var Wilnis = City{
	Name:        "Wilnis",
	Coordinates: []float64{4.899067, 52.195378},
}

var Woerden = City{
	Name:        "Woerden",
	Coordinates: []float64{4.891356, 52.086757},
}

var Vianen = City{
	Name:        "Vianen",
	Coordinates: []float64{5.09372, 51.989039},
}

var DeMeern = City{
	Name:        "De Meern",
	Coordinates: []float64{5.03658, 52.080582},
}

var Papekop = City{
	Name:        "Papekop",
	Coordinates: []float64{4.853489, 52.044843},
}

var Breukelen = City{
	Name:        "Breukelen",
	Coordinates: []float64{5.002549, 52.172152},
}

var DeBilt = City{
	Name:        "De Bilt",
	Coordinates: []float64{5.178706, 52.109701},
}

var DriebergenRijsenburg = City{
	Name:        "Driebergen-Rijsenburg",
	Coordinates: []float64{5.28267, 52.051922},
}

var IJsselstein = City{
	Name:        "IJsselstein",
	Coordinates: []float64{5.045879, 52.024638},
}

var Vleuten = City{
	Name:        "Vleuten",
	Coordinates: []float64{5.011092, 52.106979},
}

var Mijdrecht = City{
	Name:        "Mijdrecht",
	Coordinates: []float64{4.864817, 52.206018},
}

var Linschoten = City{
	Name:        "Linschoten",
	Coordinates: []float64{4.91286, 52.06433},
}

var Odijk = City{
	Name:        "Odijk",
	Coordinates: []float64{5.23439, 52.04947},
}

var Doorn = City{
	Name:        "Doorn",
	Coordinates: []float64{5.345864, 52.033299},
}

var Oudewater = City{
	Name:        "Oudewater",
	Coordinates: []float64{4.866287, 52.024527},
}

var Vreeland = City{
	Name:        "Vreeland",
	Coordinates: []float64{5.021718, 52.209893},
}

var Houten = City{
	Name:        "Houten",
	Coordinates: []float64{5.167245, 52.035909},
}

var Vinkeveen = City{
	Name:        "Vinkeveen",
	Coordinates: []float64{4.932163, 52.214764},
}

var Harmelen = City{
	Name:        "Harmelen",
	Coordinates: []float64{4.963882, 52.091534},
}

var Langbroek = City{
	Name:        "Langbroek",
	Coordinates: []float64{5.32597, 52.01114},
}

var Lopik = City{
	Name:        "Lopik",
	Coordinates: []float64{4.94915, 51.97528},
}

var Kockengen = City{
	Name:        "Kockengen",
	Coordinates: []float64{4.952222, 52.151667},
}

var Polsbroek = City{
	Name:        "Polsbroek",
	Coordinates: []float64{4.851944, 51.978056},
}
