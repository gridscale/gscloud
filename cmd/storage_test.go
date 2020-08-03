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
	"github.com/gridscale/gscloud/runtime"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

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

type mockClient struct{}

func (g mockClient) GetStorageList(ctx context.Context) ([]gsclient.Storage, error) {
	return mockStorageList, nil
}

func (g mockClient) DeleteStorage(ctx context.Context, id string) error {
	return nil
}

func Test_StorageListCmd(t *testing.T) {
	marshalledMockStorage, _ := json.Marshal(mockStorageList)
	type testCase struct {
		expectedOutput string
		jsonFlag       bool
		quietFlag      bool
	}
	buf := new(bytes.Buffer)
	headers := []string{"id", "name", "capacity", "changetime", "status"}
	rows := [][]string{
		{
			mockStorage.Properties.ObjectUUID,
			mockStorage.Properties.Name,
			strconv.Itoa(mockStorage.Properties.Capacity),
			strconv.FormatInt(int64(mockStorage.Properties.ChangeTime.Hour()), 10),
			mockStorage.Properties.Status,
		},
	}
	render.Table(buf, headers, rows)
	testCases := []testCase{
		{
			expectedOutput: buf.String(),
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
		r, w, _ := os.Pipe()
		os.Stdout = w

		jsonFlag = test.jsonFlag
		quietFlag = test.quietFlag

		mockClient := mockClient{}
		rt, _ = runtime.NewTestRuntime()
		rt.SetStorageOperator(mockClient)

		cmd := storageListCmd.Run
		cmd(new(cobra.Command), []string{})

		resetFlags()

		w.Close()
		out, _ := ioutil.ReadAll(r)
		assert.Equal(t, test.expectedOutput, string(out))
	}
}

func Test_StorageCmdDelete(t *testing.T) {
	type testCase struct {
		expectedOutput string
	}
	err := testCase{expectedOutput: ""}
	r, w, _ := os.Pipe()
	os.Stdout = w
	mockClient := mockClient{}
	rt, _ = runtime.NewTestRuntime()
	rt.SetStorageOperator(mockClient)

	cmd := storageRemoveCmd.Run
	cmd(new(cobra.Command), []string{"rm", mockStorage.Properties.ObjectUUID})

	w.Close()
	out, _ := ioutil.ReadAll(r)
	assert.Equal(t, err.expectedOutput, string(out))
}
