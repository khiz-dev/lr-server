package auth

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type TokenType string

const (
	// TokenTypeAccess -
	TokenTypeAccess TokenType = "chirpy-access"
	AuthHeader = "Authorization"
)

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	signingKey := []byte(tokenSecret)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Subject:   userID.String(),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		Issuer:    "chirpy",
	})
	return token.SignedString(signingKey)
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	claimsStruct := jwt.RegisteredClaims{}
	parsedToken, err := jwt.ParseWithClaims(tokenString, &claimsStruct, func(token *jwt.Token) (interface{}, error) {
		return []byte(tokenSecret), nil
	})
	userIDString, err := parsedToken.Claims.GetSubject()
	if err != nil {
		return uuid.Nil, err
	}

	issuer, err := parsedToken.Claims.GetIssuer()
	if err != nil {
		return uuid.Nil, err
	}
	if issuer != "chirpy" {
		return uuid.Nil, errors.New("invalid issuer")
	}

	id, err := uuid.Parse(userIDString)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid user ID: %w", err)
	}
	return id, nil
}

func GetBearerToken(headers http.Header) (string, error) {
	authHeader := headers.Get(AuthHeader)
	if authHeader == "" {
		return "", errors.New("authorization header missing")
	}
	const prefix = "Bearer "
	if len(authHeader) <= len(prefix) || authHeader[:len(prefix)] != prefix {
		return "", errors.New("invalid authorization header format")
	}
	return authHeader[len(prefix):], nil
}

func GetUserIDFromContext(ctx context.Context) (uuid.UUID, error) {
	userIDString, ok := ctx.Value("userID").(string)
	if !ok {
		return uuid.Nil, errors.New("user ID not found in context")
	}
	return uuid.Parse(userIDString)
}

