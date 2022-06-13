package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_FileExists(t *testing.T) {
	testCases := []struct {
		Filename string
		Expected bool
	}{
		{
			Filename: "utils_test.go",
			Expected: true,
		},
		{
			Filename: ".",
			Expected: false,
		},
		{
			Filename: "1e35c3fc03706c064e95f34f8ca15256f77789aa",
			Expected: false,
		},
	}

	for _, test := range testCases {
		assert.Equal(t, test.Expected, FileExists(test.Filename))
	}
}

func Test_StringLess(t *testing.T) {
	testCases := []struct {
		string1  string
		string2  string
		expected bool
	}{
		{
			string1:  "abcd",
			string2:  "abcde",
			expected: true,
		},
		{
			string1:  "abcde",
			string2:  "abcde",
			expected: false,
		},
		{
			string1:  "abcdef",
			string2:  "abcde",
			expected: false,
		},
		{
			string1:  "",
			string2:  "abcde",
			expected: true,
		},
	}

	for _, test := range testCases {
		sorter := StringSorter{test.string1, test.string2}
		assert.Equal(t, test.expected, sorter.Less(0, 1))
	}
}

func Test_contains(t *testing.T) {
	testCases := []struct {
		Slice    []string
		Item     string
		Expected bool
	}{
		{
			Slice:    []string{"test1", "test2", "test3"},
			Item:     "test3",
			Expected: true,
		},
		{
			Slice:    []string{"test1", "test2", "test3"},
			Item:     "test4",
			Expected: false,
		},
	}

	for _, test := range testCases {
		assert.Equal(t, test.Expected, Contains(test.Slice, test.Item))
	}
}
