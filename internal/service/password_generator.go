package service

import (
	"fmt"

	"github.com/eljamo/mempass/internal/config"
)

type PasswordGeneratorService interface {
	Generate() ([]string, error)
}

type DefaultPasswordGeneratorService struct {
	cfg *config.Config
}

func NewPasswordGeneratorService(cfg *config.Config) *DefaultPasswordGeneratorService {
	return &DefaultPasswordGeneratorService{cfg}
}

func (s *DefaultPasswordGeneratorService) Generate() ([]string, error) {
	np := s.cfg.NumPasswords
	if np < 1 {
		return nil, fmt.Errorf("num_passwords must be greater than 0")
	}

	rngs := NewRNGService()
	wls, err := NewWordListService(s.cfg, rngs)
	if err != nil {
		return nil, err
	}

	ts := NewTransformerService(s.cfg, rngs)
	ss := NewSeparatorService(s.cfg, rngs)
	ps := NewPaddingService(s.cfg, rngs)
	var pws []string

	for i := 0; i < np; i++ {
		sl, err := wls.GetWords(s.cfg.NumWords)
		if err != nil {
			return nil, err
		}

		slt, err := ts.Transform(sl)
		if err != nil {
			return nil, err
		}

		sls, err := ss.Separate(slt)
		if err != nil {
			return nil, err
		}

		pw, err := ps.Pad(sls)
		if err != nil {
			return nil, err
		}

		pws = append(pws, pw)
	}

	return pws, nil
}
