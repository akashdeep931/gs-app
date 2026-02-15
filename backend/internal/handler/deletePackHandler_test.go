package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"gs-app/backend/internal/store"

	"github.com/go-chi/chi"
)

func TestDeletePackHandler(t *testing.T) {
	tests := []struct {
		name           string
		sizeParam      string
		initialPacks   []int
		wantStatusCode int
		wantError      string
		wantPacks      []int
	}{
		{"non-numeric size should return 400", "abc", []int{250, 500}, http.StatusBadRequest, "size must be a valid integer", nil},
		{"not found size should return 404", "999", []int{250, 500}, http.StatusNotFound, "pack size 999 not found", nil},
		{"last pack size should return 400", "250", []int{250}, http.StatusBadRequest, "cannot remove last pack size", nil},
		{"valid removal should return 200", "250", []int{250, 500}, http.StatusOK, "", []int{500}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &APIConfig{
				PackStore: store.NewPacksStore(tt.initialPacks, discardLogger),
			}

			req := httptest.NewRequest(http.MethodDelete, "/v1/packs/"+tt.sizeParam, nil)
			rec := httptest.NewRecorder()

			r := chi.NewRouter()
			r.Delete("/v1/packs/{size}", cfg.DeletePackHandler)
			r.ServeHTTP(rec, req)

			if rec.Code != tt.wantStatusCode {
				t.Errorf("expected status %d, got %d", tt.wantStatusCode, rec.Code)
			}

			if tt.wantError != "" {
				var got map[string]string
				if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
					t.Fatalf("failed to decode error response: %v", err)
				}
				if got["error"] != tt.wantError {
					t.Errorf("expected error %q, got %q", tt.wantError, got["error"])
				}
				return
			}

			var got map[string]any
			if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
				t.Fatalf("failed to decode response: %v", err)
			}

			rawPacks, ok := got["packs"].([]any)
			if !ok {
				t.Fatalf("expected packs to be an array, got %T", got["packs"])
			}

			if len(rawPacks) != len(tt.wantPacks) {
				t.Errorf("expected %d packs, got %d", len(tt.wantPacks), len(rawPacks))
			}

			for i, want := range tt.wantPacks {
				if int(rawPacks[i].(float64)) != want {
					t.Errorf("expected packs[%d] = %d, got %v", i, want, rawPacks[i])
				}
			}
		})
	}
}
