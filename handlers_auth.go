package main

import (
	"encoding/json"
	"fmt"
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
			IsChirpyRed: user.IsChirpyRed,
		},
	})
}

func (apiCfg *apiConfig) handleRefreshToken(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		helper.ResponseBadRequest(w, "Refresh token is required")
		return
	}
	if refreshToken == "" {
		helper.ResponseBadRequest(w, "Refresh token is required")
		return
	}

	storedToken, err := apiCfg.db.GetRefreshTokenByToken(r.Context(), refreshToken)
	if err != nil {
		fmt.Println(err)
		helper.ResponseUnauthorized(w, "Invalid refresh token")
		return
	}

	if storedToken.ExpiresAt.Before(time.Now()) || storedToken.RevokedAt.Valid {
		helper.ResponseUnauthorized(w, "Refresh token has expired")
		return
	}

	// Generate a new JWT token
	defaultDuration := time.Duration(AccessTokenDefaultExpiresInSeconds) * time.Second
	jwtToken, err := auth.MakeJWT(storedToken.UserID, apiCfg.secret, defaultDuration)
	if err != nil {
		helper.ResponseInternalServerError(w, "Failed to generate JWT")
		return
	}

	type refreshTokenResponse struct {
		Token                        string `json:"token"`
		ExpiresInSeconds             int    `json:"expires_in_seconds"`
		RefreshToken                 string `json:"refresh_token"`
		RefreshTokenExpiresInSeconds int    `json:"refresh_token_expires_in_seconds"`
	}

	helper.ResponseOK(w, refreshTokenResponse{
		Token:                        jwtToken,
		ExpiresInSeconds:             AccessTokenDefaultExpiresInSeconds,
		RefreshToken:                 refreshToken,
		RefreshTokenExpiresInSeconds: RefreshTokenDefaultExpiresInSeconds,
	})
}

func (apiCfg *apiConfig) handleRevokeRefreshToken(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		helper.ResponseBadRequest(w, "Refresh token is required")
		return
	}
	if refreshToken == "" {
		helper.ResponseBadRequest(w, "Refresh token is required")
		return
	}

	storedToken, err := apiCfg.db.GetRefreshTokenByToken(r.Context(), refreshToken)
	if err != nil {
		fmt.Println(err)
		helper.ResponseInternalServerError(w, "Failed to revoke refresh token")
		return
	}
	err = apiCfg.db.RevokeRefreshToken(r.Context(), storedToken.Token)
	if err != nil {
		helper.ResponseInternalServerError(w, "Failed to revoke refresh token")
		return
	}

	helper.ResponseNoContent(w)
}
