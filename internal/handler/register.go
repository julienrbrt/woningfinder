package handler

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/render"
	"github.com/woningfinder/woningfinder/internal/customer"
	handlerErrors "github.com/woningfinder/woningfinder/internal/handler/errors"
)

// Register contains the handler for registering on WoningFinder
func (h *handler) Register(w http.ResponseWriter, r *http.Request) {
	user := &customer.User{}
	if err := render.Bind(r, user); err != nil {
		render.Render(w, r, handlerErrors.ErrorRenderer(err))
		return
	}

	// lowercase email
	user.Email = strings.ToLower(user.Email)

	if err := h.userService.CreateUser(user); err != nil {
		errorMsg := fmt.Errorf("error while creating user")
		h.logger.Sugar().Warnf("%w: %w", errorMsg, err)
		render.Render(w, r, handlerErrors.ServerErrorRenderer(errorMsg))
		return
	}

	// send account activation email
	if err := h.emailService.SendActivateAccount(user); err != nil {
		errorMsg := fmt.Errorf("error while sending activation email")
		h.logger.Sugar().Warnf("%w: %w", errorMsg, err)
		render.Render(w, r, handlerErrors.ServerErrorRenderer(errorMsg))
	}

}
