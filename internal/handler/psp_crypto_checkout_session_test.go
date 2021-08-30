package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/woningfinder/woningfinder/internal/customer"
	corporationService "github.com/woningfinder/woningfinder/internal/services/corporation"
	emailService "github.com/woningfinder/woningfinder/internal/services/email"
	userService "github.com/woningfinder/woningfinder/internal/services/user"
	"github.com/woningfinder/woningfinder/pkg/cryptocom"
	"github.com/woningfinder/woningfinder/pkg/logging"
	"github.com/woningfinder/woningfinder/pkg/stripe"
)

func Test_CreateCryptoCheckoutSession_Error(t *testing.T) {
	a := assert.New(t)
	logger := logging.NewZapLoggerWithoutSentry()

	corporationServiceMock := corporationService.NewServiceMock(nil)
	userServiceMock := userService.NewServiceMock(nil)
	emailServiceMock := emailService.NewServiceMock(nil)
	cryptoMock := cryptocom.NewClientMock(cryptocom.CryptoCheckoutSession{}, errors.New("foo"))
	handler := &handler{logger, corporationServiceMock, userServiceMock, emailServiceMock, stripe.NewClientMock(false), cryptoMock}

	// create request
	req, err := http.NewRequest(http.MethodPost, "", nil)
	a.NoError(err)

	// record response
	rr := httptest.NewRecorder()
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler.createCryptoCheckoutSession("foo@bar.com", customer.PlanBasis, w, r)
	})

	// server request
	h.ServeHTTP(rr, req)

	// verify status code
	a.Equal(http.StatusInternalServerError, rr.Code)

	// verify expected value
	a.Contains(rr.Body.String(), "error while creating crypto.com new checkout session")
}

func Test_CreateCryptoCheckoutSession(t *testing.T) {
	a := assert.New(t)
	logger := logging.NewZapLoggerWithoutSentry()

	corporationServiceMock := corporationService.NewServiceMock(nil)
	userServiceMock := userService.NewServiceMock(nil)
	emailServiceMock := emailService.NewServiceMock(nil)
	cryptoMock := cryptocom.NewClientMock(cryptocom.CryptoCheckoutSession{
		PaymentURL: "https://example.org/foo",
	}, nil)
	handler := &handler{logger, corporationServiceMock, userServiceMock, emailServiceMock, stripe.NewClientMock(false), cryptoMock}

	// create request
	req, err := http.NewRequest(http.MethodPost, "", nil)
	a.NoError(err)

	// record response
	rr := httptest.NewRecorder()
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler.createCryptoCheckoutSession("foo@bar.com", customer.PlanBasis, w, r)
	})

	// server request
	h.ServeHTTP(rr, req)

	// verify status code
	a.Equal(http.StatusOK, rr.Code)

	// verify expected value
	var response paymentProcessorResponse
	a.NoError(json.Unmarshal(rr.Body.Bytes(), &response))

	a.NotEmpty(response.CryptoPaymentURL)
}
