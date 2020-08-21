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
)

var mockServer = gsclient.Server{
	Properties: gsclient.ServerProperties{
		ObjectUUID: "xxx",
	},
}

type mockServerOp struct{}

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
	return nil
}
func (o mockServerOp) DeleteServer(ctx context.Context, id string) error {
	return nil
}
func (o mockServerOp) CreateServer(ctx context.Context, body gsclient.ServerCreateRequest) (gsclient.ServerCreateResponse, error) {
	return gsclient.ServerCreateResponse{}, nil
}
func (o mockServerOp) GetTemplateByName(ctx context.Context, name string) (gsclient.Template, error) {
	return gsclient.Template{}, nil
}
func (o mockServerOp) CreateStorage(ctx context.Context, body gsclient.StorageCreateRequest) (gsclient.CreateResponse, error) {
	return gsclient.CreateResponse{}, nil
}
func (o mockServerOp) CreateServerStorage(ctx context.Context, id string, body gsclient.ServerStorageRelationCreateRequest) error {
	return nil
}

func (o mockServerOp) UpdateServer(ctx context.Context, id string, body gsclient.ServerUpdateRequest) error {
	return nil
}

func Test_ServerCommmandDelete(t *testing.T) {
	r, w, _ := os.Pipe()
	rt, _ = runtime.NewTestRuntime()
	op := mockServerOp{}
	rt.SetServerOperator(op)
	os.Stdout = w
	cmd := serverRmCmd.Run
	cmd(new(cobra.Command), []string{"rm", mockServer.Properties.ObjectUUID})
	w.Close()
	out, _ := ioutil.ReadAll(r)
	assert.Equal(t, "", string(out))
}
