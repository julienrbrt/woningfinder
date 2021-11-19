package woonburo

import "github.com/woningfinder/woningfinder/internal/corporation"

func (c *client) FetchOffers(ch chan<- corporation.Offer) error {
	defer close(ch)

	return nil
}
