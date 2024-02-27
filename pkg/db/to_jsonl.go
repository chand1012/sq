package db

import (
	"bytes"
	"encoding/json"
	"fmt"

	"database/sql"
)

func ToJSONL(db *sql.DB, tableName string) (string, error) {
	// Query table data
	query := fmt.Sprintf("SELECT * FROM %s", tableName)
	rows, err := db.Query(query)
	if err != nil {
		return "", err
	}
	defer rows.Close()

	return RowsToJSONL(rows)
}

func RowsToJSONL(rows *sql.Rows) (string, error) {
	var buffer bytes.Buffer

	// Get column names
	cols, err := rows.Columns()
	if err != nil {
		return "", err
	}

	// Iterate through rows and build JSONL data
	for rows.Next() {
		// Scan row data
		scanValues := make([]interface{}, len(cols))
		scanPointers := make([]interface{}, len(cols))
		for i := range scanValues {
			scanPointers[i] = &scanValues[i]
		}
		err := rows.Scan(scanPointers...)
		if err != nil {
			return "", err
		}

		// Create a map to store row data, excluding null values
		rowMap := map[string]interface{}{}
		for i, col := range cols {
			if scanValues[i] != nil {
				rowMap[col] = scanValues[i]
			}
		}

		// Marshal row data to JSON
		jsonData, err := json.Marshal(rowMap)
		if err != nil {
			return "", err
		}

		// Append JSON data with newline character
		buffer.WriteString(string(jsonData))
		buffer.WriteByte('\n')
	}

	return buffer.String(), nil
}
