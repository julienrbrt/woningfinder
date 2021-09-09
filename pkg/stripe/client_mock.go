package stripe

import stripe "github.com/stripe/stripe-go/v72"

var stripeKeyTest = "sk_test_51HkWn4HWufZqidI12yfUuTsZxIdKfSlblDYcAYPda4hzMnGrDcDCLannohEiYI0TUXT1rPdx186CyhKvo67H96Ty00vP5NDSrZ"

type mockClient struct{}

// NewClientMock creates a mock client for Stripe
func NewClientMock(initLibrary bool) Client {
	stripe.Key = ""
	if initLibrary {
		stripe.Key = stripeKeyTest
	}

	return &mockClient{}
}

func (m *mockClient) WebhookSigningKey() string {
	return "foo"
}
