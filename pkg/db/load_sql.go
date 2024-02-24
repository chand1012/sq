package db

import (
	"database/sql"
	"os"

	_ "github.com/glebarez/go-sqlite"
)

func createTempDB() (*sql.DB, string, error) {
	tmpFile, err := os.CreateTemp(os.TempDir(), "sq-*.sql")
	if err != nil {
		return nil, "", err
	}
	defer tmpFile.Close()

	db, err := sql.Open("sqlite", tmpFile.Name())
	if err != nil {
		return nil, "", err
	}

	return db, tmpFile.Name(), nil
}

// load a sql file from bytes
func LoadStdin(bytes []byte) (*sql.DB, string, error) {
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

func LoadFile(fileName string) (*sql.DB, string, error) {
	db, err := sql.Open("sqlite", fileName)
	if err != nil {
		return nil, fileName, err
	}

	return db, fileName, nil
}
