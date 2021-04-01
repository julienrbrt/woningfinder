package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/jwtauth"
	"github.com/go-chi/render"
	"github.com/woningfinder/woningfinder/internal/auth"
	handlerEntity "github.com/woningfinder/woningfinder/internal/handler/entity"
)

// GetUser gets all the user information
func (h *handler) GetUser(w http.ResponseWriter, r *http.Request) {
	// extract jwt
	_, claims, err := jwtauth.FromContext(r.Context())
	if err != nil {
		render.Render(w, r, handlerEntity.ErrBadRequest)
		return
	}

	// get user from jwt claims
	userFromJTW, err := auth.ExtractUserFromJWT(claims)
	if err != nil {
		render.Render(w, r, handlerEntity.ErrBadRequest)
		return
	}

	user, err := h.userService.GetUser(userFromJTW)
	if err != nil {
		render.Render(w, r, handlerEntity.ServerErrorRenderer(err))
		return
	}

	json.NewEncoder(w).Encode(user)
}
