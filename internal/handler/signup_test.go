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
	"github.com/stripe/stripe-go"
	handlerErrors "github.com/woningfinder/woningfinder/internal/handler/errors"
	corporationService "github.com/woningfinder/woningfinder/internal/services/corporation"
	notificationService "github.com/woningfinder/woningfinder/internal/services/notification"
	paymentService "github.com/woningfinder/woningfinder/internal/services/payment"
	userService "github.com/woningfinder/woningfinder/internal/services/user"
	"github.com/woningfinder/woningfinder/pkg/email"
	"github.com/woningfinder/woningfinder/pkg/logging"
)

func Test_SignUp_ErrEmptyRequest(t *testing.T) {
	a := assert.New(t)
	logger := logging.NewZapLoggerWithoutSentry()

	corporationServiceMock := corporationService.NewServiceMock(nil)
	userServiceMock := userService.NewServiceMock(nil)
	notificationServiceMock := notificationService.NewServiceMock(nil)
	paymentServiceMock := paymentService.NewServiceMock(nil)
	handler := &handler{logger, corporationServiceMock, userServiceMock, notificationServiceMock, paymentServiceMock, "", &email.ClientMock{}}

	// create request
	req, err := http.NewRequest(http.MethodPost, "/signup", nil)
	a.NoError(err)

	// record response
	rr := httptest.NewRecorder()
	h := http.HandlerFunc(handler.SignUp)

	// server request
	h.ServeHTTP(rr, req)

	// verify status code
	a.Equal(http.StatusBadRequest, rr.Code)

	// verify expected value
	a.Contains(rr.Body.String(), "Bad request")
}

func Test_SignUp_ErrUserService(t *testing.T) {
	a := assert.New(t)
	logger := logging.NewZapLoggerWithoutSentry()

	corporationServiceMock := corporationService.NewServiceMock(nil)
	userServiceMock := userService.NewServiceMock(errors.New("foo"))
	notificationServiceMock := notificationService.NewServiceMock(nil)
	paymentServiceMock := paymentService.NewServiceMock(nil)
	handler := &handler{logger, corporationServiceMock, userServiceMock, notificationServiceMock, paymentServiceMock, "", &email.ClientMock{}}

	// create request
	data, err := ioutil.ReadFile("testdata/signup-request-pro.json")
	a.NoError(err)

	req, err := http.NewRequest(http.MethodPost, "/signup", strings.NewReader(string(data)))
	req.Header.Set("Content-Type", "application/json")
	a.NoError(err)

	// record response
	rr := httptest.NewRecorder()
	h := http.HandlerFunc(handler.SignUp)

	// server request
	h.ServeHTTP(rr, req)

	// verify status code
	a.Equal(http.StatusInternalServerError, rr.Code)

	// verify expected value
	expected, err := json.Marshal(handlerErrors.ServerErrorRenderer(errors.New("error while creating user")))
	a.NoError(err)
	a.Equal(string(expected), strings.Trim(rr.Body.String(), "\n"))
}

func Test_SignUp_InvalidPlan(t *testing.T) {
	a := assert.New(t)
	logger := logging.NewZapLoggerWithoutSentry()

	corporationServiceMock := corporationService.NewServiceMock(nil)
	userServiceMock := userService.NewServiceMock(nil)
	notificationServiceMock := notificationService.NewServiceMock(nil)
	paymentServiceMock := paymentService.NewServiceMock(nil)
	handler := &handler{logger, corporationServiceMock, userServiceMock, notificationServiceMock, paymentServiceMock, "", &email.ClientMock{}}

	// create request
	data, err := ioutil.ReadFile("testdata/signup-invalid-plan-request.json")
	a.NoError(err)

	req, err := http.NewRequest(http.MethodPost, "/signup", strings.NewReader(string(data)))
	req.Header.Set("Content-Type", "application/json")
	a.NoError(err)

	// record response
	rr := httptest.NewRecorder()
	h := http.HandlerFunc(handler.SignUp)

	// server request
	h.ServeHTTP(rr, req)

	// verify status code
	a.Equal(http.StatusBadRequest, rr.Code)
}

