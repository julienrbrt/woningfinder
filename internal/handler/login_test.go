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
	handlerErrors "github.com/woningfinder/woningfinder/internal/handler/errors"
	corporationService "github.com/woningfinder/woningfinder/internal/services/corporation"
	emailService "github.com/woningfinder/woningfinder/internal/services/email"
	userService "github.com/woningfinder/woningfinder/internal/services/user"
	"github.com/woningfinder/woningfinder/pkg/cryptocom"
	"github.com/woningfinder/woningfinder/pkg/logging"
	"github.com/woningfinder/woningfinder/pkg/stripe"
)

func Test_Login_ErrBadRequest(t *testing.T) {
	a := assert.New(t)
	logger := logging.NewZapLoggerWithoutSentry()

	corporationServiceMock := corporationService.NewServiceMock(nil)
	userServiceMock := userService.NewServiceMock(nil)
	emailServiceMock := emailService.NewServiceMock(nil)
	cryptoMock := cryptocom.NewClientMock(cryptocom.CryptoCheckoutSession{}, nil)
	handler := &handler{logger, corporationServiceMock, userServiceMock, emailServiceMock, stripe.NewClientMock(false), cryptoMock}

	// create request
	req, err := http.NewRequest(http.MethodPost, "/login", nil)
	a.NoError(err)

	// record response
	rr := httptest.NewRecorder()
	h := http.HandlerFunc(handler.Login)

	// server request
	h.ServeHTTP(rr, req)

	// verify status code
	a.Equal(http.StatusBadRequest, rr.Code)

	// verify expected value
	a.Contains(rr.Body.String(), "Bad request")
}

func Test_Login_ErrUserService(t *testing.T) {
	a := assert.New(t)
	logger := logging.NewZapLoggerWithoutSentry()

	corporationServiceMock := corporationService.NewServiceMock(nil)
	userServiceMock := userService.NewServiceMock(errors.New("foo"))
	emailServiceMock := emailService.NewServiceMock(nil)
	cryptoMock := cryptocom.NewClientMock(cryptocom.CryptoCheckoutSession{}, nil)
	handler := &handler{logger, corporationServiceMock, userServiceMock, emailServiceMock, stripe.NewClientMock(false), cryptoMock}

	// request data
	data, err := ioutil.ReadFile("testdata/login-request.json")
	a.NoError(err)

	// create request
	req, err := http.NewRequest(http.MethodPost, "/login", strings.NewReader(string(data)))
	req.Header.Set("Content-Type", "application/json")
	a.NoError(err)

	// record response
	rr := httptest.NewRecorder()
	h := http.HandlerFunc(handler.Login)

	// server request
	h.ServeHTTP(rr, authenticateRequest(req))

	// verify status code
	a.Equal(http.StatusNotFound, rr.Code)

	// verify expected value
	expected, err := json.Marshal(handlerErrors.ErrNotFound)
	a.NoError(err)
	a.Equal(string(expected), strings.Trim(rr.Body.String(), "\n"))
}

func Test_Login_ErremailService(t *testing.T) {
	a := assert.New(t)
	logger := logging.NewZapLoggerWithoutSentry()

	corporationServiceMock := corporationService.NewServiceMock(nil)
	userServiceMock := userService.NewServiceMock(nil)
	emailServiceMock := emailService.NewServiceMock(errors.New("foo"))
	cryptoMock := cryptocom.NewClientMock(cryptocom.CryptoCheckoutSession{}, nil)
	handler := &handler{logger, corporationServiceMock, userServiceMock, emailServiceMock, stripe.NewClientMock(false), cryptoMock}

	// request data
	data, err := ioutil.ReadFile("testdata/login-request.json")
	a.NoError(err)

	// create request
	req, err := http.NewRequest(http.MethodPost, "/login", strings.NewReader(string(data)))
	req.Header.Set("Content-Type", "application/json")
	a.NoError(err)

	// record response
	rr := httptest.NewRecorder()
	h := http.HandlerFunc(handler.Login)

	// server request
	h.ServeHTTP(rr, authenticateRequest(req))

	// verify status code
	a.Equal(http.StatusInternalServerError, rr.Code)

	// verify expected value
	expected, err := json.Marshal(handlerErrors.ServerErrorRenderer(errors.New("error while sending login email")))
	a.NoError(err)
	a.Equal(string(expected), strings.Trim(rr.Body.String(), "\n"))
}

func Test_Login(t *testing.T) {
	a := assert.New(t)
	logger := logging.NewZapLoggerWithoutSentry()

	corporationServiceMock := corporationService.NewServiceMock(nil)
	userServiceMock := userService.NewServiceMock(nil)
	emailServiceMock := emailService.NewServiceMock(nil)
	cryptoMock := cryptocom.NewClientMock(cryptocom.CryptoCheckoutSession{}, nil)
	handler := &handler{logger, corporationServiceMock, userServiceMock, emailServiceMock, stripe.NewClientMock(false), cryptoMock}

	// request data
	data, err := ioutil.ReadFile("testdata/login-request.json")
	a.NoError(err)

	// create request
	req, err := http.NewRequest(http.MethodPost, "/login", strings.NewReader(string(data)))
	req.Header.Set("Content-Type", "application/json")
	a.NoError(err)

	// record response
	rr := httptest.NewRecorder()
	h := http.HandlerFunc(handler.Login)

	// server request
	h.ServeHTTP(rr, authenticateRequest(req))

	// verify status code
	a.Equal(http.StatusOK, rr.Code)

	// verify body
	a.Empty(rr.Body.String())
}
