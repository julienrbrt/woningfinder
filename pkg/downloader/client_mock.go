package downloader

import "net/url"

type mockClient struct {
	err    error
	output string
}

// NewClientMock creates a mock client for the file downloader
func NewClientMock(err error, output string) Client {
	return &mockClient{
		err:    err,
		output: output,
	}
}

func (m *mockClient) Download(originalName string, fileURL *url.URL) (string, error) {
	if m.err != nil {
		return m.output, m.err
	}

	return m.output, nil
}

func (m *mockClient) Get(originalName string) ([]byte, error) {
	return nil, m.err
}
