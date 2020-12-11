package corporation

import (
	"errors"
)

// Client defines a Housing Corporation client
type Client interface {
	Login(username, password string) error
	FetchOffer() ([]Offer, error)
	ApplyOffer(offer Offer) error
}

// ErrAuthFailed is retuned when the authentication to the housing coporation has failed
var ErrAuthFailed = errors.New("authentication failed")
