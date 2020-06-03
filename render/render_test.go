package render

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWithByteBuffer(t *testing.T) {
	out := new(bytes.Buffer)
	Table(out, []string{"test", "version", "text"}, [][]string{{"1", "version 1", "empty"}, {"a2", "b2", "c3"}})

	countedLines := strings.Count(out.String(), "\n")
	assert.Equal(t, countedLines, 3)

	fields := strings.Fields(out.String())
	assert.Equal(t, fields[0], "TEST")
}
