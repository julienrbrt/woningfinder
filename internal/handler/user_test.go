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

func Test_GetUserInfo_ErrUnauthorized(t *testing.T) {
	a := assert.New(t)
	logger := logging.NewZapLoggerWithoutSentry()

	corporationServiceMock := corporationService.NewServiceMock(nil)
	userServiceMock := userService.NewServiceMock(nil)
	emailServiceMock := emailService.NewServiceMock(nil)
	imgClientMock := downloader.NewClientMock(nil)
	handler := &handler{logger, corporationServiceMock, userServiceMock, emailServiceMock, imgClientMock}

	// create request
	req, err := http.NewRequest(http.MethodGet, "/me", nil)
	a.NoError(err)

	// record response
	rr := httptest.NewRecorder()
	h := http.HandlerFunc(handler.GetUserInfo)

	// server request
	h.ServeHTTP(rr, req)

	// verify status code
	a.Equal(http.StatusBadRequest, rr.Code)

	// verify expected value
	a.Contains(rr.Body.String(), "Bad request")
}

func Test_GetUserInfo_ErrUserService(t *testing.T) {
	a := assert.New(t)
	logger := logging.NewZapLoggerWithoutSentry()

	corporationServiceMock := corporationService.NewServiceMock(nil)
	userServiceMock := userService.NewServiceMock(errors.New("foo"))
	emailServiceMock := emailService.NewServiceMock(nil)
	imgClientMock := downloader.NewClientMock(nil)
	handler := &handler{logger, corporationServiceMock, userServiceMock, emailServiceMock, imgClientMock}

	// create request
	req, err := http.NewRequest(http.MethodGet, "/me", nil)
	a.NoError(err)

	// record response
	rr := httptest.NewRecorder()
	h := http.HandlerFunc(handler.GetUserInfo)

	// server request
	h.ServeHTTP(rr, authenticateRequest(req))

	// verify status code
	a.Equal(http.StatusInternalServerError, rr.Code)

	// verify expected value
	expected, err := json.Marshal(handlerErrors.ServerErrorRenderer(errors.New("failed to get user information")))
	a.NoError(err)
	a.Equal(string(expected), strings.Trim(rr.Body.String(), "\n"))
}

func Test_GetUserInfo(t *testing.T) {
	a := assert.New(t)
	logger := logging.NewZapLoggerWithoutSentry()

	corporationServiceMock := corporationService.NewServiceMock(nil)
	userServiceMock := userService.NewServiceMock(nil)
	emailServiceMock := emailService.NewServiceMock(nil)
	imgClientMock := downloader.NewClientMock(nil)
	handler := &handler{logger, corporationServiceMock, userServiceMock, emailServiceMock, imgClientMock}

	// create request
	req, err := http.NewRequest(http.MethodGet, "/me", nil)
	a.NoError(err)

	// record response
	rr := httptest.NewRecorder()
	h := http.HandlerFunc(handler.GetUserInfo)

	// server request
	h.ServeHTTP(rr, authenticateRequest(req))

	// verify status code
	a.Equal(http.StatusOK, rr.Code)

	data, err := ioutil.ReadFile("testdata/user-info-response.json")
	a.NoError(err)

	// verify body
	a.Equal(string(data), rr.Body.String())
}

func Test_UpdateUserInfo_BadRequestError(t *testing.T) {
	a := assert.New(t)
	logger := logging.NewZapLoggerWithoutSentry()

	corporationServiceMock := corporationService.NewServiceMock(nil)
	userServiceMock := userService.NewServiceMock(errors.New("foo"))
	emailServiceMock := emailService.NewServiceMock(nil)
	imgClientMock := downloader.NewClientMock(nil)
	handler := &handler{logger, corporationServiceMock, userServiceMock, emailServiceMock, imgClientMock}

	// create request
	req, err := http.NewRequest(http.MethodPost, "/me", nil)
	a.NoError(err)

	// record response
	rr := httptest.NewRecorder()
	h := http.HandlerFunc(handler.UpdateUserInfo)

	// server request
	h.ServeHTTP(rr, authenticateRequest(req))

	// verify status code
	a.Equal(http.StatusBadRequest, rr.Code)

	// verify expected value
	a.Contains(rr.Body.String(), "Bad request")
}

func Test_UpdateUserInfo_ErrUserService(t *testing.T) {
	a := assert.New(t)
	logger := logging.NewZapLoggerWithoutSentry()

	corporationServiceMock := corporationService.NewServiceMock(nil)
	userServiceMock := userService.NewServiceMock(errors.New("foo"))
	emailServiceMock := emailService.NewServiceMock(nil)
	imgClientMock := downloader.NewClientMock(nil)
	handler := &handler{logger, corporationServiceMock, userServiceMock, emailServiceMock, imgClientMock}

	// request data
	data, err := ioutil.ReadFile("testdata/update-user-info-request.json")
	a.NoError(err)

	// create request
	req, err := http.NewRequest(http.MethodPost, "/me", strings.NewReader(string(data)))
	req.Header.Set("Content-Type", "application/json")
	a.NoError(err)

	// record response
	rr := httptest.NewRecorder()
	h := http.HandlerFunc(handler.UpdateUserInfo)

	// server request
	h.ServeHTTP(rr, authenticateRequest(req))

	// verify status code
	a.Equal(http.StatusInternalServerError, rr.Code)

	// verify expected value
	expected, err := json.Marshal(handlerErrors.ServerErrorRenderer(errors.New("failed to update user information")))
	a.NoError(err)
	a.Equal(string(expected), strings.Trim(rr.Body.String(), "\n"))
}

