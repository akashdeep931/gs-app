package handler

import (
	"encoding/json"
	"gs-app/backend/internal/middleware"
	"gs-app/backend/internal/model"
	"net/http"
	"strings"
)

func (cfg *APIConfig) AddPackHandler(w http.ResponseWriter, r *http.Request) {
	logger := middleware.GetLogger(r)
	logger.Println("Attempting to add pack size")

	var req model.AddPackRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Printf("Invalid request body: %v", err)

		respondWithJSON(w, r, http.StatusBadRequest, map[string]string{
			"error": "invalid JSON body",
		})
		return
	}

	if req.Size <= 0 {
		logger.Printf("Invalid pack size: %d", req.Size)

		respondWithJSON(w, r, http.StatusBadRequest, map[string]string{
			"error": "size must be a positive integer",
		})
		return
	}

	if err := cfg.PackStore.Add(req.Size); err != nil {
		logger.Printf("Failed to add pack size: %v", err)

		status := http.StatusBadRequest
		if strings.Contains(err.Error(), "already exists") {
			status = http.StatusConflict
		}

		respondWithJSON(w, r, status, map[string]string{
			"error": err.Error(),
		})
		return
	}

	logger.Printf("Pack size %d added", req.Size)

	respondWithJSON(w, r, http.StatusCreated, map[string]any{
		"packs": cfg.PackStore.GetAll(),
	})
}
