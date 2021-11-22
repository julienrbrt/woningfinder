package city_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	bootstrapCorporation "github.com/woningfinder/woningfinder/internal/bootstrap/corporation"
	"github.com/woningfinder/woningfinder/internal/corporation/city"
	"github.com/woningfinder/woningfinder/pkg/logging"
)

func TestCity_Merge(t *testing.T) {
	a := assert.New(t)

	a.Equal(((&city.City{Name: "Hengelo OV"}).Merge()), city.Hengelo)
	a.Equal(((&city.City{Name: "Hengelo"}).Merge()), city.Hengelo)
	expected := &city.City{Name: "a city"}
	a.Equal(*expected, expected.Merge())
}

func TestNewSuggester_Suggest(t *testing.T) {
	a := assert.New(t)

	cities := bootstrapCorporation.CreateConnectorProvider(logging.NewZapLoggerWithoutSentry(), nil).GetCities()
	suggester := city.NewSuggester(cities)

	a.Len(suggester.Suggest("Losser"), 0)
	a.True(len(suggester.Suggest("Enschede")) > 0)
}
