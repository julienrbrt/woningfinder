package middleware

import (
	"github.com/woningfinder/woningfinder/internal/networking"
)

// CreateDefaultHeadersMiddleware adds a list of headers to the request
func CreateDefaultHeadersMiddleware(headers map[string]string) networking.ClientMiddleware {
	return func(base networking.Client) networking.Client {
		return networking.ClientFunc(func(r *networking.Request) (*networking.Response, error) {
			if len(headers) > 0 {
				for key, value := range headers {
					if _, ok := r.Headers[key]; !ok {
						r.SetHeader(key, value)
					}
				}
			}
			return base.Send(r)
		})
	}
}
