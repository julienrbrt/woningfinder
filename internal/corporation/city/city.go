package city

import (
	"strings"
	"time"
)

// cityTable defines the city name and it's corresponding struct
var CityTable = map[string]City{
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
	Vinkeveen.Name:            Vinkeveen,
	Harmelen.Name:             Harmelen,
	Amsterdam.Name:            Amsterdam,
	Amstelveen.Name:           Amstelveen,
	Aalsmeer.Name:             Aalsmeer,
	Diemen.Name:               Diemen,
	Zaandam.Name:              Zaandam,
	Hoofddorp.Name:            Hoofddorp,
	Krommenie.Name:            Krommenie,
	Kudelstaart.Name:          Kudelstaart,
	Landsmeer.Name:            Landsmeer,
	Marken.Name:               Marken,
	NieuwVennep.Name:          NieuwVennep,
	Oostzaan.Name:             Oostzaan,
	OuderkerkAanDeAmstel.Name: OuderkerkAanDeAmstel,
	Purmerend.Name:            Purmerend,
	Uithoorn.Name:             Uithoorn,
	Vijfhuizen.Name:           Vijfhuizen,
	Wormer.Name:               Wormer,
	Zwanenburg.Name:           Zwanenburg,
	Lopik.Name:                Lopik,
	Langbroek.Name:            Langbroek,
	Zuidoostbeemster.Name:     Zuidoostbeemster,
	Zaandijk.Name:             Zaandijk,
	Badhoevedorp.Name:         Badhoevedorp,
	DeKwakel.Name:             DeKwakel,
	Lisserbroek.Name:          Lisserbroek,
	Purmerland.Name:           Purmerland,
	Kockengen.Name:            Kockengen,
	Polsbroek.Name:            Polsbroek,
	Hagestein.Name:            Hagestein,
	Kampen.Name:               Kampen,
	IJsselmuiden.Name:         IJsselmuiden,
	Wanneperveen.Name:         Wanneperveen,
	Vollenhove.Name:           Vollenhove,
	KoogaandeZaan.Name:        KoogaandeZaan,
	Leersum.Name:              Leersum,
	Abcoude.Name:              Abcoude,
	Maarn.Name:                Maarn,
	Leerdam.Name:              Leerdam,
	Kamerik.Name:              Kamerik,
	Assendelft.Name:           Assendelft,
	Wormerveer.Name:           Wormerveer,
	Wijthmen.Name:             Wijthmen,
	Jisp.Name:                 Jisp,
	Rijsenhout.Name:           Rijsenhout,
	Willemsoord.Name:          Willemsoord,
	Zegveld.Name:              Zegveld,
	Middenbeemster.Name:       Middenbeemster,
	Monnickendam.Name:         Monnickendam,
	Westzaan.Name:             Westzaan,
	Beinsdorp.Name:            Beinsdorp,
	Cothen.Name:               Cothen,
	Giethoorn.Name:            Giethoorn,
	Cruquius.Name:             Cruquius,
	Aadorp.Name:               Aadorp,
	Bornerbroek.Name:          Bornerbroek,
	Wierden.Name:              Wierden,
	Cruquius.Name:             Cruquius,
	Hertogenbosch.Name:        Hertogenbosch,
	Maastricht.Name:           Maastricht,
	Heerenveen.Name:           Heerenveen,
	Deventer.Name:             Deventer,
	Terborg.Name:              Terborg,
	Haaften.Name:              Haaften,
	Loosdrecht.Name:           Loosdrecht,
	Deurne.Name:               Deurne,
	Soest.Name:                Soest,
	OudBeijerland.Name:        OudBeijerland,
	Uden.Name:                 Uden,
	Doetinchem.Name:           Doetinchem,
	HendrikIdoAmbacht.Name:    HendrikIdoAmbacht,
	CapelleAanDenIJssel.Name:  CapelleAanDenIJssel,
	Kortenhoef.Name:           Kortenhoef,
	Hoogeveen.Name:            Hoogeveen,
	Voorhout.Name:             Voorhout,
	Voorburg.Name:             Voorburg,
	Amersfoort.Name:           Amersfoort,
	Dieren.Name:               Dieren,
	Ede.Name:                  Ede,
	Breda.Name:                Breda,
	Arnhem.Name:               Arnhem,
	Hoogvliet.Name:            Hoogvliet,
	Oss.Name:                  Oss,
	Heemstede.Name:            Heemstede,
	Echt.Name:                 Echt,
	Delft.Name:                Delft,
	Terneuzen.Name:            Terneuzen,
	Roosendaal.Name:           Roosendaal,
	Harderwijk.Name:           Harderwijk,
	Hardenberg.Name:           Hardenberg,
	Zoetermeer.Name:           Zoetermeer,
	Winschoten.Name:           Winschoten,
	Hilversum.Name:            Hilversum,
	Groesbeek.Name:            Groesbeek,
	Naaldwijk.Name:            Naaldwijk,
	Heerlen.Name:              Heerlen,
	Sliedrecht.Name:           Sliedrecht,
	Leerbroek.Name:            Leerbroek,
	Groningen.Name:            Groningen,
	Kaatsheuvel.Name:          Kaatsheuvel,
	Prinsenbeek.Name:          Prinsenbeek,
	Eindhoven.Name:            Eindhoven,
	Papendrecht.Name:          Papendrecht,
	Oegstgeest.Name:           Oegstgeest,
	Goes.Name:                 Goes,
	DenHaag.Name:              DenHaag,
	Roermond.Name:             Roermond,
	Leiden.Name:               Leiden,
	Tilburg.Name:              Tilburg,
	Ermelo.Name:               Ermelo,
	Vlissingen.Name:           Vlissingen,
	Gorinchem.Name:            Gorinchem,
	Waalwijk.Name:             Waalwijk,
	Rotterdam.Name:            Rotterdam,
	Oldenzaal.Name:            Oldenzaal,
	Schiedam.Name:             Schiedam,
	Alkmaar.Name:              Alkmaar,
	Blaricum.Name:             Blaricum,
	Best.Name:                 Best,
	Heerhugowaard.Name:        Heerhugowaard,
	Borculo.Name:              Borculo,
	Hulst.Name:                Hulst,
	Voorst.Name:               Voorst,
	Haarlem.Name:              Haarlem,
}

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
	// TODO try equal fold for every member of the list

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
	}

	return *c
}

// SuggestedCityDistrict permit to get city suggested districts
func SuggestedCityDistrict(name string) map[string][]string {
	city, ok := CityTable[name]
	if !ok || len(city.SuggestedDistrict) == 0 {
		return nil
	}

	return city.SuggestedDistrict
}
