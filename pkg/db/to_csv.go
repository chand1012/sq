package db

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"reflect"
	"strings"

	"database/sql"
)

func ToCSV(db *sql.DB, tableName string) (string, error) {
	// Create a buffer to hold the CSV data
	var buffer bytes.Buffer
	w := csv.NewWriter(&buffer)

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

	// Write table header
	err = w.Write(cols)
	if err != nil {
		return "", err
	}

	// Iterate through rows and write data
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

		// Convert scanned values to strings
		row := make([]string, len(cols))
		for i, v := range scanValues {
			val := reflect.ValueOf(v)
			if val.Kind() == reflect.String {
				row[i] = val.String()
			} else {
				row[i] = fmt.Sprint(v)
			}
		}

		// Write row data
		err = w.Write(row)
		if err != nil {
			return "", err
		}
	}

	w.Flush()
	return buffer.String(), nil
}

func RowsToCSV(rows *sql.Rows) (string, error) {
	var buffer strings.Builder
	writer := csv.NewWriter(&buffer)

	// Get column names
	cols, err := rows.Columns()
	if err != nil {
		return "", err
	}

	// Write header row
	err = writer.Write(cols)
	if err != nil {
		return "", err
	}

	// Write data rows
	for rows.Next() {
		// Scan row data into interface slice
		scanValues := make([]interface{}, len(cols))
		scanPointers := make([]interface{}, len(cols))
		for i := range scanValues {
			scanPointers[i] = &scanValues[i]
		}
		err = rows.Scan(scanPointers...)
		if err != nil {
			return "", err
		}

		// Convert scanned values to strings
		row := make([]string, len(cols))
		for i, v := range scanValues {
			val := reflect.ValueOf(v)
			if val.Kind() == reflect.String {
				row[i] = val.String()
			} else {
				row[i] = fmt.Sprint(v)
			}
		}

		// Write row to CSV
		err = writer.Write(row)
		if err != nil {
			return "", err
		}
	}

	writer.Flush()
	return buffer.String(), nil
}
