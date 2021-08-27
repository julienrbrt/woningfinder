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
	corporationService "github.com/woningfinder/woningfinder/internal/services/corporation"
	emailService "github.com/woningfinder/woningfinder/internal/services/email"
	userService "github.com/woningfinder/woningfinder/internal/services/user"
	"github.com/woningfinder/woningfinder/pkg/email"
	"github.com/woningfinder/woningfinder/pkg/logging"
)

func Test_PaymentProcessor_InvalidRequest(t *testing.T) {
	a := assert.New(t)
	logger := logging.NewZapLoggerWithoutSentry()

	corporationServiceMock := corporationService.NewServiceMock(nil)
	userServiceMock := userService.NewServiceMock(nil)
	emailServiceMock := emailService.NewServiceMock(nil)
	handler := &handler{logger, corporationServiceMock, userServiceMock, emailServiceMock, "", &email.ClientMock{}}

	data, err := ioutil.ReadFile("testdata/payment-processor-invalid-payment-method-request.json")
	a.NoError(err)

	// create request
	req, err := http.NewRequest(http.MethodPost, "/payment", strings.NewReader(string(data)))
	req.Header.Set("Content-Type", "application/json")
	a.NoError(err)

	// record response
	rr := httptest.NewRecorder()
	h := http.HandlerFunc(handler.PaymentProcessor)

	// server request
	h.ServeHTTP(rr, req)

	// verify status code
	a.Equal(http.StatusBadRequest, rr.Code)

	// verify expected value
	a.Contains(rr.Body.String(), "invalid payment method: invalid")
}

func Test_PaymentProcessor_ErrUserService(t *testing.T) {
	a := assert.New(t)
	logger := logging.NewZapLoggerWithoutSentry()

	corporationServiceMock := corporationService.NewServiceMock(nil)
	userServiceMock := userService.NewServiceMock(errors.New("foo"))
	emailServiceMock := emailService.NewServiceMock(nil)
	handler := &handler{logger, corporationServiceMock, userServiceMock, emailServiceMock, "", &email.ClientMock{}}

	data, err := ioutil.ReadFile("testdata/payment-processor-stripe-request.json")
	a.NoError(err)

	// create request
	req, err := http.NewRequest(http.MethodPost, "/payment", strings.NewReader(string(data)))
	req.Header.Set("Content-Type", "application/json")
	a.NoError(err)

	// record response
	rr := httptest.NewRecorder()
	h := http.HandlerFunc(handler.PaymentProcessor)

	// server request
	h.ServeHTTP(rr, req)

	// verify status code
	a.Equal(http.StatusNotFound, rr.Code)
}

func Test_PaymentProcessor_Stripe(t *testing.T) {
	a := assert.New(t)
	logger := logging.NewZapLoggerWithoutSentry()

	corporationServiceMock := corporationService.NewServiceMock(nil)
	userServiceMock := userService.NewServiceMock(nil)
	emailServiceMock := emailService.NewServiceMock(nil)
	handler := &handler{logger, corporationServiceMock, userServiceMock, emailServiceMock, "", &email.ClientMock{}}

	data, err := ioutil.ReadFile("testdata/payment-processor-stripe-request.json")
	a.NoError(err)

	// create request
	req, err := http.NewRequest(http.MethodPost, "/payment", strings.NewReader(string(data)))
	req.Header.Set("Content-Type", "application/json")
	a.NoError(err)

	// init stripe library
	stripe.Key = stripeKeyTest

	// record response
	rr := httptest.NewRecorder()
	h := http.HandlerFunc(handler.PaymentProcessor)

	// server request
	h.ServeHTTP(rr, req)

	// verify status code
	a.Equal(http.StatusOK, rr.Code)

	// verify expected value
	var response paymentProcessorResponse

	a.NoError(json.Unmarshal(rr.Body.Bytes(), &response))
	a.NotEmpty(response.StripeSessionID)
	a.Empty(response.CryptoPaymentURL)
}

func Test_PaymentProcessor_Crypto(t *testing.T) {
	a := assert.New(t)
	logger := logging.NewZapLoggerWithoutSentry()

	corporationServiceMock := corporationService.NewServiceMock(nil)
	userServiceMock := userService.NewServiceMock(nil)
	emailServiceMock := emailService.NewServiceMock(nil)
	handler := &handler{logger, corporationServiceMock, userServiceMock, emailServiceMock, "", &email.ClientMock{}}

	data, err := ioutil.ReadFile("testdata/payment-processor-crypto-request.json")
	a.NoError(err)

	// create request
	req, err := http.NewRequest(http.MethodPost, "/payment", strings.NewReader(string(data)))
	req.Header.Set("Content-Type", "application/json")
	a.NoError(err)

	// record response
	rr := httptest.NewRecorder()
	h := http.HandlerFunc(handler.PaymentProcessor)

	// server request
	h.ServeHTTP(rr, req)

	// verify status code
	a.Equal(http.StatusOK, rr.Code)

	// verify expected value
	var response paymentProcessorResponse

	a.NoError(json.Unmarshal(rr.Body.Bytes(), &response))
	a.NotEmpty(response.CryptoPaymentURL)
	a.Empty(response.StripeSessionID)
}
