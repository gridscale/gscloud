package runtime

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/gridscale/gsclient-go/v3"
	"github.com/kirsle/configdir"
)

// Runtime holds all run-time infos.
type Runtime struct {
	accountName string
	client      interface{}
	config      Config
}

// KubernetesOperator amalgamates operations for Kubernetes PaaS.
type KubernetesOperator interface {
	RenewK8sCredentials(ctx context.Context, id string) error
	GetPaaSService(ctx context.Context, id string) (gsclient.PaaSService, error)
}

// PaaSOperator return an operation to Get a PaaS.
func (r *Runtime) PaaSOperator() gsclient.PaaSOperator {
	if UnderTest() {
		return r.client.(gsclient.PaaSOperator)
	}
	return r.client.(*gsclient.Client)
}

// SetPaaSOperator set operation to Create PaaS.
func (r *Runtime) SetPaaSOperator(op gsclient.PaaSOperator) {
	if !UnderTest() {
		panic("unexpected use")
	}
	r.client = op
}

// Account is the current selected account.
func (r *Runtime) Account() string {
	return r.accountName
}

// Client provides access to the API client.
func (r *Runtime) Client() *gsclient.Client {
	return r.client.(*gsclient.Client)
}

// Config allows access to configuration.
func (r *Runtime) Config() *Config {
	return &r.config
}

// ServerIPRelationOperator return an operation to remove a storage.
func (r *Runtime) ServerIPRelationOperator() gsclient.ServerIPRelationOperator {
	if UnderTest() {
		return r.client.(gsclient.ServerIPRelationOperator)
	}
	return r.client.(*gsclient.Client)
}

// SetServerIPRelationOperator set operation to delete storages.
func (r *Runtime) SetServerIPRelationOperator(op gsclient.ServerIPRelationOperator) {
	if !UnderTest() {
		panic("unexpected use")
	}
	r.client = op
}

// StorageOperator return an operation to remove a storage.
func (r *Runtime) StorageOperator() gsclient.StorageOperator {
	if UnderTest() {
		return r.client.(gsclient.StorageOperator)
	}
	return r.client.(*gsclient.Client)
}

// SetStorageOperator set operation to delete storages.
func (r *Runtime) SetStorageOperator(op gsclient.StorageOperator) {
	if !UnderTest() {
		panic("unexpected use")
	}
	r.client = op
}

// TemplateOperator return an operation to remove a storage.
func (r *Runtime) TemplateOperator() gsclient.TemplateOperator {
	if UnderTest() {
		return r.client.(gsclient.TemplateOperator)
	}
	return r.client.(*gsclient.Client)
}

