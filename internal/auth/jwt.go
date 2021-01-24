package auth

import (
	"fmt"
	"time"

	"github.com/lestrrat-go/jwx/jwt"

	"github.com/go-chi/jwtauth"
	"github.com/woningfinder/woningfinder/internal/domain/entity"
)

const (
	userIDKey    = "user_id"
	userEmailKey = "user_email"
)

// CreateJWTAuthenticationToken builds a jwt authentication token
func CreateJWTAuthenticationToken(secret string) *jwtauth.JWTAuth {
	return jwtauth.New("HS256", []byte(secret), nil)
}

// CreateJWTUserToken builds an authentication token valid 12 hours for a given user
func CreateJWTUserToken(tokenAuth *jwtauth.JWTAuth, user *entity.User) (jwt.Token, string, error) {
	claims := map[string]interface{}{userIDKey: user.ID, userEmailKey: user.Email}
	jwtauth.SetExpiryIn(claims, time.Hour*12)
	jwtauth.SetIssuedNow(claims)
	token, tokenString, err := tokenAuth.Encode(claims)
	if err != nil {
		return nil, "", fmt.Errorf("error building user %s token: %w", user.Email, err)
	}

	return token, tokenString, nil
}

// ExtractUserFromJWT extracts an user from its JWT claims
func ExtractUserFromJWT(claims map[string]interface{}) (*entity.User, error) {
	userID, ok := claims[userIDKey].(uint)
	if !ok {
		return nil, fmt.Errorf("error extracting %s from claims, got %v", userIDKey, claims[userIDKey])
	}

	userEmail, ok := claims[userEmailKey].(string)
	if !ok {
		return nil, fmt.Errorf("error extracting %s from claims, got %v", userEmailKey, claims[userIDKey])
	}

	return &entity.User{ID: uint(userID), Email: userEmail}, nil
}
