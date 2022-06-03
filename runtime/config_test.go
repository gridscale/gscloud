package runtime

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_CommandWithoutConfig(t *testing.T) {
	testCases := []struct {
		Args     []string
		Expected bool
	}{
		{
			Args:     []string{"gscloud", "version"},
			Expected: true,
		},
		{
			Args:     []string{"gscloud", "server", "create"},
			Expected: false,
		},
		{
			Args:     []string{"gscloud", "server", "create", "--account", "completion"},
			Expected: false,
		},
	}

	for _, test := range testCases {
		assert.Equal(t, test.Expected, CommandWithoutConfig(test.Args))
	}
}
