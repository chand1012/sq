package db

import (
	"database/sql"
	"encoding/json"

	_ "github.com/glebarez/go-sqlite"

	"github.com/chand1012/sq/pkg/constants"
	"github.com/chand1012/sq/pkg/utils"
)

func FromJSON(b []byte, tableName string) (*sql.DB, string, error) {

	if tableName == "" {
		tableName = constants.TableName
	}

	db, tempName, err := createTempDB()
	if err != nil {
		return nil, tempName, err
	}

	// convert it to a slice of maps
	var data []map[string]any
	err = json.Unmarshal(b, &data)
	if err != nil {
		return nil, tempName, err
	}

	// map to store the column types
	typeMap := make(map[string]string)
	// loop through all the records and infer the column types
	// if a column type's second inference is different from the first, it's type is set to text
	for _, record := range data {
		for k, v := range record {
			if typeMap[k] == "" {
				typeMap[k] = reflectType(v)
			} else if typeMap[k] != reflectType(v) {
				typeMap[k] = "TEXT"
			}
		}
	}

	columns, types := utils.BreakOutMap(typeMap)

	// // preprocess the column names
	// columns = processColumnNames(columns)

	createQuery := genCreateTableQuery(tableName, columns, types)

	_, err = db.Exec(createQuery)
	if err != nil {
		return nil, tempName, err
	}

	insertQuery := genInsertQuery(tableName, columns)

	stmt, err := db.Prepare(insertQuery)
	if err != nil {
		return nil, tempName, err
	}

	for _, record := range data {
		var values []any
		for _, column := range columns {
			// if the column is not present in the record, insert a NULL
			if record[column] == nil {
				values = append(values, nil)
				continue
			}
			values = append(values, record[column])
		}
		_, err = stmt.Exec(values...)
		if err != nil {
			return nil, tempName, err
		}
	}

	return db, tempName, nil
}
