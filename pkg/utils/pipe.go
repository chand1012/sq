package utils

import (
	"io"
	"os"
)

func ReadStdin() ([]byte, error) {
	bytes, err := io.ReadAll(os.Stdin)
	if err != nil {
		return nil, err // Return the error to the caller
	}
	return bytes, nil // Return the bytes to the caller
}
