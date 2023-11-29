package service

import (
	"errors"

	"github.com/eljamo/libpass/v2/config"
	"github.com/eljamo/libpass/v2/internal/stringcheck"
)

type SeparatorService interface {
	Separate(slice []string) ([]string, error)
}

type DefaultSeparatorService struct {
	cfg    *config.Config
	rngSvc RNGService
}

func NewSeparatorService(cfg *config.Config, rngSvc RNGService) (*DefaultSeparatorService, error) {
	svc := &DefaultSeparatorService{cfg, rngSvc}

	if err := svc.validate(); err != nil {
		return nil, err
	}

	return svc, nil
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
	if s.cfg.SeparatorCharacter == config.Random {
		sa := s.cfg.SeparatorAlphabet
		num, err := s.rngSvc.GenerateWithMax(len(sa))

		if err != nil {
			return "", err
		}

		return string(sa[num]), nil
	}

	return s.cfg.SeparatorCharacter, nil
}

func (s *DefaultSeparatorService) validate() error {
	if s.cfg.SeparatorCharacter != config.Random && len(s.cfg.SeparatorCharacter) > 1 {
		return errors.New("separator_character must be a single character if specified")
	}

	if s.cfg.SeparatorCharacter == config.Random {
		sa := s.cfg.SeparatorAlphabet
		if len(sa) == 0 {
			return errors.New("separator_alphabet cannot be empty")
		}

		chk := stringcheck.HasElementWithLengthGreaterThanOne(sa)
		if chk {
			return errors.New("separator_alphabet cannot contain elements with a length greater than 1")
		}
	}

	return nil
}
