package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	handlerErrors "github.com/woningfinder/woningfinder/internal/handler/errors"
	corporationService "github.com/woningfinder/woningfinder/internal/services/corporation"
	emailService "github.com/woningfinder/woningfinder/internal/services/email"
	userService "github.com/woningfinder/woningfinder/internal/services/user"
	"github.com/woningfinder/woningfinder/pkg/cryptocom"
	"github.com/woningfinder/woningfinder/pkg/logging"
	"github.com/woningfinder/woningfinder/pkg/stripe"
)

func Test_GetOffering(t *testing.T) {
	a := assert.New(t)
	logger := logging.NewZapLoggerWithoutSentry()

	corporationServiceMock := corporationService.NewServiceMock(nil)
	userServiceMock := userService.NewServiceMock(nil)
	emailServiceMock := emailService.NewServiceMock(nil)
	cryptoMock := cryptocom.NewClientMock(cryptocom.CryptoCheckoutSession{}, nil)
	handler := &handler{logger, corporationServiceMock, userServiceMock, emailServiceMock, stripe.NewClientMock(false), cryptoMock}

	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest(http.MethodGet, "/offering", nil)
	a.NoError(err)

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	h := http.HandlerFunc(handler.GetOffering)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	h.ServeHTTP(rr, req)

	// verify status code
	a.Equal(rr.Code, http.StatusOK)

	// verify expected value
	data, err := ioutil.ReadFile("testdata/offering-response.json")
	a.NoError(err)

	a.Equal(string(data), rr.Body.String())
}

func Test_GetOffering_ErrCorporationService(t *testing.T) {
	a := assert.New(t)
	logger := logging.NewZapLoggerWithoutSentry()

	corporationServiceMock := corporationService.NewServiceMock(errors.New("foo"))
	userServiceMock := userService.NewServiceMock(nil)
	emailServiceMock := emailService.NewServiceMock(nil)
	cryptoMock := cryptocom.NewClientMock(cryptocom.CryptoCheckoutSession{}, nil)
	handler := &handler{logger, corporationServiceMock, userServiceMock, emailServiceMock, stripe.NewClientMock(false), cryptoMock}

	// create request
	req, err := http.NewRequest(http.MethodGet, "/offering", nil)
	a.NoError(err)

	// record response
	rr := httptest.NewRecorder()
	h := http.HandlerFunc(handler.GetOffering)

	// server request
	h.ServeHTTP(rr, req)

	// verify status code
	a.Equal(http.StatusInternalServerError, rr.Code)

	// verify expected value
	expected, err := json.Marshal(handlerErrors.ServerErrorRenderer(fmt.Errorf("error while getting offering")))
	a.NoError(err)
	a.Equal(string(expected), strings.Trim(rr.Body.String(), "\n"))
}
