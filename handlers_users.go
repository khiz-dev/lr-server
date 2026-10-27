package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/khiz-dev/goserve/internal/helper"
)


func (apiCfg *apiConfig) handleCreateUser(w http.ResponseWriter, r *http.Request) {
    type createUserRequest struct {
        Email string `json:"email"`
    }

    type createUserResponse struct {
        ID        string `json:"id"`
        Email     string `json:"email"`
        CreatedAt string `json:"created_at"`
        UpdatedAt string `json:"updated_at"`
    }

    var req createUserRequest

    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusBadRequest)
        helper.ResponseBadRequest(w, "Invalid request payload")
        return
    }

    if req.Email == "" {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusBadRequest)
        helper.ResponseBadRequest(w, "Email is required")
        return
    }

    insertedUser, err := apiCfg.db.CreateUser(r.Context(), req.Email)
    if err != nil {
        fmt.Println("Failed to create user:", err)
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusInternalServerError)
        helper.ResponseInternalServerError(w, "Failed to create user")
        return
    }

    response := createUserResponse{
        ID:        insertedUser.ID.String(),
        Email:     insertedUser.Email,
        CreatedAt: insertedUser.CreatedAt.String(),
        UpdatedAt: insertedUser.UpdatedAt.String(),
    }

    helper.ResponseCreated(w, response)
}


func (apiCfg *apiConfig) handleUserReset(w http.ResponseWriter, r *http.Request) {
	if os.Getenv("ENV") != "dev" {
		w.WriteHeader(http.StatusForbidden)
		helper.ResponseBadRequest(w, "User reset is only allowed in development environment")
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	err := apiCfg.db.ResetUsers(r.Context())
	if err != nil {
		fmt.Println("Failed to reset users:", err)
		w.WriteHeader(http.StatusInternalServerError)
		helper.ResponseInternalServerError(w, "Failed to reset users")
		return
	}
	w.WriteHeader(http.StatusOK)
	helper.ResponseOK(w, "Users reset successfully")
}