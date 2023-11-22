package json_merge

import (
	"encoding/json"
	"fmt"
)

func Merge(json1, json2 string) (string, error) {
	var m1, m2 map[string]any

	if err := json.Unmarshal([]byte(json1), &m1); err != nil {
		return "", fmt.Errorf("failed to unmarshal json1: %w", err)
	}
	if err := json.Unmarshal([]byte(json2), &m2); err != nil {
		return "", fmt.Errorf("failed to unmarshal json2: %w", err)
	}

	merged := mergeMaps(m1, m2)
	mergedJSON, err := json.Marshal(merged)
	if err != nil {
		return "", fmt.Errorf("failed to marshal merged JSON: %w", err)
	}

	return string(mergedJSON), nil
}

func mergeMaps(m1, m2 map[string]any) map[string]any {
	for k, v := range m2 {
		m1[k] = v
	}
	return m1
}
