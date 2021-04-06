package handler

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/jwtauth"
	"github.com/stretchr/testify/assert"
	"github.com/woningfinder/woningfinder/internal/auth"
	"github.com/woningfinder/woningfinder/internal/domain/entity"
	corporationService "github.com/woningfinder/woningfinder/internal/services/corporation"
	paymentService "github.com/woningfinder/woningfinder/internal/services/payment"
	userService "github.com/woningfinder/woningfinder/internal/services/user"
	"github.com/woningfinder/woningfinder/pkg/email"
	"github.com/woningfinder/woningfinder/pkg/logging"
)

func Test_GetCorporationCredentials_ErrUnauthorized(t *testing.T) {
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
	h := http.HandlerFunc(handler.GetCorporationCredentials)

	// server request
	h.ServeHTTP(rr, req)

	// verify status code
	a.Equal(http.StatusBadRequest, rr.Code)

	// verify expected value
	a.Contains(rr.Body.String(), "Bad request")
}

func Test_GetCorporationCredentials_ErrUserService(t *testing.T) {
	a := assert.New(t)
	logger := logging.NewZapLoggerWithoutSentry()

	corporationServiceMock := corporationService.NewServiceMock(nil)
	userServiceMock := userService.NewServiceMock(errors.New("foo"))
	paymentServiceMock := paymentService.NewServiceMock(nil)
	handler := &handler{logger, corporationServiceMock, userServiceMock, paymentServiceMock, "", &email.ClientMock{}}

	// create request
	req, err := http.NewRequest(http.MethodPost, endpoint, nil)
	a.NoError(err)

	// record response
	rr := httptest.NewRecorder()
	h := http.HandlerFunc(handler.GetCorporationCredentials)

	// server request
	h.ServeHTTP(rr, authenticateRequest(req))

	// verify status code
	a.Equal(http.StatusInternalServerError, rr.Code)

	// verify expected value
	expected, err := json.Marshal(entity.ServerErrorRenderer(errors.New("failed getting housing corporation relevant for you")))
	a.NoError(err)
	a.Equal(string(expected), strings.Trim(rr.Body.String(), "\n"))
}

func Test_GetCorporationCredentials(t *testing.T) {
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
	h := http.HandlerFunc(handler.GetCorporationCredentials)

	// server request
	h.ServeHTTP(rr, authenticateRequest(req))

	// verify status code
	a.Equal(http.StatusOK, rr.Code)

	data, err := ioutil.ReadFile("testdata/corporation-credentials-response.json")
	a.NoError(err)

	// verify body
	a.Equal(string(data), rr.Body.String())
}

func Test_UpdateCorporationCredentials_ErrUnauthorized(t *testing.T) {
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
	h := http.HandlerFunc(handler.UpdateCorporationCredentials)

	// server request
	h.ServeHTTP(rr, req)

	// verify status code
	a.Equal(http.StatusBadRequest, rr.Code)

	// verify expected value
	a.Contains(rr.Body.String(), "Bad request")
}

func Test_UpdateCorporationCredentials_ErrEmptyRequest(t *testing.T) {
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
	h := http.HandlerFunc(handler.UpdateCorporationCredentials)

	// server request
	h.ServeHTTP(rr, authenticateRequest(req))

	// verify status code
	a.Equal(http.StatusBadRequest, rr.Code)

	// verify expected value
	a.Contains(rr.Body.String(), "Bad request")
}

func Test_UpdateCorporationCredentials_ErrUserService(t *testing.T) {
	a := assert.New(t)
	logger := logging.NewZapLoggerWithoutSentry()

	corporationServiceMock := corporationService.NewServiceMock(nil)
	userServiceMock := userService.NewServiceMock(errors.New("foo"))
	paymentServiceMock := paymentService.NewServiceMock(nil)
	handler := &handler{logger, corporationServiceMock, userServiceMock, paymentServiceMock, "", &email.ClientMock{}}

	// create request
	data, err := ioutil.ReadFile("testdata/corporation-credentials-request.json")
	a.NoError(err)

	req, err := http.NewRequest(http.MethodPost, endpoint, strings.NewReader(string(data)))
	req.Header.Set("Content-Type", "application/json")
	a.NoError(err)

	// record response
	rr := httptest.NewRecorder()
	h := http.HandlerFunc(handler.UpdateCorporationCredentials)

	// server request
	h.ServeHTTP(rr, authenticateRequest(req))

	// verify status code
	a.Equal(http.StatusInternalServerError, rr.Code)

	// verify expected value
	expected, err := json.Marshal(entity.ServerErrorRenderer(errors.New("failed creating corporation credentials")))
	a.NoError(err)
	a.Equal(string(expected), strings.Trim(rr.Body.String(), "\n"))
}

func Test_UpdateCorporationCredentials(t *testing.T) {
	a := assert.New(t)
	logger := logging.NewZapLoggerWithoutSentry()

	corporationServiceMock := corporationService.NewServiceMock(nil)
	userServiceMock := userService.NewServiceMock(nil)
	paymentServiceMock := paymentService.NewServiceMock(nil)
	handler := &handler{logger, corporationServiceMock, userServiceMock, paymentServiceMock, "", &email.ClientMock{}}

	// create request
	data, err := ioutil.ReadFile("testdata/corporation-credentials-request.json")
	a.NoError(err)

	req, err := http.NewRequest(http.MethodPost, endpoint, strings.NewReader(string(data)))
	req.Header.Set("Content-Type", "application/json")
	a.NoError(err)

	// record response
	rr := httptest.NewRecorder()
	h := http.HandlerFunc(handler.UpdateCorporationCredentials)

	// server request
	h.ServeHTTP(rr, authenticateRequest(req))

	// verify status code
	a.Equal(http.StatusOK, rr.Code)

	// verify expected value
	a.Empty(rr.Body.String())
}

// authenticateRequest mock the login of an user
func authenticateRequest(req *http.Request) *http.Request {
	ctx := req.Context()
	user := entity.User{ID: 3, Email: "foo@bar.com"}
	token, _, _ := auth.CreateJWTUserToken(auth.CreateJWTAuthenticationToken("foo"), &user)
	ctx = jwtauth.NewContext(ctx, token, nil)

	return req.WithContext(ctx)
}
