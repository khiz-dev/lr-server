package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/khiz-dev/goserve/internal/auth"
	"github.com/khiz-dev/goserve/internal/helper"
)

func (apiCfg *apiConfig) handleLogin(w http.ResponseWriter, r *http.Request) {
    type loginRequest struct {
        Email    string `json:"email"`
        Password string `json:"password"`
		ExpiresInSeconds int `json:"expires_in_seconds"`
    }
	type loginResponse struct {
		Token string `json:"token"`
		ExpiresInSeconds int `json:"expires_in_seconds"`
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

	defaultExpiresInSeconds := DefaultExpiresInSeconds
	if req.ExpiresInSeconds > 0 && req.ExpiresInSeconds <= DefaultExpiresInSeconds {
		defaultExpiresInSeconds = req.ExpiresInSeconds
	}
	defaultDuration := time.Duration(defaultExpiresInSeconds) * time.Second

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

    helper.ResponseOK(w, loginResponse{
        Token: jwtToken,
        ExpiresInSeconds: defaultExpiresInSeconds,
        User: User{
            ID:    user.ID.String(),
            Email: user.Email,
            CreatedAt: user.CreatedAt.String(),
            UpdatedAt: user.UpdatedAt.String(),
        },
    })
}