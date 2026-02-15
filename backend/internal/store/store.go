package store

import (
	"fmt"
	"log"
	"slices"
	"sync"
)

type PackStore struct {
	mu     sync.RWMutex
	sizes  map[int]bool
	logger *log.Logger
}

func NewPacksStore(initial []int, logger *log.Logger) *PackStore {
	s := &PackStore{sizes: make(map[int]bool), logger: logger}

	for _, size := range initial {
		s.sizes[size] = true
	}

	return s
}

func (s *PackStore) GetSizes() []int {
	s.logger.Println("PackStore - Getting all pack sizes")

	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]int, 0, len(s.sizes))

	for size := range s.sizes {
		result = append(result, size)
	}

	slices.Sort(result)

	s.logger.Printf("PackStore - %d sizes retrieved", len(result))

	return result
}

func (s *PackStore) Add(size int) error {
	s.logger.Printf("PackStore - Adding pack size %d", size)

	if size <= 0 {
		err := fmt.Errorf("pack size must be a positive integer")

		s.logger.Printf("PackStore - Error adding pack size: %v", err)
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if s.sizes[size] {
		err := fmt.Errorf("pack size %d already exists", size)

		s.logger.Printf("PackStore - Error adding pack size: %v", err)
		return err
	}

	s.sizes[size] = true

	s.logger.Printf("PackStore - Pack size %d added", size)
	return nil
}

func (s *PackStore) Remove(size int) error {
	s.logger.Printf("PackStore - Removing pack size %d", size)

	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.sizes[size] {
		err := fmt.Errorf("pack size %d not found", size)

		s.logger.Printf("PackStore - Error removing pack size: %v", err)
		return err
	}

	if len(s.sizes) <= 1 {
		err := fmt.Errorf("cannot remove last pack size")

		s.logger.Printf("PackStore - Error removing pack size: %v", err)
		return err
	}

	delete(s.sizes, size)

	s.logger.Printf("PackStore - Pack size %d removed", size)
	return nil
}

func (s *PackStore) Set(sizes []int) error {
	s.logger.Printf("PackStore - Setting pack sizes to %v", sizes)

	if len(sizes) == 0 {
		err := fmt.Errorf("pack sizes cannot be empty")

		s.logger.Printf("PackStore - Error setting pack sizes: %v", err)
		return err
	}

	newSizes := make(map[int]bool)
	for _, size := range sizes {
		if size <= 0 {
			err := fmt.Errorf("all pack sizes must be positive integers")

			s.logger.Printf("PackStore - Error setting pack sizes: %v", err)
			return err
		}
		newSizes[size] = true
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.sizes = newSizes

	s.logger.Printf("PackStore - Pack sizes set to %v", sizes)
	return nil
}
