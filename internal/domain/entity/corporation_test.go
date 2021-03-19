package entity_test

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/woningfinder/woningfinder/internal/domain/entity"
)

var corporationInfo = entity.Corporation{
	APIEndpoint: &url.URL{Scheme: "https", Host: "example.com"},
	Name:        "OnsHuis",
	URL:         "https://example.com",
	Cities: []entity.City{
		{Name: "Enschede"},
		{Name: "Hengelo"},
	},
	SelectionMethod: []entity.SelectionMethod{
		{
			Method: entity.SelectionRandom,
		},
	},
}

func Test_Corporation_IsValid(t *testing.T) {
	a := assert.New(t)
	a.Nil(corporationInfo.IsValid())
}

func Test_Corporation_IsValid_Invalid(t *testing.T) {
	a := assert.New(t)
	corp := entity.Corporation{
		Name: "Corporation",
		URL:  "https://example.com",
	}
	a.Error(corp.IsValid())
}
