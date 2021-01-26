package handler

import (
	"net/http"
	"time"

	"github.com/go-chi/jwtauth"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/httprate"
	"github.com/go-chi/render"
	handlerEntity "github.com/woningfinder/woningfinder/internal/handler/entity"
	customMiddleware "github.com/woningfinder/woningfinder/internal/handler/middleware"
	"github.com/woningfinder/woningfinder/internal/services/corporation"
	"github.com/woningfinder/woningfinder/internal/services/payment"
	"github.com/woningfinder/woningfinder/internal/services/user"
	"github.com/woningfinder/woningfinder/pkg/logging"
)

type handler struct {
	logger                   *logging.Logger
	corporationService       corporation.Service
	userService              user.Service
	paymentService           payment.Service
	paymentWebhookSigningKey string
}

// NewHandler creates a WoningFinder API router
func NewHandler(logger *logging.Logger, corporationService corporation.Service, userService user.Service, paymentService payment.Service, paymentWebhookSigningKey string, jwtAuth *jwtauth.JWTAuth) http.Handler {
	handler := &handler{logger, corporationService, userService, paymentService, paymentWebhookSigningKey}

	// router configuration
	r := chi.NewRouter()
	// add middlewares (order matters!)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(customMiddleware.CreateDefaultHeadersMiddleware(map[string]string{"Content-Type": "application/json"}))
	r.Use(customMiddleware.CreateZapMiddleware(logger))
	r.Use(middleware.Recoverer)                                                         // recovers from panics and returns 500
	r.Use(httprate.Limit(10, 10*time.Second, httprate.KeyByIP, httprate.KeyByEndpoint)) // Overall rate-limiter, keyed by IP and URL path (aka endpoint). This means each user (by IP) will receive a unique limit counter per endpoint.

	// register default routes
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		render.Render(w, r, handlerEntity.ErrNotFound)
	})
	r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusMethodNotAllowed)
		render.Render(w, r, handlerEntity.ErrMethodNotAllowed)
	})

	// register routes

	// public routes
	r.Group(func(r chi.Router) {
		r.Get("/cities", handler.GetCities)
		r.Post("/signup", handler.SignUp)
		r.Post("/stripe-webhook", handler.ProcessPayment)
	})

	// protected routes
	r.Group(func(r chi.Router) {
		// Seek, verify and validate JWT tokens
		r.Use(customMiddleware.JWTVerifierMiddleware(jwtAuth))
		// Handle valid / invalid tokens.
		r.Use(customMiddleware.CreateJWTValidatorMiddleware)

		r.Route("/corporation-credentials", func(r chi.Router) {
			r.Get("/", handler.GetCorporationCredentials)
			r.Post("/", handler.UpdateCorporationCredentials)
			r.Delete("/", handler.DeleteCorporationCredentials)
		})
	})

	return r
}
