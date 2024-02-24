package utils

import (
	"errors"
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
