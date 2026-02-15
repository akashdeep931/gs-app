package handler

import (
	"encoding/json"
	"net/http"
)

func respondWithJSON(w http.ResponseWriter, r *http.Request, statusCode int, data any) {
	w.Header().Set("Content-Type", "application/json")

	origin := r.Header.Get("Origin")
	if origin != "" {
		w.Header().Set("Access-Control-Allow-Origin", origin)
	}

	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}
