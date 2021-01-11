package corporation

import (
	"errors"

	"github.com/woningfinder/woningfinder/internal/domain/entity"
)

// ErrAuthFailed is retuned when the authentication to the housing coporation has failed
var ErrAuthFailed = errors.New("authentication failed")

// Client defines a Housing Corporation client
type Client interface {
	Login(username, password string) error
	FetchOffer() ([]entity.Offer, error)
	ReactToOffer(offer entity.Offer) error
}
