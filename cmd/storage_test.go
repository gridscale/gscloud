package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/gridscale/gsclient-go/v3"
	"github.com/gridscale/gscloud/render"
	"github.com/gridscale/gscloud/runtime"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

var changeTime, _ = time.Parse(time.RFC3339, "2020-07-02T16:15:00+02:00")

var mockStorage = gsclient.Storage{
	Properties: gsclient.StorageProperties{
		ObjectUUID: "xxx-xxx-xxx",
		Name:       "test",
		Capacity:   10,
		Status:     "active",
		ChangeTime: gsclient.GSTime{
			Time: changeTime,
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

func (g mockClient) CloneStorage(ctx context.Context, id string) (gsclient.CreateResponse, error) {
	return gsclient.CreateResponse{}, nil
}

func (g mockClient) CreateStorage(ctx context.Context, body gsclient.StorageCreateRequest) (gsclient.CreateResponse, error) {
	return gsclient.CreateResponse{}, nil
}

func (g mockClient) GetDeletedStorages(ctx context.Context) ([]gsclient.Storage, error) {
	return []gsclient.Storage{}, nil
}

func (g mockClient) GetStorage(ctx context.Context, id string) (gsclient.Storage, error) {
	return gsclient.Storage{}, nil
}

func (g mockClient) GetStorageEventList(ctx context.Context, id string) ([]gsclient.Event, error) {
	return []gsclient.Event{}, nil
}

func (g mockClient) GetStoragesByLocation(ctx context.Context, id string) ([]gsclient.Storage, error) {
	return []gsclient.Storage{}, nil
}

func (g mockClient) UpdateStorage(ctx context.Context, id string, body gsclient.StorageUpdateRequest) error {
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
			"xxx-xxx-xxx",
			"test",
			"10",
			changeTime.Local().Format(time.RFC3339),
			"active",
		},
	}
	render.AsTable(buf, headers, rows, render.Options{})
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

		cmd := storageLsCmd.Run
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

	cmd := storageRmCmd.Run
	cmd(new(cobra.Command), []string{"rm", mockStorage.Properties.ObjectUUID})

	w.Close()
	out, _ := ioutil.ReadAll(r)
	assert.Equal(t, err.expectedOutput, string(out))
}
