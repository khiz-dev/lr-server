package helper

import (
	"encoding/json"
	"net/http"
)

func ResponseOK(w http.ResponseWriter, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(payload)
	if err != nil {
		ResponseInternalServerError(w, "Invalid response payload")
	}
}

func ResponseNoContent(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}

func ResponseCreated(w http.ResponseWriter, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err := json.NewEncoder(w).Encode(payload)
	if err != nil {
		ResponseInternalServerError(w, "Invalid response payload")
	}
}


func ResponseNotFound(w http.ResponseWriter, payload any) {
	type notFoundResponse struct {
		Error string `json:"error"`
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(notFoundResponse{Error: payload.(string)})
}

func ResponseUnauthorized(w http.ResponseWriter, payload any) {
	type unauthorizedResponse struct {
		Error string `json:"error"`
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	json.NewEncoder(w).Encode(unauthorizedResponse{Error: payload.(string)})
}

func ResponseNoPermission(w http.ResponseWriter, payload any) {
	type noPermissionResponse struct {
		Error string `json:"error"`
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusForbidden)
	json.NewEncoder(w).Encode(noPermissionResponse{Error: payload.(string)})
}

func ResponseBadRequest(w http.ResponseWriter, payload any) {
	type badResponse struct {
		Error string `json:"error"`
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(badResponse{Error: payload.(string)})
}

func ResponseInternalServerError(w http.ResponseWriter, payload any) {
	type internalErrorResponse struct {
		Error string `json:"error"`
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(internalErrorResponse{Error: payload.(string)})
}