package compare

// DeepEqual compares two maps recursively
func DeepEqual(a, b map[string]interface{}) bool {
	// Check if the maps have the same number of keys
	if len(a) != len(b) {
		return false
	}

	// Check if all keys in 'a' exist in 'b' and have the same values
	for key, aValue := range a {
		bValue, ok := b[key]
		if !ok || !deepValueEqual(aValue, bValue) {
			return false
		}
	}

	return true
}

func deepValueEqual(a, b interface{}) bool {
	// If both values are maps, recursively compare them
	aMap, aIsMap := a.(map[string]interface{})
	bMap, bIsMap := b.(map[string]interface{})
	if aIsMap && bIsMap {
		return DeepEqual(aMap, bMap)
	}

	// Otherwise, use equality comparison
	return a == b
}
