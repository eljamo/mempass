package service

import (
	"crypto/rand"
	"errors"
	"math/big"
)

type RNGService interface {
	GenerateWithMax(max int) (int, error)
	Generate() (int, error)
	GenerateSlice(length int) ([]int, error)
	GenerateSliceWithMax(length, max int) ([]int, error)
}

type DefaultRNGService struct{}

func NewRNGService() *DefaultRNGService {
	return &DefaultRNGService{}
}

var maxInt = int(^uint(0) >> 1)
var maxDigit = 10

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

func (s *DefaultRNGService) Generate() (int, error) {
	return s.GenerateWithMax(maxInt)
}

func (s *DefaultRNGService) GenerateSliceWithMax(length, max int) ([]int, error) {
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

func (s *DefaultRNGService) GenerateSlice(length int) ([]int, error) {
	return s.GenerateSliceWithMax(length, maxInt)
}
