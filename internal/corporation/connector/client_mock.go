package connector

import (
	"github.com/julienrbrt/woningfinder/internal/corporation"
)

type clientMock struct {
	offers []corporation.Offer
	err    error
}

func NewClientMock(offers []corporation.Offer, err error) Client {
	return &clientMock{
		offers: offers,
		err:    err,
	}
}

func (c *clientMock) Login(_, _ string) error {
	return c.err
}

func (c *clientMock) FetchOffers(ch chan<- corporation.Offer) error {
	if c.err != nil {
		return c.err
	}

	for _, offer := range c.offers {
		ch <- offer
	}

	return nil
}

func (c *clientMock) React(_ corporation.Offer) error {
	return c.err
}
