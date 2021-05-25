package handler

import (
	_ "embed"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/render"
	"github.com/woningfinder/woningfinder/internal/customer"
	handlerErrors "github.com/woningfinder/woningfinder/internal/handler/errors"
	"github.com/woningfinder/woningfinder/pkg/util"
)

type waitinglistFormRequest struct {
	Email    string `json:"email"`
	CityName string `json:"city"`
	AntiSpam int    `json:"phone"` // using phone to lure bots
}

// Bind permits go-chi router to verify the user input and marshal it
func (w *waitinglistFormRequest) Bind(r *http.Request) error {
	if w.Email == "" || w.CityName == "" || w.AntiSpam == 0 {
		return errors.New("fields cannot be empty")
	}

	if !util.IsEmailValid(w.Email) {
		return fmt.Errorf("email invalid")
	}

	// verify anti spam
	// the anti spam is build the sum of all bytes of the sent email and message + 374
	var sum = 374
	bytes := append([]byte(w.Email), []byte(w.CityName)...)
	for _, b := range bytes {
		sum += int(b)
	}

	if w.AntiSpam != sum {
		return errors.New("waiting list form cannot be processed")
	}

	return nil
}

// Render permits go-chi router to render the user
func (*waitinglistFormRequest) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (h *handler) WaitingListForm(w http.ResponseWriter, r *http.Request) {
	waitingListRequest := &waitinglistFormRequest{}
	if err := render.Bind(r, waitingListRequest); err != nil {
		render.Render(w, r, handlerErrors.ErrorRenderer(err))
		return
	}

	// save waiting list
	if err := h.userService.CreateWaitingList(&customer.WaitingList{Email: waitingListRequest.Email, CityName: waitingListRequest.CityName}); err != nil {
		errorMsg := fmt.Errorf("error while adding %s to the waiting list", waitingListRequest.Email)
		h.logger.Sugar().Warnf("%w: %w", errorMsg, err)
		render.Render(w, r, handlerErrors.ServerErrorRenderer(errorMsg))
		return
	}

	// returns 200 by default
}
