package city

import "strings"

// cityTable defines the city name and it's corresponding struct
var cityTable = map[string]City{
	Aalten.Name:               Aalten,
	Dinxperlo.Name:            Dinxperlo,
	Bredevoort.Name:           Bredevoort,
	Neede.Name:                Neede,
	Borne.Name:                Borne,
	Bussum.Name:               Bussum,
	Wehl.Name:                 Wehl,
	Enschede.Name:             Enschede,
	Haaksbergen.Name:          Haaksbergen,
	Hengelo.Name:              Hengelo,
	Losser.Name:               Losser,
	Overdinkel.Name:           Overdinkel,
	DeLutte.Name:              DeLutte,
	Groenlo.Name:              Groenlo,
	Ulft.Name:                 Ulft,
	Winterswijk.Name:          Winterswijk,
	Zwolle.Name:               Zwolle,
	Almelo.Name:               Almelo,
	Hertme.Name:               Hertme,
	Apeldoorn.Name:            Apeldoorn,
	Heerenberg.Name:           Heerenberg,
	Zenderen.Name:             Zenderen,
	Utrecht.Name:              Utrecht,
	Zeist.Name:                Zeist,
	Bilthoven.Name:            Bilthoven,
	Bunnik.Name:               Bunnik,
	Nieuwegein.Name:           Nieuwegein,
	Maarssen.Name:             Maarssen,
	WijkBijDuurstede.Name:     WijkBijDuurstede,
	DenDoler.Name:             DenDoler,
	Maartensdijk.Name:         Maartensdijk,
	Baambrugge.Name:           Baambrugge,
	Wilnis.Name:               Wilnis,
	Woerden.Name:              Woerden,
	Vianen.Name:               Vianen,
	DeMeern.Name:              DeMeern,
	Papekop.Name:              Papekop,
	Breukelen.Name:            Breukelen,
	DeBilt.Name:               DeBilt,
	IJsselstein.Name:          IJsselstein,
	Vleuten.Name:              Vleuten,
	DriebergenRijsenburg.Name: DriebergenRijsenburg,
}

// Merge cities that are supposed to be the same but that housing corporation name differently
func MergeCity(city City) City {
	switch {
	case strings.Contains(city.Name, "Hengelo"):
		return Hengelo
	case strings.Contains(city.Name, "Winterswijk"):
		return Winterswijk
	case strings.EqualFold(city.Name, "s-Heerenberg"):
		return Heerenberg
	}

	return city
}

// HasSuggestedCityDistrict permit to see if a city has an city districts
// This method is used for checking if we support a city districts before making the request to Mapbox (see in the connector logic).
func HasSuggestedCityDistrict(cityName string) bool {
	city, ok := cityTable[cityName]

	return ok && len(city.District) > 0
}

func SuggestedCityDistrictFromName(cityName string) (map[string][]string, bool) {
	city, ok := cityTable[cityName]
	if !ok {
		return nil, false
	}

	return city.District, true
}
