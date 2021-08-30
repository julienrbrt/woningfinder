package cryptocom

type mockClient struct {
	err    error
	output CryptoCheckoutSession
}

// NewClientMock creates a mock client for Crypto.com
func NewClientMock(output CryptoCheckoutSession, err error) Client {
	return &mockClient{
		output: output,
		err:    err,
	}
}

func (m *mockClient) CreatePayment(session CryptoCheckoutSession) (*CryptoCheckoutSession, error) {
	return &m.output, m.err
}

func (m *mockClient) VerifyEvent(signatureHeader, eventRaw string) bool {
	return true
}
