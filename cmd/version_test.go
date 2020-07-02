package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func Test_VersionCmdOutput(t *testing.T) {
	// Create a pipe
	r, w, _ := os.Pipe()
	// Change the standard output to the write side of the pipe
	os.Stdout = w

	// set GitCommit and Version to dummy values
	GitCommit = "test_commit"
	Version = "test_version"

	expectedOutput := fmt.Sprintf("Version:\t%s\nGit commit:\t%s\n", Version, GitCommit)
	//run the cmd
	versionCmdRun(new(cobra.Command), []string{})
	// close the write side of the pipe so that we can start reading from
	// the read side of the pipe
	w.Close()
	// Read the standard output's result from the read side of the pipe
	out, _ := ioutil.ReadAll(r)
	assert.Equal(t, expectedOutput, string(out))
}
