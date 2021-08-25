package handler

import (
	"net/http"
)

// CryptoWebhook is called via the Crypto.com webhook and confirm that a user has paid
func (h *handler) CryptoWebhook(w http.ResponseWriter, r *http.Request) {

	// returns 200 by default
}
