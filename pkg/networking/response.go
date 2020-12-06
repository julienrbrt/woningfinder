package networking

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

// Response defines a response.
type Response struct {
	Request     *Request
	StatusCode  int
	Headers     map[string][]string
	Body        io.ReadCloser
	bodyBytes   []byte
	RawRequest  *http.Request
	RawResponse *http.Response
}

// CopyBody reads the response body to bytes.
func (r *Response) CopyBody() ([]byte, error) {
	if r.Body == nil {
		return nil, nil
	}
	if r.bodyBytes == nil {
		body, err := ioutil.ReadAll(r.Body)
		defer r.Body.Close()

		if err != nil {
			return nil, fmt.Errorf("failed to read response body for copying response from request %q: %w", r.Request, err)
		}
		r.SetBody(ioutil.NopCloser(bytes.NewReader(body)))
		r.bodyBytes = body
	}

	return r.bodyBytes, nil
}

// SetBody sets the reponse's body.
func (r *Response) SetBody(body io.ReadCloser) {
	r.Body = body
	if r.RawResponse != nil {
		r.RawResponse.Body = body
	}
}

// ReadJSONBody unmarshal a json body.
func (r *Response) ReadJSONBody(data interface{}) error {
	if r.Body == nil {
		return fmt.Errorf("unable to read json body from nil body")
	}
	body, err := r.CopyBody()
	if err != nil {
		return fmt.Errorf("unable to read the response body: %w", err)
	}
	if err := json.Unmarshal(body, data); err != nil {
		return fmt.Errorf("failed to unmarshal response from %q to JSON: %w", r.Request, err)
	}
	return nil
}
