package handler

import (
	"errors"
	"net/http"

	"github.com/go-chi/render"
	"github.com/woningfinder/woningfinder/internal/domain/entity"
)

func (h *handler) SignUp(w http.ResponseWriter, r *http.Request) {
	user := &entity.User{}
	if err := render.Bind(r, user); err != nil {
		render.Render(w, r, ErrorRenderer(err))
		return
	}

	if err := h.userService.CreateUser(user); err != nil {
		render.Render(w, r, ServerErrorRenderer(errors.New("error while registering an account")))
		return
	}

	// returns 200 by default
}
