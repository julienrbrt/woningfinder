package spaces

import "net/url"

type mockClient struct {
	err    error
	output string
}

// NewClientMock creates a mock client for DigitalOcean Spaces
func NewClientMock(err error, output string) Client {
	return &mockClient{
		err:    err,
		output: output,
	}
}

func (m *mockClient) UploadPicture(prefix, originalName string, pictureURL *url.URL) (string, error) {
	if m.err != nil {
		return m.output, m.err
	}

	return m.output, nil
}
