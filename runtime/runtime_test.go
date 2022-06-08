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
		ExpectedAccount      ProjectEntry
		ExpectedErrorIsNil   bool
	}

	testAccount := ProjectEntry{
		Name:   "test",
		UserID: "test",
		Token:  "test",
		URL:    "test.example.com",
	}

	testCases := []testCase{
		{
			Configuration:        Config{[]ProjectEntry{testAccount}},
			AccountName:          testAccount.Name,
			Environment:          []string{},
			ExpectedRuntimeIsNil: false,
			ExpectedAccount:      testAccount,
			ExpectedErrorIsNil:   true,
		},
		{
			Configuration:        Config{[]ProjectEntry{testAccount}},
			AccountName:          "default",
			Environment:          []string{},
			ExpectedRuntimeIsNil: true,
			ExpectedAccount:      ProjectEntry{},
			ExpectedErrorIsNil:   false,
		},
		{
			Configuration:        Config{[]ProjectEntry{}},
			AccountName:          "default",
			Environment:          []string{},
			ExpectedRuntimeIsNil: false,
			ExpectedAccount:      ProjectEntry{},
			ExpectedErrorIsNil:   true,
		},
		{
			Configuration:        Config{[]ProjectEntry{testAccount}},
			AccountName:          testAccount.Name,
			Environment:          []string{"GRIDSCALE_UUID=envUserId", "GRIDSCALE_TOKEN=envToken", "GRIDSCALE_URL=env.example.com"},
			ExpectedRuntimeIsNil: false,
			ExpectedAccount:      ProjectEntry{Name: testAccount.Name, UserID: "envUserId", Token: "envToken", URL: "env.example.com"},
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
			assert.Equal(t, test.ExpectedAccount, rt.account)
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
