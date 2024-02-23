package db

import (
	"database/sql"
	"os"

	_ "github.com/glebarez/go-sqlite"
)

// load a sql file from bytes
func LoadSQLFileFromBytes(bytes []byte) (*sql.DB, string, error) {
	tmpFile, err := os.CreateTemp(os.TempDir(), "sq-*.sql")
	if err != nil {
		return nil, "", err
	}
	defer tmpFile.Close()

	_, err = tmpFile.Write(bytes)
	if err != nil {
		return nil, "", err
	}

	db, err := sql.Open("sqlite", tmpFile.Name())
	if err != nil {
		return nil, "", err
	}

	return db, tmpFile.Name(), nil
}

func LoadSQLFile(fileName string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", fileName)
	if err != nil {
		return nil, err
	}

	return db, nil
}
