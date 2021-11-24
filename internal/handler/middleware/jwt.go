package middleware

import (
	"net/http"

	jwtauth "github.com/go-chi/jwtauth/v5"
	"github.com/go-chi/render"
	handlerErrors "github.com/julienrbrt/woningfinder/internal/handler/errors"
	"github.com/lestrrat-go/jwx/jwt"
)

// JWTVerifierMiddleware verify if the JWT token is present in the request
// It can be present in the header or in the query (jwt)
func JWTVerifierMiddleware(token *jwtauth.JWTAuth) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return jwtauth.Verify(token, jwtauth.TokenFromQuery, jwtauth.TokenFromHeader)(next)
	}
}

// CreateJWTValidatorMiddleware validates the JWT token
func CreateJWTValidatorMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, _, err := jwtauth.FromContext(r.Context())
		if err != nil {
			render.Render(w, r, handlerErrors.ErrUnauthorized)
			return
		}

		if token == nil || jwt.Validate(token) != nil {
			render.Render(w, r, handlerErrors.ErrUnauthorized)
			return
		}

		// Token is authenticated, pass it through
		next.ServeHTTP(w, r)
	})
}
