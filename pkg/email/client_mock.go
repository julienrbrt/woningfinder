package email

type mockClient struct {
	err error
}

// NewClientMock creates a mock client for email
func NewClientMock(err error) Client {
	return &mockClient{
		err: err,
	}
}

func (c *mockClient) Send(_, _, _ string) error {
	return c.err
}
