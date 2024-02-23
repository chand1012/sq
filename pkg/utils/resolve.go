package utils

import (
	"bytes"
	"encoding/json"
	"unicode/utf8"
)

func IsSQLiteFile(data []byte) bool {
	// SQLite database file header (magic string)
	const sqliteHeader = "SQLite format 3\x00"
	// Convert the first 16 bytes of data to a string for comparison.
	// Ensure data has at least 16 bytes to avoid slicing beyond its length.
	if len(data) < len(sqliteHeader) {
		return false
	}
	header := string(data[:len(sqliteHeader)])
	return header == sqliteHeader
}

func IsValidJSON(b []byte) bool {
	if !utf8.Valid(b) {
		return false
	}
	var js json.RawMessage
	return json.Unmarshal(b, &js) == nil
}

func IsValidJSONL(b []byte) bool {
	if !utf8.Valid(b) {
		return false
	}

	lines := bytes.Split(b, []byte("\n"))
	for _, line := range lines {
		// Skip empty lines
		if len(line) == 0 {
			continue
		}
		if !IsValidJSON(line) {
			return false
		}
	}
	return true // All lines are valid JSON
}
