package db

import "strconv"

func guessType(s string) string {
	// First, try to parse it as an integer
	if _, err := strconv.ParseInt(s, 10, 64); err == nil {
		return "INTEGER"
	}

	// If it's not an integer, try to parse it as a float
	if _, err := strconv.ParseFloat(s, 64); err == nil {
		return "REAL"
	}

	// If it's neither, then return "Text"
	return "TEXT"
}

// same as above, just takes generics and uses reflection
func reflectType(v any) string {
	switch v.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return "INTEGER"
	case float32, float64:
		return "REAL"
	default:
		return "TEXT"
	}
}
