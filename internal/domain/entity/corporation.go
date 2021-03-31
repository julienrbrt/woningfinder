package entity

import (
	"net/url"
	"time"
)

// Corporation defines a housing corporations basic data
// That data is shared between every housing corporations
type Corporation struct {
	APIEndpoint     *url.URL  `pg:"-" json:",omitempty"`
	CreatedAt       time.Time `pg:"default:now()"`
	DeletedAt       time.Time `pg:",soft_delete" json:"-"`
	Name            string    `pg:",pk"`
	URL             string
	Cities          []City            `pg:"many2many:corporation_cities"`
	SelectionMethod []SelectionMethod `pg:"-"`
	SelectionTime   time.Time
}

// City defines a city where a housing corporation operates or when an house offer lies
type City struct {
	CreatedAt time.Time `pg:"default:now()"`
	Name      string    `pg:",pk"`
	District  []string  `pg:"-" json:",omitempty"`
}

// CorporationCity defines the many-to-many relationship table
type CorporationCity struct {
	CorporationName string
	CityName        string
}
