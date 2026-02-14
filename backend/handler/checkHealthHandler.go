package handler

import (
	"gs-app/backend/middleware"
	"net/http"
)

func CheckHealthHandler(w http.ResponseWriter, r *http.Request) {
	logger := middleware.GetLogger(r)
	logger.Println("Checking service health")

	respondWithJSON(w, r, http.StatusOK, map[string]string{"status": "ok"})

	logger.Println("Service is healthy")
}
