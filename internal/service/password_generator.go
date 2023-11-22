package service

import (
	"errors"
	"fmt"
	"sync"

	"github.com/eljamo/mempass/internal/config"
)

type PasswordGeneratorService interface {
	Generate() ([]string, error)
}

type DefaultPasswordGeneratorService struct {
	cfg            *config.Config
	transformerSvc TransformerService
	separatorSvc   SeparatorService
	paddingSvc     PaddingService
	wordListSvc    WordListService
}

func NewPasswordGeneratorService(
	cfg *config.Config,
	transformerSvc TransformerService,
	separatorSvc SeparatorService,
	paddingSvc PaddingService,
	wordListSvc WordListService,
) (*DefaultPasswordGeneratorService, error) {
	np := cfg.NumPasswords
	if np < 1 {
		return nil, errors.New("num_passwords must be greater than 0")
	}

	return &DefaultPasswordGeneratorService{
		cfg,
		transformerSvc,
		separatorSvc,
		paddingSvc,
		wordListSvc,
	}, nil
}

func (s *DefaultPasswordGeneratorService) Generate() ([]string, error) {
	np := s.cfg.NumPasswords

	var wg sync.WaitGroup
	wg.Add(np)

	pws := make([]string, np)
	errChan := make(chan error, np)

	for i := 0; i < np; i++ {
		go func(i int) {
			defer wg.Done()

			sl, err := s.wordListSvc.GetWords()
			if err != nil {
				errChan <- err
				return
			}

			slt, err := s.transformerSvc.Transform(sl)
			if err != nil {
				errChan <- err
				return
			}

			sls, err := s.separatorSvc.Separate(slt)
			if err != nil {
				errChan <- err
				return
			}

			pw, err := s.paddingSvc.Pad(sls)
			if err != nil {
				errChan <- err
				return
			}

			pws[i] = pw
		}(i)
	}

	wg.Wait()
	close(errChan)

	if len(errChan) > 0 {
		return nil, fmt.Errorf("%w", <-errChan)
	}

	return pws, nil
}
