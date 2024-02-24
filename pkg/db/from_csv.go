package db

import (
	"bytes"
	"database/sql"
	"encoding/csv"

	"github.com/chand1012/sq/pkg/utils"
	_ "github.com/glebarez/go-sqlite"
)

func FromCSV(b []byte, tableName string) (*sql.DB, string, error) {
	if tableName == "" {
		tableName = "sq_table"
	}
	// create a new temp database
	db, tempName, err := createTempDB()
	if err != nil {
		return nil, "", err
	}

	// load the csv using the csv reader
	csvReader := csv.NewReader(bytes.NewReader(b))
	// get the headers for column names
	headers, err := csvReader.Read()
	if err != nil {
		return nil, tempName, err
	}
	// also get the first row for type inference
	firstRow, err := csvReader.Read()
	if err != nil {
		return nil, tempName, err
	}
	columnTypes := make([]string, len(headers))
	for i, v := range firstRow {
		// attempt to infer the column type
		columnTypes[i] = guessType(v)
	}

	// construct the query that will create the table
	createTableQuery := genCreateTableQuery(tableName, headers, columnTypes)

	// execute the query to create the table
	_, err = db.Exec(createTableQuery)
	if err != nil {
		return nil, tempName, err
	}

	// construct the query that will insert the data
	insertQuery := genInsertQuery(tableName, headers)

	// prepare the insert statement
	stmt, err := db.Prepare(insertQuery)
	if err != nil {
		return nil, tempName, err
	}

	// iterate over the csv and insert the data
	for {
		record, err := csvReader.Read()
		if err != nil {
			break
		}
		// convert the records into a slice of any
		_, err = stmt.Exec(utils.StringSliceToAnySlice(record)...)
		if err != nil {
			return nil, tempName, err
		}
	}

	return db, tempName, nil
}
