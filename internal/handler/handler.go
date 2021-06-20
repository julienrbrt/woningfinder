package handler

import (
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/httprate"
	"github.com/go-chi/jwtauth"
	"github.com/go-chi/render"
	handlerErrors "github.com/woningfinder/woningfinder/internal/handler/errors"
	customMiddleware "github.com/woningfinder/woningfinder/internal/handler/middleware"
	"github.com/woningfinder/woningfinder/internal/services/corporation"
	"github.com/woningfinder/woningfinder/internal/services/notification"
	"github.com/woningfinder/woningfinder/internal/services/payment"
	"github.com/woningfinder/woningfinder/internal/services/user"
	"github.com/woningfinder/woningfinder/pkg/email"
	"github.com/woningfinder/woningfinder/pkg/logging"
)

type handler struct {
	logger                   *logging.Logger
	corporationService       corporation.Service
	userService              user.Service
	notificationService      notification.Service
	paymentService           payment.Service
	paymentWebhookSigningKey string
	emailClient              email.Client
}

// NewHandler creates a WoningFinder API router
func NewHandler(logger *logging.Logger, corporationService corporation.Service, userService user.Service, notificationService notification.Service, paymentService payment.Service, paymentWebhookSigningKey string, jwtAuth *jwtauth.JWTAuth, emailClient email.Client) http.Handler {
	handler := &handler{
		logger:                   logger,
		corporationService:       corporationService,
		userService:              userService,
		notificationService:      notificationService,
		paymentService:           paymentService,
		paymentWebhookSigningKey: paymentWebhookSigningKey,
		emailClient:              emailClient,
	}

	// router configuration
	r := chi.NewRouter()
	// add middlewares (order matters!)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(customMiddleware.CreateDefaultHeadersMiddleware(map[string]string{"Content-Type": "application/json"}))
	r.Use(customMiddleware.CreateZapMiddleware(logger))
	r.Use(middleware.StripSlashes)                                                                             //strip any trailing slash from the request
	r.Use(middleware.Recoverer)                                                                                // recovers from panics and returns 500
	r.Use(httprate.Limit(10, 10*time.Second, httprate.WithKeyFuncs(httprate.KeyByIP, httprate.KeyByEndpoint))) // overall rate-limiter, keyed by IP and URL path (aka endpoint). This means each user (by IP) will receive a unique limit counter per endpoint.

	// register default routes
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		render.Render(w, r, handlerErrors.ErrNotFound)
	})
	r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusMethodNotAllowed)
		render.Render(w, r, handlerErrors.ErrMethodNotAllowed)
	})

	// register routes

	// public routes
	r.Group(func(r chi.Router) {
		r.Get("/offering", handler.GetOffering)
		r.Post("/contact", handler.ContactForm)
		r.Post("/waitinglist", handler.WaitingListForm)
		r.Post("/login", handler.Login)
		r.Post("/signup", handler.SignUp)
		r.Post("/payment", handler.PaymentProcessor)
		r.Post("/stripe-webhook", handler.PaymentValidator)
	})

	// protected routes
	r.Group(func(r chi.Router) {
		// Seek, verify and validate JWT tokens
		r.Use(customMiddleware.JWTVerifierMiddleware(jwtAuth))
		// Handle valid / invalid tokens.
		r.Use(customMiddleware.CreateJWTValidatorMiddleware)

		r.Route("/me", func(r chi.Router) {
			r.Get("/", handler.GetUser)
			r.Get("/corporation-credentials", handler.GetCorporationCredentials)
			r.Post("/corporation-credentials", handler.UpdateCorporationCredentials)
		})
	})

	return r
}
