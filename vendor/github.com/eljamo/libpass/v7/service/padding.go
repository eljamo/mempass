package service

import (
	"fmt"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/eljamo/libpass/v7/config"
	"github.com/eljamo/libpass/v7/config/option"
	"github.com/eljamo/libpass/v7/internal/validator"
)

// Defines the interface for a service that provides functionality to pad a
// slice of strings.
type PaddingService interface {
	// Takes a slice of strings and applies extra characters or digits
	// and joins the slice, or returns an error
	Pad(slice []string) (string, error)
}

// Implements the PaddingService interface. It provides methods to add padding
// to strings based on predefined configuration settings.
type DefaultPaddingService struct {
	cfg    *config.Settings
	rngSvc RNGService
}

// Creates a new instance of DefaultPaddingService with the provided
// configuration and RNGService. It returns an error if the provided
// configuration is invalid.
func NewPaddingService(cfg *config.Settings, rngSvc RNGService) (*DefaultPaddingService, error) {
	svc := &DefaultPaddingService{cfg, rngSvc}

	if err := svc.validate(); err != nil {
		return nil, err
	}

	return svc, nil
}

// Takes a slice of strings and applies padding based on the service's
// configuration. It returns the padded string or an error if padding cannot
// be applied.
func (s *DefaultPaddingService) Pad(slice []string) (string, error) {
	pwd, err := s.digits(slice)
	if err != nil {
		return "", err
	}

	// Remove any separator characters which remain on the edges of the slice
	pwe := s.removeEdgeSeparatorCharacter(pwd)
	// Remove any whitespace characters which remain on the edges of the slice
	pwt := strings.TrimSpace(strings.Join(pwe, ""))

	pws, err := s.symbols(pwt)
	if err != nil {
		return "", err
	}

	return pws, nil
}

// Adds random digits before and after the given slice based on the
// configuration. It returns a slice with the added digits or an error if the
// digits cannot be generated.
func (s *DefaultPaddingService) digits(slice []string) ([]string, error) {
	paddedSlice := make([]string, 0, s.cfg.PaddingDigitsBefore+len(slice)+s.cfg.PaddingDigitsAfter)
	rdb, err := s.generateRandomDigits(s.cfg.PaddingDigitsBefore)
	if err != nil {
		return nil, err
	}

	paddedSlice = append(paddedSlice, rdb...)
	paddedSlice = append(paddedSlice, slice...)

	rda, err := s.generateRandomDigits(s.cfg.PaddingDigitsAfter)
	if err != nil {
		return nil, err
	}

	paddedSlice = append(paddedSlice, rda...)

	return paddedSlice, nil
}

// Generates a specified number of random digits. It returns a slice of the
// generated digits as strings or an error if the generation fails.
func (s *DefaultPaddingService) generateRandomDigits(num int) ([]string, error) {
	digits := make([]string, 0, num)
	for i := 0; i < num; i++ {
		num, err := s.rngSvc.GenerateDigit()
		if err != nil {
			return nil, err
		}

		digits = append(digits, strconv.Itoa(num))
	}

	return digits, nil
}

// Removes separator characters from the edges of the given slice. It handles
// both specified and random separator characters based on the configuration.
func (s *DefaultPaddingService) removeEdgeSeparatorCharacter(slice []string) []string {
	if len(slice) == 0 {
		return slice
	}

	if s.cfg.SeparatorCharacter == option.PaddingCharacterRandom {
		return s.removeRandomEdgeSeparatorCharacter(slice)
	}

	start, end := 0, len(slice)
	if slice[start] == s.cfg.SeparatorCharacter {
		start++
	}
	if end > start && slice[end-1] == s.cfg.SeparatorCharacter {
		end--
	}

	return slice[start:end]
}

// Specifically handles the removal of random edge separator characters. It is
// used when the separator character is set to random in the configuration.
func (s *DefaultPaddingService) removeRandomEdgeSeparatorCharacter(slice []string) []string {
	if len(slice) == 0 {
		return slice
	}

	if len(s.cfg.SeparatorAlphabet) == 0 {
		return slice
	}

	start, end := 0, len(slice)
	if validator.IsElementInSlice(s.cfg.SeparatorAlphabet, slice[start]) {
		start++
	}
	if end > start && validator.IsElementInSlice(s.cfg.SeparatorAlphabet, slice[end-1]) {
		end--
	}

	return slice[start:end]
}

