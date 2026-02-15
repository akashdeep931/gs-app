package handler

import (
	"encoding/json"
	"gs-app/backend/internal/middleware"
	"gs-app/backend/internal/model"
	calculator "gs-app/backend/internal/packSizeCalculator"
	"net/http"
)

func (cfg *APIConfig) CalculateHandler(w http.ResponseWriter, r *http.Request) {
	logger := middleware.GetLogger(r)
	logger.Println("Attempting packs calculation")

	var req model.CalculateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Printf("Invalid request body: %v", err)

		respondWithJSON(w, r, http.StatusBadRequest, map[string]string{
			"error": "invalid JSON body",
		})
		return
	}

	if req.Items <= 0 {
		logger.Printf("Invalid items count: %d", req.Items)

		respondWithJSON(w, r, http.StatusBadRequest, map[string]string{
			"error": "items must be a positive integer",
		})
		return
	}

	packSizes := cfg.PackStore.GetSizes()
	packs := calculator.CalculatePackSize(req.Items, packSizes)

	itemsShipped := 0
	for size, qty := range packs {
		itemsShipped += size * qty
	}

	logger.Printf("Shipping %d items for order of %d", itemsShipped, req.Items)

	respondWithJSON(w, r, http.StatusOK, map[string]any{
		"packs":         packs,
		"items_shipped": itemsShipped,
	})
}
