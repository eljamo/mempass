package service

import (
	"fmt"

	"github.com/eljamo/libpass/v7/config"
	"github.com/eljamo/libpass/v7/config/option"
	"github.com/eljamo/libpass/v7/internal/validator"
)

// Defines the interface for a service that can separate elements of a string
// slice.
type SeparatorService interface {
	// Separate takes a slice of strings and inserts a separator character
	// between each element of the slice or returns an error if the slice
	// cannot be separated
	Separate(slice []string) ([]string, error)
}

// Implements the SeparatorService, providing functionality to separate string
// slices.
type DefaultSeparatorService struct {
	cfg    *config.Settings
	rngSvc RNGService
}

// Creates a new instance of DefaultSeparatorService. It validates the provided
// configuration and returns an error if the configuration is invalid.
func NewSeparatorService(cfg *config.Settings, rngSvc RNGService) (*DefaultSeparatorService, error) {
	svc := &DefaultSeparatorService{cfg, rngSvc}

	if err := svc.validate(); err != nil {
		return nil, err
	}

	return svc, nil
}

// Separate takes a slice of strings and inserts a separator character between
// each element of the slice. The separator character is determined based on the
// configuration. It returns the modified slice or an error if the separator
// character cannot be determined.
func (s *DefaultSeparatorService) Separate(slice []string) ([]string, error) {
	char, err := s.getSeparatorCharacter()
	if err != nil {
		return nil, fmt.Errorf("failed to get separator character: %w", err)
	}

	separatedSlice := make([]string, 0, len(slice))
	for _, element := range slice {
		separatedSlice = append(separatedSlice, char, element)
	}
	separatedSlice = append(separatedSlice, char)

	return separatedSlice, nil
}

// Returns the separator character based on the service configuration. It either
// returns a predefined character or a random character from a specified
// alphabet. Returns an error if it fails to return a random character.
func (s *DefaultSeparatorService) getSeparatorCharacter() (string, error) {
	if s.cfg.SeparatorCharacter == option.SeparatorCharacterRandom {
		num, err := s.rngSvc.GenerateWithMax(len(s.cfg.SeparatorAlphabet))
		if err != nil {
			return "", fmt.Errorf("failed to generate random number for separator character: %w", err)
		}

		return s.cfg.SeparatorAlphabet[num], nil
	}

	return s.cfg.SeparatorCharacter, nil
}

// Checks the configuration of the DefaultSeparatorService for  correctness.
// It ensures that the separator character is either a single  character or a
// valid random character from the alphabet. Returns an error if the
// configuration is invalid.
func (s *DefaultSeparatorService) validate() error {
	if s.cfg.SeparatorCharacter != option.SeparatorCharacterRandom && len(s.cfg.SeparatorCharacter) > 1 {
		return fmt.Errorf("%s must be a single character if specified", option.ConfigKeySeparatorCharacter)
	}

	if s.cfg.SeparatorCharacter == option.SeparatorCharacterRandom {
		if len(s.cfg.SeparatorAlphabet) == 0 {
			return fmt.Errorf("%s cannot be empty", option.ConfigKeySeparatorAlphabet)
		}

		chk := validator.HasElementWithLengthGreaterThanOne(s.cfg.SeparatorAlphabet)
		if chk {
			return fmt.Errorf("%s cannot contain elements with a length greater than 1", option.ConfigKeySeparatorAlphabet)
		}
	}

	return nil
}
