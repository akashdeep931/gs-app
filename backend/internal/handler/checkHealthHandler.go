package handler

import (
	"gs-app/backend/internal/middleware"
	"net/http"
)

func (cfg *APIConfig) CheckHealthHandler(w http.ResponseWriter, r *http.Request) {
	logger := middleware.GetLogger(r)
	logger.Println("Checking service health")

	respondWithJSON(w, r, http.StatusOK, map[string]string{"status": "ok"})

	logger.Println("Service is healthy")
}
