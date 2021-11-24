package domijn_test

import (
	"testing"

	"github.com/julienrbrt/woningfinder/internal/corporation/connector"
	"github.com/julienrbrt/woningfinder/internal/corporation/connector/domijn"
	"github.com/julienrbrt/woningfinder/pkg/logging"
	"github.com/julienrbrt/woningfinder/pkg/mapbox"
	"github.com/stretchr/testify/assert"
)

func Test_Login_Failed(t *testing.T) {
	a := assert.New(t)
	mockMapbox := mapbox.NewClientMock(nil, "district")

	client, err := domijn.NewClient(logging.NewZapLoggerWithoutSentry(), mockMapbox)
	a.NoError(err)

	err = client.Login("example@example.com", "unexesting")
	a.Error(connector.ErrAuthFailed, err)
}
