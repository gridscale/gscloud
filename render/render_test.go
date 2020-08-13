package render

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_AsTable(t *testing.T) {
	out := new(bytes.Buffer)
	AsTable(out, []string{"test", "version", "text"}, [][]string{{"1", "version 1", "empty"}, {"a2", "b2", "c3"}}, Options{})

	countedLines := strings.Count(out.String(), "\n")
	assert.Equal(t, countedLines, 3)

	fields := strings.Fields(out.String())
	assert.Equal(t, fields[0], "TEST")
}

func Test_AsTableWithoutHeader(t *testing.T) {
	out := new(bytes.Buffer)
	opts := Options{
		NoHeader: true,
	}
	AsTable(out, []string{"a", "b"}, [][]string{{"1", "2"}, {"1", "2"}}, opts)

	countedLines := strings.Count(out.String(), "\n")
	assert.Equal(t, countedLines, 2)
}

func Test_AsJSON(t *testing.T) {
	type someStruct struct {
		Test string `json:"test"`
	}
	val := someStruct{Test: "test_value"}
	expectedOutput := "[{\"test\":\"test_value\"}]\n"
	buffer := new(bytes.Buffer)
	AsJSON(buffer, val)
	assert.Equal(t, expectedOutput, buffer.String())
}
