package service

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/eljamo/mempass/internal/config"
)

func TestDefaultTransformerService_Transform(t *testing.T) {
	rngs := NewMockRNGService()
	erngs := NewMockEvenRNGService()
	tests := []struct {
		name      string
		cfg       *config.Config
		rng       RNGService
		input     []string
		expected  []string
		expectErr bool
	}{
		{
			name:     "Alternate",
			cfg:      &config.Config{CaseTransform: config.ALTERNATE},
			rng:      rngs,
			input:    []string{"hello", "world"},
			expected: []string{"hello", "WORLD"},
		},
		{
			name:     "Capitalisation",
			cfg:      &config.Config{CaseTransform: config.CAPITALISE},
			rng:      rngs,
			input:    []string{"hello", "world"},
			expected: []string{"Hello", "World"},
		},
		{
			name:     "Inversion",
			cfg:      &config.Config{CaseTransform: config.INVERT},
			rng:      rngs,
			input:    []string{"hello", "world"},
			expected: []string{"hELLO", "wORLD"},
		},
		{
			name:     "Lower",
			cfg:      &config.Config{CaseTransform: config.LOWER},
			rng:      rngs,
			input:    []string{"HELLO", "WORLD"},
			expected: []string{"hello", "world"},
		},
		{
			name:     "Upper",
			cfg:      &config.Config{CaseTransform: config.UPPER},
			rng:      rngs,
			input:    []string{"hello", "world"},
			expected: []string{"HELLO", "WORLD"},
		},
		{
			name:     "Random",
			cfg:      &config.Config{CaseTransform: config.RANDOM},
			rng:      rngs,
			input:    []string{"hello", "world"},
			expected: []string{"hello", "world"},
		},
		{
			name:     "Random with even RNG",
			cfg:      &config.Config{CaseTransform: config.RANDOM},
			rng:      erngs,
			input:    []string{"hello", "world"},
			expected: []string{"HELLO", "WORLD"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewTransformerService(tt.cfg, tt.rng)
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
	validCaseTransforms := []string{
		config.ALTERNATE,
		config.CAPITALISE,
		config.INVERT,
		config.LOWER,
		config.RANDOM,
		config.UPPER,
		config.NONE,
	}

	for _, validTransform := range validCaseTransforms {
		t.Run(fmt.Sprintf("Valid case transform: %s", validTransform), func(t *testing.T) {
			cfg := DefaultTransformerService{
				cfg: &config.Config{CaseTransform: validTransform},
			}
			if err := cfg.validate(); err != nil {
				t.Errorf("validate() with valid case transform %s returned error: %v", validTransform, err)
			}
		})
	}

	t.Run("Invalid case transform", func(t *testing.T) {
		cfg := DefaultTransformerService{
			cfg: &config.Config{CaseTransform: "invalid_case_transform"},
		}
		if err := cfg.validate(); err == nil {
			t.Error("validate() with invalid case transform did not return an error")
		}
	})
}
