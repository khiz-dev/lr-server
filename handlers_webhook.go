package main

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/khiz-dev/goserve/internal/helper"
)

func (apiCfg *apiConfig) handleWebhook(w http.ResponseWriter, r *http.Request) {
	// Implement webhook handling logic here
	type webhookPayload struct {
		Event string `json:"event"`
		Data  struct {
			UserID string `json:"user_id"`
		} `json:"data"`
	}

	var v webhookPayload
	err := json.NewDecoder(r.Body).Decode(&v)
	if err != nil {
		helper.ResponseBadRequest(w, "Invalid request payload")
		return
	}

	userID, err := uuid.Parse(v.Data.UserID)
	if err != nil {
		helper.ResponseBadRequest(w, "Invalid user ID")
		return
	}

	user, err := apiCfg.db.GetUserByID(r.Context(), userID)
	if err != nil {
		helper.ResponseNotFound(w, "User not found")
		return
	}

	if v.Event != "user.upgraded" {
		helper.ResponseNoContent(w)
		return
	}

	err = apiCfg.db.MarkChirpRed(r.Context(), user.ID)
	if err != nil {
		helper.ResponseBadRequest(w, "Failed to mark chirp red")
		return
	}
	helper.ResponseNoContent(w)
}