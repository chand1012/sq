package file_types

import (
	"testing"
)

func TestIsSQLiteFile(t *testing.T) {
	testCases := []struct {
		data     []byte
		expected bool
	}{
		{[]byte("SQLite format 3\x00"), true},
		{[]byte(""), false},
	}

	for _, testCase := range testCases {
		actual := IsSQLiteFile(testCase.data)
		if actual != testCase.expected {
			t.Errorf("Expected %v, but got %v", testCase.expected, actual)
		}
	}
}

func TestIsValidJSON(t *testing.T) {
	testCases := []struct {
		data     []byte
		expected bool
	}{
		{[]byte("{\"key\": \"value\"}"), true},
		{[]byte("{"), false},
	}

	for _, testCase := range testCases {
		actual := IsValidJSON(testCase.data)
		if actual != testCase.expected {
			t.Errorf("Expected %v, but got %v", testCase.expected, actual)
		}
	}
}

func TestIsValidJSONL(t *testing.T) {
	testCases := []struct {
		data     []byte
		expected bool
	}{
		{[]byte("{\"key\": \"value\"}\n{\"key2\": \"value2\"}"), true},
		{[]byte("{\"key\": \"value\"}\n{\"key2\": "), false},
	}

	for _, testCase := range testCases {
		actual := IsValidJSONL(testCase.data)
		if actual != testCase.expected {
			t.Errorf("Expected %v, but got %v", testCase.expected, actual)
		}
	}
}
