package utils

import (
	"errors"
	"regexp"
)

func GetJSONKeys(j any) ([]string, error) {
	var keys []string
	switch t := j.(type) {
	case map[string]interface{}:
		for key := range t {
			keys = append(keys, key)
		}
	default:
		return nil, errors.New("not a map, no keys")
	}
	return keys, nil
}

func IsJSONArray(data []byte) bool {
	// check to make sure the entire data is
	// wrapped in brackets
	// use regex
	jsonArrayRegex := regexp.MustCompile(`^\[.*\]$`)
	return jsonArrayRegex.MatchString(string(data))
}
