package handler

import (
	"gs-app/backend/internal/middleware"
	"net/http"
)

func (cfg *APIConfig) GetPacksHandler(w http.ResponseWriter, r *http.Request) {
	logger := middleware.GetLogger(r)
	logger.Println("Fetching pack sizes")

	sizes := cfg.PackStore.GetSizes()

	logger.Printf("Returning %d pack sizes", len(sizes))

	respondWithJSON(w, r, http.StatusOK, map[string]any{
		"packs": sizes,
	})
}
