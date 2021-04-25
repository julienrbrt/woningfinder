package city_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/woningfinder/woningfinder/internal/corporation"
	"github.com/woningfinder/woningfinder/internal/corporation/city"
)

func TestCity_Merge(t *testing.T) {
	a := assert.New(t)

	mergedCity := city.Merge(corporation.City{Name: "Hengelo OV"})
	a.Equal(mergedCity, city.Hengelo)

	mergedCity = city.Merge(corporation.City{Name: "Hengelo"})
	a.Equal(mergedCity, city.Hengelo)

	expected := corporation.City{Name: "a city"}
	mergedCity = city.Merge(expected)
	a.Equal(mergedCity, expected)
}
