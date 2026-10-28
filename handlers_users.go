package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/khiz-dev/goserve/internal/auth"
	"github.com/khiz-dev/goserve/internal/database"
	"github.com/khiz-dev/goserve/internal/helper"
)

type User struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}


func (apiCfg *apiConfig) handleCreateUser(w http.ResponseWriter, r *http.Request) {
    type createUserRequest struct {
        Email string `json:"email"`
		Password string `json:"password"`
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

    if req.Password == "" {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusBadRequest)
        helper.ResponseBadRequest(w, "Password is required")
        return
    }
	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		helper.ResponseInternalServerError(w, "Failed to hash password")
		return
	}
    insertedUser, err := apiCfg.db.CreateUser(r.Context(), database.CreateUserParams{
		Email: req.Email,
		HashedPassword: sql.NullString{String: hashedPassword, Valid: true},
	})
    if err != nil {
        fmt.Println("Failed to create user:", err)
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusInternalServerError)
        helper.ResponseInternalServerError(w, "Failed to create user")
        return
    }

    helper.ResponseCreated(w, User{
        ID:        insertedUser.ID.String(),
        Email:     insertedUser.Email,
        CreatedAt: insertedUser.CreatedAt.String(),
        UpdatedAt: insertedUser.UpdatedAt.String(),
    })
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

func (apiCfg *apiConfig) handleUserUpdate(w http.ResponseWriter, r *http.Request) {
	type updateUserRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var req updateUserRequest

	userID, err := auth.GetUserIDFromContext(r.Context())
	if err != nil {
		helper.ResponseUnauthorized(w, "User is not authenticated")
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helper.ResponseBadRequest(w, "Invalid request payload")
		return
	}

	if req.Email == "" {
		helper.ResponseBadRequest(w, "Empty email cannot be updated")
		return
	}

	if req.Password == "" {
		helper.ResponseBadRequest(w, "Empty password cannot be updated")
		return
	}

	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		helper.ResponseInternalServerError(w, "Failed to hash password")
		return
	}
	updatedUser, err := apiCfg.db.UpdateUser(r.Context(), database.UpdateUserParams{
		ID:             userID,
		Email:          req.Email,
		HashedPassword: sql.NullString{String: hashedPassword, Valid: true},
	})
	if err != nil {
		helper.ResponseInternalServerError(w, "Failed to update user")
		return
	}

	helper.ResponseOK(w, User{
		ID:        updatedUser.ID.String(),
		Email:     updatedUser.Email,
		CreatedAt: updatedUser.CreatedAt.String(),
		UpdatedAt: updatedUser.UpdatedAt.String(),
	})
}