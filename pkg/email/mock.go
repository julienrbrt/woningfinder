package email

type ClientMock struct {
	Err error
}

func (c *ClientMock) Send(_, _, _ string) error {
	return c.Err
}
