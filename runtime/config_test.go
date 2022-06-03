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
		{
			Args:     []string{"gscloud", "--json", "completion"},
			Expected: true,
		},
	}

	for _, test := range testCases {
		assert.Equal(t, test.Expected, CommandWithoutConfig(test.Args))
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
		assert.Equal(t, test.Expected, contains(test.Slice, test.Item))
	}
}
