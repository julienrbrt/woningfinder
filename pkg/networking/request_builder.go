package networking

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"
)

// RequestBuilder permits to build a networking.Request as a http.Request
// This permits to build the request ready to be used by the Go HTTP client and checking if the given request can be correctly sent.
func RequestBuilder(r *Request) (*http.Request, error) {
	if r.Host == nil {
		return nil, fmt.Errorf("no host provided on request %q", r)
	}
	body, err := r.CopyBody()
	if err != nil {
		return nil, fmt.Errorf("failed to get a copy of the request body of request %q to create an HTTP request object: %w", r, err)
	}

	req, err := http.NewRequestWithContext(r.GetContext(), r.GetMethod(), createURL(r), bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request object for %q: %w", r, err)
	}
	for k, v := range r.Headers {
		req.Header.Set(k, v)
	}

	return req, nil
}

func createURL(r *Request) string {
	newURL := *r.Host
	newURL.Path = strings.TrimRight(path.Clean(fmt.Sprintf("/%s/%s", r.Host.Path, r.Path)), "/")
	newURL.RawQuery = r.Query.Encode()

	return strings.ReplaceAll(newURL.String(), url.PathEscape(url.PathEscape("/")), url.PathEscape("/"))
}
