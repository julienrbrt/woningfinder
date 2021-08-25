package handler

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/render"
	"github.com/woningfinder/woningfinder/internal/customer"
	handlerErrors "github.com/woningfinder/woningfinder/internal/handler/errors"
	userService "github.com/woningfinder/woningfinder/internal/services/user"
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

		if errors.Is(err, userService.ErrUserAlreadyExist) {
			render.Render(w, r, handlerErrors.ErrorRenderer(fmt.Errorf("%s: %s", errorMsg, err.Error())))
			return
		}

		h.logger.Sugar().Warnf("%s: %w", errorMsg, err)
		render.Render(w, r, handlerErrors.ServerErrorRenderer(errorMsg))
		return
	}

	// send welcome email
	if err := h.emailService.SendWelcome(user); err != nil {
		// just logging error
		h.logger.Sugar().Warnf("error while sending activation email: %w", err)
	}
}
