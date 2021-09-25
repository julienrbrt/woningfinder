package city_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/woningfinder/woningfinder/internal/city"
)

func TestCity(t *testing.T) {
	a := assert.New(t)

	a.Equal(city.Hengelo.Name, "Hengelo OV")
	a.Equal(len(city.Almelo.Districts()), len(city.Almelo.District))
	a.Empty(city.Almelo.Neighbourhoods())
}
