package main

import (
	"encoding/json"
	"net/http"

	"github.com/khiz-dev/goserve/internal/auth"
	"github.com/khiz-dev/goserve/internal/helper"
)

func (apiCfg *apiConfig) handleLogin(w http.ResponseWriter, r *http.Request) {
    type loginRequest struct {
        Email    string `json:"email"`
        Password string `json:"password"`
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

    helper.ResponseOK(w, User{
        ID:    user.ID.String(),
        Email: user.Email,
		CreatedAt: user.CreatedAt.String(),
		UpdatedAt: user.UpdatedAt.String(),
    })
}