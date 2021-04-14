package city

import "github.com/woningfinder/woningfinder/internal/corporation"

// https://nl.wikipedia.org/wiki/Wijken_en_buurten_in_Bussum
var Bussum = corporation.City{
	Name: "Bussum",
	District: []corporation.CityDistrict{
		{Name: "Ooster"},
		{Name: "Wester"},
		{Name: "Midden"},
		{Name: "Boslaan"},
		{Name: "Dondersstraat"},
		{Name: "Bloemenbuurt"},
		{Name: "Verbindingslaan"},
		{Name: "Brink"},
		{Name: "Kom van Bieghel"},
		{Name: "Vondellaan"},
		{Name: "Bijlstraat"},
		{Name: "Cereslaan"},
		{Name: "Laarderwegkwartier"},
		{Name: "Batterijlaan"},
		{Name: "Lomanplein"},
		{Name: "Spiegelzicht"},
		{Name: "Bredius"},
		{Name: "Nijverheidswerf"},
		{Name: "Godelindebuurt"},
		{Name: "Waltherlaan"},
		{Name: "Prins Hendrikkwartie"},
		{Name: "Koedijk"},
		{Name: "Meijerkamp"},
		{Name: "Hooftlaan"},
		{Name: "Schimmellaan"},
		{Name: "Franse Kamp"},
	},
}