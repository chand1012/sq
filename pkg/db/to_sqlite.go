package db

import (
	"database/sql"
	"io"
)

func createDB(fileName string) (*sql.DB, string, error) {
	db, err := sql.Open("sqlite", fileName)
	if err != nil {
		return nil, fileName, err
	}

	return db, fileName, nil
}

func getNextRow(rows *sql.Rows) ([]any, error) {
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	values := make([]interface{}, len(columns))
	valuePtrs := make([]interface{}, len(columns))

	for i := range columns {
		valuePtrs[i] = &values[i]
	}

	// call next to get the next row
	if !rows.Next() {
		return nil, io.EOF
	}

	err = rows.Scan(valuePtrs...)
	if err != nil {
		return nil, err
	}

	return values, nil
}

func RowsToSQLite(rows *sql.Rows, tableName string, fileName string) error {
	d, _, err := createDB(fileName)
	if err != nil {
		return err
	}
	defer d.Close()
	columns, err := rows.Columns()
	if err != nil {
		return err
	}

	firstRow, err := getNextRow(rows)
	if err != nil {
		return err
	}

	var columnTypes []string
	for _, value := range firstRow {
		t := reflectType(value)
		columnTypes = append(columnTypes, t)
	}

	createQuery := genCreateTableQuery(tableName, columns, columnTypes)
	_, err = d.Exec(createQuery)
	if err != nil {
		return err
	}

	insertQuery := genInsertQuery(tableName, columns)
	stmt, err := d.Prepare(insertQuery)
	if err != nil {
		return err
	}
	// we already pulled the first row, so we need to insert it
	_, err = stmt.Exec(firstRow...)
	for {
		row, err := getNextRow(rows)
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		if row == nil {
			break
		}

		_, err = stmt.Exec(row...)
		if err != nil {
			return err
		}
	}

	return err
}
