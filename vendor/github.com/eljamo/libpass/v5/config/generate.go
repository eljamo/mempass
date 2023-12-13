package config

import (
	"encoding/json"
	"fmt"

	"github.com/eljamo/libpass/v5/internal/merger"
)

func mapToJSONString(m map[string]any) (string, error) {
	mj, err := json.Marshal(m)
	if err != nil {
		return "", err
	}

	return string(mj), nil
}

func mergeMaps(ms ...map[string]any) (string, error) {
	mm := merger.Map(ms...)

	return mapToJSONString(mm)
}

func jsonToSettings(js string) (*Settings, error) {
	var cfg Settings
	if err := json.Unmarshal([]byte(js), &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func Generate(ms ...map[string]any) (*Settings, error) {
	if len(ms) == 0 {
		return nil, fmt.Errorf("no config provided")
	}

	js, err := mergeMaps(ms...)
	if err != nil {
		return nil, err
	}

	return jsonToSettings(js)
}
