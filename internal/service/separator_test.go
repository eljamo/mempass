package service

import (
	"slices"
	"testing"

	"github.com/eljamo/mempass/internal/config"
)

func TestNewSeparatorService(t *testing.T) {
	t.Parallel()

	mockRNGService := &MockRNGService{}

	tests := []struct {
		name    string
		cfg     *config.Config
		wantErr bool
	}{
		{
			name:    "Valid configuration",
			cfg:     &config.Config{SeparatorCharacter: "*", SeparatorAlphabet: []string{"!", "@", "#", "$", "%"}},
			wantErr: false,
		},
		{
			name:    "Invalid configuration - invalid separator character",
			cfg:     &config.Config{SeparatorCharacter: "invalid"},
			wantErr: true,
		},
		{
			name:    "Valid configuration - separator alphabet",
			cfg:     &config.Config{SeparatorCharacter: config.Random, SeparatorAlphabet: []string{""}},
			wantErr: false,
		},

		{
			name:    "Valid configuration - separator alphabet",
			cfg:     &config.Config{SeparatorCharacter: config.Random, SeparatorAlphabet: []string{"a"}},
			wantErr: false,
		},
		{
			name:    "Invalid configuration - empty separator alphabet",
			cfg:     &config.Config{SeparatorCharacter: config.Random, SeparatorAlphabet: []string{"aaa"}},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			_, err := NewSeparatorService(tt.cfg, mockRNGService)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewSeparatorService() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSeparatorService_Separate(t *testing.T) {
	t.Parallel()

	rngs := &MockRNGService{}
	erngs := &MockEvenRNGService{}

	tests := []struct {
		name      string
		cfg       *config.Config
		rngSvc    RNGService
		input     []string
		expected  []string
		expectErr error
	}{
		{
			name:     "With fixed separator",
			cfg:      &config.Config{SeparatorCharacter: "-"},
			rngSvc:   rngs,
			input:    []string{"a", "b", "c"},
			expected: []string{"-", "a", "-", "b", "-", "c", "-"},
		},
		{
			name:     "With empty slice",
			cfg:      &config.Config{SeparatorCharacter: "-"},
			rngSvc:   rngs,
			input:    []string{},
			expected: []string{"-"},
		},
		{
			name:     "With random separator",
			cfg:      &config.Config{SeparatorCharacter: config.Random, SeparatorAlphabet: []string{"!", "-", "="}},
			rngSvc:   rngs,
			input:    []string{"a", "b", "c"},
			expected: []string{"-", "a", "-", "b", "-", "c", "-"},
		},
		{
			name:     "With random separator with RNG returning a even number",
			cfg:      &config.Config{SeparatorCharacter: config.Random, SeparatorAlphabet: []string{"!", "-", "="}},
			rngSvc:   erngs,
			input:    []string{"a", "b", "c"},
			expected: []string{"=", "a", "=", "b", "=", "c", "="},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			service, err := NewSeparatorService(tt.cfg, tt.rngSvc)
			if err != nil {
				t.Errorf("unexpected error with service init: %s", err)
				return
			}

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
