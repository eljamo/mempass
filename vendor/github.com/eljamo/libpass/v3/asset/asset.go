package asset

import (
	"embed"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/eljamo/libpass/v3/config"
)

//go:embed preset/* word_list/*
var Files embed.FS

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

func loadFileData(filePath string, readerFunc func(string) ([]byte, error)) (string, error) {
	data, err := readerFunc(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read file '%s': %w", filePath, err)
	}

	var parsed any
	if err := json.Unmarshal(data, &parsed); err != nil {
		return "", fmt.Errorf("invalid JSON content in '%s': %w", filePath, err)
	}

	jsonData, err := json.Marshal(parsed)
	if err != nil {
		return "", fmt.Errorf("error marshaling JSON data from '%s': %w", filePath, err)
	}

	return string(jsonData), nil
}

func LoadJSONFile(filePath string) (string, error) {
	return loadFileData(filePath, os.ReadFile)
}

func GetWordList(key string) ([]string, error) {
	fileName, ok := keyToFile(key, config.WordListKey)
	if !ok {
		return nil, fmt.Errorf("invalid word list key '%s'", key)
	}

	filePath := fmt.Sprintf("%s/%s", config.WordListKey, fileName)
	data, err := Files.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read embedded text file '%s': %w", filePath, err)
	}

	return strings.Split(string(data), "\n"), nil
}

func GetJSONPreset(key string) (string, error) {
	fileName, ok := keyToFile(key, config.PresetKey)
	if !ok {
		return "", fmt.Errorf("invalid JSON preset key '%s'", key)
	}

	filePath := fmt.Sprintf("%s/%s", config.PresetKey, fileName)
	return loadFileData(filePath, Files.ReadFile)
}
