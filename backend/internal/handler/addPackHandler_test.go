package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"gs-app/backend/internal/store"
)

func TestAddPackHandler(t *testing.T) {
	tests := []struct {
		name           string
		body           string
		initialPacks   []int
		wantStatusCode int
		wantError      string
		wantPacks      []int
	}{
		{"invalid JSON should return 400", `not json`, []int{250}, http.StatusBadRequest, "invalid JSON body", nil},
		{"empty body should return 400", ``, []int{250}, http.StatusBadRequest, "invalid JSON body", nil},
		{"zero size should return 400", `{"size": 0}`, []int{250}, http.StatusBadRequest, "size must be a positive integer", nil},
		{"negative size should return 400", `{"size": -5}`, []int{250}, http.StatusBadRequest, "size must be a positive integer", nil},
		{"duplicate size should return 409", `{"size": 250}`, []int{250}, http.StatusConflict, "pack size 250 already exists", nil},
		{"valid size should return 201", `{"size": 500}`, []int{250}, http.StatusCreated, "", []int{250, 500}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &APIConfig{
				PackStore: store.NewPacksStore(tt.initialPacks, discardLogger),
			}

			req := httptest.NewRequest(http.MethodPost, "/v1/packs", bytes.NewBufferString(tt.body))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			cfg.AddPackHandler(rec, req)

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
