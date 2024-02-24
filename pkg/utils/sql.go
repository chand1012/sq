package utils

import (
	"encoding/json"

	"github.com/krasun/gosqlparser"
)

type sqlQuery struct {
	Table string
}

func GetTableName(query string) (string, error) {
	parsedQuery, err := gosqlparser.Parse(query)
	if err != nil {
		return "", err
	}

	j, err := json.Marshal(parsedQuery)
	if err != nil {
		return "", err
	}

	// unmarsal the json into a struct
	var q sqlQuery
	err = json.Unmarshal(j, &q)

	return q.Table, err
}
