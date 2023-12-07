package asset

import (
	"bufio"
	"embed"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/eljamo/libpass/v4/config"
)

//go:embed preset/* word_list/*
var files embed.FS
var fileMap = map[string]map[string]string{
	config.PresetKey: {
		config.AppleID:       "appleid.json",
		config.Default:       "default.json",
		config.NTLM:          "ntlm.json",
		config.SecurityQ:     "securityq.json",
		config.Web16:         "web16.json",
		config.Web16XKPasswd: "web16_xkpasswd.json",
		config.Web32:         "web32.json",
		config.WiFi:          "wifi.json",
		config.XKCD:          "xkcd.json",
		config.XKCDXKPasswd:  "xkcd_xkpasswd.json",
	},
	config.WordListKey: {
		config.All:           "all.txt",
		config.DoctorWho:     "doctor_who.txt",
		config.EN:            "en.txt",
		config.ENSmall:       "en_small.txt",
		config.GameOfThrones: "game_of_thrones.txt",
		config.HarryPotter:   "harry_potter.txt",
		config.MiddleEarth:   "middle_earth.txt",
		config.Pokemon:       "pokemon.txt",
		config.StarTrek:      "star_trek.txt",
		config.StarWars:      "star_wars.txt",
	},
}

func keyToFile(key, fileType string) (string, bool) {
	file, ok := fileMap[fileType][strings.ToUpper(key)]

	return file, ok
}

func loadJSONFileData(filePath string, readerFunc func(string) ([]byte, error)) (string, error) {
	data, err := readerFunc(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read file (%s): %w", filePath, err)
	}

	var parsed any
	if err := json.Unmarshal(data, &parsed); err != nil {
		return "", fmt.Errorf("invalid JSON content in (%s): %w", filePath, err)
	}

	jsonData, err := json.Marshal(parsed)
	if err != nil {
		return "", fmt.Errorf("marshaling JSON data for %s failed: %w", filePath, err)
	}

	return string(jsonData), nil
}

// LoadJSONFile reads a JSON file from the given file path and returns its
// content as a string. It handles file reading and JSON unmarshalling. In case
// of any error during these operations, an error is returned
func LoadJSONFile(filePath string) (string, error) {
	return loadJSONFileData(filePath, os.ReadFile)
}

func getWordListFilePath(key string) (string, error) {
	fileName, ok := keyToFile(key, config.WordListKey)
	if !ok {
		return "", fmt.Errorf("invalid %s value (%s)", config.WordListKey, key)
	}

	return path.Join(config.WordListKey, fileName), nil
}

// GetWordList retrieves a list of words from an embedded file identified by the
// given key. The method reads the file content, splits it by newline characters
// and returns the result as a slice of strings. If the file cannot be found or
// read, an error is returned.
func GetWordList(key string) ([]string, error) {
	path, err := getWordListFilePath(key)
	if err != nil {
		return nil, err
	}

	data, err := files.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read embedded text file (%s): %w", path, err)
	}

	return strings.Split(string(data), "\n"), nil
}

// GetFilteredWordList reads a word list from an embedded file identified by the
// given key, and filters the words based on the specified minimum and maximum
// length. It returns a slice of strings that meet the length criteria. If the
// file cannot be opened or read, or if an error occurs during scanning, an
// error is returned.
func GetFilteredWordList(key string, minLen int, maxLen int) ([]string, error) {
	path, err := getWordListFilePath(key)
	if err != nil {
		return nil, err
	}

	file, err := files.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var wl []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) >= minLen && len(line) <= maxLen {
			wl = append(wl, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return wl, nil
}

// GetJSONPreset reads a JSON preset file identified by the given key from
// embedded files. It returns the content of the JSON file as a string. If the
// file is not found, cannot be read, or contains invalid JSON, an error is
// returned.
func GetJSONPreset(key string) (string, error) {
	fileName, ok := keyToFile(key, config.PresetKey)
	if !ok {
		return "", fmt.Errorf("invalid %s value (%s)", config.PresetKey, key)
	}

	filePath := path.Join(config.PresetKey, fileName)

	return loadJSONFileData(filePath, files.ReadFile)
}
