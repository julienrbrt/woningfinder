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

	handlerErrors "github.com/julienrbrt/woningfinder/internal/handler/errors"
	corporationService "github.com/julienrbrt/woningfinder/internal/services/corporation"
	emailService "github.com/julienrbrt/woningfinder/internal/services/email"
	userService "github.com/julienrbrt/woningfinder/internal/services/user"
	"github.com/julienrbrt/woningfinder/pkg/downloader"
	"github.com/julienrbrt/woningfinder/pkg/logging"
	"github.com/stretchr/testify/assert"
)

func Test_GetOffering(t *testing.T) {
	a := assert.New(t)
	logger := logging.NewZapLoggerWithoutSentry()

	corporationServiceMock := corporationService.NewServiceMock(nil)
	userServiceMock := userService.NewServiceMock(nil)
	emailServiceMock := emailService.NewServiceMock(nil)
	imgClientMock := downloader.NewClientMock(nil, "")
	handler := &handler{logger, corporationServiceMock, userServiceMock, emailServiceMock, imgClientMock}

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
	imgClientMock := downloader.NewClientMock(nil, "")
	handler := &handler{logger, corporationServiceMock, userServiceMock, emailServiceMock, imgClientMock}

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
