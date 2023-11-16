package asset

import (
	"embed"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/eljamo/mempass/internal/config"
)

//go:embed preset/* word_list/*
var Files embed.FS

func keyToTXTFile(key string) (string, bool) {
	fileMap := map[string]string{
		config.ALL:             "all.txt",
		config.DOCTOR_WHO:      "doctor_who.txt",
		config.EN:              "en.txt",
		config.EN_SMALL:        "en_small.txt",
		config.GAME_OF_THRONES: "game_of_thrones.txt",
		config.HARRY_POTTER:    "harry_potter.txt",
		config.MIDDLE_EARTH:    "middle_earth.txt",
		config.POKEMON:         "pokemon.txt",
		config.STAR_TREK:       "star_trek.txt",
		config.STAR_WARS:       "star_wars.txt",
	}

	file, ok := fileMap[strings.ToUpper(key)]
	return file, ok
}

func keyToJSONFile(key string) (string, bool) {
	fileMap := map[string]string{
		config.APPLEID:        "appleid.json",
		config.DEFAULT:        "default.json",
		config.NTLM:           "ntlm.json",
		config.SECURITYQ:      "securityq.json",
		config.WEB16:          "web16.json",
		config.WEB16_XKPASSWD: "web16_xkpasswd.json",
		config.WEB32:          "web32.json",
		config.WIFI:           "wifi.json",
		config.XKCD:           "xkcd.json",
		config.XKCD_XKPASSWD:  "xkcd_xkpasswd.json",
	}

	file, ok := fileMap[strings.ToUpper(key)]
	return file, ok
}

func readJSONFileData(filePath string, readerFunc func(string) ([]byte, error)) (string, error) {
	data, err := readerFunc(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %s", err)
	}

	var parsed any
	if err := json.Unmarshal(data, &parsed); err != nil {
		return "", fmt.Errorf("invalid JSON content: %s", err)
	}

	return string(data), nil
}

func loadEmbeddedJSONFile(filePath string) (string, error) {
	return readJSONFileData(filePath, Files.ReadFile)
}

func LoadJSONFile(filePath string) (string, error) {
	return readJSONFileData(filePath, os.ReadFile)
}

func loadTXTFile(filePath string) ([]byte, error) {
	data, err := Files.ReadFile(filePath)
	if err != nil {
		return make([]byte, 0), fmt.Errorf("failed to read embedded file: %s", err)
	}

	return data, nil
}

func GetWordList(key string) ([]byte, error) {
	fileName, ok := keyToTXTFile(key)
	if !ok {
		return make([]byte, 0), fmt.Errorf("invalid %s name: %s", config.WORD_LIST_KEY, key)
	}

	content, err := loadTXTFile(fmt.Sprintf("%s/%s", config.WORD_LIST_KEY, fileName))
	if err != nil {
		return make([]byte, 0), err
	}

	return content, nil
}

func GetJSONPreset(key string) (string, error) {
	fileName, ok := keyToJSONFile(key)
	if !ok {
		return "", fmt.Errorf("invalid %s name: %s", config.PRESET_KEY, key)
	}

	content, err := loadEmbeddedJSONFile(fmt.Sprintf("%s/%s", config.PRESET_KEY, fileName))
	if err != nil {
		return "", err
	}

	return content, nil
}
