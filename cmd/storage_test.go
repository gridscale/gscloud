package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/gridscale/gsclient-go/v3"
	"github.com/gridscale/gscloud/render"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

// mockStorage is a mock storage data
var mockStorage = gsclient.Storage{
	Properties: gsclient.StorageProperties{
		ObjectUUID: "xxx-xxx-xxx",
		Name:       "test",
		Capacity:   10,
		Status:     "active",
		ChangeTime: gsclient.GSTime{
			Time: time.Now(),
		},
	},
}
var mockStorageList = []gsclient.Storage{
	mockStorage,
}

// mockStorageGetter implements storageGetter interface,
// it is used for mocking data
type mockStorageGetter struct {
}

// GetStorageList returns a mocking list of storages
func (g mockStorageGetter) GetStorageList(ctx context.Context) ([]gsclient.Storage, error) {
	return mockStorageList, nil
}

func (g mockStorageGetter) DeleteStorage(ctx context.Context, id string) error {
	return nil
}

func Test_StorageCmdOutput(t *testing.T) {
	marshalledMockStorage, _ := json.Marshal(mockStorageList)
	type testCase struct {
		expectedOutput string
		jsonFlag       bool
		quietFlag      bool
	}
	mockRes := new(bytes.Buffer)
	headers := []string{"id", "name", "capacity", "changetime", "status"}
	mockRows := [][]string{
		{
			mockStorage.Properties.ObjectUUID,
			mockStorage.Properties.Name,
			strconv.Itoa(mockStorage.Properties.Capacity),
			strconv.FormatInt(int64(mockStorage.Properties.ChangeTime.Hour()), 10),
			mockStorage.Properties.Status,
		},
	}
	render.Table(mockRes, headers, mockRows)
	testCases := []testCase{
		{
			expectedOutput: mockRes.String(),
		},
		{
			jsonFlag:       true,
			expectedOutput: string(marshalledMockStorage) + "\n",
		},
		{
			quietFlag:      true,
			expectedOutput: mockStorage.Properties.ObjectUUID + "\n",
		},
	}
	for _, test := range testCases {
		// Create a pipe
		r, w, _ := os.Pipe()
		// Change the standard output to the write side of the pipe
		os.Stdout = w

		// Set the flags to test values
		jsonFlag = test.jsonFlag
		quietFlag = test.quietFlag

		mockClient := mockStorageGetter{}
		cmd := produceStorageCmdRunFunc(mockClient, storageListAction)
		cmd(new(cobra.Command), []string{})

		// reset the flags back to the default values
		resetFlags()

		// close the write side of the pipe so that we can start reading from
		// the read side of the pipe
		w.Close()
		out, _ := ioutil.ReadAll(r)
		assert.Equal(t, test.expectedOutput, string(out))
	}
}

func Test_StorageCmdDeleteOutput(t *testing.T) {
	type testCase struct {
		expectedOutput string
	}
	err := testCase{expectedOutput: ""}
	r, w, _ := os.Pipe()
	// Change the standard output to the write side of the pipe
	os.Stdout = w
	mockClient := mockStorageGetter{}
	cmd := produceStorageCmdRunFunc(mockClient, storageDeleteAction)
	cmd(new(cobra.Command), []string{"rm", mockStorage.Properties.ObjectUUID})
	// close the write side of the pipe so that we can start reading from
	// the read side of the pipe
	w.Close()
	out, _ := ioutil.ReadAll(r)
	assert.Equal(t, err.expectedOutput, string(out))
}
