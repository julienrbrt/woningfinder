package connector

import (
	"github.com/woningfinder/woningfinder/internal/corporation"
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

func (c *clientMock) GetOffers() ([]corporation.Offer, error) {
	if c.err != nil {
		return nil, c.err
	}

	return c.offers, nil
}

func (c *clientMock) React(_ corporation.Offer) error {
	return c.err
}
