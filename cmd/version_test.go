package cmd

import (
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func Test_VersionCmdOutput(t *testing.T) {
	r, w, _ := os.Pipe()
	os.Stdout = w

	GitCommit = "foo"
	Version = "bar"
	expectedOutput := fmt.Sprintf("Version:\t%s\nGit commit:\t%s\n", Version, GitCommit)

	versionCmd.Run(new(cobra.Command), []string{})
	w.Close()

	out, _ := io.ReadAll(r)
	assert.Equal(t, expectedOutput, string(out))
}
