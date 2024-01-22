package config

import (
	"encoding/json"

	"github.com/eljamo/libpass/v7/internal/merger"
)

func mapToJSON(m map[string]any) ([]byte, error) {
	mj, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}

	return mj, nil
}

func mergeMaps(ms ...map[string]any) ([]byte, error) {
	mm := merger.Map(ms...)

	return mapToJSON(mm)
}

func jsonToSettings(s *Settings, js []byte) error {
	if err := json.Unmarshal(js, &s); err != nil {
		return err
	}

	return nil
}

// NewSettings creates a Settings struct from the given maps of unmarshalled JSON.
// If no maps are given, the default settings are returned.
func New(ms ...map[string]any) (*Settings, error) {
	settings := DefaultSettings()
	if len(ms) == 0 {
		return settings, nil
	}

	js, err := mergeMaps(ms...)
	if err != nil {
		return nil, err
	}

	err = jsonToSettings(settings, js)
	if err != nil {
		return nil, err
	}

	return settings, nil
}
