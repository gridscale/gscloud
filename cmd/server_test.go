package cmd

import (
	"context"
	"io/ioutil"
	"os"
	"testing"

	"github.com/gridscale/gsclient-go/v3"
	"github.com/gridscale/gscloud/runtime"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var mockServer = gsclient.Server{
	Properties: gsclient.ServerProperties{
		ObjectUUID: "xxx",
	},
}

type mockServerOp struct {
	mock.Mock
}

func (o mockServerOp) GetServerList(ctx context.Context) ([]gsclient.Server, error) {
	return nil, nil
}

func (o mockServerOp) StartServer(ctx context.Context, id string) error {
	return nil
}

func (o mockServerOp) StopServer(ctx context.Context, id string) error {
	return nil
}

func (o mockServerOp) ShutdownServer(ctx context.Context, id string) error {
	args := o.Called(id)
	return args.Error(0)
}

func (o mockServerOp) DeleteServer(ctx context.Context, id string) error {
	args := o.Called(id)
	return args.Error(0)
}

func (o mockServerOp) CreateServer(ctx context.Context, body gsclient.ServerCreateRequest) (gsclient.ServerCreateResponse, error) {
	return gsclient.ServerCreateResponse{}, nil
}

func (o mockServerOp) CreateStorage(ctx context.Context, body gsclient.StorageCreateRequest) (gsclient.CreateResponse, error) {
	return gsclient.CreateResponse{}, nil
}

func (o mockServerOp) UpdateServer(ctx context.Context, id string, body gsclient.ServerUpdateRequest) error {
	return nil
}

func (o mockServerOp) GetDeletedServers(ctx context.Context) ([]gsclient.Server, error) {
	return []gsclient.Server{}, nil
}

func (o mockServerOp) GetServer(ctx context.Context, id string) (gsclient.Server, error) {
	return gsclient.Server{}, nil
}

func (o mockServerOp) GetServerEventList(ctx context.Context, id string) ([]gsclient.Event, error) {
	return []gsclient.Event{}, nil
}

func (o mockServerOp) GetServerMetricList(ctx context.Context, id string) ([]gsclient.ServerMetric, error) {
	return []gsclient.ServerMetric{}, nil
}

func (o mockServerOp) GetServersByLocation(ctx context.Context, id string) ([]gsclient.Server, error) {
	return []gsclient.Server{}, nil
}

func (o mockServerOp) IsServerOn(ctx context.Context, id string) (bool, error) {
	return false, nil
}

func Test_ServerCommmandDelete(t *testing.T) {
	r, w, _ := os.Pipe()
	rt, _ = runtime.NewTestRuntime()
	op := mockServerOp{}
	op.On("DeleteServer", mock.Anything).Return(nil)
	rt.SetServerOperator(op)
	os.Stdout = w
	cmd := serverRmCmd.Run
	cmd(new(cobra.Command), []string{"rm", mockServer.Properties.ObjectUUID})
	w.Close()
	out, _ := ioutil.ReadAll(r)
	assert.Equal(t, "", string(out))
}
