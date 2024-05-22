package service

import (
	"fmt"
	"slices"
	"strings"
	"sync"
	"unicode"

	"github.com/eljamo/libpass/v7/config"
	"github.com/eljamo/libpass/v7/config/option"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// Setup a pool of casers to be used for capitalisation
var caserPool = &sync.Pool{
	New: func() any {
		c := cases.Title(language.English)
		return &c
	},
}

// Defines an interface for transforming a slice of strings
type TransformerService interface {
	// Transform takes a slice of strings and transforms each element or returns
	// an error if the transformation fails.
	Transform(slice []string) ([]string, error)
}

// Implements the TransformerService, providing functionality to transform
// string slices based on a predefined configuration.
type DefaultTransformerService struct {
	cfg    *config.Settings
	rngSvc RNGService
}

// Creates a new valid instance of DefaultTransformerService with the given
// configuration and RNG service
func NewTransformerService(cfg *config.Settings, rngSvc RNGService) (*DefaultTransformerService, error) {
	svc := &DefaultTransformerService{cfg, rngSvc}

	if err := svc.validate(); err != nil {
		return nil, err
	}

	return svc, nil
}

// Transform takes a slice of strings and transforms each element
// according to the configured transformation rule.
// Returns the transformed slice or an error if the transformation fails.
//
// Transform Types:
//   - Alternate
//   - AlternateLettercase
//   - Capitalise
//   - CapitaliseInvert
//   - Invert
//   - Lower
//   - LowerVowelUpperConsonant
//   - Random
//   - Sentence
//   - Upper
func (s *DefaultTransformerService) Transform(slice []string) ([]string, error) {
	switch s.cfg.CaseTransform {
	case option.CaseTransformAlternate:
		return s.alternate(slice), nil
	case option.CaseTransformAlternateLettercase:
		return s.alternateLettercase(slice)
	case option.CaseTransformCapitalise:
		return s.capitalise(slice), nil
	case option.CaseTransformCapitaliseInvert:
		return s.capitaliseInvert(slice)
	case option.CaseTransformInvert: // Same as CapitaliseInvert but reserved to maintain compatibility with xkpasswd.net configs
		return s.capitaliseInvert(slice)
	case option.CaseTransformLower:
		return s.lower(slice), nil
	case option.CaseTransformLowerVowelUpperConsonant:
		return s.lowerVowelUpperConsonant(slice)
	case option.CaseTransformRandom:
		return s.random(slice)
	case option.CaseTransformSentence:
		return s.sentence(slice), nil
	case option.CaseTransformUpper:
		return s.upper(slice), nil
	}

	return slice, nil
}

// alternate applies alternating casing to each element of the slice.
//
// Example Output: string[]{"hello", "WORLD"}
func (s *DefaultTransformerService) alternate(slice []string) []string {
	for i, w := range slice {
		if i%2 == 0 {
			slice[i] = strings.ToLower(w)
		} else {
			slice[i] = strings.ToUpper(w)
		}
	}

	return slice
}

// alternateLettercase takes a slice of strings and alternates the casing of
// each letter within each string. Starting with lowercase, it switches
// between lowercase and uppercase for each subsequent letter.
// The function returns a new slice of strings with the applied transformations
// or an error if an issue occurs during string building.
//
// Example Output: string[]{"hElLo", "WoRlD"}, nil
func (s *DefaultTransformerService) alternateLettercase(slice []string) ([]string, error) {
	result := make([]string, 0, len(slice))
	for _, str := range slice {
		var sb strings.Builder
		upper := false
		for _, r := range str {
			var err error
			if unicode.IsLetter(r) {
				if upper {
					r = unicode.ToUpper(r)
				} else {
					r = unicode.ToLower(r)
				}
				upper = !upper
			}
			_, err = sb.WriteRune(r)
			if err != nil {
				return nil, fmt.Errorf("failed to write rune to string builder: %w", err)
			}
		}
		result = append(result, sb.String())
	}

	return result, nil
}

// Capitialises each element in the slice
//
// Example Output: string[]{"Hello", "World"}
func (s *DefaultTransformerService) capitalise(slice []string) []string {
	caser := caserPool.Get().(*cases.Caser)
	for i, w := range slice {
		slice[i] = caser.String(w)
	}

	caserPool.Put(caser)

	return slice
}

// Inverts the casing of a capitialised string in the slice
//
// Example output: string[]{"hELLO", "wORLD"}, nil
func (s *DefaultTransformerService) capitaliseInvert(slice []string) ([]string, error) {
	for i, w := range slice {
		var sb strings.Builder
		for j, r := range w {
			if j == 0 {
				_, err := sb.WriteRune(unicode.ToLower(r))
				if err != nil {
					return nil, fmt.Errorf("failed to write rune to string builder: %w", err)
				}
			} else {
				_, err := sb.WriteRune(unicode.ToUpper(r))
				if err != nil {
					return nil, fmt.Errorf("failed to write rune to string builder: %w", err)
				}
			}
		}
		slice[i] = sb.String()
	}

	return slice, nil
}

func (s *DefaultTransformerService) lower(slice []string) []string {
	for i, w := range slice {
		slice[i] = strings.ToLower(w)
	}

	return slice
}

func isVowel(r rune) bool {
	return strings.ContainsRune("aeiouAEIOU", r)
}

// lowerVowelUpperConsonant processes a slice of strings, transforming each string
// by applying lowercase to vowels and uppercase to consonants.
// It iterates through each rune in a string, checks if it is a vowel using
// the isVowel function, and accordingly changes its case.
// The function returns the transformed slice of strings or an error if any
// occurs during the string building process.
//
// Example Output: string[]{"hEllO", "wOrld"}, nil
func (s *DefaultTransformerService) lowerVowelUpperConsonant(slice []string) ([]string, error) {
	result := make([]string, 0, len(slice))
	for _, str := range slice {
		var sb strings.Builder
		for _, r := range str {
			if isVowel(r) {
				_, err := sb.WriteRune(unicode.ToLower(r))
				if err != nil {
					return nil, fmt.Errorf("failed to write rune to string builder: %w", err)
				}
			} else {
				_, err := sb.WriteRune(unicode.ToUpper(r))
				if err != nil {
					return nil, fmt.Errorf("failed to write rune to string builder: %w", err)
				}
			}
		}
		result = append(result, sb.String())
	}

	return result, nil
}

func (s *DefaultTransformerService) random(slice []string) ([]string, error) {
	for i, w := range slice {
		r, err := s.rngSvc.Generate()
		if err != nil {
			return nil, fmt.Errorf("failed to generate random number for random case: %w", err)
		}

		if r%2 == 0 {
			slice[i] = strings.ToUpper(w)
		} else {
			slice[i] = strings.ToLower(w)
		}
	}

	return slice, nil
}

// sentence applies sentence casing across each element of the slice
//
// Example Output: string[]{"Hello", "world"}
func (s *DefaultTransformerService) sentence(slice []string) []string {
	caser := caserPool.Get().(*cases.Caser)
	for i, w := range slice {
		if i == 0 {
			slice[i] = caser.String(w)
		} else {
			slice[i] = strings.ToLower(w)
		}
	}

	caserPool.Put(caser)

	return slice
}

func (s *DefaultTransformerService) validate() error {
	if !slices.Contains(option.TransformTypes, s.cfg.CaseTransform) {
		return fmt.Errorf("not a valid %s type (%s)", option.ConfigKeyCaseTransform, s.cfg.CaseTransform)
	}

	return nil
}

func (s *DefaultTransformerService) upper(slice []string) []string {
	for i, w := range slice {
		slice[i] = strings.ToUpper(w)
	}

	return slice
}
