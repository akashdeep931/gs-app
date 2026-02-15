package handler

import (
	"encoding/json"
	"gs-app/backend/internal/middleware"
	"gs-app/backend/internal/model"
	"net/http"
)

func (cfg *APIConfig) SetPacksHandler(w http.ResponseWriter, r *http.Request) {
	logger := middleware.GetLogger(r)
	logger.Println("Attempting to set pack sizes")

	var req model.SetPacksRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Printf("Invalid request body: %v", err)

		respondWithJSON(w, r, http.StatusBadRequest, map[string]string{
			"error": "invalid JSON body",
		})
		return
	}

	if err := cfg.PackStore.Set(req.Sizes); err != nil {
		logger.Printf("Failed to set pack sizes: %v", err)

		respondWithJSON(w, r, http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
		return
	}

	logger.Printf("Pack sizes updated to %v", req.Sizes)

	respondWithJSON(w, r, http.StatusOK, map[string]any{
		"packs": cfg.PackStore.GetAll(),
	})
}
