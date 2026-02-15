package packSizesCalculator

import (
	"testing"
)

func TestCalculatePackSizes(t *testing.T) {
	defaultPacks := []int{250, 500, 1000, 2000, 5000}

	tests := []struct {
		name       string
		orderSize  int
		packSizes  []int
		totalPacks map[int]int
		total      int
	}{
		{"1 item should return smallest pack", 1, defaultPacks, map[int]int{250: 1}, 250},
		{"250 items should return one 250 pack", 250, defaultPacks, map[int]int{250: 1}, 250},
		{"251 items should round up to 500 pack", 251, defaultPacks, map[int]int{500: 1}, 500},
		{"501 items should return 500 and 250 packs", 501, defaultPacks, map[int]int{500: 1, 250: 1}, 750},
		{"12001 items should return multiple packs", 12001, defaultPacks, map[int]int{5000: 2, 2000: 1, 250: 1}, 12250},

		{"0 items should return nil", 0, defaultPacks, nil, 0},
		{"negative items should return nil", -5, defaultPacks, nil, 0},
		{"empty pack sizes should return nil", 100, []int{}, nil, 0},

		{"5000 items should return exact match", 5000, defaultPacks, map[int]int{5000: 1}, 5000},
		{"1500 items should return exact match", 1500, defaultPacks, map[int]int{1000: 1, 500: 1}, 1500},
		{"10000 items should return exact match", 10000, defaultPacks, map[int]int{5000: 2}, 10000},

		{"7 items with custom packs 3 and 5 should return 8", 7, []int{3, 5}, map[int]int{3: 1, 5: 1}, 8},
		{"1 item with single pack size should return 1", 1, []int{1}, map[int]int{1: 1}, 1},
		{"1 item with single large pack should return 500", 1, []int{500}, map[int]int{500: 1}, 500},

		{"100000 items should return 20 packs of 5000", 100000, defaultPacks, map[int]int{5000: 20}, 100000},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := CalculatePackSizes(tt.orderSize, tt.packSizes)

			if tt.totalPacks == nil {
				if res != nil {
					t.Errorf("expected nil, got %v", res)
				}
				return
			}

			if len(res) != len(tt.totalPacks) {
				t.Errorf("expected %v, got %v", tt.totalPacks, res)
				return
			}

			for k, v := range tt.totalPacks {
				if res[k] != v {
					t.Errorf("expected %v, got %v", tt.totalPacks, res)
					return
				}
			}

			total := 0
			for size, qty := range res {
				total += size * qty
			}
			if total != tt.total {
				t.Errorf("expected total %d, got %d", tt.total, total)
			}
		})
	}
}
