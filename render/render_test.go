package render

import (
	"bytes"
	"strings"
	"testing"
)

func TestWithByteBuffer(t *testing.T) {
	out := new(bytes.Buffer)
	Table(out, []string{"test", "version", "text"}, [][]string{{"1", "version 1", "empty"}, {"a2", "b2", "c3"}})

	countedLines := strings.Count(out.String(), "\n")
	if countedLines != 4 {
		t.Fail()
	}

	fields := strings.Fields(out.String())
	if fields[0] != "TEST" {
		t.Fail()
	}

}
