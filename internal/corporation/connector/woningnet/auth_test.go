package woningnet_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	bootstrapCorporation "github.com/woningfinder/woningfinder/internal/bootstrap/corporation"
	"github.com/woningfinder/woningfinder/internal/corporation/connector"
	"github.com/woningfinder/woningfinder/internal/corporation/connector/woningnet"
	"github.com/woningfinder/woningfinder/pkg/logging"
	"github.com/woningfinder/woningfinder/pkg/mapbox"
)

func Test_Login_Failed(t *testing.T) {
	a := assert.New(t)

	client := bootstrapCorporation.CreateWoningNetClient(logging.NewZapLoggerWithoutSentry(), mapbox.NewClientMock(nil, "district"), woningnet.HengeloBorneInfo)

	err := client.Login("example@example.com", "unexesting")
	a.Error(connector.ErrAuthFailed, err)
}
