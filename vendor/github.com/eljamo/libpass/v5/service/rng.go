package service

import (
	"crypto/rand"
	"errors"
	"math/big"
)

const (
	maxInt   = int(^uint(0) >> 1) // Maximum value for an int variable for the build architecture
	maxDigit = 10                 // Maximum digit value, used in GenerateDigit
)

// RNGService defines an interface for random number generation.
// It provides methods for generating random integers and slices of integers.
type RNGService interface {
	// Generates a random integer up to the specified maximum value.
	GenerateWithMax(max int) (int, error)
	// Generates a random integer
	Generate() (int, error)
	// Generates a single digit (0-9).
	GenerateDigit() (int, error)
	// Generates a slice of random integers
	GenerateSlice(length int) ([]int, error)
	// Generates a slice of random integers, each up to the specified maximum value.
	GenerateSliceWithMax(length int, max int) ([]int, error)
}

// DefaultRNGService is a struct implementing the RNGService interface.
type DefaultRNGService struct{}

// Creates a new instance of DefaultRNGService.
func NewRNGService() *DefaultRNGService {
	return &DefaultRNGService{}
}

// Generates a random integer up to the specified maximum value.
func (s *DefaultRNGService) GenerateWithMax(max int) (int, error) {
	if max < 1 {
		return 0, errors.New("rng max cannot be less than 1")
	}

	n, err := rand.Int(rand.Reader, big.NewInt(int64(max)))
	if err != nil {
		return 0, err
	}

	return int(n.Int64()), nil
}

// Generates a random integer with the maximum possible value for int.
func (s *DefaultRNGService) Generate() (int, error) {
	return s.GenerateWithMax(maxInt)
}

// GenerateDigit generates a single digit (0-9).
func (s *DefaultRNGService) GenerateDigit() (int, error) {
	return s.GenerateWithMax(maxDigit)
}

// Generates a slice of random integers, each up to the specified maximum value.
func (s *DefaultRNGService) GenerateSliceWithMax(length int, max int) ([]int, error) {
	if length < 0 {
		return nil, errors.New("rng slice length cannot be less than 0")
	}

	if max < 1 {
		return nil, errors.New("rng max cannot be less than 1")
	}

	if length == 0 {
		return []int{}, nil
	}

	slice := make([]int, length)
	for i := 0; i < length; i++ {
		n, err := s.GenerateWithMax(max)
		if err != nil {
			return nil, err
		}
		slice[i] = n
	}

	return slice, nil
}

// Generates a slice of random integers with the maximum possible value for int.
func (s *DefaultRNGService) GenerateSlice(length int) ([]int, error) {
	return s.GenerateSliceWithMax(length, maxInt)
}
