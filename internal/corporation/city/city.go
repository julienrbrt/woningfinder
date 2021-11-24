package city

import (
	"strings"
	"time"
)

// City defines a city where a housing corporation operates or when an house offer lies
type City struct {
	CreatedAt         time.Time           `pg:"default:now()" json:"-"`
	Name              string              `pg:",pk" json:"name"`
	Latitude          float64             `json:"latitude,omitempty"`
	Longitude         float64             `json:"longitude,omitempty"`
	District          []string            `pg:"-" json:"district,omitempty"`
	SuggestedDistrict map[string][]string `pg:"-" json:"suggested_district,omitempty"`
}

// Merge cities that are supposed to be the same but that housing corporation name differently
func (c *City) Merge() City {
	switch {
	case strings.Contains(c.Name, "Amsterdam"):
		return Amsterdam
	case strings.Contains(c.Name, "Hengelo"):
		return Hengelo
	case strings.Contains(c.Name, "Winterswijk"):
		return Winterswijk
	case strings.EqualFold(c.Name, "s-Heerenberg"):
		return Heerenberg
	case strings.Contains(c.Name, "IJsselstein"):
		return IJsselstein
	case strings.Contains(c.Name, "Hertogenbosch"):
		return Hertogenbosch
	case strings.Contains(c.Name, "Gravenhage"):
		return DenHaag
	case strings.EqualFold(c.Name, "Zoetermeer"):
		return Zoetermeer
	case strings.Contains(c.Name, "Bodegraven"):
		return Bodegraven
	case strings.Contains(c.Name, "Graveland"):
		return Graveland
	case strings.Contains(c.Name, "Elst"):
		return Elst
	case strings.Contains(c.Name, "Nieuwerkerk a/d"):
		return NieuwerkerkAanDenIJssel
	case strings.Contains(c.Name, "Rijswijk"):
		return Rijswijk
	case strings.Contains(c.Name, "Gravenzande"):
		return Gravenzande
	case strings.EqualFold(c.Name, "Haren Gn"):
		return Haren
	case strings.Contains(c.Name, "Bleskensgraaf"):
		return Bleskensgraaf
	}

	return *c
}

// Suggester permit to get city suggested districts
type Suggester interface {
	Suggest(name string) map[string][]string
}

type suggester struct {
	cities map[string]City
}

func NewSuggester(cities map[string]City) Suggester {
	return &suggester{cities}
}

func (s *suggester) Suggest(name string) map[string][]string {
	city, ok := s.cities[name]
	if !ok || len(city.SuggestedDistrict) == 0 {
		return nil
	}

	return city.SuggestedDistrict
}
