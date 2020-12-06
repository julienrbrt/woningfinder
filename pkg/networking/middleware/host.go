package middleware

import (
	"net/url"

	"github.com/woningfinder/woningfinder/pkg/networking"
)

// CreateHostMiddleware permits to specify the host of the request
func CreateHostMiddleware(host *url.URL) networking.ClientMiddleware {
	return func(base networking.Client) networking.Client {
		return networking.ClientFunc(func(r *networking.Request) (*networking.Response, error) {
			r.Host = host
			return base.Send(r)
		})
	}
}
