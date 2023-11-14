package service

import (
	"errors"

	"github.com/eljamo/mempass/internal/config"
)

type SeparatorService interface {
	Separate(slice []string) ([]string, error)
}

type DefaultSeparatorService struct {
	cfg *config.Config
	rng RNGService
}

func NewSeparatorService(cfg *config.Config, rng RNGService) *DefaultSeparatorService {
	return &DefaultSeparatorService{cfg, rng}
}

func (s *DefaultSeparatorService) Separate(slice []string) ([]string, error) {
	char, err := s.getSeparatorCharacter()
	if err != nil {
		return nil, err
	}

	separatedSlice := make([]string, 0, len(slice))
	for _, element := range slice {
		separatedSlice = append(separatedSlice, char, element)
	}
	separatedSlice = append(separatedSlice, char)

	return separatedSlice, nil
}

func (s *DefaultSeparatorService) getSeparatorCharacter() (string, error) {
	if s.cfg.SeparatorCharacter == config.RANDOM {
		sa := s.cfg.SeparatorAlphabet
		if len(sa) == 0 {
			return "", errors.New("separator_alphabet cannot be empty")
		}

		num, err := s.rng.GenerateN(len(sa))
		if err != nil {
			return "", err
		}

		return string(sa[num]), nil
	}

	return s.cfg.SeparatorCharacter, nil
}
