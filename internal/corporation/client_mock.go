package corporation

import "github.com/woningfinder/woningfinder/internal/domain/entity"

type clientMock struct {
	offers []entity.Offer
	err    error
}

func NewClientMock(offers []entity.Offer, err error) Client {
	return &clientMock{
		offers: offers,
		err:    err,
	}
}

func (c *clientMock) Login(_, _ string) error {
	return c.err
}

func (c *clientMock) FetchOffer() ([]entity.Offer, error) {
	if c.err != nil {
		return nil, c.err
	}

	return c.offers, nil
}

func (c *clientMock) ReactToOffer(_ entity.Offer) error {
	return c.err
}
