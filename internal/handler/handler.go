package handler

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	customMiddleware "github.com/woningfinder/woningfinder/internal/handler/middleware"
	"github.com/woningfinder/woningfinder/internal/services/corporation"
	"github.com/woningfinder/woningfinder/internal/services/user"
	"github.com/woningfinder/woningfinder/pkg/logging"
)

type handler struct {
	logger *logging.Logger
	corporation.CorporationService
	user.UserService
}

// NewHandler creates a WoningFinder API router
func NewHandler(logger *logging.Logger, corporationService corporation.CorporationService, userService user.UserService) http.Handler {
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
	r.NotFound(handler.NotFound)
	r.MethodNotAllowed(handler.MethodNotAllowed)

	// register routes
	r.Get("/cities", handler.Cities)
	r.Post("/signup", handler.SignUp)

	return r
}
