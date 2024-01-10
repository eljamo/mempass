package config

import (
	"encoding/json"
	"fmt"

	"github.com/eljamo/libpass/v5/internal/merger"
)

func mapToJSONString(m map[string]any) ([]byte, error) {
	mj, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}

	return mj, nil
}

func mergeMaps(ms ...map[string]any) ([]byte, error) {
	mm := merger.Map(ms...)

	return mapToJSONString(mm)
}

func jsonToSettings(js []byte) (*Settings, error) {
	var cfg Settings
	if err := json.Unmarshal(js, &cfg); err != nil {
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
