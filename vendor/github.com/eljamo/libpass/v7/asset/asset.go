package asset

import (
	"bufio"
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"strings"

	"github.com/eljamo/libpass/v7/config/option"
)

//go:embed preset/* word_list/*
var files embed.FS

var fileMap = map[string]map[string]string{
	option.ConfigKeyPreset: {
		option.PresetAppleID:       "appleid.json",
		option.PresetDefault:       "default.json",
		option.PresetNTLM:          "ntlm.json",
		option.PresetSecurityQ:     "securityq.json",
		option.PresetWeb16:         "web16.json",
		option.PresetWeb16XKPasswd: "web16_xkpasswd.json",
		option.PresetWeb32:         "web32.json",
		option.PresetWiFi:          "wifi.json",
		option.PresetXKCD:          "xkcd.json",
		option.PresetXKCDXKPasswd:  "xkcd_xkpasswd.json",
	},
	option.ConfigKeyWordList: {
		option.WordList40k:           "40k.txt",
		option.WordListAll:           "all.txt",
		option.WordListDoctorWho:     "doctor_who.txt",
		option.WordListEN:            "en.txt",
		option.WordListENSmall:       "en_small.txt",
		option.WordListGameOfThrones: "game_of_thrones.txt",
		option.WordListHarryPotter:   "harry_potter.txt",
		option.WordListMiddleEarth:   "middle_earth.txt",
		option.WordListPokemon:       "pokemon.txt",
		option.WordListStarTrek:      "star_trek.txt",
		option.WordListStarWars:      "star_wars.txt",
	},
}

var (
	ErrReadFile        = errors.New("failed to read file")
	ErrJSON            = errors.New("invalid JSON content")
	ErrInvalidWordList = errors.New("invalid word list")
	ErrInvalidPreset   = errors.New("invalid preset")
)

func keyToFile(key, fileType string) (string, bool) {
	file, ok := fileMap[fileType][strings.ToUpper(key)]

	return file, ok
}

func loadJSONFileData(filePath string, readerFunc func(string) ([]byte, error)) (map[string]any, error) {
	data, err := readerFunc(filePath)
	if err != nil {
		return nil, errors.Join(ErrReadFile, fmt.Errorf("error reading file (%s): %w", filePath, err))
	}

	var jmap map[string]any
	if err := json.Unmarshal(data, &jmap); err != nil {
		return nil, errors.Join(ErrJSON, fmt.Errorf("error unmarshaling JSON (%s): %w", filePath, err))
	}

	return jmap, nil
}

// LoadJSONFile reads a JSON file from the given file path and returns its
// content as a map. In case of any error during these operations, an error
// is returned
func LoadJSONFile(filePath string) (map[string]any, error) {
	return loadJSONFileData(filePath, os.ReadFile)
}

func getWordListFilePath(key string) (string, error) {
	fileName, ok := keyToFile(key, option.ConfigKeyWordList)
	if !ok {
		return "", errors.Join(ErrInvalidWordList, fmt.Errorf("invalid %s value (%s)", option.ConfigKeyWordList, key))
	}

	return path.Join(option.ConfigKeyWordList, fileName), nil
}

// GetWordList retrieves a list of words from an embedded file identified by the
// given key. The method reads the file content, splits it by newline characters
// and returns the result as a slice of strings. If the file cannot be found or
// read, an error is returned.
func GetWordList(key string) ([]string, error) {
	filePath, err := getWordListFilePath(key)
	if err != nil {
		return nil, err
	}

	data, err := files.ReadFile(filePath)
	if err != nil {
		return nil, errors.Join(ErrReadFile, fmt.Errorf("failed to read embedded text file (%s): %w", filePath, err))
	}

	return strings.Split(string(data), "\n"), nil
}

// readAndFilterWords reads from an io.Reader, and filters the words based on the specified minimum and maximum length.
func readAndFilterWords(filePath string, minLen int, maxLen int, fs embed.FS) ([]string, error) {
	file, err := fs.Open(filePath)
	if err != nil {
		return nil, errors.Join(ErrReadFile, fmt.Errorf("failed to open embedded text file (%s): %w", filePath, err))
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Printf("failed to close file (%s): %v", filePath, err)
		}
	}()

	var wl []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) >= minLen && len(line) <= maxLen {
			wl = append(wl, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, errors.Join(ErrReadFile, fmt.Errorf("failed to scan embedded text file (%s): %w", filePath, err))
	}

	return wl, nil
}

// GetFilteredWordList reads a word list from an embedded file identified by the
// given key, and filters the words based on the specified minimum and maximum
// length. It returns a slice of strings that meet the length criteria. If the
// file cannot be opened or read, or if an error occurs during scanning, an
// error is returned.
func GetFilteredWordList(key string, minLen int, maxLen int) ([]string, error) {
	filePath, err := getWordListFilePath(key)
	if err != nil {
		return nil, err
	}

	return readAndFilterWords(filePath, minLen, maxLen, files)
}

func getPresetFilePath(key string) (string, error) {
	fileName, ok := keyToFile(key, option.ConfigKeyPreset)
	if !ok {
		return "", errors.Join(ErrInvalidPreset, fmt.Errorf("invalid %s value (%s)", option.ConfigKeyPreset, key))
	}

	return path.Join(option.ConfigKeyPreset, fileName), nil
}

// GetJSONPreset reads a JSON preset file identified by the given key from
// embedded files. It returns the content of the JSON file as a map, if not an
// error is returned.
func GetJSONPreset(key string) (map[string]any, error) {
	filePath, err := getPresetFilePath(key)
	if err != nil {
		return nil, err
	}

	return loadJSONFileData(filePath, files.ReadFile)
}
