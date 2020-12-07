package corporation

// mockClient mocks the housing corporation client
type mockClient struct {
	offers []Offer
	err    error
}

// CreateMockClient initialize a housing coporation mock client
func CreateMockClient(offers []Offer, err error) mockClient {
	return mockClient{
		offers: offers,
		err:    err,
	}
}

func (m *mockClient) Login(username, password string) error {
	return m.err
}

func (m *mockClient) FetchOffer(minimumPrice float64) ([]Offer, error) {
	return m.offers, m.err
}

func (m *mockClient) ApplyOffer(offer Offer) error {
	return m.err
}
