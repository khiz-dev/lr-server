package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/uuid"

	"github.com/khiz-dev/goserve/internal/auth"
	"github.com/khiz-dev/goserve/internal/database"
	"github.com/khiz-dev/goserve/internal/helper"
)

type Chirp struct {
	ID        string `json:"id"`
	Body      string `json:"body"`
	UserID    string `json:"user_id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func (apiCfg *apiConfig) handleCreateChirp(w http.ResponseWriter, r *http.Request) {
	type createChirpRequest struct {
		Body   string `json:"body"`
	}

	req := createChirpRequest{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		helper.ResponseBadRequest(w, "Invalid request payload")
		return
	}

	err = validateChirp(req.Body)
	if err != nil {
		helper.ResponseBadRequest(w, err.Error())
		return
	}
	
	userID, err := auth.GetUserIDFromContext(r.Context())
	if err != nil {
		helper.ResponseBadRequest(w, "Invalid user ID")
		return
	}
	chrip, err := apiCfg.db.CreateChirp(context.Background(), database.CreateChirpParams{
		Body: req.Body,
		UserID: userID,
	})
	if err != nil {
		helper.ResponseInternalServerError(w, err.Error())
		return
	}
	resp := Chirp{
		ID:        chrip.ID.String(),
		Body:      chrip.Body,
		UserID:    chrip.UserID.String(),
		CreatedAt: chrip.CreatedAt.String(),
		UpdatedAt: chrip.UpdatedAt.String(),
	}
	helper.ResponseCreated(w, resp)
}

func (apiCfg *apiConfig) handleGetAllChrips(w http.ResponseWriter, r *http.Request) {
	authorIDParam	 := r.URL.Query().Get("author_id")
	authorID, err := uuid.Parse(authorIDParam)
	if authorIDParam != "" && err != nil {
		helper.ResponseBadRequest(w, "Invalid author ID")
		return
	}
	
	var chrips []database.Chrip
	if err == nil && authorIDParam != "" {
		chrips, err = apiCfg.db.GetChripsByUserID(context.Background(), authorID)
		if err != nil {
			helper.ResponseInternalServerError(w, err.Error())
			return
		}
	} else {
		chrips, err = apiCfg.db.GetAllChrips(context.Background())
		if err != nil {
			helper.ResponseInternalServerError(w, err.Error())
			return
		}
	}
	resp := make([]Chirp, len(chrips))
	for i, c := range chrips {
		resp[i] = Chirp{
			ID:        c.ID.String(),
			Body:      c.Body,
			UserID:    c.UserID.String(),
			CreatedAt: c.CreatedAt.String(),
			UpdatedAt: c.UpdatedAt.String(),
		}
	}
	helper.ResponseOK(w, resp)
}

func (apiCfg *apiConfig) handleGetChirpByID(w http.ResponseWriter, r *http.Request) {
	chirpID := r.PathValue("chirpID")
	id, err := uuid.Parse(chirpID)
	if err != nil {
		helper.ResponseBadRequest(w, err.Error())
		return
	}

	chrip, err := apiCfg.db.GetChirpByID(context.Background(), id)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		helper.ResponseNotFound(w, "Chirp not found")
		return
	}
	if err != nil {
		helper.ResponseInternalServerError(w, err.Error())
		return
	}
	resp := Chirp{
		ID:        chrip.ID.String(),
		Body:      chrip.Body,
		UserID:    chrip.UserID.String(),
		CreatedAt: chrip.CreatedAt.String(),
		UpdatedAt: chrip.UpdatedAt.String(),
	}
	helper.ResponseOK(w, resp)
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

func (apiCfg *apiConfig) handleDeleteChirp(w http.ResponseWriter, r *http.Request) {
	chirpID := r.PathValue("chirpID")
	id, err := uuid.Parse(chirpID)
	if err != nil {
		helper.ResponseBadRequest(w, err.Error())
		return
	}

	userID, err := auth.GetUserIDFromContext(r.Context())
	if err != nil {
		helper.ResponseUnauthorized(w, "Unauthorized")
		return
	}

	chirp, err := apiCfg.db.GetChirpByID(context.Background(), id)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		helper.ResponseNotFound(w, "Chirp not found")
		return
	}
	if err != nil {
		helper.ResponseInternalServerError(w, err.Error())
		return
	}
	if chirp.UserID != userID {
		helper.ResponseNoPermission(w, "No permission to delete this chirp")
		return
	}

	err = apiCfg.db.DeleteChirp(context.Background(), id)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		helper.ResponseNotFound(w, "Chirp not found")
		return
	}
	if err != nil {
		helper.ResponseInternalServerError(w, err.Error())
		return
	}
	helper.ResponseNoContent(w)
}