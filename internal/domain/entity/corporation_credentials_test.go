package entity_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/woningfinder/woningfinder/internal/domain/entity"
)

func Test_CorporationCredentials_IsValid(t *testing.T) {
	a := assert.New(t)

	credentials := entity.CorporationCredentials{
		Corporation: entity.Corporation{
			Name: "De Woonplaats",
		},
		Login:    "foo",
		Password: "bar",
	}

	a.Nil(credentials.IsValid())
}

func Test_CorporationCredentials_IsValid_Invalid(t *testing.T) {
	a := assert.New(t)

	credentials := entity.CorporationCredentials{
		Login:    "foo",
		Password: "bar",
	}
	a.Error(credentials.IsValid())

	credentials = entity.CorporationCredentials{
		Corporation: entity.Corporation{
			Name: "De Woonplaats",
		},
	}
	a.Error(credentials.IsValid())

}
