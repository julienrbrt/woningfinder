package handler

import (
	"encoding/json"
	"errors"
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

func Test_Register_ErrEmptyRequest(t *testing.T) {
	a := assert.New(t)
	logger := logging.NewZapLoggerWithoutSentry()

	corporationServiceMock := corporationService.NewServiceMock(nil)
	userServiceMock := userService.NewServiceMock(nil)
	emailServiceMock := emailService.NewServiceMock(nil)
	imgClientMock := downloader.NewClientMock(nil, "")
	handler := &handler{logger, corporationServiceMock, userServiceMock, emailServiceMock, imgClientMock}

	// create request
	req, err := http.NewRequest(http.MethodPost, "/register", nil)
	a.NoError(err)

	// record response
	rr := httptest.NewRecorder()
	h := http.HandlerFunc(handler.Register)

	// server request
	h.ServeHTTP(rr, req)

	// verify status code
	a.Equal(http.StatusBadRequest, rr.Code)

	// verify expected value
	a.Contains(rr.Body.String(), "Bad request")
}

func Test_Register_ErrUserService(t *testing.T) {
	a := assert.New(t)
	logger := logging.NewZapLoggerWithoutSentry()

	corporationServiceMock := corporationService.NewServiceMock(nil)
	userServiceMock := userService.NewServiceMock(errors.New("foo"))
	emailServiceMock := emailService.NewServiceMock(nil)
	imgClientMock := downloader.NewClientMock(nil, "")
	handler := &handler{logger, corporationServiceMock, userServiceMock, emailServiceMock, imgClientMock}

	// create request
	data, err := ioutil.ReadFile("testdata/register-request.json")
	a.NoError(err)

	req, err := http.NewRequest(http.MethodPost, "/register", strings.NewReader(string(data)))
	req.Header.Set("Content-Type", "application/json")
	a.NoError(err)

	// record response
	rr := httptest.NewRecorder()
	h := http.HandlerFunc(handler.Register)

	// server request
	h.ServeHTTP(rr, req)

	// verify status code
	a.Equal(http.StatusInternalServerError, rr.Code)

	// verify expected value
	expected, err := json.Marshal(handlerErrors.ServerErrorRenderer(errors.New("error while creating user")))
	a.NoError(err)
	a.Equal(string(expected), strings.Trim(rr.Body.String(), "\n"))
}

func Test_Register_ErrEmailService(t *testing.T) {
	a := assert.New(t)
	logger := logging.NewZapLoggerWithoutSentry()

	corporationServiceMock := corporationService.NewServiceMock(nil)
	userServiceMock := userService.NewServiceMock(nil)
	emailServiceMock := emailService.NewServiceMock(errors.New("foo"))
	imgClientMock := downloader.NewClientMock(nil, "")
	handler := &handler{logger, corporationServiceMock, userServiceMock, emailServiceMock, imgClientMock}

	// create request
	data, err := ioutil.ReadFile("testdata/register-request.json")
	a.NoError(err)

	req, err := http.NewRequest(http.MethodPost, "/register", strings.NewReader(string(data)))
	req.Header.Set("Content-Type", "application/json")
	a.NoError(err)

	// record response
	rr := httptest.NewRecorder()
	h := http.HandlerFunc(handler.Register)

	// server request
	h.ServeHTTP(rr, req)

	// verify status code
	a.Equal(http.StatusOK, rr.Code)
}

func Test_Register_InvalidHousingType(t *testing.T) {
	a := assert.New(t)
	logger := logging.NewZapLoggerWithoutSentry()

	corporationServiceMock := corporationService.NewServiceMock(nil)
	userServiceMock := userService.NewServiceMock(nil)
	emailServiceMock := emailService.NewServiceMock(nil)
	imgClientMock := downloader.NewClientMock(nil, "")
	handler := &handler{logger, corporationServiceMock, userServiceMock, emailServiceMock, imgClientMock}

	// create request
	data, err := ioutil.ReadFile("testdata/register-invalid-housing-type-request.json")
	a.NoError(err)

	req, err := http.NewRequest(http.MethodPost, "/register", strings.NewReader(string(data)))
	req.Header.Set("Content-Type", "application/json")
	a.NoError(err)

	// record response
	rr := httptest.NewRecorder()
	h := http.HandlerFunc(handler.Register)

	// server request
	h.ServeHTTP(rr, req)

	// verify status code
	a.Equal(http.StatusBadRequest, rr.Code)
}

func Test_Register(t *testing.T) {
	a := assert.New(t)
	logger := logging.NewZapLoggerWithoutSentry()

	corporationServiceMock := corporationService.NewServiceMock(nil)
	userServiceMock := userService.NewServiceMock(nil)
	emailServiceMock := emailService.NewServiceMock(nil)
	imgClientMock := downloader.NewClientMock(nil, "")
	handler := &handler{logger, corporationServiceMock, userServiceMock, emailServiceMock, imgClientMock}

	// create request
	data, err := ioutil.ReadFile("testdata/register-request.json")
	a.NoError(err)

	req, err := http.NewRequest(http.MethodPost, "/register", strings.NewReader(string(data)))
	req.Header.Set("Content-Type", "application/json")
	a.NoError(err)

	// record response
	rr := httptest.NewRecorder()
	h := http.HandlerFunc(handler.Register)

	// server request
	h.ServeHTTP(rr, req)

	// verify status code
	a.Equal(http.StatusOK, rr.Code)
	a.Empty(rr.Body.String())
}
