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
	"github.com/woningfinder/woningfinder/pkg/logging"
	"github.com/woningfinder/woningfinder/pkg/stripe"
)

func Test_WaitingListForm_ErrEmptyRequest(t *testing.T) {
	a := assert.New(t)
	logger := logging.NewZapLoggerWithoutSentry()

	corporationServiceMock := corporationService.NewServiceMock(nil)
	userServiceMock := userService.NewServiceMock(nil)
	emailServiceMock := emailService.NewServiceMock(nil)
	handler := &handler{logger, corporationServiceMock, userServiceMock, emailServiceMock, stripe.NewClientMock(false)}

	// create request
	req, err := http.NewRequest(http.MethodPost, "/waitinglist", nil)
	a.NoError(err)

	// record response
	rr := httptest.NewRecorder()
	h := http.HandlerFunc(handler.WaitingListForm)

	// server request
	h.ServeHTTP(rr, req)

	// verify status code
	a.Equal(http.StatusBadRequest, rr.Code)

	// verify expected value
	a.Contains(rr.Body.String(), "Bad request")
}

func Test_WaitingListForm_Spam(t *testing.T) {
	a := assert.New(t)
	logger := logging.NewZapLoggerWithoutSentry()

	corporationServiceMock := corporationService.NewServiceMock(nil)
	userServiceMock := userService.NewServiceMock(nil)
	emailServiceMock := emailService.NewServiceMock(nil)
	handler := &handler{logger, corporationServiceMock, userServiceMock, emailServiceMock, stripe.NewClientMock(false)}

	// create request
	data, err := ioutil.ReadFile("testdata/waitinglist-request-spam.json")
	a.NoError(err)

	req, err := http.NewRequest(http.MethodPost, "/waitinglist", strings.NewReader(string(data)))
	req.Header.Set("Content-Type", "application/json")
	a.NoError(err)

	// record response
	rr := httptest.NewRecorder()
	h := http.HandlerFunc(handler.WaitingListForm)

	// server request
	h.ServeHTTP(rr, req)

	// verify status code
	a.Equal(http.StatusBadRequest, rr.Code)

	// verify expected value
	a.Contains(rr.Body.String(), "Bad request")
}

func Test_WaitingListForm_MalformedEmail(t *testing.T) {
	a := assert.New(t)
	logger := logging.NewZapLoggerWithoutSentry()

	corporationServiceMock := corporationService.NewServiceMock(nil)
	userServiceMock := userService.NewServiceMock(nil)
	emailServiceMock := emailService.NewServiceMock(nil)
	handler := &handler{logger, corporationServiceMock, userServiceMock, emailServiceMock, stripe.NewClientMock(false)}

	// create request
	data, err := ioutil.ReadFile("testdata/waitinglist-request-bad-email.json")
	a.NoError(err)

	req, err := http.NewRequest(http.MethodPost, "/waitinglist", strings.NewReader(string(data)))
	req.Header.Set("Content-Type", "application/json")
	a.NoError(err)

	// record response
	rr := httptest.NewRecorder()
	h := http.HandlerFunc(handler.WaitingListForm)

	// server request
	h.ServeHTTP(rr, req)

	// verify status code
	a.Equal(http.StatusBadRequest, rr.Code)

	// verify expected value
	a.Contains(rr.Body.String(), "Bad request")
}

func Test_WaitingListForm_ErrUSerService(t *testing.T) {
	a := assert.New(t)
	logger := logging.NewZapLoggerWithoutSentry()
	corporationServiceMock := corporationService.NewServiceMock(nil)
	userServiceMock := userService.NewServiceMock(errors.New("err"))
	emailServiceMock := emailService.NewServiceMock(nil)
	handler := &handler{logger, corporationServiceMock, userServiceMock, emailServiceMock, stripe.NewClientMock(false)}

	// create request
	data, err := ioutil.ReadFile("testdata/waitinglist-request.json")
	a.NoError(err)

	req, err := http.NewRequest(http.MethodPost, "/waitinglist", strings.NewReader(string(data)))
	req.Header.Set("Content-Type", "application/json")
	a.NoError(err)

	// record response
	rr := httptest.NewRecorder()
	h := http.HandlerFunc(handler.WaitingListForm)

	// server request
	h.ServeHTTP(rr, req)

	// verify status code
	a.Equal(http.StatusInternalServerError, rr.Code)

	// verify expected value
	expected, err := json.Marshal(handlerErrors.ServerErrorRenderer(errors.New("error while adding john.snow@woningfinder.nl to the waiting list")))
	a.NoError(err)
	a.Equal(string(expected), strings.Trim(rr.Body.String(), "\n"))
}

func Test_WaitingListForm_ErrEmailService(t *testing.T) {
	a := assert.New(t)
	logger := logging.NewZapLoggerWithoutSentry()

	corporationServiceMock := corporationService.NewServiceMock(nil)
	userServiceMock := userService.NewServiceMock(nil)
	emailServiceMock := emailService.NewServiceMock(errors.New("foo"))
	handler := &handler{logger, corporationServiceMock, userServiceMock, emailServiceMock, stripe.NewClientMock(false)}

	// create request
	data, err := ioutil.ReadFile("testdata/waitinglist-request.json")
	a.NoError(err)

	req, err := http.NewRequest(http.MethodPost, "/waitinglist", strings.NewReader(string(data)))
	req.Header.Set("Content-Type", "application/json")
	a.NoError(err)

	// record response
	rr := httptest.NewRecorder()
	h := http.HandlerFunc(handler.WaitingListForm)

	// server request
	h.ServeHTTP(rr, req)

	// verify status code
	a.Equal(http.StatusOK, rr.Code)
}

func Test_WaitingListForm_Success(t *testing.T) {
	a := assert.New(t)
	logger := logging.NewZapLoggerWithoutSentry()

	corporationServiceMock := corporationService.NewServiceMock(nil)
	userServiceMock := userService.NewServiceMock(nil)
	emailServiceMock := emailService.NewServiceMock(nil)
	handler := &handler{logger, corporationServiceMock, userServiceMock, emailServiceMock, stripe.NewClientMock(false)}

	// create request
	data, err := ioutil.ReadFile("testdata/waitinglist-request.json")
	a.NoError(err)

	req, err := http.NewRequest(http.MethodPost, "/waitinglist", strings.NewReader(string(data)))
	req.Header.Set("Content-Type", "application/json")
	a.NoError(err)

	// record response
	rr := httptest.NewRecorder()
	h := http.HandlerFunc(handler.WaitingListForm)

	// server request
	h.ServeHTTP(rr, req)

	// verify status code
	a.Equal(http.StatusOK, rr.Code)
}
