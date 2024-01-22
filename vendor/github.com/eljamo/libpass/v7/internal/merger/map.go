package merger

// Map takes multiple map[string]any arguments and merges them into a single
// map. If the same key appears in multiple input maps, the value from the last
// map containing the key will be used in the resulting merged map.
func Map[K comparable, V any](maps ...map[K]V) map[K]V {
	mergedMap := make(map[K]V)

	for _, currentMap := range maps {
		for key, value := range currentMap {
			mergedMap[key] = value
		}
	}

	return mergedMap
}
