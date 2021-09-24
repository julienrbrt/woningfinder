package city

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
