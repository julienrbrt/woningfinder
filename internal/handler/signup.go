package handler

import (
	"fmt"
	"net/http"

	"github.com/go-chi/render"
	"github.com/woningfinder/woningfinder/internal/customer"
	handlerErrors "github.com/woningfinder/woningfinder/internal/handler/errors"
)

// SignUp contains the handler for registering on WoningFinder
func (h *handler) SignUp(w http.ResponseWriter, r *http.Request) {
	user := &customer.User{}
	if err := render.Bind(r, user); err != nil {
		render.Render(w, r, handlerErrors.ErrorRenderer(err))
		return
	}

	if err := h.userService.CreateUser(user); err != nil {
		errorMsg := fmt.Errorf("error while creating user")
		h.logger.Sugar().Warnf("%w: %w", errorMsg, err)
		render.Render(w, r, handlerErrors.ServerErrorRenderer(errorMsg))
		return
	}

	// process payment by creating a session id from Stripe
	h.createCheckoutSession(user.Email, user.Plan.Name, w, r)
}
