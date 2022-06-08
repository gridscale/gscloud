package runtime

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NewRuntime(t *testing.T) {
	type testCase struct {
		Configuration        Config
		AccountName          string
		Environment          []string
		ExpectedRuntimeIsNil bool
		ExpectedAccount      AccountEntry
		ExpectedErrorIsNil   bool
	}

	testAccount := AccountEntry{
		Name:   "test",
		UserID: "test",
		Token:  "test",
		URL:    "test.example.com",
	}

	testCases := []testCase{
		{
			Configuration:        Config{[]AccountEntry{testAccount}},
			AccountName:          testAccount.Name,
			Environment:          []string{},
			ExpectedRuntimeIsNil: false,
			ExpectedAccount:      testAccount,
			ExpectedErrorIsNil:   true,
		},
		{
			Configuration:        Config{[]AccountEntry{testAccount}},
			AccountName:          "default",
			Environment:          []string{},
			ExpectedRuntimeIsNil: true,
			ExpectedAccount:      AccountEntry{},
			ExpectedErrorIsNil:   false,
		},
		{
			Configuration:        Config{[]AccountEntry{}},
			AccountName:          "default",
			Environment:          []string{},
			ExpectedRuntimeIsNil: false,
			ExpectedAccount:      AccountEntry{},
			ExpectedErrorIsNil:   true,
		},
		{
			Configuration:        Config{[]AccountEntry{testAccount}},
			AccountName:          testAccount.Name,
			Environment:          []string{"GRIDSCALE_UUID=envUserId", "GRIDSCALE_TOKEN=envToken", "GRIDSCALE_URL=env.example.com"},
			ExpectedRuntimeIsNil: false,
			ExpectedAccount:      AccountEntry{Name: testAccount.Name, UserID: "envUserId", Token: "envToken", URL: "env.example.com"},
			ExpectedErrorIsNil:   true,
		},
	}

	for _, test := range testCases {
		oldEnviron := os.Environ()
		resetEnv(test.Environment)

		rt, err := NewRuntime(test.Configuration, test.AccountName, false)

		assert.Equal(t, test.ExpectedErrorIsNil, err == nil)
		assert.Equal(t, test.ExpectedRuntimeIsNil, rt == nil)

		if rt != nil {
			for _, ac := range rt.config.Accounts {
				if ac.Name == rt.accountName {
					assert.Equal(t, test.ExpectedAccount, ac)
					break
				}
			}
		}

		resetEnv(oldEnviron)
	}
}

func resetEnv(environ []string) {
	os.Clearenv()

	for _, s := range environ {
		splitString := strings.Split(s, "=")

		if len(splitString) >= 2 {
			os.Setenv(splitString[0], splitString[1])
		}
	}
}
