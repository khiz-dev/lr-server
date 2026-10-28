package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/khiz-dev/goserve/internal/auth"
	"github.com/khiz-dev/goserve/internal/database"
	"github.com/khiz-dev/goserve/internal/helper"
)

func (apiCfg *apiConfig) handleLogin(w http.ResponseWriter, r *http.Request) {
	type loginRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	type loginResponse struct {
		Token                       string `json:"token"`
		RefreshToken                string `json:"refresh_token"`
		TokenExpiryInSeconds        int    `json:"expires_in_seconds"`
		RefreshTokenExpiryInSeconds int    `json:"refresh_token_expires_in_seconds"`
		User
	}

	var req loginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helper.ResponseBadRequest(w, "Invalid request payload")
		return
	}

	if req.Email == "" || req.Password == "" {
		helper.ResponseBadRequest(w, "Email and password are required")
		return
	}

	defaultDuration := time.Duration(AccessTokenDefaultExpiresInSeconds) * time.Second

	user, err := apiCfg.db.GetUserByEmail(r.Context(), req.Email)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		helper.ResponseBadRequest(w, "Invalid email or password")
		return
	}

	samePassword, err := auth.ComparePassword(user.HashedPassword.String, req.Password)
	if err != nil || !samePassword {
		helper.ResponseUnauthorized(w, "Invalid email or password")
		return
	}

	jwtToken, err := auth.MakeJWT(user.ID, apiCfg.secret, defaultDuration)
	if err != nil {
		helper.ResponseInternalServerError(w, "Failed to generate JWT")
		return
	}

	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		helper.ResponseInternalServerError(w, "Failed to generate refresh token")
	}

	expiresAt := time.Now().Add(time.Duration(RefreshTokenDefaultExpiresInSeconds) * time.Second)
	_, err = apiCfg.db.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		UserID:    user.ID,
		Token:     refreshToken,
		ExpiresAt: expiresAt,
	})
	if err != nil {
		helper.ResponseInternalServerError(w, "Failed to store refresh token")
		return
	}

	helper.ResponseOK(w, loginResponse{
		Token:                       jwtToken,
		TokenExpiryInSeconds:        AccessTokenDefaultExpiresInSeconds,
		RefreshToken:                refreshToken,
		RefreshTokenExpiryInSeconds: RefreshTokenDefaultExpiresInSeconds,
		User: User{
			ID:        user.ID.String(),
			Email:     user.Email,
			CreatedAt: user.CreatedAt.String(),
			UpdatedAt: user.UpdatedAt.String(),
		},
	})
}

func (apiCfg *apiConfig) handleRefreshToken(w http.ResponseWriter, r *http.Request) {
	refreshToken := r.Header.Get(auth.AuthHeader)
	if refreshToken == "" {
		helper.ResponseBadRequest(w, "Refresh token is required")
		return
	}
}
