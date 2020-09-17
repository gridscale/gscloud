package runtime

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_SelectAccount(t *testing.T) {
	testAccount := AccountEntry{
		Name: "test",
	}
	testConfig := Config{Accounts: []AccountEntry{testAccount}}
	rt, err := NewRuntime(testConfig, "test")
	assert.Nil(t, err)
	assert.Equal(t, rt.Account(), "test")
}

