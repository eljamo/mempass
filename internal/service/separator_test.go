package service

import (
	"errors"
	"slices"
	"testing"

	"github.com/eljamo/mempass/internal/config"
)

func TestSeparatorService_Separate(t *testing.T) {
	rngs := NewMockRNGService()
	erngs := NewMockEvenRNGService()
	tests := []struct {
		name      string
		cfg       *config.Config
		rng       RNGService
		input     []string
		expected  []string
		expectErr error
	}{
		{
			name:     "With fixed separator",
			cfg:      &config.Config{SeparatorCharacter: "-"},
			rng:      rngs,
			input:    []string{"a", "b", "c"},
			expected: []string{"-", "a", "-", "b", "-", "c", "-"},
		},
		{
			name:     "With empty slice",
			cfg:      &config.Config{SeparatorCharacter: "-"},
			rng:      rngs,
			input:    []string{},
			expected: []string{"-"},
		},
		{
			name:     "With random separator",
			cfg:      &config.Config{SeparatorCharacter: config.RANDOM, SeparatorAlphabet: []string{"!", "-", "="}},
			rng:      rngs,
			input:    []string{"a", "b", "c"},
			expected: []string{"-", "a", "-", "b", "-", "c", "-"},
		},
		{
			name:     "With random separator with RNG returning a even number",
			cfg:      &config.Config{SeparatorCharacter: config.RANDOM, SeparatorAlphabet: []string{"!", "-", "="}},
			rng:      erngs,
			input:    []string{"a", "b", "c"},
			expected: []string{"=", "a", "=", "b", "=", "c", "="},
		},
		{
			name:      "With random separator but empty alphabet",
			cfg:       &config.Config{SeparatorCharacter: config.RANDOM, SeparatorAlphabet: []string{}},
			rng:       rngs,
			input:     []string{"a", "b", "c"},
			expectErr: errors.New("separator_alphabet cannot be empty"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewSeparatorService(tt.cfg, tt.rng)
			got, err := service.Separate(tt.input)

			if tt.expectErr != nil {
				if err == nil {
					t.Errorf("expected error but got none")
				} else if err.Error() != tt.expectErr.Error() {
					t.Errorf("expected error %q but got %q", tt.expectErr, err)
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %s", err)
				return
			}

			if !slices.Equal(got, tt.expected) {
				t.Errorf("expected %v, but got %v", tt.expected, got)
			}
		})
	}
}
