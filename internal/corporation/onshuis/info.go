// Package onshuis implements OnsHuis API
package onshuis

import "github.com/woningfinder/woningfinder/internal/corporation"

// Info contains the information of OnsHuis
var Info = corporation.Corporation{
	Name: "OnsHuis",
	URL:  "https://mijnonshuis.com",
	Cities: []corporation.City{
		{
			Name:   "Enschede",
			Region: "Overijssel",
		},
		{
			Name:   "Hengelo",
			Region: "Overijssel",
		},
	},
	SelectionMethod: []corporation.SelectionMethod{
		{
			Method: corporation.SelectionRandom,
		},
		{
			Method: corporation.SelectionRandom,
		},
	},
}
