package woningnet_test

import (
	"testing"

	bootstrapCorporation "github.com/julienrbrt/woningfinder/internal/bootstrap/corporation"
	"github.com/julienrbrt/woningfinder/internal/corporation/connector"
	"github.com/julienrbrt/woningfinder/internal/corporation/connector/woningnet"
	"github.com/julienrbrt/woningfinder/pkg/logging"
	"github.com/julienrbrt/woningfinder/pkg/mapbox"
	"github.com/stretchr/testify/assert"
)

func Test_Login_Failed(t *testing.T) {
	a := assert.New(t)

	client := bootstrapCorporation.CreateWoningNetClient(logging.NewZapLoggerWithoutSentry(), mapbox.NewClientMock(nil, "district"), woningnet.UtrechtInfo)

	err := client.Login("example@example.com", "unexesting")
	a.Error(connector.ErrAuthFailed, err)
}
