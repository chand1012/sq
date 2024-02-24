package file_types

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