// SetTemplateOperator set operation to delete storages.
func (r *Runtime) SetTemplateOperator(op gsclient.TemplateOperator) {
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
func (r *Runtime) SSHKeyOperator() gsclient.SSHKeyOperator {
	if UnderTest() {
		return r.client.(gsclient.SSHKeyOperator)
	}
	return r.client.(*gsclient.Client)
}

// SetSSHKeyOperator set operation to manipulate SSH keys.
func (r *Runtime) SetSSHKeyOperator(op gsclient.SSHKeyOperator) {
	if !UnderTest() {
		panic("unexpected use")
	}
	r.client = op
}

// ServerOperator return operation for server objects.
func (r *Runtime) ServerOperator() gsclient.ServerOperator {
	if UnderTest() {
		return r.client.(gsclient.ServerOperator)
	}
	return r.client.(*gsclient.Client)
}

// SetServerOperator set operation for server objects.
func (r *Runtime) SetServerOperator(op gsclient.ServerOperator) {
	if !UnderTest() {
		panic("unexpected use")
	}
	r.client = op
}

// ISOImageOperator return operation for server objects.
func (r *Runtime) ISOImageOperator() gsclient.ISOImageOperator {
	if UnderTest() {
		return r.client.(gsclient.ISOImageOperator)
	}
	return r.client.(*gsclient.Client)
}

// SetISOImageOperator set operation for ISO image objects.
func (r *Runtime) SetISOImageOperator(op gsclient.ISOImageOperator) {
	if !UnderTest() {
		panic("unexpected use")
	}
	r.client = op
}

// NetworkOperator return operations for network objects.
func (r *Runtime) NetworkOperator() gsclient.NetworkOperator {
	if UnderTest() {
		return r.client.(gsclient.NetworkOperator)
	}
	return r.client.(*gsclient.Client)
}

// SetNetworkOperator set operations to work on network objects.
func (r *Runtime) SetNetworkOperator(op gsclient.NetworkOperator) {
	if !UnderTest() {
		panic("unexpected use")
	}
	r.client = op
}

// IPOperator return operations to manipulate IP addresses.
func (r *Runtime) IPOperator() gsclient.IPOperator {
	if UnderTest() {
		return r.client.(gsclient.IPOperator)
	}
	return r.client.(*gsclient.Client)
}

// SetIPOperator set operations to manipulate IP addresses.
func (r *Runtime) SetIPOperator(op gsclient.IPOperator) {
	if !UnderTest() {
		panic("unexpected use")
	}
	r.client = op
}

// ServerStorageRelationOperator return an operation to associate server objects with storages.
func (r *Runtime) ServerStorageRelationOperator() gsclient.ServerStorageRelationOperator {
	if UnderTest() {
		return r.client.(gsclient.ServerStorageRelationOperator)
	}
	return r.client.(*gsclient.Client)
}

// SetServerStorageRelationOperator set operation to delete storages.
func (r *Runtime) SetServerStorageRelationOperator(op gsclient.ServerStorageRelationOperator) {
	if !UnderTest() {
		panic("unexpected use")
	}
	r.client = op
}

// NewRuntime creates a new runtime for a given account. Usually there should be
// only one runtime instance in the program.
func NewRuntime(conf Config, accountName string) (*Runtime, error) {
	var ac AccountEntry
	var accountIndex = -1

	for i, a := range conf.Accounts {
		if accountName == a.Name {
			ac = a
			accountIndex = i
			break
		}
	}

	if len(conf.Accounts) > 0 && accountIndex == -1 {
		if !CommandWithoutConfig(os.Args) {
			return nil, fmt.Errorf("account '%s' does not exist", accountName)
		}
	}

	if len(conf.Accounts) >= accountIndex && accountIndex != -1 {
		ac = LoadEnvVariables(ac)
		conf.Accounts[accountIndex] = ac
	}

	client := newClient(ac)
	rt := &Runtime{
		accountName: ac.Name,
		client:      client,
		config:      conf,
	}
	return rt, nil
}

// LoadEnvVariables loads UserId, Token and URL from their respective environment variable
func LoadEnvVariables(defaultAc AccountEntry) AccountEntry {
	userIDEnv, userIDEnvExists := os.LookupEnv("GRIDSCALE_UUID")
	if userIDEnvExists {
		defaultAc.UserID = userIDEnv
	}

	tokenEnv, tokenEnvExists := os.LookupEnv("GRIDSCALE_TOKEN")
	if tokenEnvExists {
		defaultAc.Token = tokenEnv
	}

	apiURLEnv, apiURLEnvExists := os.LookupEnv("GRIDSCALE_URL")
	if apiURLEnvExists {
		defaultAc.URL = apiURLEnv
	}
	return defaultAc
}

// NewTestRuntime creates a pretty useless runtime instance. Except maybe if
// used for testing.
func NewTestRuntime() (*Runtime, error) {
	rt := &Runtime{
		accountName: "test",
		client:      nil,
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

// newClient creates new gsclient from a given instance of AccountEntry
func newClient(account AccountEntry) *gsclient.Client {
	if account.URL == "" {
		config := gsclient.DefaultConfiguration(account.UserID, account.Token)
		return gsclient.NewClient(config)
	}
	config := gsclient.NewConfiguration(
		account.URL,
		account.UserID,
		account.Token,
		false,
		true,
		500,
		0, // no retries
	)
	return gsclient.NewClient(config)
}
