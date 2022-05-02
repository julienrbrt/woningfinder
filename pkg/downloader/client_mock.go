package downloader

import "net/url"

type mockClient struct {
	err error
}

// NewClientMock creates a mock client for the file downloader
func NewClientMock(err error) Client {
	return &mockClient{
		err: err,
	}
}

func (m *mockClient) Download(originalName string, fileURL *url.URL) error {
	return m.err
}

func (m *mockClient) Get(originalName string) ([]byte, error) {
	return nil, m.err
}
