// Package dewoonplaats implements De Woonplaats API
package dewoonplaats

import "github.com/woningfinder/woningfinder/internal/corporation"

// Info contains the information of De Woonplaats
var Info = corporation.Corporation{
	Name: "De Woonplaats",
	URL:  "https://dewoonplaats.nl",
	Cities: []corporation.City{
		{
			Name:   "Enschede",
			Region: "Overijssel",
		},
		{
			Name:   "Zwolle",
			Region: "Overijssel",
		},
		{
			Name:   "Aatlen",
			Region: "Gelderland",
		},
		{
			Name:   "Dinxperlo",
			Region: "Gelderland",
		},
		{
			Name:   "Winterswijk",
			Region: "Gelderland",
		},
		{
			Name:   "Neede",
			Region: "Gelderland",
		},
		{
			Name:   "Wehl",
			Region: "Gelderland",
		},
	},
	SelectionMethod: []corporation.SelectionMethod{
		{
			Method: corporation.SelectionRandom,
		},
		{
			Method: corporation.SelectionFirstComeFirstServed,
		},
	},
}