func Test_UpdateUserInfo(t *testing.T) {
	a := assert.New(t)
	logger := logging.NewZapLoggerWithoutSentry()

	corporationServiceMock := corporationService.NewServiceMock(nil)
	userServiceMock := userService.NewServiceMock(nil)
	emailServiceMock := emailService.NewServiceMock(nil)
	imgClientMock := downloader.NewClientMock(nil)
	handler := &handler{logger, corporationServiceMock, userServiceMock, emailServiceMock, imgClientMock}

	// request data
	data, err := ioutil.ReadFile("testdata/update-user-info-request.json")
	a.NoError(err)

	// create request
	req, err := http.NewRequest(http.MethodPost, "/me", strings.NewReader(string(data)))
	req.Header.Set("Content-Type", "application/json")
	a.NoError(err)

	// record response
	rr := httptest.NewRecorder()
	h := http.HandlerFunc(handler.UpdateUserInfo)

	// server request
	h.ServeHTTP(rr, authenticateRequest(req))

	// verify status code
	a.Equal(http.StatusOK, rr.Code)
}

func Test_DeleteUser_ErrEmptyRequest(t *testing.T) {
	a := assert.New(t)
	logger := logging.NewZapLoggerWithoutSentry()

	corporationServiceMock := corporationService.NewServiceMock(nil)
	userServiceMock := userService.NewServiceMock(errors.New("foo"))
	emailServiceMock := emailService.NewServiceMock(nil)
	imgClientMock := downloader.NewClientMock(nil)
	handler := &handler{logger, corporationServiceMock, userServiceMock, emailServiceMock, imgClientMock}

	// create request
	req, err := http.NewRequest(http.MethodPost, "/me/delete", nil)
	a.NoError(err)

	// record response
	rr := httptest.NewRecorder()
	h := http.HandlerFunc(handler.DeleteUser)

	// server request
	h.ServeHTTP(rr, authenticateRequest(req))

	// verify status code
	a.Equal(http.StatusBadRequest, rr.Code)

	// verify expected value
	a.Contains(rr.Body.String(), "Bad request")
}

func Test_DeleteUser_BadRequestError(t *testing.T) {
	a := assert.New(t)
	logger := logging.NewZapLoggerWithoutSentry()

	corporationServiceMock := corporationService.NewServiceMock(nil)
	userServiceMock := userService.NewServiceMock(nil)
	emailServiceMock := emailService.NewServiceMock(nil)
	imgClientMock := downloader.NewClientMock(nil)
	handler := &handler{logger, corporationServiceMock, userServiceMock, emailServiceMock, imgClientMock}

	// request data
	data, err := ioutil.ReadFile("testdata/delete-user-invalid-request.json")
	a.NoError(err)

	// create request
	req, err := http.NewRequest(http.MethodPost, "/me/delete", strings.NewReader(string(data)))
	req.Header.Set("Content-Type", "application/json")
	a.NoError(err)

	// record response
	rr := httptest.NewRecorder()
	h := http.HandlerFunc(handler.DeleteUser)

	// server request
	h.ServeHTTP(rr, authenticateRequest(req))

	// verify status code
	a.Equal(http.StatusBadRequest, rr.Code)
}

func Test_DeleteUser_ErrUserService(t *testing.T) {
	a := assert.New(t)
	logger := logging.NewZapLoggerWithoutSentry()

	corporationServiceMock := corporationService.NewServiceMock(nil)
	userServiceMock := userService.NewServiceMock(errors.New("foo"))
	emailServiceMock := emailService.NewServiceMock(nil)
	imgClientMock := downloader.NewClientMock(nil)
	handler := &handler{logger, corporationServiceMock, userServiceMock, emailServiceMock, imgClientMock}

	// request data
	data, err := ioutil.ReadFile("testdata/delete-user-request.json")
	a.NoError(err)

	// create request
	req, err := http.NewRequest(http.MethodPost, "/me/delete", strings.NewReader(string(data)))
	req.Header.Set("Content-Type", "application/json")
	a.NoError(err)

	// record response
	rr := httptest.NewRecorder()
	h := http.HandlerFunc(handler.DeleteUser)

	// server request
	h.ServeHTTP(rr, authenticateRequest(req))

	// verify status code
	a.Equal(http.StatusInternalServerError, rr.Code)

	// verify expected value
	expected, err := json.Marshal(handlerErrors.ServerErrorRenderer(errors.New("failed to delete user")))
	a.NoError(err)
	a.Equal(string(expected), strings.Trim(rr.Body.String(), "\n"))
}

func Test_DeleteUser(t *testing.T) {
	a := assert.New(t)
	logger := logging.NewZapLoggerWithoutSentry()

	corporationServiceMock := corporationService.NewServiceMock(nil)
	userServiceMock := userService.NewServiceMock(nil)
	emailServiceMock := emailService.NewServiceMock(nil)
	imgClientMock := downloader.NewClientMock(nil)
	handler := &handler{logger, corporationServiceMock, userServiceMock, emailServiceMock, imgClientMock}

	// request data
	data, err := ioutil.ReadFile("testdata/delete-user-request.json")
	a.NoError(err)

	// create request
	req, err := http.NewRequest(http.MethodPost, "/me/delete", strings.NewReader(string(data)))
	req.Header.Set("Content-Type", "application/json")
	a.NoError(err)

	// record response
	rr := httptest.NewRecorder()
	h := http.HandlerFunc(handler.DeleteUser)

	// server request
	h.ServeHTTP(rr, authenticateRequest(req))

	// verify status code
	a.Equal(http.StatusOK, rr.Code)
}
