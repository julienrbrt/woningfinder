package handler

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	handlerEntity "github.com/woningfinder/woningfinder/internal/handler/entity"
	corporationService "github.com/woningfinder/woningfinder/internal/services/corporation"
	paymentService "github.com/woningfinder/woningfinder/internal/services/payment"
	userService "github.com/woningfinder/woningfinder/internal/services/user"
	"github.com/woningfinder/woningfinder/pkg/email"
	"github.com/woningfinder/woningfinder/pkg/logging"
)

func Test_GetUser_ErrUnauthorized(t *testing.T) {
	a := assert.New(t)
	logger := logging.NewZapLoggerWithoutSentry()

	corporationServiceMock := corporationService.NewServiceMock(nil)
	userServiceMock := userService.NewServiceMock(nil)
	paymentServiceMock := paymentService.NewServiceMock(nil)
	handler := &handler{logger, corporationServiceMock, userServiceMock, paymentServiceMock, "", &email.ClientMock{}}

	// create request
	req, err := http.NewRequest(http.MethodPost, endpoint, nil)
	a.NoError(err)

	// record response
	rr := httptest.NewRecorder()
	h := http.HandlerFunc(handler.GetUser)

	// server request
	h.ServeHTTP(rr, req)

	// verify status code
	a.Equal(http.StatusBadRequest, rr.Code)

	// verify expected value
	a.Contains(rr.Body.String(), "Bad request")
}

func Test_GetUser_ErrUserService(t *testing.T) {
	a := assert.New(t)
	logger := logging.NewZapLoggerWithoutSentry()
	expectedErr := errors.New("foo")

	corporationServiceMock := corporationService.NewServiceMock(nil)
	userServiceMock := userService.NewServiceMock(expectedErr)
	paymentServiceMock := paymentService.NewServiceMock(nil)
	handler := &handler{logger, corporationServiceMock, userServiceMock, paymentServiceMock, "", &email.ClientMock{}}

	// create request
	req, err := http.NewRequest(http.MethodPost, endpoint, nil)
	a.NoError(err)

	// record response
	rr := httptest.NewRecorder()
	h := http.HandlerFunc(handler.GetUser)

	// server request
	h.ServeHTTP(rr, authenticateRequest(req))

	// verify status code
	a.Equal(http.StatusInternalServerError, rr.Code)

	// verify expected value
	expected, err := json.Marshal(handlerEntity.ServerErrorRenderer(expectedErr))
	a.NoError(err)
	a.Equal(string(expected), strings.Trim(rr.Body.String(), "\n"))
}

func Test_GetUser(t *testing.T) {
	a := assert.New(t)
	logger := logging.NewZapLoggerWithoutSentry()

	corporationServiceMock := corporationService.NewServiceMock(nil)
	userServiceMock := userService.NewServiceMock(nil)
	paymentServiceMock := paymentService.NewServiceMock(nil)
	handler := &handler{logger, corporationServiceMock, userServiceMock, paymentServiceMock, "", &email.ClientMock{}}

	// create request
	req, err := http.NewRequest(http.MethodPost, endpoint, nil)
	a.NoError(err)

	// record response
	rr := httptest.NewRecorder()
	h := http.HandlerFunc(handler.GetUser)

	// server request
	h.ServeHTTP(rr, authenticateRequest(req))

	// verify status code
	a.Equal(http.StatusOK, rr.Code)

	data, err := ioutil.ReadFile("testdata/user-info-response.json")
	a.NoError(err)

	// verify body
	a.Equal(string(data), rr.Body.String())
}