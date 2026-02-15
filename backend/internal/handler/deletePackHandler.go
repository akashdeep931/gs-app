package handler

import (
	"gs-app/backend/internal/middleware"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi"
)

func (cfg *APIConfig) DeletePackHandler(w http.ResponseWriter, r *http.Request) {
	logger := middleware.GetLogger(r)
	logger.Println("Attempting to delete pack size")

	sizeParam := chi.URLParam(r, "size")

	size, err := strconv.Atoi(sizeParam)
	if err != nil {
		logger.Printf("Invalid size parameter: %s", sizeParam)

		respondWithJSON(w, r, http.StatusBadRequest, map[string]string{
			"error": "size must be a valid integer",
		})
		return
	}

	if err := cfg.PackStore.Remove(size); err != nil {
		logger.Printf("Failed to remove pack size: %v", err)

		status := http.StatusBadRequest
		if strings.Contains(err.Error(), "not found") {
			status = http.StatusNotFound
		}

		respondWithJSON(w, r, status, map[string]string{
			"error": err.Error(),
		})
		return
	}

	logger.Printf("Pack size %d removed", size)

	respondWithJSON(w, r, http.StatusOK, map[string]any{
		"packs": cfg.PackStore.GetSizes(),
	})
}
