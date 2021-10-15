package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/woningfinder/woningfinder/internal/customer"
	corporationService "github.com/woningfinder/woningfinder/internal/services/corporation"
	emailService "github.com/woningfinder/woningfinder/internal/services/email"
	userService "github.com/woningfinder/woningfinder/internal/services/user"
	"github.com/woningfinder/woningfinder/pkg/logging"
	"github.com/woningfinder/woningfinder/pkg/stripe"
)

func Test_CreateStripeCheckoutSession_ErrAPIKeyMissing(t *testing.T) {
	a := assert.New(t)
	logger := logging.NewZapLoggerWithoutSentry()

	corporationServiceMock := corporationService.NewServiceMock(nil)
	userServiceMock := userService.NewServiceMock(nil)
	emailServiceMock := emailService.NewServiceMock(nil)
	handler := &handler{logger, corporationServiceMock, userServiceMock, emailServiceMock, stripe.NewClientMock(false)}

	// create request
	req, err := http.NewRequest(http.MethodPost, "", nil)
	a.NoError(err)

	// record response
	rr := httptest.NewRecorder()
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler.createStripeCheckoutSession("foo@bar.com", customer.PlanPro, w, r)
	})

	// server request
	h.ServeHTTP(rr, req)

	// verify status code
	a.Equal(http.StatusInternalServerError, rr.Code)

	// verify expected value
	a.Contains(rr.Body.String(), "error while creating stripe new checkout session")
}

func Test_CreateStripeCheckoutSession_InvalidPlan(t *testing.T) {
	a := assert.New(t)
	logger := logging.NewZapLoggerWithoutSentry()

	corporationServiceMock := corporationService.NewServiceMock(nil)
	userServiceMock := userService.NewServiceMock(nil)
	emailServiceMock := emailService.NewServiceMock(nil)
	handler := &handler{logger, corporationServiceMock, userServiceMock, emailServiceMock, stripe.NewClientMock(true)}

	// create request
	req, err := http.NewRequest(http.MethodPost, "", nil)
	a.NoError(err)

	// record response
	rr := httptest.NewRecorder()
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler.createStripeCheckoutSession("foo@bar.com", customer.Plan{}, w, r)
	})

	// server request
	h.ServeHTTP(rr, req)

	// verify status code
	a.Equal(http.StatusInternalServerError, rr.Code)

	// verify expected value
	a.Contains(rr.Body.String(), "error while creating stripe new checkout session")
}

func Test_CreateStripeCheckoutSession(t *testing.T) {
	a := assert.New(t)
	logger := logging.NewZapLoggerWithoutSentry()

	corporationServiceMock := corporationService.NewServiceMock(nil)
	userServiceMock := userService.NewServiceMock(nil)
	emailServiceMock := emailService.NewServiceMock(nil)
	handler := &handler{logger, corporationServiceMock, userServiceMock, emailServiceMock, stripe.NewClientMock(true)}

	// create request
	req, err := http.NewRequest(http.MethodPost, "", nil)
	a.NoError(err)

	// record response
	rr := httptest.NewRecorder()
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler.createStripeCheckoutSession("foo@bar.com", customer.PlanPro, w, r)
	})

	// server request
	h.ServeHTTP(rr, req)

	// verify status code
	a.Equal(http.StatusOK, rr.Code)

	// verify expected value
	var response subscriptionResponse
	a.NoError(json.Unmarshal(rr.Body.Bytes(), &response))

	a.NotEmpty(response.StripeSessionID)
}
