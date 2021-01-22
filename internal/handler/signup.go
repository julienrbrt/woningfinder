package handler

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/woningfinder/woningfinder/internal/domain/entity"
	handlerEntity "github.com/woningfinder/woningfinder/internal/handler/entity"
)

func (h *handler) SignUp(w http.ResponseWriter, r *http.Request) {
	user := &entity.User{}
	if err := render.Bind(r, user); err != nil {
		render.Render(w, r, handlerEntity.ErrorRenderer(err))
		return
	}

	if err := h.userService.CreateUser(user); err != nil {
		render.Render(w, r, handlerEntity.ServerErrorRenderer(err))
		return
	}

	// returns 200 by default
}
