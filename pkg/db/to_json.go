package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
)

func ToJSON(db *sql.DB, tableName string) (string, error) {
	// Create an empty slice to store row data as maps
	var data []map[string]interface{}

	// Query table data
	query := fmt.Sprintf("SELECT * FROM %s", tableName)
	rows, err := db.Query(query)
	if err != nil {
		return "", err
	}
	defer rows.Close()

	// Get column information
	cols, err := rows.Columns()
	if err != nil {
		return "", err
	}

	// Iterate through rows and build data slice
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

		// Create a map to store row data
		rowMap := make(map[string]interface{})
		for i, col := range cols {
			rowMap[col] = scanValues[i]
		}

		// Append row data to the slice
		data = append(data, rowMap)
	}

	// Marshal data slice to JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	return string(jsonData), nil
}

func RowsToJSON(rows *sql.Rows) (string, error) {
	var data []map[string]interface{}

	// Get column names
	cols, err := rows.Columns()
	if err != nil {
		return "", err
	}

	// Iterate through rows and build data slice
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

		// Append row data to the slice
		data = append(data, rowMap)
	}

	// Marshal data slice to JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	return string(jsonData), nil
}
