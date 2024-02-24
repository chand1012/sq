package db

import (
	"bytes"
	"database/sql"
	"encoding/json"

	_ "github.com/glebarez/go-sqlite"

	"github.com/chand1012/sq/pkg/constants"
	"github.com/chand1012/sq/pkg/utils"
)

func FromJSONL(b []byte, tableName string) (*sql.DB, string, error) {
	if tableName == "" {
		tableName = constants.TableName
	}

	db, tempName, err := createTempDB()
	if err != nil {
		return nil, "", err
	}

	// separate the jsonl into lines
	lines := bytes.Split(b, []byte("\n"))

	typeMap := make(map[string]string)
	// loop through each lines and get all the keys
	for _, line := range lines {
		if len(line) == 0 {
			continue
		}
		// convert each line to a map
		var data map[string]any
		err = json.Unmarshal(line, &data)
		if err != nil {
			return nil, tempName, err
		}
		// loop through the map and get all the keys
		for k, v := range data {
			if typeMap[k] == "" {
				typeMap[k] = reflectType(v)
			} else if typeMap[k] != reflectType(v) {
				typeMap[k] = "TEXT"
			}
		}
	}

	columns, types := utils.BreakOutMap(typeMap)

	// preprocess the column names
	columns = processColumnNames(columns)

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

	// iterate over the jsonl and insert the data
	for _, line := range lines {
		// convert the records into a slice of any
		var data map[string]any
		err = json.Unmarshal(line, &data)
		if len(line) == 0 {
			continue
		}
		if err != nil {
			return nil, tempName, err
		}
		// convert the map to a slice of any
		args := make([]any, len(columns))
		for i, col := range columns {
			args[i] = data[col]
		}
		_, err = stmt.Exec(args...)
		if err != nil {
			return nil, tempName, err
		}
	}

	return db, tempName, nil
}
