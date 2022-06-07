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
