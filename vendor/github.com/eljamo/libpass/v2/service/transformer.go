package service

import (
	"fmt"
	"slices"
	"strings"
	"unicode"

	"github.com/eljamo/libpass/v2/config"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type TransformerService interface {
	Transform(slice []string) ([]string, error)
}

type DefaultTransformerService struct {
	cfg    *config.Config
	rngSvc RNGService
}

func NewTransformerService(cfg *config.Config, rngSvc RNGService) (*DefaultTransformerService, error) {
	svc := &DefaultTransformerService{cfg, rngSvc}

	if err := svc.validate(); err != nil {
		return nil, err
	}

	return svc, nil
}

func (s *DefaultTransformerService) Transform(slice []string) ([]string, error) {
	switch s.cfg.CaseTransform {
	case config.Alternate:
		return s.alternate(slice), nil
	case config.AlternateLettercase:
		return alternateLettercase(slice)
	case config.Capitalise:
		return s.capitalise(slice), nil
	case config.CapitaliseInvert:
		return s.capitaliseInvert(slice)
	case config.Invert: // Same as CapitaliseInvert but reserved to maintain compatibility with xkpasswd.net configs
		return s.capitaliseInvert(slice)
	case config.Lower:
		return s.lower(slice), nil
	case config.LowerVowelUpperConsonant:
		return lowerVowelUpperConsonant(slice)
	case config.Random:
		return s.random(slice)
	case config.Sentence:
		return s.sentence(slice), nil
	case config.Upper:
		return s.upper(slice), nil
	case config.None:
	default:
		return slice, nil
	}

	return slice, nil
}

var caser = cases.Title(language.English)

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

func alternateLettercase(slice []string) ([]string, error) {
	var result []string
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
				return nil, err
			}
		}
		result = append(result, sb.String())
	}

	return result, nil
}

func (s *DefaultTransformerService) capitalise(slice []string) []string {
	for i, w := range slice {
		slice[i] = caser.String(w)
	}

	return slice
}

func (s *DefaultTransformerService) capitaliseInvert(slice []string) ([]string, error) {
	for i, w := range slice {
		var sb strings.Builder
		for j, r := range w {
			if j == 0 {
				_, err := sb.WriteRune(unicode.ToLower(r))
				if err != nil {
					return nil, err
				}
			} else {
				_, err := sb.WriteRune(unicode.ToUpper(r))
				if err != nil {
					return nil, err
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

func lowerVowelUpperConsonant(slice []string) ([]string, error) {
	var result []string
	for _, str := range slice {
		var sb strings.Builder
		for _, r := range str {
			if isVowel(r) {
				_, err := sb.WriteRune(unicode.ToLower(r))
				if err != nil {
					return nil, err
				}
			} else {
				_, err := sb.WriteRune(unicode.ToUpper(r))
				if err != nil {
					return nil, err
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
			return nil, err
		}

		if r%2 == 0 {
			slice[i] = strings.ToUpper(w)
		} else {
			slice[i] = strings.ToLower(w)
		}
	}

	return slice, nil
}

func (s *DefaultTransformerService) sentence(slice []string) []string {
	for i, w := range slice {
		if i == 0 {
			slice[i] = caser.String(w)
		} else {
			slice[i] = strings.ToLower(w)
		}
	}

	return slice
}

func (s *DefaultTransformerService) validate() error {
	if !slices.Contains(config.TransformType, s.cfg.CaseTransform) {
		return fmt.Errorf("not a valid case_transform type (%s)", s.cfg.CaseTransform)
	}

	return nil
}

func (s *DefaultTransformerService) upper(slice []string) []string {
	for i, w := range slice {
		slice[i] = strings.ToUpper(w)
	}

	return slice
}
