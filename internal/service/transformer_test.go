package service

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/eljamo/mempass/internal/config"
)

func TestNewTransformerService(t *testing.T) {
	t.Parallel()

	mockRNGService := &MockEvenRNGService{}

	validTransformType := config.Upper
	invalidTransformType := "invalid"

	tests := []struct {
		name          string
		caseTransform string
		wantErr       bool
	}{
		{
			name:          "Valid configuration",
			caseTransform: validTransformType,
			wantErr:       false,
		},
		{
			name:          "Invalid configuration",
			caseTransform: invalidTransformType,
			wantErr:       true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			cfg := &config.Config{CaseTransform: tt.caseTransform}
			_, err := NewTransformerService(cfg, mockRNGService)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewTransformerService() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDefaultTransformerService_Transform(t *testing.T) {
	t.Parallel()

	rngs := &MockRNGService{}
	erngs := &MockEvenRNGService{}

	tests := []struct {
		name      string
		cfg       *config.Config
		rngSvc    RNGService
		input     []string
		expected  []string
		expectErr bool
	}{
		{
			name:     "Alternate",
			cfg:      &config.Config{CaseTransform: config.Alternate},
			rngSvc:   rngs,
			input:    []string{"hello", "world"},
			expected: []string{"hello", "WORLD"},
		},
		{
			name:     "Alternate Lettercase",
			cfg:      &config.Config{CaseTransform: config.AlternateLettercase},
			rngSvc:   rngs,
			input:    []string{"hello", "world"},
			expected: []string{"hElLo", "wOrLd"},
		},
		{
			name:     "Capitalisation",
			cfg:      &config.Config{CaseTransform: config.Capitalise},
			rngSvc:   rngs,
			input:    []string{"hello", "world"},
			expected: []string{"Hello", "World"},
		},
		{
			name:     "Capitalisation Inversed",
			cfg:      &config.Config{CaseTransform: config.CapitaliseInvert},
			rngSvc:   rngs,
			input:    []string{"hello", "world"},
			expected: []string{"hELLO", "wORLD"},
		},
		{
			name:     "Inversion",
			cfg:      &config.Config{CaseTransform: config.Invert},
			rngSvc:   rngs,
			input:    []string{"hello", "world"},
			expected: []string{"hELLO", "wORLD"},
		},
		{
			name:     "Lower",
			cfg:      &config.Config{CaseTransform: config.Lower},
			rngSvc:   rngs,
			input:    []string{"HELLO", "WORLD"},
			expected: []string{"hello", "world"},
		},
		{
			name:     "Lower Vowel Upper Consonant",
			cfg:      &config.Config{CaseTransform: config.LowerVowelUpperConsonant},
			rngSvc:   rngs,
			input:    []string{"hello", "world"},
			expected: []string{"HeLLo", "WoRLD"},
		},
		{
			name:     "Sentence",
			cfg:      &config.Config{CaseTransform: config.Sentence},
			rngSvc:   rngs,
			input:    []string{"hello", "world"},
			expected: []string{"Hello", "world"},
		},
		{
			name:     "Upper",
			cfg:      &config.Config{CaseTransform: config.Upper},
			rngSvc:   rngs,
			input:    []string{"hello", "world"},
			expected: []string{"HELLO", "WORLD"},
		},
		{
			name:     "Random",
			cfg:      &config.Config{CaseTransform: config.Random},
			rngSvc:   rngs,
			input:    []string{"hello", "world"},
			expected: []string{"hello", "world"},
		},
		{
			name:     "Random with even RNG",
			cfg:      &config.Config{CaseTransform: config.Random},
			rngSvc:   erngs,
			input:    []string{"hello", "world"},
			expected: []string{"HELLO", "WORLD"},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			service, err := NewTransformerService(tt.cfg, tt.rngSvc)
			if err != nil {
				t.Errorf("unexpected error with service init: %s", err)
			}

			got, err := service.Transform(tt.input)

			if (err != nil) != tt.expectErr {
				t.Errorf("DefaultTransformerService.Transform() error = %v, expectErr %v", err, tt.expectErr)
				return
			}
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("DefaultTransformerService.Transform() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestDefaultTransformerService_Validate(t *testing.T) {
	t.Parallel()

	validCaseTransforms := []string{
		config.Alternate,
		config.AlternateLettercase,
		config.Capitalise,
		config.CapitaliseInvert,
		config.Invert,
		config.Lower,
		config.LowerVowelUpperConsonant,
		config.Random,
		config.Upper,
		config.None,
	}

	for _, validTransform := range validCaseTransforms {
		validTransform := validTransform
		t.Run(fmt.Sprintf("Valid case transform: %s", validTransform), func(t *testing.T) {
			t.Parallel()

			cfg := DefaultTransformerService{
				cfg: &config.Config{CaseTransform: validTransform},
			}
			if err := cfg.validate(); err != nil {
				t.Errorf("validate() with valid case transform %s returned error: %v", validTransform, err)
			}
		})
	}

	t.Run("Invalid case transform", func(t *testing.T) {
		t.Parallel()

		cfg := DefaultTransformerService{
			cfg: &config.Config{CaseTransform: "invalid_case_transform"},
		}
		if err := cfg.validate(); err == nil {
			t.Error("validate() with invalid case transform did not return an error")
		}
	})
}
