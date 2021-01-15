package entity_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/woningfinder/woningfinder/internal/corporation/dewoonplaats"
	"github.com/woningfinder/woningfinder/internal/domain/entity"
)

func Test_Corporation_IsValid(t *testing.T) {
	a := assert.New(t)
	a.Nil(dewoonplaats.Info.IsValid())
}

func Test_Corporation_IsValid_Invalid(t *testing.T) {
	a := assert.New(t)
	corp := entity.Corporation{
		Name: "Corporation",
		URL:  "https://example.com",
	}
	a.Error(corp.IsValid())
}
