package runtime

import (
	"os"
	"strings"
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

func Test_SelectAccountSettings(t *testing.T) {
	defer resetEnv(os.Environ())

	os.Setenv("GRIDSCALE_UUID", "passedUserId")
	os.Setenv("GRIDSCALE_TOKEN", "passedToken")
	os.Setenv("GRIDSCALE_URL", "passed.example.com")

	testAccount := AccountEntry{
		Name:   "test",
		UserID: "test",
		Token:  "test",
		URL:    "test.example.com",
	}
	testConfig := Config{Accounts: []AccountEntry{testAccount}}
	rt, err := NewRuntime(testConfig, "test")

	assert.Nil(t, err)
	assert.Equal(t, "passedUserId", rt.config.Accounts[0].UserID)
	assert.Equal(t, "passedToken", rt.config.Accounts[0].Token)
	assert.Equal(t, "passed.example.com", rt.config.Accounts[0].URL)
}

func resetEnv(environ []string) {
	os.Clearenv()

	for _, s := range environ {
		splitString := strings.Split(s, "=")

		os.Setenv(splitString[0], splitString[1])
	}
}
