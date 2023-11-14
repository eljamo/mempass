package service

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/eljamo/mempass/asset"
	"github.com/eljamo/mempass/internal/config"
)

type WordListService interface {
	GetWords(num int) ([]string, error)
}

type DefaultWordListService struct {
	cfg      *config.Config
	rng      RNGService
	wordList [][]byte
}

func NewWordListService(cfg *config.Config, rng RNGService) (*DefaultWordListService, error) {
	wl, err := getWordList(cfg.WordList, cfg.WordLengthMin, cfg.WordLengthMax)
	if err != nil {
		return &DefaultWordListService{}, err
	}

	return &DefaultWordListService{
		cfg,
		rng,
		wl,
	}, nil
}

func getWordList(wordList string, wordMinLength int, wordMaxLength int) ([][]byte, error) {
	if wordMaxLength < wordMinLength {
		return nil, fmt.Errorf("word_length_max (%d) must be greater than or equal to word_length_min (%d)", wordMaxLength, wordMinLength)
	}

	wl, err := asset.GetWordList(wordList)
	if err != nil {
		return nil, err
	}

	aw, err := bytes.Split(wl, []byte("\n")), nil
	if err != nil {
		return nil, err
	}

	var fw [][]byte

	for _, word := range aw {
		if len(word) >= wordMinLength && len(word) <= wordMaxLength {
			fw = append(fw, word)
		}
	}

	if len(fw) == 0 {
		return nil, fmt.Errorf("no words found in word list %s with minimum length of %d and maximum length of %d", wordList, wordMinLength, wordMaxLength)
	}

	return fw, nil
}

func (s *DefaultWordListService) getWord() (string, error) {
	num, err := s.rng.GenerateN(len(s.wordList))
	if err != nil {
		return "", err
	}
	w := string(s.wordList[int64(num)])

	return w, nil
}

func (s *DefaultWordListService) GetWords(num int) ([]string, error) {
	if num < 2 {
		return nil, errors.New("num_words must be greater than or equal to 2")
	}

	wl := make([]string, 0, num)

	for {
		if len(wl) == num {
			break
		}

		word, err := s.getWord()
		if err != nil {
			return wl, err
		}

		wl = append(wl, word)
	}

	return wl, nil
}
