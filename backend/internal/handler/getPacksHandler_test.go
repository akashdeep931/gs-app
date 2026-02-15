package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"gs-app/backend/internal/store"
)

func TestGetPacksHandler(t *testing.T) {
	tests := []struct {
		name           string
		packSizes      []int
		wantStatusCode int
		wantSizes      []int
	}{
		{
			"should return default pack sizes",
			[]int{250, 500, 1000, 2000, 5000},
			http.StatusOK,
			[]int{250, 500, 1000, 2000, 5000},
		},
		{
			"should return single pack size",
			[]int{100},
			http.StatusOK,
			[]int{100},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &APIConfig{
				PackStore: store.NewPacksStore(tt.packSizes, discardLogger),
			}

			req := httptest.NewRequest(http.MethodGet, "/v1/packs", nil)
			rec := httptest.NewRecorder()

			cfg.GetPacksHandler(rec, req)

			if rec.Code != tt.wantStatusCode {
				t.Errorf("expected status %d, got %d", tt.wantStatusCode, rec.Code)
			}

			var got map[string]any
			if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
				t.Fatalf("failed to decode response body: %v", err)
			}

			rawSizes, ok := got["pack_sizes"].([]any)
			if !ok {
				t.Fatalf("expected pack_sizes to be an array, got %T", got["pack_sizes"])
			}

			if len(rawSizes) != len(tt.wantSizes) {
				t.Errorf("expected %d sizes, got %d", len(tt.wantSizes), len(rawSizes))
			}

			for i, want := range tt.wantSizes {
				if int(rawSizes[i].(float64)) != want {
					t.Errorf("expected sizes[%d] = %d, got %v", i, want, rawSizes[i])
				}
			}
		})
	}
}
