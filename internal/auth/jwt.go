package auth

import (
	"fmt"
	"strconv"
	"time"

	"github.com/go-chi/jwtauth"
	"github.com/lestrrat-go/jwx/jwt"
	"github.com/woningfinder/woningfinder/internal/customer"
)

const (
	userIDKey    = "user_id"
	userEmailKey = "user_email"
	loginTime    = 2
)

// CreateJWTAuthenticationToken builds a jwt authentication token
func CreateJWTAuthenticationToken(secret string) *jwtauth.JWTAuth {
	return jwtauth.New("HS256", []byte(secret), nil)
}

// CreateJWTUserToken builds an authentication token valid 2h for a given user
func CreateJWTUserToken(jwtAuth *jwtauth.JWTAuth, user *customer.User) (jwt.Token, string, error) {
	claims := map[string]interface{}{userIDKey: user.ID, userEmailKey: user.Email}
	jwtauth.SetExpiryIn(claims, time.Hour*loginTime)
	jwtauth.SetIssuedNow(claims)
	token, tokenString, err := jwtAuth.Encode(claims)
	if err != nil {
		return nil, "", fmt.Errorf("error building user %s token: %w", user.Email, err)
	}

	return token, tokenString, nil
}

// ExtractUserFromJWT extracts an user from its JWT claims
func ExtractUserFromJWT(claims map[string]interface{}) (*customer.User, error) {
	userID, err := strconv.ParseUint(fmt.Sprintf("%v", claims[userIDKey]), 10, 64)
	if err != nil {
		return nil, fmt.Errorf("error when parsing %v to uint: %w", claims[userIDKey], err)
	}

	userEmail, ok := claims[userEmailKey].(string)
	if !ok {
		return nil, fmt.Errorf("error extracting %s from claims, got %v", userEmailKey, claims[userIDKey])
	}

	return &customer.User{ID: uint(userID), Email: userEmail}, nil
}
