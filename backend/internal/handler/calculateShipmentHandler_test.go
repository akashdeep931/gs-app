package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"gs-app/backend/internal/store"
)

func TestCalculateShipmentHandler(t *testing.T) {
	defaultPacks := []int{250, 500, 1000, 2000, 5000}

	tests := []struct {
		name           string
		body           string
		packSizes      []int
		wantStatusCode int
		wantError      string
		wantPacks      map[string]int
		wantShipped    int
	}{
		{"invalid JSON should return 400", `not json`, defaultPacks, http.StatusBadRequest, "invalid JSON body", nil, 0},
		{"empty body should return 400", ``, defaultPacks, http.StatusBadRequest, "invalid JSON body", nil, 0},
		{"zero items should return 400", `{"items": 0}`, defaultPacks, http.StatusBadRequest, "items must be a positive integer", nil, 0},
		{"negative items should return 400", `{"items": -5}`, defaultPacks, http.StatusBadRequest, "items must be a positive integer", nil, 0},

		{"1 item should return smallest pack", `{"items": 1}`, defaultPacks, http.StatusOK, "", map[string]int{"250": 1}, 250},
		{"250 items should return one 250 pack", `{"items": 250}`, defaultPacks, http.StatusOK, "", map[string]int{"250": 1}, 250},
		{"251 items should round up to 500 pack", `{"items": 251}`, defaultPacks, http.StatusOK, "", map[string]int{"500": 1}, 500},
		{"501 items should return 500 and 250 packs", `{"items": 501}`, defaultPacks, http.StatusOK, "", map[string]int{"500": 1, "250": 1}, 750},
		{"12001 items should return multiple packs", `{"items": 12001}`, defaultPacks, http.StatusOK, "", map[string]int{"5000": 2, "2000": 1, "250": 1}, 12250},
		{"custom packs 3 and 5 with 7 items should return 8 shipped", `{"items": 7}`, []int{3, 5}, http.StatusOK, "", map[string]int{"3": 1, "5": 1}, 8},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &APIConfig{
				PackStore: store.NewPacksStore(tt.packSizes, discardLogger),
			}

			req := httptest.NewRequest(http.MethodPost, "/v1/calculate", bytes.NewBufferString(tt.body))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			cfg.CalculateShipmentHandler(rec, req)

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

			var got map[string]json.RawMessage
			if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
				t.Fatalf("failed to decode response: %v", err)
			}

			var gotPacks map[string]int
			if err := json.Unmarshal(got["packs"], &gotPacks); err != nil {
				t.Fatalf("failed to decode packs: %v", err)
			}

			if len(gotPacks) != len(tt.wantPacks) {
				t.Errorf("expected %v, got %v", tt.wantPacks, gotPacks)
			}
			for size, qty := range tt.wantPacks {
				if gotPacks[size] != qty {
					t.Errorf("expected packs[%s] = %d, got %d", size, qty, gotPacks[size])
				}
			}

			var gotShipped float64
			if err := json.Unmarshal(got["items_shipped"], &gotShipped); err != nil {
				t.Fatalf("failed to decode items_shipped: %v", err)
			}
			if int(gotShipped) != tt.wantShipped {
				t.Errorf("expected items_shipped %d, got %d", tt.wantShipped, int(gotShipped))
			}
		})
	}
}
