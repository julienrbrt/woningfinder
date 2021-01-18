package middleware

import (
	"net/http"
)

// CreateDefaultHeadersMiddleware adds a list of headers to the request
func CreateDefaultHeadersMiddleware(headers map[string]string) func(next http.Handler) http.Handler {
	return defaultHeaders{headers: headers}.middleware
}

type defaultHeaders struct {
	headers map[string]string
}

func (m defaultHeaders) middleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if len(m.headers) > 0 {
			for key, value := range m.headers {
				if _, ok := r.Header[key]; !ok {
					w.Header().Set(key, value)
				}
			}
		}

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
