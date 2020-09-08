package runtime

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/gridscale/gsclient-go/v3"
	"github.com/kirsle/configdir"
	"github.com/spf13/viper"
)

// Runtime holds all run-time infos.
type Runtime struct {
	account string
	client  interface{}
}

// StorageOperator represents operations on storages.
type StorageOperator interface {
	DeleteStorage(ctx context.Context, id string) error
	GetStorageList(ctx context.Context) ([]gsclient.Storage, error)
}

// TemplateOperator represents operations on templates.
type TemplateOperator interface {
	GetTemplateList(ctx context.Context) ([]gsclient.Template, error)
}

// KubernetesOperator amalgamates operations for Kubernetes PaaS.
type KubernetesOperator interface {
	RenewK8sCredentials(ctx context.Context, id string) error
	GetPaaSService(ctx context.Context, id string) (gsclient.PaaSService, error)
}

// SSHKeyOperator is used for manipulating SSH keys.
type SSHKeyOperator interface {
	GetSshkeyList(ctx context.Context) ([]gsclient.Sshkey, error)
	CreateSshkey(ctx context.Context, body gsclient.SshkeyCreateRequest) (gsclient.CreateResponse, error)
	DeleteSshkey(ctx context.Context, id string) error
}

// ServerOperator contains all operations regarding server objects.
type ServerOperator interface {
	GetServerList(ctx context.Context) ([]gsclient.Server, error)
	StartServer(ctx context.Context, id string) error
	StopServer(ctx context.Context, id string) error
	ShutdownServer(ctx context.Context, id string) error
	DeleteServer(ctx context.Context, id string) error
	CreateServer(ctx context.Context, body gsclient.ServerCreateRequest) (gsclient.ServerCreateResponse, error)
	GetTemplateByName(ctx context.Context, name string) (gsclient.Template, error)
	CreateStorage(ctx context.Context, body gsclient.StorageCreateRequest) (gsclient.CreateResponse, error)
	CreateServerStorage(ctx context.Context, id string, body gsclient.ServerStorageRelationCreateRequest) error
	UpdateServer(ctx context.Context, id string, body gsclient.ServerUpdateRequest) error
}

// NetworkOperator interface that amalgamates all operations regarding network objects.
type NetworkOperator interface {
	GetNetworkList(ctx context.Context) ([]gsclient.Network, error)
	DeleteNetwork(ctx context.Context, id string) error
}

// Account is the current selected account.
func (r *Runtime) Account() string {
	return r.account
}

// StorageOperator return an operation to remove a storage.
func (r *Runtime) StorageOperator() StorageOperator {
	if UnderTest() {
		return r.client.(StorageOperator)
	}
	return r.client.(*gsclient.Client)
}

// SetStorageOperator set operation to delete storages.
func (r *Runtime) SetStorageOperator(op StorageOperator) {
	if !UnderTest() {
		panic("unexpected use")
	}
	r.client = op
}

// TemplateOperator return an operation to remove a storage.
func (r *Runtime) TemplateOperator() TemplateOperator {
	if UnderTest() {
		return r.client.(TemplateOperator)
	}
	return r.client.(*gsclient.Client)
}

// SetTemplateOperator set operation to delete storages.
func (r *Runtime) SetTemplateOperator(op TemplateOperator) {
	if !UnderTest() {
		panic("unexpected use")
	}
	r.client = op
}

// KubernetesOperator return operation relating to Kubernetes managed services.
func (r *Runtime) KubernetesOperator() KubernetesOperator {
	if UnderTest() {
		return r.client.(KubernetesOperator)
	}
	return r.client.(*gsclient.Client)
}

// SetKubernetesOperator set Kubernetes PaaS operation.
func (r *Runtime) SetKubernetesOperator(op KubernetesOperator) {
	if !UnderTest() {
		panic("unexpected use")
	}
	r.client = op
}

// SSHKeyOperator return operation to manipulate SSH keys.
func (r *Runtime) SSHKeyOperator() SSHKeyOperator {
	if UnderTest() {
		return r.client.(SSHKeyOperator)
	}
	return r.client.(*gsclient.Client)
}

// SetSSHKeyOperator set operation to manipulate SSH keys.
func (r *Runtime) SetSSHKeyOperator(op SSHKeyOperator) {
	if !UnderTest() {
		panic("unexpected use")
	}
	r.client = op
}

// ServerOperator return operation for server objects.
func (r *Runtime) ServerOperator() ServerOperator {
	if UnderTest() {
		return r.client.(ServerOperator)
	}
	return r.client.(*gsclient.Client)
}

// SetServerOperator set operation for server objects.
func (r *Runtime) SetServerOperator(op ServerOperator) {
	if !UnderTest() {
		panic("unexpected use")
	}
	r.client = op
}

// NetworkOperator return operations for network objects.
func (r *Runtime) NetworkOperator() NetworkOperator {
	if UnderTest() {
		return r.client.(NetworkOperator)
	}
	return r.client.(*gsclient.Client)
}

// SetNetworkOperator set operations to work on network objects.
func (r *Runtime) SetNetworkOperator(op NetworkOperator) {
	if !UnderTest() {
		panic("unexpected use")
	}
	r.client = op
}

// NewRuntime creates a new runtime for a given account. Usually there should be
// only one runtime instance in the program.
func NewRuntime(account string) (*Runtime, error) {
	client, err := newClient(account)
	if err != nil {
		return nil, err
	}

	rt := &Runtime{
		account: account,
		client:  client,
	}
	return rt, nil
}

// NewTestRuntime creates a pretty useless runtime instance. Except maybe if
// used for testing.
func NewTestRuntime() (*Runtime, error) {
	rt := &Runtime{
		account: "test",
		client:  nil,
	}
	return rt, nil
}

// UnderTest returns true when gscloud is running within 'Go test'.
func UnderTest() bool {
	return strings.HasSuffix(os.Args[0], ".test")
}

// CachePath returns the local cache path of the current user.
func CachePath() string {
	return configdir.LocalCache("gscloud")
}

func newClient(account string) (*gsclient.Client, error) {
	var ac AccountEntry
	var accountInConfig = false

	conf := &Config{}
	err := viper.Unmarshal(conf)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	for _, a := range conf.Accounts {
		if account == a.Name {
			ac = a
			accountInConfig = true
			break
		}
	}

	if len(conf.Accounts) > 0 && !accountInConfig {
		return nil, fmt.Errorf("account '%s' does not exist", account)
	}

	config := gsclient.NewConfiguration(
		ac.URL,
		ac.UserID,
		ac.Token,
		false,
		true,
		500,
		0, // no retries
	)
	return gsclient.NewClient(config), nil
}
