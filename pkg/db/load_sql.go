package db

import (
	"database/sql"
	"errors"
	"os"

	_ "github.com/glebarez/go-sqlite"
)

func createTempDB() (*sql.DB, string, error) {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		return nil, "", err
	}

	return db, "", nil
}

// check if the file is a valid sqlite db
func IsValidDB(fileName string) bool {
	db, err := sql.Open("sqlite", fileName)
	if err != nil {
		return false
	}
	defer db.Close()

	// try to query the db
	_, err = db.Query("SELECT * FROM sqlite_master")
	return err == nil
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
	if !IsValidDB(fileName) {
		return nil, "", errors.New("file is not a valid SQLite database")
	}
	db, err := sql.Open("sqlite", fileName)
	if err != nil {
		return nil, fileName, err
	}

	return db, fileName, nil
}
