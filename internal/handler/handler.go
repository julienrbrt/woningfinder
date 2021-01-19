package handler

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	customMiddleware "github.com/woningfinder/woningfinder/internal/handler/middleware"
	"github.com/woningfinder/woningfinder/internal/services/corporation"
	"github.com/woningfinder/woningfinder/internal/services/user"
	"github.com/woningfinder/woningfinder/pkg/logging"
)

type handler struct {
	logger             *logging.Logger
	corporationService corporation.Service
	userService        user.Service
}

// NewHandler creates a WoningFinder API router
func NewHandler(logger *logging.Logger, corporationService corporation.Service, userService user.Service) http.Handler {
	handler := &handler{logger, corporationService, userService}

	// router configuration
	r := chi.NewRouter()
	// add middlewares (order matters!)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(customMiddleware.CreateDefaultHeadersMiddleware(map[string]string{"Content-Type": "application/json"}))
	r.Use(customMiddleware.CreateZapMiddleware(logger))
	r.Use(middleware.Recoverer)

	// register default routes
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		render.Render(w, r, ErrNotFound)
	})
	r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusMethodNotAllowed)
		render.Render(w, r, ErrMethodNotAllowed)
	})

	// register routes
	r.Get("/cities", handler.GetCities)
	r.Post("/signup", handler.SignUp)
	r.Route("/corporation-credentials", func(r chi.Router) {
		r.Get("/", handler.GetCorporationCredentials)
		r.Post("/", handler.UpdateCorporationCredentials)
		r.Delete("/", handler.DeleteCorporationCredentials)
	})

	return r
}
