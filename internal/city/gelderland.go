package city

// https://nl.wikipedia.org/wiki/Wijken_en_buurten_in_Aalten
var Aalten = City{
	Name: "Aalten",
}

var Dinxperlo = City{
	Name: "Dinxperlo",
}

var Bredevoort = City{
	Name: "Bredevoort",
}

// https://nl.wikipedia.org/wiki/Wijken_en_buurten_in_Apeldoorn
var Apeldoorn = City{
	Name: "Apeldoorn",
	District: map[string][]string{
		"Centrum": nil,
		"West":    nil,
		"Zuid":    nil,
		"Oost":    nil,
		"Noord":   nil,
	},
}

// https://nl.wikipedia.org/wiki/Wijken_en_buurten_in_Berkelland
var Neede = City{
	Name: "Neede",
}

// https://nl.wikipedia.org/wiki/Wijken_en_buurten_in_Doetinchem
var Wehl = City{
	Name: "Wehl",
}

// https://nl.wikipedia.org/wiki/Wijken_en_buurten_in_Montferland
var Heerenberg = City{
	Name: "'s-Heerenberg",
}

// https://nl.wikipedia.org/wiki/Wijken_en_buurten_in_Oost_Gelre
var Groenlo = City{
	Name: "Groenlo",
}

// https://nl.wikipedia.org/wiki/Wijken_en_buurten_in_Oude_IJsselstreek
var Ulft = City{
	Name: "Ulft",
}

// https://nl.wikipedia.org/wiki/Wijken_en_buurten_in_Winterswijk
var Winterswijk = City{
	Name: "Winterswijk",
	District: map[string][]string{
		"Centrale deel": nil,
		"Noordoost":     nil,
		"Noordwest":     nil,
		"Zuidoost":      nil,
		"Zuidwest":      nil,
	},
}
