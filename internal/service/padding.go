package service

import (
	"fmt"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/eljamo/mempass/internal/config"
)

type PaddingService interface {
	Pad(slice []string) (string, error)
}

type DefaultPaddingService struct {
	cfg *config.Config
	rng RNGService
}

func NewPaddingService(cfg *config.Config, rng RNGService) *DefaultPaddingService {
	return &DefaultPaddingService{cfg, rng}
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
	if before < 0 || after < 0 {
		return nil, fmt.Errorf("padding_digits_before and padding_digits_after must be greater than or equal to 0")
	}

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
		num, err := s.rng.GenerateN(10)
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
	case config.FIXED:
		return s.fixed(pw, char, s.cfg.PaddingCharactersBefore, s.cfg.PaddingCharactersAfter)
	case config.ADAPTIVE:
		return s.adaptive(pw, char, s.cfg.PadToLength)
	case config.NONE:
		return pw, nil
	}

	return pw, nil
}

func (s *DefaultPaddingService) getPaddingCharacter() (string, error) {
	if s.cfg.PaddingCharacter == config.RANDOM {
		num, err := s.rng.GenerateN(len(s.cfg.SymbolAlphabet))
		if err != nil {
			return "", err
		}
		return string(s.cfg.SymbolAlphabet[num]), nil
	}

	return s.cfg.PaddingCharacter, nil
}

func (s *DefaultPaddingService) fixed(pw string, char string, before int, after int) (string, error) {
	if before < 0 || after < 0 {
		return "", fmt.Errorf("padding_characters_before and padding_characters_after must be greater than or equal to 0")
	}

	paddingBefore := strings.Repeat(char, before)
	paddingAfter := strings.Repeat(char, after)

	return paddingBefore + pw + paddingAfter, nil
}

func (s *DefaultPaddingService) adaptive(pw string, char string, padLen int) (string, error) {
	if padLen < 0 {
		return "", fmt.Errorf("pad_to_length must be greater than or equal to 0")
	}

	pwLen := utf8.RuneCountInString(pw)
	if padLen <= pwLen {
		return pw, nil
	}

	diff := padLen - pwLen
	padding := strings.Repeat(char, diff)

	return pw + padding, nil
}
