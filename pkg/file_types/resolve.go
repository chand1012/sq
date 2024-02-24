package file_types

import (
	"path/filepath"
	"strings"
)

// enum of the valid types
// json, csv, and jsonl for now
type FileType int

const (
	JSON FileType = iota
	CSV
	JSONL
	SQLite
	Unknown
)

func Resolve(b []byte) FileType {
	if IsSQLiteFile(b) {
		return SQLite
	}
	if IsValidCSV(b) {
		return CSV
	}
	if IsValidJSON(b) {
		return JSON
	}
	if IsValidJSONL(b) {
		return JSONL
	}
	return Unknown
}

func ResolveByPath(fileName string) FileType {
	// Get the file extension (lowercase)
	ext := strings.ToLower(filepath.Ext(fileName))

	switch ext {
	case ".csv":
		return CSV
	case ".sqlite", ".db", ".sqlite3", ".db3", ".sdb", ".dat":
		return SQLite
	case ".json":
		return JSON
	case ".jsonl":
		return JSONL
	default:
		return Unknown
	}
}

func (ft FileType) String() string {
	switch ft {
	case JSON:
		return "json"
	case CSV:
		return "csv"
	case JSONL:
		return "jsonl"
	case SQLite:
		return "sqlite"
	default:
		return "unknown"
	}
}
