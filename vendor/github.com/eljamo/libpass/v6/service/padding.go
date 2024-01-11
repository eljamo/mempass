package service

import (
	"errors"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/eljamo/libpass/v6/config"
	"github.com/eljamo/libpass/v6/config/option"
	"github.com/eljamo/libpass/v6/internal/validator"
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
	before, after := s.cfg.PaddingDigitsBefore, s.cfg.PaddingDigitsAfter
	paddedSlice := make([]string, 0, before+len(slice)+after)

	rdb, err := s.generateRandomDigits(before)
	if err != nil {
		return nil, err
	}

	paddedSlice = append(paddedSlice, rdb...)
	paddedSlice = append(paddedSlice, slice...)

	rda, err := s.generateRandomDigits(after)
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

	sc := s.cfg.SeparatorCharacter
	if sc == option.Random {
		return s.removeRandomEdgeSeparatorCharacter(slice)
	}

	start, end := 0, len(slice)
	if slice[start] == sc {
		start++
	}
	if end > start && slice[end-1] == sc {
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

	sa := s.cfg.SeparatorAlphabet
	if len(sa) == 0 {
		return slice
	}

	start, end := 0, len(slice)
	if validator.IsElementInSlice(sa, slice[start]) {
		start++
	}
	if end > start && validator.IsElementInSlice(sa, slice[end-1]) {
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
	case option.Fixed:
		return s.fixed(pw, char)
	case option.Adaptive:
		return s.adaptive(pw, char), nil
	case option.None:
		return pw, nil
	}

	return pw, nil
}

// Retrieves the character to be used for padding. It selects a random character
// from the symbol alphabet if the padding character is set to random.
func (s *DefaultPaddingService) getPaddingCharacter() (string, error) {
	if s.cfg.PaddingCharacter == option.Random {
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
	padLen := s.cfg.PadToLength
	pwLen := utf8.RuneCountInString(pw)
	if padLen <= pwLen {
		return pw
	}

	diff := padLen - pwLen
	padding := strings.Repeat(char, diff)

	return pw + padding
}

// Checks the service's configuration for any invalid values. It ensures the
// integrity of the padding settings before processing the padding operations.
func (s *DefaultPaddingService) validate() error {
	if s.cfg.PaddingType != option.None && s.cfg.PaddingCharacter != option.Random && len(s.cfg.PaddingCharacter) > 1 {
		return errors.New("padding_character must be a single character if specified")
	}

	if s.cfg.PaddingCharacter == option.Random {
		sa := s.cfg.SymbolAlphabet
		if len(sa) == 0 {
			return errors.New("symbol_alphabet cannot be empty")
		}

		chk := validator.HasElementWithLengthGreaterThanOne(sa)
		if chk {
			return errors.New("symbol_alphabet cannot contain elements with a length greater than 1")
		}
	}

	if s.cfg.PaddingDigitsBefore < 0 || s.cfg.PaddingDigitsAfter < 0 {
		return errors.New("padding_digits_before and padding_digits_after must be greater than or equal to 0")
	}

	if s.cfg.PaddingCharactersBefore < 0 || s.cfg.PaddingCharactersAfter < 0 {
		return errors.New("padding_characters_before and padding_characters_after must be greater than or equal to 0")
	}

	if s.cfg.PadToLength < 0 {
		return errors.New("pad_to_length must be greater than or equal to 0")
	}

	return nil
}
