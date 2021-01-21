package mapbox

type mockClient struct {
	err    error
	output string
}

// NewClientMock creates a mock client for Mapbox
func NewClientMock(err error, output string) Client {
	return &mockClient{
		err:    err,
		output: output,
	}
}

func (m *mockClient) CityDistrictFromCoords(_, _ string) (string, error) {
	if m.err != nil {
		return "", m.err
	}

	return m.output, nil
}

func (m *mockClient) CityDistrictFromAddress(_ string) (string, error) {
	if m.err != nil {
		return "", m.err
	}

	return m.output, nil
}