// Applies symbol-based padding to the provided string as per the service
// configuration. It handles different padding types (fixed, adaptive) and
// returns the padded string or an error.
func (s *DefaultPaddingService) symbols(pw string) (string, error) {
	char, err := s.getPaddingCharacter()
	if err != nil {
		return "", err
	}

	switch s.cfg.PaddingType {
	case option.PaddingTypeFixed:
		return s.fixed(pw, char)
	case option.PaddingTypeAdaptive:
		return s.adaptive(pw, char), nil
	case option.PaddingTypeNone:
		return pw, nil
	}

	return pw, nil
}

// Retrieves the character to be used for padding. It selects a random character
// from the symbol alphabet if the padding character is set to random.
func (s *DefaultPaddingService) getPaddingCharacter() (string, error) {
	if s.cfg.PaddingCharacter == option.PaddingCharacterRandom {
		num, err := s.rngSvc.GenerateWithMax(len(s.cfg.SymbolAlphabet))
		if err != nil {
			return "", err
		}
		return string(s.cfg.SymbolAlphabet[num]), nil
	}

	return s.cfg.PaddingCharacter, nil
}

// Applies a fixed number of padding characters before and after the input
// string. It returns the resulting string with the specified fixed padding.
func (s *DefaultPaddingService) fixed(pw string, char string) (string, error) {
	paddingBefore := strings.Repeat(char, s.cfg.PaddingCharactersBefore)
	paddingAfter := strings.Repeat(char, s.cfg.PaddingCharactersAfter)

	return paddingBefore + pw + paddingAfter, nil
}

// Applies padding to the input string to meet a specified total length.
// The padding is added at the end of the string and uses the provided
// padding character.
func (s *DefaultPaddingService) adaptive(pw string, char string) string {
	pwLen := utf8.RuneCountInString(pw)
	if s.cfg.PadToLength <= pwLen {
		return pw
	}

	diff := s.cfg.PadToLength - pwLen
	padding := strings.Repeat(char, diff)

	return pw + padding
}

// Checks the service's configuration for any invalid values. It ensures the
// integrity of the padding settings before processing the padding operations.
func (s *DefaultPaddingService) validate() error {
	if s.cfg.PaddingType != option.PaddingTypeNone && s.cfg.PaddingCharacter != option.PaddingCharacterRandom && len(s.cfg.PaddingCharacter) > 1 {
		return fmt.Errorf("%s must be a single character if specified", option.ConfigKeyPaddingCharacter)
	}

	if s.cfg.PaddingCharacter == option.PaddingCharacterRandom {
		if len(s.cfg.SymbolAlphabet) == 0 {
			return fmt.Errorf("%s cannot be empty", option.ConfigKeySymbolAlphabet)
		}

		chk := validator.HasElementWithLengthGreaterThanOne(s.cfg.SymbolAlphabet)
		if chk {
			return fmt.Errorf("%s cannot contain elements with a length greater than 1", option.ConfigKeySymbolAlphabet)
		}
	}

	if s.cfg.PaddingDigitsBefore < 0 || s.cfg.PaddingDigitsAfter < 0 {
		return fmt.Errorf("%s and %s must be greater than or equal to 0", option.ConfigKeyPaddingDigitsBefore, option.ConfigKeyPaddingDigitsAfter)
	}

	if s.cfg.PaddingCharactersBefore < 0 || s.cfg.PaddingCharactersAfter < 0 {
		return fmt.Errorf("%s and %s must be greater than or equal to 0", option.ConfigKeyPaddingCharactersBefore, option.ConfigKeyPaddingCharactersAfter)
	}

	if s.cfg.PadToLength < 0 {
		return fmt.Errorf("%s must be greater than or equal to 0", option.ConfigKeyPadToLength)
	}

	return nil
}
