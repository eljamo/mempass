package service

import (
	"fmt"

	"github.com/eljamo/libpass/v7/config"
	"github.com/eljamo/libpass/v7/config/option"
)

// PasswordGeneratorService defines the interface for a service that generates
// passwords. It requires an implementation of the Generate method that returns
// a slice of strings representing generated passwords and an error if the
// generation process fails.
type PasswordGeneratorService interface {
	// Generate creates a list of passwords and returns the list or an error
	Generate() ([]string, error)
}

// DefaultPasswordGeneratorService implements the PasswordGeneratorService
// interface providing a concrete implementation for password generation. It
// combines various services like transformers, separators, padders, and word
// list services to generate passwords based on provided configuration.
type DefaultPasswordGeneratorService struct {
	cfg            *config.Settings
	transformerSvc TransformerService
	separatorSvc   SeparatorService
	paddingSvc     PaddingService
	wordListSvc    WordListService
}

const (
	numPasswordMin int = 1
	numPasswordMax int = 10
)

// NewCustomPasswordGeneratorService constructs a new instance of
// DefaultPasswordGeneratorService with the provided services and
// configuration. It validates the configuration and returns an
// error if the configuration is invalid (e.g., number of passwords
// is less than 1).
func NewCustomPasswordGeneratorService(
	cfg *config.Settings,
	transformerSvc TransformerService,
	separatorSvc SeparatorService,
	paddingSvc PaddingService,
	wordListSvc WordListService,
) (*DefaultPasswordGeneratorService, error) {
	if cfg.NumPasswords < numPasswordMin || cfg.NumPasswords > numPasswordMax {
		return nil, fmt.Errorf(
			"%s (%d) must be between %d and %d",
			option.ConfigKeyNumPasswords,
			cfg.NumPasswords,
			numPasswordMin,
			numPasswordMax,
		)
	}

	return &DefaultPasswordGeneratorService{
		cfg,
		transformerSvc,
		separatorSvc,
		paddingSvc,
		wordListSvc,
	}, nil
}

// NewPasswordGeneratorService constructs a DefaultPasswordGeneratorService with default
// implementations for its dependent services (transformer, separator, padding, and word list services).
// It initializes each service with the provided configuration and random number generator service.
func NewPasswordGeneratorService(
	cfg *config.Settings,
) (*DefaultPasswordGeneratorService, error) {
	rngs := NewRNGService()
	wls, err := NewWordListService(cfg, rngs)
	if err != nil {
		return nil, err
	}

	ts, err := NewTransformerService(cfg, rngs)
	if err != nil {
		return nil, err
	}

	ss, err := NewSeparatorService(cfg, rngs)
	if err != nil {
		return nil, err
	}

	ps, err := NewPaddingService(cfg, rngs)
	if err != nil {
		return nil, err
	}

	return NewCustomPasswordGeneratorService(cfg, ts, ss, ps, wls)
}

// Generate creates a list of passwords using the services provided to the
// DefaultPasswordGeneratorService instance and returns the list of generated
// passwords or the first error if one or more is encountered.
func (s *DefaultPasswordGeneratorService) Generate() ([]string, error) {
	pws := make([]string, s.cfg.NumPasswords)

	for i := 0; i < s.cfg.NumPasswords; i++ {
		// Get a list of words from the wordList service
		sl, err := s.wordListSvc.GetWords()
		if err != nil {
			return nil, err
		}

		// Transform the casing of words or letters using the transformer service
		slt, err := s.transformerSvc.Transform(sl)
		if err != nil {
			return nil, err
		}

		// Separate the transformed list using the separator service using special characters
		sls, err := s.separatorSvc.Separate(slt)
		if err != nil {
			return nil, err
		}

		// Pad the password with digits and special characters using the padding service
		pw, err := s.paddingSvc.Pad(sls)
		if err != nil {
			return nil, err
		}

		pws[i] = pw
	}

	return pws, nil
}