func Test_SignUp_InvalidHousingType(t *testing.T) {
	a := assert.New(t)
	logger := logging.NewZapLoggerWithoutSentry()

	corporationServiceMock := corporationService.NewServiceMock(nil)
	userServiceMock := userService.NewServiceMock(nil)
	notificationServiceMock := notificationService.NewServiceMock(nil)
	paymentServiceMock := paymentService.NewServiceMock(nil)
	handler := &handler{logger, corporationServiceMock, userServiceMock, notificationServiceMock, paymentServiceMock, "", &email.ClientMock{}}

	// create request
	data, err := ioutil.ReadFile("testdata/signup-invalid-housing-type-request.json")
	a.NoError(err)

	req, err := http.NewRequest(http.MethodPost, "/signup", strings.NewReader(string(data)))
	req.Header.Set("Content-Type", "application/json")
	a.NoError(err)

	// record response
	rr := httptest.NewRecorder()
	h := http.HandlerFunc(handler.SignUp)

	// server request
	h.ServeHTTP(rr, req)

	// verify status code
	a.Equal(http.StatusBadRequest, rr.Code)
}

func Test_SignUp_Basis(t *testing.T) {
	a := assert.New(t)
	logger := logging.NewZapLoggerWithoutSentry()

	corporationServiceMock := corporationService.NewServiceMock(nil)
	userServiceMock := userService.NewServiceMock(nil)
	notificationServiceMock := notificationService.NewServiceMock(nil)
	paymentServiceMock := paymentService.NewServiceMock(nil)
	handler := &handler{logger, corporationServiceMock, userServiceMock, notificationServiceMock, paymentServiceMock, "", &email.ClientMock{}}

	// create request
	data, err := ioutil.ReadFile("testdata/signup-request-basis.json")
	a.NoError(err)

	req, err := http.NewRequest(http.MethodPost, "/signup", strings.NewReader(string(data)))
	req.Header.Set("Content-Type", "application/json")
	a.NoError(err)

	// init stripe library
	// we do that because the signup handler directly talks to stripe
	stripe.Key = stripeKeyTest

	// record response
	rr := httptest.NewRecorder()
	h := http.HandlerFunc(handler.SignUp)

	// server request
	h.ServeHTTP(rr, req)

	// verify status code
	a.Equal(http.StatusOK, rr.Code)

	// verify expected value
	var result createCheckoutSessionResponse
	a.NoError(json.Unmarshal(rr.Body.Bytes(), &result))
	a.NotEmpty(result.SessionID)
}

func Test_SignUp_Pro(t *testing.T) {
	a := assert.New(t)
	logger := logging.NewZapLoggerWithoutSentry()

	corporationServiceMock := corporationService.NewServiceMock(nil)
	userServiceMock := userService.NewServiceMock(nil)
	notificationServiceMock := notificationService.NewServiceMock(nil)
	paymentServiceMock := paymentService.NewServiceMock(nil)
	handler := &handler{logger, corporationServiceMock, userServiceMock, notificationServiceMock, paymentServiceMock, "", &email.ClientMock{}}

	// create request
	data, err := ioutil.ReadFile("testdata/signup-request-pro.json")
	a.NoError(err)

	req, err := http.NewRequest(http.MethodPost, "/signup", strings.NewReader(string(data)))
	req.Header.Set("Content-Type", "application/json")
	a.NoError(err)

	// init stripe library
	// we do that because the signup handler directly talks to stripe
	stripe.Key = stripeKeyTest

	// record response
	rr := httptest.NewRecorder()
	h := http.HandlerFunc(handler.SignUp)

	// server request
	h.ServeHTTP(rr, req)

	// verify status code
	a.Equal(http.StatusOK, rr.Code)

	// verify expected value
	var result createCheckoutSessionResponse
	a.NoError(json.Unmarshal(rr.Body.Bytes(), &result))
	a.NotEmpty(result.SessionID)
}
