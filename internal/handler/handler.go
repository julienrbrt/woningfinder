package handler

import (
	"net/http"
	"time"

	chi "github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httprate"
	jwtauth "github.com/go-chi/jwtauth/v5"
	"github.com/go-chi/render"
	handlerErrors "github.com/julienrbrt/woningfinder/internal/handler/errors"
	customMiddleware "github.com/julienrbrt/woningfinder/internal/handler/middleware"
	"github.com/julienrbrt/woningfinder/internal/services/corporation"
	emailService "github.com/julienrbrt/woningfinder/internal/services/email"
	"github.com/julienrbrt/woningfinder/internal/services/user"
	"github.com/julienrbrt/woningfinder/pkg/downloader"
	"github.com/julienrbrt/woningfinder/pkg/logging"
)

type handler struct {
	logger             *logging.Logger
	corporationService corporation.Service
	userService        user.Service
	emailService       emailService.Service
	imgClient          downloader.Client
}

// NewHandler creates a WoningFinder API router
func NewHandler(logger *logging.Logger, jwtAuth *jwtauth.JWTAuth, corporationService corporation.Service, userService user.Service, emailService emailService.Service, imgClient downloader.Client) http.Handler {
	handler := &handler{
		logger:             logger,
		corporationService: corporationService,
		userService:        userService,
		emailService:       emailService,
		imgClient:          imgClient,
	}

	// router configuration
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(customMiddleware.CreateDefaultHeadersMiddleware(map[string]string{
		"Content-Type":                 "application/json",
		"Access-Control-Allow-Origin":  "https://woningfinder.nl",
		"Access-Control-Allow-Method":  "GET, POST, PUT, DELETE, OPTIONS",
		"Access-Control-Allow-Headers": "Accept, Authorization, Content-Type, X-CSRF-Token",
		"Access-Control-Max-Age":       "86400",
	}))
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
		r.Post("/register", handler.Register)
		r.Get("/match/{img}", handler.GetOfferImage)
	})

	// protected routes
	r.Group(func(r chi.Router) {
		// Seek, verify and validate JWT tokens
		r.Use(customMiddleware.JWTVerifierMiddleware(jwtAuth))
		// Handle valid / invalid tokens.
		r.Use(customMiddleware.CreateJWTValidatorMiddleware)

		r.Route("/me", func(r chi.Router) {
			r.Get("/", handler.GetUserInfo)
			r.Post("/", handler.UpdateUserInfo)
			r.Post("/housing-preferences", handler.UpdateHousingPreferences)
			r.Get("/corporation-credentials", handler.GetCorporationCredentials)
			r.Post("/corporation-credentials", handler.UpdateCorporationCredentials)
			r.Post("/delete", handler.DeleteUser)
		})
	})

	return r
}
