package store

import (
	"io"
	"log"
	"slices"
	"testing"
)

var discardLogger = log.New(io.Discard, "", 0)

func TestPacksStore(t *testing.T) {
	tests := []struct {
		name    string
		initial []int
		want    []int
	}{
		{"unsorted sizes should return sorted", []int{500, 250, 1000}, []int{250, 500, 1000}},
		{"many unsorted sizes should return sorted", []int{5000, 250, 2000, 500, 1000}, []int{250, 500, 1000, 2000, 5000}},
		{"duplicate sizes should be deduplicated", []int{250, 250, 500}, []int{250, 500}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewPacksStore(tt.initial, discardLogger)
			res := s.GetAll()

			if !slices.Equal(res, tt.want) {
				t.Errorf("expected %v, got %v", tt.want, res)
			}
		})
	}
}

func TestAdd(t *testing.T) {
	tests := []struct {
		name      string
		initial   []int
		add       int
		wantErr   bool
		wantSizes []int
	}{
		{"new size should be added", []int{250, 500}, 1000, false, []int{250, 500, 1000}},
		{"duplicate size should return error", []int{250, 500}, 250, true, []int{250, 500}},
		{"zero size should return error", []int{250}, 0, true, []int{250}},
		{"negative size should return error", []int{250}, -10, true, []int{250}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewPacksStore(tt.initial, discardLogger)
			err := s.Add(tt.add)

			if (err != nil) != tt.wantErr {
				t.Fatalf("expected error=%v, got %v", tt.wantErr, err)
			}

			res := s.GetAll()

			if !slices.Equal(res, tt.wantSizes) {
				t.Errorf("expected %v, got %v", tt.wantSizes, res)
			}
		})
	}
}

func TestRemove(t *testing.T) {
	tests := []struct {
		name      string
		initial   []int
		remove    int
		wantErr   bool
		wantSizes []int
	}{
		{"existing size should be removed", []int{250, 500, 1000}, 500, false, []int{250, 1000}},
		{"non-existent size should return error", []int{250, 500}, 9999, true, []int{250, 500}},
		{"last remaining size should return error", []int{250}, 250, true, []int{250}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewPacksStore(tt.initial, discardLogger)
			err := s.Remove(tt.remove)

			if (err != nil) != tt.wantErr {
				t.Fatalf("expected error=%v, got %v", tt.wantErr, err)
			}

			res := s.GetAll()

			if !slices.Equal(res, tt.wantSizes) {
				t.Errorf("expected %v, got %v", tt.wantSizes, res)
			}
		})
	}
}

func TestSet(t *testing.T) {
	tests := []struct {
		name      string
		initial   []int
		set       []int
		wantErr   bool
		wantSizes []int
	}{
		{"new sizes should replace existing", []int{250, 500}, []int{100, 200, 300}, false, []int{100, 200, 300}},
		{"duplicate sizes should be deduplicated", []int{250}, []int{100, 100, 200}, false, []int{100, 200}},
		{"empty sizes should return error", []int{250}, []int{}, true, []int{250}},
		{"negative size should return error", []int{250}, []int{100, -5, 200}, true, []int{250}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewPacksStore(tt.initial, discardLogger)
			err := s.Set(tt.set)

			if (err != nil) != tt.wantErr {
				t.Fatalf("expected error=%v, got %v", tt.wantErr, err)
			}

			res := s.GetAll()

			if !slices.Equal(res, tt.wantSizes) {
				t.Errorf("expected %v, got %v", tt.wantSizes, res)
			}
		})
	}
}
