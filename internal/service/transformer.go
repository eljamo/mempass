package service

import (
	"fmt"
	"slices"
	"strings"

	"github.com/eljamo/mempass/internal/config"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type TransformerService interface {
	Transform(slice []string) ([]string, error)
}

type DefaultTransformerService struct {
	cfg *config.Config
	rng RNGService
}

func NewTransformerService(cfg *config.Config, rng RNGService) *DefaultTransformerService {
	return &DefaultTransformerService{
		cfg,
		rng,
	}
}

func (s *DefaultTransformerService) Transform(slice []string) ([]string, error) {
	err := s.validate()
	if err != nil {
		return []string{}, err
	}

	switch s.cfg.CaseTransform {
	case config.ALTERNATE:
		return s.alternate(slice), nil
	case config.CAPITALISE:
		return s.capitalise(slice), nil
	case config.INVERT:
		return s.invert(slice), nil
	case config.LOWER:
		return s.lower(slice), nil
	case config.RANDOM:
		return s.random(slice)
	case config.UPPER:
		return s.upper(slice), nil
	case config.NONE:
	default:
		return slice, nil
	}

	return slice, err
}

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

func (s *DefaultTransformerService) capitalise(slice []string) []string {
	caser := cases.Title(language.English)
	for i := range slice {
		slice[i] = caser.String(slice[i])
	}

	return slice
}

func (s *DefaultTransformerService) invert(slice []string) []string {
	for i, w := range slice {
		if len(w) > 1 {
			slice[i] = string(w[0]) + strings.ToUpper(w[1:])
		}
	}

	return slice
}

func (s *DefaultTransformerService) lower(slice []string) []string {
	for i, w := range slice {
		slice[i] = strings.ToLower(w)
	}

	return slice
}

func (s *DefaultTransformerService) random(slice []string) ([]string, error) {
	for i, w := range slice {
		r, err := s.rng.Generate()
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

func (s *DefaultTransformerService) validate() error {
	if !slices.Contains(config.TransformType, s.cfg.CaseTransform) {
		return fmt.Errorf("Invalid case_transform type: %s", s.cfg.CaseTransform)
	}

	return nil
}

func (s *DefaultTransformerService) upper(slice []string) []string {
	for i, w := range slice {
		slice[i] = strings.ToUpper(w)
	}

	return slice
}
