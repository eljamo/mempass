package json_merge

import (
	"encoding/json"
	"fmt"
)

func Merge(json1, json2 string) (string, error) {
	m1, err := unmarshalJSON(json1)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal json1: %w", err)
	}

	m2, err := unmarshalJSON(json2)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal json2: %w", err)
	}

	mm := mergeMaps(m1, m2)
	mj, err := json.Marshal(mm)
	if err != nil {
		return "", fmt.Errorf("failed to marshal merged JSON: %w", err)
	}

	return string(mj), nil
}

func unmarshalJSON(jsonStr string) (map[string]any, error) {
	var m map[string]any
	if err := json.Unmarshal([]byte(jsonStr), &m); err != nil {
		return nil, err
	}
	return m, nil
}

func mergeMaps(m1, m2 map[string]any) map[string]any {
	for k, v := range m2 {
		m1[k] = v
	}
	return m1
}
