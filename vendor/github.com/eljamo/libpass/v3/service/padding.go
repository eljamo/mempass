package service

import (
	"errors"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/eljamo/libpass/v3/config"
	"github.com/eljamo/libpass/v3/internal/stringcheck"
)

type PaddingService interface {
	Pad(slice []string) (string, error)
}

type DefaultPaddingService struct {
	cfg    *config.Config
	rngSvc RNGService
}

func NewPaddingService(cfg *config.Config, rngSvc RNGService) (*DefaultPaddingService, error) {
	svc := &DefaultPaddingService{cfg, rngSvc}

	if err := svc.validate(); err != nil {
		return nil, err
	}

	return svc, nil
}

func (s *DefaultPaddingService) Pad(slice []string) (string, error) {
	pwd, err := s.digits(slice)
	if err != nil {
		return "", err
	}

	pwe := s.removeEdgeSeparatorCharacter(pwd)
	pwt := strings.TrimSpace(strings.Join(pwe, ""))

	pws, err := s.symbols(pwt)
	if err != nil {
		return "", err
	}

	return pws, nil
}

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

func (s *DefaultPaddingService) generateRandomDigits(count int) ([]string, error) {
	digits := make([]string, 0, count)
	for i := 0; i < count; i++ {
		num, err := s.rngSvc.GenerateWithMax(maxDigit)
		if err != nil {
			return nil, err
		}

		digits = append(digits, strconv.Itoa(num))
	}

	return digits, nil
}

func (s *DefaultPaddingService) removeEdgeSeparatorCharacter(slice []string) []string {
	if len(slice) == 0 {
		return slice
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

func (s *DefaultPaddingService) symbols(pw string) (string, error) {
	char, err := s.getPaddingCharacter()
	if err != nil {
		return "", err
	}

	switch s.cfg.PaddingType {
	case config.Fixed:
		return s.fixed(pw, char)
	case config.Adaptive:
		return s.adaptive(pw, char)
	case config.None:
		return pw, nil
	}

	return pw, nil
}

func (s *DefaultPaddingService) getPaddingCharacter() (string, error) {
	if s.cfg.PaddingCharacter == config.Random {
		num, err := s.rngSvc.GenerateWithMax(len(s.cfg.SymbolAlphabet))
		if err != nil {
			return "", err
		}
		return string(s.cfg.SymbolAlphabet[num]), nil
	}

	return s.cfg.PaddingCharacter, nil
}

func (s *DefaultPaddingService) fixed(pw string, char string) (string, error) {
	paddingBefore := strings.Repeat(char, s.cfg.PaddingCharactersBefore)
	paddingAfter := strings.Repeat(char, s.cfg.PaddingCharactersAfter)

	return paddingBefore + pw + paddingAfter, nil
}

func (s *DefaultPaddingService) adaptive(pw string, char string) (string, error) {
	padLen := s.cfg.PadToLength
	pwLen := utf8.RuneCountInString(pw)
	if padLen <= pwLen {
		return pw, nil
	}

	diff := padLen - pwLen
	padding := strings.Repeat(char, diff)

	return pw + padding, nil
}

func (s *DefaultPaddingService) validate() error {
	if s.cfg.PaddingType != config.None && s.cfg.PaddingCharacter != config.Random && len(s.cfg.PaddingCharacter) > 1 {
		return errors.New("padding_character must be a single character if specified")
	}

	if s.cfg.PaddingCharacter == config.Random {
		sa := s.cfg.SymbolAlphabet
		if len(sa) == 0 {
			return errors.New("symbol_alphabet cannot be empty")
		}

		chk := stringcheck.HasElementWithLengthGreaterThanOne(sa)
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
