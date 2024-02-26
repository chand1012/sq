package db

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
)

func guessType(s string) string {
	// assume all numbers are floats
	if _, err := strconv.ParseFloat(s, 64); err == nil {
		return "REAL"
	}

	// Otherwise, it's a string
	return "TEXT"
}

// same as above, just takes generics and uses reflection
func reflectType(v any) string {
	switch v.(type) {
	case float32, float64, int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return "REAL"
	default:
		return "TEXT"
	}
}

// makes the column names all lowercase and replaces spaces with underscores
func processColumnNames(columnNames []string) []string {
	var processedColumnNames []string
	for _, name := range columnNames {
		processedColumnNames = append(processedColumnNames, strings.ReplaceAll(strings.ToLower(name), " ", "_"))
	}
	return processedColumnNames
}

func GetColumnNames(db *sql.DB, tableName string) ([]string, error) {
	// Build the query to retrieve column names
	query := "SELECT * FROM '" + tableName + "' LIMIT 1"

	// Execute the query
	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error querying column names: %w", err)
	}
	defer rows.Close()

	// Extract and print each column name
	columns, err := rows.Columns()
	if err != nil {
		return nil, fmt.Errorf("error getting column names: %w", err)
	}

	// Return the column names
	return columns, nil
}
