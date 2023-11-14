package service

import (
	"crypto/rand"
	"errors"
	"math/big"
)

type RNGService interface {
	GenerateN(max int) (int, error)
	Generate() (int, error)
}

type DefaultRNGService struct{}

func NewRNGService() *DefaultRNGService {
	return &DefaultRNGService{}
}

func (s *DefaultRNGService) GenerateN(max int) (int, error) {
	if max < 0 {
		return 0, errors.New("rng max cannot be less than 0")
	}

	n, err := rand.Int(rand.Reader, big.NewInt(int64(max)))
	if err != nil {
		return 0, err
	}

	return int(n.Int64()), nil
}

func (s *DefaultRNGService) Generate() (int, error) {
	max := int(^uint(0) >> 1)

	return s.GenerateN(max)
}
