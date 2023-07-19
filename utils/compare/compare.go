package compare

import "regexp"

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

	// if a is a string, use a as a regex and check if b is a match
	aString, aIsString := a.(string)
	if aIsString {
		match, _ := regexp.MatchString(makeRegex(aString), b.(string))
		if match {
			return true
		}
	}

	// if b is a string, use b as a regex and check if a is a match
	bString, bIsString := b.(string)
	if bIsString {
		match, _ := regexp.MatchString(makeRegex(bString), a.(string))
		if match {
			return true
		}
	}

	// else, just use equality comparison
	return a == b
}

func makeRegex(str string) string {
	// if string doesn't start with ^, and doesn't end with $, add them
	if str[0] != '^' {
		str = "^" + str
	}
	if str[len(str)-1] != '$' {
		str = str + "$"
	}

	return str
}
