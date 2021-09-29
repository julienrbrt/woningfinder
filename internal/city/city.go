package city

import (
	"strings"
	"time"
)

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
	Mijdrecht.Name:            Mijdrecht,
	Linschoten.Name:           Linschoten,
	Odijk.Name:                Odijk,
	Doorn.Name:                Doorn,
	Oudewater.Name:            Oudewater,
	Vreeland.Name:             Vreeland,
	Denekamp.Name:             Denekamp,
	Ootmarsum.Name:            Ootmarsum,
	Vriezenveen.Name:          Vriezenveen,
	Houten.Name:               Houten,
}

// City defines a city where a housing corporation operates or when an house offer lies
type City struct {
	CreatedAt         time.Time           `pg:"default:now()" json:"-"`
	Name              string              `pg:",pk" json:"name"`
	District          []string            `pg:"-" json:"district,omitempty"`
	SuggestedDistrict map[string][]string `pg:"-" json:"suggested_district,omitempty"`
	Coordinates       []float64           `pg:"-" json:"coordinates,omitempty"` // as longitude latitude for mapbox
}

// Merge cities that are supposed to be the same but that housing corporation name differently
func (c *City) Merge() City {
	switch {
	case strings.Contains(c.Name, "Hengelo"):
		return Hengelo
	case strings.Contains(c.Name, "Winterswijk"):
		return Winterswijk
	case strings.EqualFold(c.Name, "s-Heerenberg"):
		return Heerenberg
	}

	return *c
}

// GetCoordinates gets the city coordinates
func GetCoordinates(name string) []float64 {
	city, ok := cityTable[name]
	if !ok {
		return nil
	}

	return city.Coordinates
}

// SuggestedCityDistrict permit to get city suggested districts
func SuggestedCityDistrict(name string) map[string][]string {
	city, ok := cityTable[name]
	if !ok || len(city.SuggestedDistrict) == 0 {
		return nil
	}

	return city.SuggestedDistrict
}
