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
		ProjectName          string
		Environment          []string
		ExpectedRuntimeIsNil bool
		ExpectedProject      ProjectEntry
		ExpectedErrorIsNil   bool
	}

	testProject := ProjectEntry{
		Name:   "test",
		UserID: "test",
		Token:  "test",
		URL:    "test.example.com",
	}

	testCases := []testCase{
		{
			Configuration:        Config{[]ProjectEntry{testProject}},
			ProjectName:          testProject.Name,
			Environment:          []string{},
			ExpectedRuntimeIsNil: false,
			ExpectedProject:      testProject,
			ExpectedErrorIsNil:   true,
		},
		{
			Configuration:        Config{[]ProjectEntry{testProject}},
			ProjectName:          "default",
			Environment:          []string{},
			ExpectedRuntimeIsNil: true,
			ExpectedProject:      ProjectEntry{},
			ExpectedErrorIsNil:   false,
		},
		{
			Configuration:        Config{[]ProjectEntry{}},
			ProjectName:          "default",
			Environment:          []string{},
			ExpectedRuntimeIsNil: false,
			ExpectedProject:      ProjectEntry{},
			ExpectedErrorIsNil:   true,
		},
		{
			Configuration:        Config{[]ProjectEntry{testProject}},
			ProjectName:          testProject.Name,
			Environment:          []string{"GRIDSCALE_UUID=envUserId", "GRIDSCALE_TOKEN=envToken", "GRIDSCALE_URL=env.example.com"},
			ExpectedRuntimeIsNil: false,
			ExpectedProject:      ProjectEntry{Name: testProject.Name, UserID: "envUserId", Token: "envToken", URL: "env.example.com"},
			ExpectedErrorIsNil:   true,
		},
	}

	for _, test := range testCases {
		oldEnviron := os.Environ()
		resetEnv(test.Environment)

		rt, err := NewRuntime(test.Configuration, test.ProjectName, false)

		assert.Equal(t, test.ExpectedErrorIsNil, err == nil)
		assert.Equal(t, test.ExpectedRuntimeIsNil, rt == nil)

		if rt != nil {
			for _, ac := range rt.config.Projects {
				if ac.Name == rt.ProjectName {
					assert.Equal(t, test.ExpectedProject, ac)
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
