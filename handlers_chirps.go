package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"

	"github.com/khiz-dev/goserve/internal/database"
	"github.com/khiz-dev/goserve/internal/helper"
)

func (apiCfg *apiConfig) handleCreateChirp(w http.ResponseWriter, r *http.Request) {
	type createChirpRequest struct {
		Body   string `json:"body"`
		UserID string `json:"user_id"`
	}
	type createChirpResponse struct {
		ID        string `json:"id"`
		Body      string `json:"body"`
		UserID    string `json:"user_id"`
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
	}

	req := createChirpRequest{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		helper.ResponseBadRequest(w, "Invalid request payload")
		return
	}

	err = validateChirp(req.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		helper.ResponseBadRequest(w, err.Error())
		return
	}
	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		helper.ResponseBadRequest(w, "Invalid user ID")
		return
	}
	chrip, err := apiCfg.db.CreateChirp(context.Background(), database.CreateChirpParams{
		Body: req.Body,
		UserID: userID,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		helper.ResponseInternalServerError(w, err.Error())
		return
	}
	resp := createChirpResponse{
		ID:        chrip.ID.String(),
		Body:      chrip.Body,
		UserID:    chrip.UserID.String(),
		CreatedAt: chrip.CreatedAt.String(),
		UpdatedAt: chrip.UpdatedAt.String(),
	}
	helper.ResponseCreated(w, resp)
}


func chirpHandler(w http.ResponseWriter, r *http.Request) {
	type body struct {
		Body string `json:"body"`
	}
	type errorResponse struct {
		Error string `json:"error"`
	}
	type successResponse struct {
		CleanedBody string `json:"cleaned_body"`
	}
	b := body{}
	err := json.NewDecoder(r.Body).Decode(&b)
	if err != nil {
		resp := errorResponse{Error: err.Error()}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)
		return
	}
	if len(b.Body) > 140 {
		resp := errorResponse{Error: "Chirp body exceeds 140 characters"}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)
		return
	}
	resp := successResponse{CleanedBody: helper.SanitizeChirp(b.Body)}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}