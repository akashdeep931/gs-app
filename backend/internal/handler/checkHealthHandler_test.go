package handler

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"gs-app/backend/internal/store"
)

var discardLogger = log.New(io.Discard, "", 0)

func TestCheckHealthHandler(t *testing.T) {
	tests := []struct {
		name           string
		wantStatusCode int
		wantBody       map[string]string
	}{
		{"should return 200 with ok status", http.StatusOK, map[string]string{"status": "ok"}},
	}

	cfg := &APIConfig{
		PackStore: store.NewPacksStore([]int{250}, discardLogger),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/v1/health", nil)
			rec := httptest.NewRecorder()

			cfg.CheckHealthHandler(rec, req)

			if rec.Code != tt.wantStatusCode {
				t.Errorf("expected status %d, got %d", tt.wantStatusCode, rec.Code)
			}

			var got map[string]string
			if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
				t.Fatalf("failed to decode response body: %v", err)
			}

			for k, v := range tt.wantBody {
				if got[k] != v {
					t.Errorf("expected body[%q] = %q, got %q", k, v, got[k])
				}
			}
		})
	}
}
