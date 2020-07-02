package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/gridscale/gsclient-go/v3"
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

func Test_StorageCmdOutput(t *testing.T) {
	marshalledMockStorage, _ := json.Marshal(mockStorageList)
	type testCase struct {
		expectedOutput string
		jsonFlag       bool
		idFlag         bool
		quietFlag      bool
	}
	testCases := []testCase{
		{
			expectedOutput: fmt.Sprintf("NAME  CAPACITY  CHANGETIME  STATUS  \n%s  %d        %s          %s  \n",
				mockStorage.Properties.Name,
				mockStorage.Properties.Capacity,
				strconv.FormatInt(int64(mockStorage.Properties.ChangeTime.Hour()), 10),
				mockStorage.Properties.Status,
			),
		},
		{
			jsonFlag:       true,
			idFlag:         false,
			quietFlag:      false,
			expectedOutput: string(marshalledMockStorage) + "\n",
		},
		{
			idFlag:    true,
			quietFlag: false,
			expectedOutput: fmt.Sprintf("NAME  CAPACITY  CHANGETIME  STATUS  ID           \n%s  %d        %s          %s  %s  \n",
				mockStorage.Properties.Name,
				mockStorage.Properties.Capacity,
				strconv.FormatInt(int64(mockStorage.Properties.ChangeTime.Hour()), 10),
				mockStorage.Properties.Status,
				mockStorage.Properties.ObjectUUID,
			),
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
		idFlag = test.idFlag
		quietFlag = test.quietFlag

		mockClient := mockStorageGetter{}
		cmd := produceStorageCmdRunFunc(mockClient)
		cmd(new(cobra.Command), []string{})

		// reset the flags back to the default values
		resetFlags()

		// close the write side of the pipe so that we can start reading from
		// the read side of the pipe
		w.Close()
		// Read the standard output's result from the read side of the pipe
		out, _ := ioutil.ReadAll(r)
		assert.Equal(t, test.expectedOutput, string(out))
	}
}
