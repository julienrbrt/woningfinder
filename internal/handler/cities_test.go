package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/woningfinder/woningfinder/pkg/logging"

	"github.com/stretchr/testify/assert"
	handlerEntity "github.com/woningfinder/woningfinder/internal/handler/entity"
	corporationService "github.com/woningfinder/woningfinder/internal/services/corporation"
	paymentService "github.com/woningfinder/woningfinder/internal/services/payment"
	userService "github.com/woningfinder/woningfinder/internal/services/user"
)

func Test_GetCities(t *testing.T) {
	a := assert.New(t)
	logger := logging.NewZapLoggerWithoutSentry()

	corporationServiceMock := corporationService.NewServiceMock(nil)
	userServiceMock := userService.NewServiceMock(nil)
	paymentServiceMock := paymentService.NewServiceMock(nil)
	handler := &handler{logger, corporationServiceMock, userServiceMock, paymentServiceMock, ""}

	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest(http.MethodGet, "/get-cities", nil)
	a.NoError(err)

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	h := http.HandlerFunc(handler.GetCities)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	h.ServeHTTP(rr, req)

	// verify status code
	a.Equal(rr.Code, http.StatusOK)

	// verify expected value
	expected, err := json.Marshal(corporationService.ExpectedMockGetCities)
	a.NoError(err)
	a.Equal(string(expected), strings.Trim(rr.Body.String(), "\n"))
}

func Test_GetCities_ErrCorporationService(t *testing.T) {
	a := assert.New(t)
	logger := logging.NewZapLoggerWithoutSentry()

	corporationServiceMock := corporationService.NewServiceMock(errors.New("foo"))
	userServiceMock := userService.NewServiceMock(nil)
	paymentServiceMock := paymentService.NewServiceMock(nil)
	handler := &handler{logger, corporationServiceMock, userServiceMock, paymentServiceMock, ""}

	// create request
	req, err := http.NewRequest(http.MethodGet, "/get-cities", nil)
	a.NoError(err)

	// record response
	rr := httptest.NewRecorder()
	h := http.HandlerFunc(handler.GetCities)

	// server request
	h.ServeHTTP(rr, req)

	// verify status code
	a.Equal(http.StatusInternalServerError, rr.Code)

	// verify expected value
	expected, err := json.Marshal(handlerEntity.ServerErrorRenderer(fmt.Errorf("error while getting cities")))
	a.NoError(err)
	a.Equal(string(expected), strings.Trim(rr.Body.String(), "\n"))
}
