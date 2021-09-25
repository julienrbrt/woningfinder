package city

// https://nl.wikipedia.org/wiki/Wijken_en_buurten_in_Utrecht
var Utrecht = City{
	Name: "Utrecht",
	District: map[string][]string{
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
		"Vleuten-De Meern ": nil,
	},
}

// https://nl.wikipedia.org/wiki/Wijken_en_buurten_in_Zeist
var Zeist = City{
	Name: "Zeist",
	District: map[string][]string{
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
			"	Nijenheim",
			"	Crosestein",
			"	Vogelwijk",
			"	Brugakker en De Clomp",
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

// https://nl.wikipedia.org/wiki/Wijken_en_buurten_in_De_Bilt
var Bilthoven = City{
	Name: "Bilthoven",
}

var Bunnik = City{
	Name: "Bunnik",
}

// https://nl.wikipedia.org/wiki/Wijken_en_buurten_in_Nieuwegein
var Nieuwegein = City{
	Name: "Nieuwegein",
}

// https://allecijfers.nl/gemeente-overzicht/stichtse-vecht/
var Maarssen = City{
	Name: "Maarssen",
}

// https://allecijfers.nl/gemeente-overzicht/stichtse-vecht/
var WijkBijDuurstede = City{
	Name: "Wijk bij Duurstede",
}

var DenDoler = City{
	Name: "Den Dolder",
}

var Maartensdijk = City{
	Name: "Maartensdijk",
}

var Baambrugge = City{
	Name: "Baambrugge",
}

var Wilnis = City{
	Name: "Wilnis",
}

var Woerden = City{
	Name: "Woerden",
}

var Vianen = City{
	Name: "Vianen",
}

var DeMeern = City{
	Name: "De Meern",
}

var Papekop = City{
	Name: "Papekop",
}

var Breukelen = City{
	Name: "Breukelen",
}

var DeBilt = City{
	Name: "De Bilt",
}

var DriebergenRijsenburg = City{
	Name: "Driebergen-Rijsenburg",
}

var IJsselstein = City{
	Name: "IJsselstein",
}

var Vleuten = City{
	Name: "Vleuten",
}
