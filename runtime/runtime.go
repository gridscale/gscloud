package runtime

import (
	"context"
	"fmt"
	"os"

	"github.com/gridscale/gsclient-go/v3"
	"github.com/gridscale/gscloud/utils"
	"github.com/kirsle/configdir"
)

// Runtime holds all run-time infos.
type Runtime struct {
	ProjectName string
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
	if utils.UnderTest() {
		return r.client.(gsclient.PaaSOperator)
	}
	return r.client.(*gsclient.Client)
}

// SetPaaSOperator set operation to Create PaaS.
func (r *Runtime) SetPaaSOperator(op gsclient.PaaSOperator) {
	if !utils.UnderTest() {
		panic("unexpected use")
	}
	r.client = op
}

// Project is the current selected project.
func (r *Runtime) Project() string {
	return r.ProjectName
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
	if utils.UnderTest() {
		return r.client.(gsclient.ServerIPRelationOperator)
	}
	return r.client.(*gsclient.Client)
}

// SetServerIPRelationOperator set operation to delete storages.
func (r *Runtime) SetServerIPRelationOperator(op gsclient.ServerIPRelationOperator) {
	if !utils.UnderTest() {
		panic("unexpected use")
	}
	r.client = op
}

// StorageOperator return an operation to remove a storage.
func (r *Runtime) StorageOperator() gsclient.StorageOperator {
	if utils.UnderTest() {
		return r.client.(gsclient.StorageOperator)
	}
	return r.client.(*gsclient.Client)
}

// SetStorageOperator set operation to delete storages.
func (r *Runtime) SetStorageOperator(op gsclient.StorageOperator) {
	if !utils.UnderTest() {
		panic("unexpected use")
	}
	r.client = op
}

// TemplateOperator return an operation to remove a storage.
func (r *Runtime) TemplateOperator() gsclient.TemplateOperator {
	if utils.UnderTest() {
		return r.client.(gsclient.TemplateOperator)
	}
	return r.client.(*gsclient.Client)
}

// SetTemplateOperator set operation to delete storages.
func (r *Runtime) SetTemplateOperator(op gsclient.TemplateOperator) {
	if !utils.UnderTest() {
		panic("unexpected use")
	}
	r.client = op
}

// KubernetesOperator return operation relating to Kubernetes managed services.
func (r *Runtime) KubernetesOperator() KubernetesOperator {
	if utils.UnderTest() {
		return r.client.(KubernetesOperator)
	}
	return r.client.(*gsclient.Client)
}

// SetKubernetesOperator set Kubernetes PaaS operation.
func (r *Runtime) SetKubernetesOperator(op KubernetesOperator) {
	if !utils.UnderTest() {
		panic("unexpected use")
	}
	r.client = op
}

// SSHKeyOperator return operation to manipulate SSH keys.
func (r *Runtime) SSHKeyOperator() gsclient.SSHKeyOperator {
	if utils.UnderTest() {
		return r.client.(gsclient.SSHKeyOperator)
	}
	return r.client.(*gsclient.Client)
}

// SetSSHKeyOperator set operation to manipulate SSH keys.
func (r *Runtime) SetSSHKeyOperator(op gsclient.SSHKeyOperator) {
	if !utils.UnderTest() {
		panic("unexpected use")
	}
	r.client = op
}

// ServerOperator return operation for server objects.
func (r *Runtime) ServerOperator() gsclient.ServerOperator {
	if utils.UnderTest() {
		return r.client.(gsclient.ServerOperator)
	}
	return r.client.(*gsclient.Client)
}

// SetServerOperator set operation for server objects.
func (r *Runtime) SetServerOperator(op gsclient.ServerOperator) {
	if !utils.UnderTest() {
		panic("unexpected use")
	}
	r.client = op
}

// ISOImageOperator return operation for server objects.
func (r *Runtime) ISOImageOperator() gsclient.ISOImageOperator {
	if utils.UnderTest() {
		return r.client.(gsclient.ISOImageOperator)
	}
	return r.client.(*gsclient.Client)
}

// SetISOImageOperator set operation for ISO image objects.
func (r *Runtime) SetISOImageOperator(op gsclient.ISOImageOperator) {
	if !utils.UnderTest() {
		panic("unexpected use")
	}
	r.client = op
}

// NetworkOperator return operations for network objects.
func (r *Runtime) NetworkOperator() gsclient.NetworkOperator {
	if utils.UnderTest() {
		return r.client.(gsclient.NetworkOperator)
	}
	return r.client.(*gsclient.Client)
}

// SetNetworkOperator set operations to work on network objects.
func (r *Runtime) SetNetworkOperator(op gsclient.NetworkOperator) {
	if !utils.UnderTest() {
		panic("unexpected use")
	}
	r.client = op
}

// IPOperator return operations to manipulate IP addresses.
func (r *Runtime) IPOperator() gsclient.IPOperator {
	if utils.UnderTest() {
		return r.client.(gsclient.IPOperator)
	}
	return r.client.(*gsclient.Client)
}

// SetIPOperator set operations to manipulate IP addresses.
func (r *Runtime) SetIPOperator(op gsclient.IPOperator) {
	if !utils.UnderTest() {
		panic("unexpected use")
	}
	r.client = op
}

// ServerStorageRelationOperator return an operation to associate server objects with storages.
func (r *Runtime) ServerStorageRelationOperator() gsclient.ServerStorageRelationOperator {
	if utils.UnderTest() {
		return r.client.(gsclient.ServerStorageRelationOperator)
	}
	return r.client.(*gsclient.Client)
}

// SetServerStorageRelationOperator set operation to delete storages.
func (r *Runtime) SetServerStorageRelationOperator(op gsclient.ServerStorageRelationOperator) {
	if !utils.UnderTest() {
		panic("unexpected use")
	}
	r.client = op
}

// NewRuntime creates a new runtime for a given project. Usually there should be
// only one runtime instance in the program.
func NewRuntime(conf Config, projectName string, commandWithoutConfig bool) (*Runtime, error) {
	var ac ProjectEntry
	var projectIndex = -1

	for i, a := range conf.Projects {
		if projectName == a.Name {
			ac = a
			projectIndex = i
			break
		}
	}

	if projectIndex == -1 {
		if len(conf.Projects) > 0 && !commandWithoutConfig {
			return nil, fmt.Errorf("project '%s' does not exist", projectName)
		}
	} else {
		ac = LoadEnvVariables(ac)
		conf.Projects[projectIndex] = ac
	}

	client := newClient(ac)
	rt := &Runtime{
		ProjectName: ac.Name,
		client:      client,
		config:      conf,
	}
	return rt, nil
}

// LoadEnvVariables loads UserId, Token and URL from their respective environment variable
func LoadEnvVariables(defaultAc ProjectEntry) ProjectEntry {
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
		ProjectName: "test",
		client:      nil,
	}
	return rt, nil
}

// CachePath returns the local cache path of the current user.
func CachePath() string {
	return configdir.LocalCache("gscloud")
}

// newClient creates new gsclient from a given instance of ProjectEntry
func newClient(project ProjectEntry) *gsclient.Client {
	if project.URL == "" {
		config := gsclient.DefaultConfiguration(project.UserID, project.Token)
		return gsclient.NewClient(config)
	}
	config := gsclient.NewConfiguration(
		project.URL,
		project.UserID,
		project.Token,
		false,
		true,
		500,
		0, // no retries
	)
	return gsclient.NewClient(config)
}
