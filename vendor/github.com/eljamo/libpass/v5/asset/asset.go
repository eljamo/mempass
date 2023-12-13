package asset

import (
	"bufio"
	"embed"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/eljamo/libpass/v5/config/option"
)

//go:embed preset/* word_list/*
var files embed.FS
var fileMap = map[string]map[string]string{
	option.PresetKey: {
		option.AppleID:       "appleid.json",
		option.Default:       "default.json",
		option.NTLM:          "ntlm.json",
		option.SecurityQ:     "securityq.json",
		option.Web16:         "web16.json",
		option.Web16XKPasswd: "web16_xkpasswd.json",
		option.Web32:         "web32.json",
		option.WiFi:          "wifi.json",
		option.XKCD:          "xkcd.json",
		option.XKCDXKPasswd:  "xkcd_xkpasswd.json",
	},
	option.WordListKey: {
		option.All:           "all.txt",
		option.DoctorWho:     "doctor_who.txt",
		option.EN:            "en.txt",
		option.ENSmall:       "en_small.txt",
		option.GameOfThrones: "game_of_thrones.txt",
		option.HarryPotter:   "harry_potter.txt",
		option.MiddleEarth:   "middle_earth.txt",
		option.Pokemon:       "pokemon.txt",
		option.StarTrek:      "star_trek.txt",
		option.StarWars:      "star_wars.txt",
	},
}

func keyToFile(key, fileType string) (string, bool) {
	file, ok := fileMap[fileType][strings.ToUpper(key)]

	return file, ok
}

func loadJSONFileData(filePath string, readerFunc func(string) ([]byte, error)) (map[string]any, error) {
	data, err := readerFunc(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file (%s): %w", filePath, err)
	}

	var jmap map[string]any
	if err := json.Unmarshal(data, &jmap); err != nil {
		return nil, fmt.Errorf("invalid JSON content in (%s): %w", filePath, err)
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
	fileName, ok := keyToFile(key, option.WordListKey)
	if !ok {
		return "", fmt.Errorf("invalid %s value (%s)", option.WordListKey, key)
	}

	return path.Join(option.WordListKey, fileName), nil
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
// embedded files. It returns the content of the JSON file as a map, if not an
// error is returned.
func GetJSONPreset(key string) (map[string]any, error) {
	fileName, ok := keyToFile(key, option.PresetKey)
	if !ok {
		return nil, fmt.Errorf("invalid %s value (%s)", option.PresetKey, key)
	}

	filePath := path.Join(option.PresetKey, fileName)

	return loadJSONFileData(filePath, files.ReadFile)
}
