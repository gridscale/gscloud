package gsclient

import (
	"context"
	"errors"
	"net/http"
	"path"
)

//ServerList JSON struct of a list of servers
type ServerList struct {
	//Array of servers
	List map[string]ServerProperties `json:"servers"`
}

//DeletedServerList JSON struct of a list of deleted servers
type DeletedServerList struct {
	//Array of deleted servers
	List map[string]ServerProperties `json:"deleted_servers"`
}

//Server JSON struct of a single server
type Server struct {
	//Properties of a server
	Properties ServerProperties `json:"server"`
}

//ServerProperties JSON struct of properties of a server
type ServerProperties struct {
	//The UUID of an object is always unique, and refers to a specific object.
	ObjectUUID string `json:"object_uuid"`

	//The human-readable name of the object. It supports the full UTF-8 charset, with a maximum of 64 characters.
	Name string `json:"name"`

	//Indicates the amount of memory in GB.
	Memory int `json:"memory"`

	//Number of server cores.
	Cores int `json:"cores"`

	//Specifies the hardware settings for the virtual machine.
	HardwareProfile string `json:"hardware_profile"`

	//Status indicates the status of the object. it could be in-provisioning or active
	Status string `json:"status"`

	//Helps to identify which datacenter an object belongs to.
	LocationUUID string `json:"location_uuid"`

	//The power status of the server.
	Power bool `json:"power"`

	//The price for the current period since the last bill.
	CurrentPrice float64 `json:"current_price"`

	//Which Availability-Zone the Server is placed.
	AvailabilityZone string `json:"availability_zone"`

	//If the server should be auto-started in case of a failure (default=true).
	AutoRecovery bool `json:"auto_recovery"`

	//Legacy-Hardware emulation instead of virtio hardware.
	//If enabled, hotplugging cores, memory, storage, network, etc. will not work,
	//but the server will most likely run every x86 compatible operating system.
	//This mode comes with a performance penalty, as emulated hardware does not benefit from the virtio driver infrastructure.
	Legacy bool `json:"legacy"`

	//The token used by the panel to open the websocket VNC connection to the server console.
	ConsoleToken string `json:"console_token"`

	//Total minutes of memory used.
	UsageInMinutesMemory int `json:"usage_in_minutes_memory"`

	//Total minutes of cores used.
	UsageInMinutesCores int `json:"usage_in_minutes_cores"`

	//List of labels.
	Labels []string `json:"labels"`

	//The information about other object which are related to this server. the object could be ip, storage, network, and isoimage
	Relations ServerRelations `json:"relations"`

	//Defines the date and time the object was initially created.
	CreateTime GSTime `json:"create_time"`

	//Defines the date and time of the last object change.
	ChangeTime GSTime `json:"change_time"`
}

//ServerRelations JSON struct of a list of server relations
type ServerRelations struct {
	//Array of object (ServerIsoImageRelationProperties)
	IsoImages []ServerIsoImageRelationProperties `json:"isoimages"`

	//Array of object (ServerNetworkRelationProperties)
	Networks []ServerNetworkRelationProperties `json:"networks"`

	//Array of object (ServerIPRelationProperties)
	PublicIPs []ServerIPRelationProperties `json:"public_ips"`

	//Array of object (ServerStorageRelationProperties)
	Storages []ServerStorageRelationProperties `json:"storages"`
}

//ServerCreateRequest JSON struct of a request for creating a server
type ServerCreateRequest struct {
	//The human-readable name of the object. It supports the full UTF-8 charset, with a maximum of 64 characters.
	Name string `json:"name"`

	//The amount of server memory in GB.
	Memory int `json:"memory"`

	//The number of server cores.
	Cores int `json:"cores"`

	//Specifies the hardware settings for the virtual machine.
	//Allowed values: nil, DefaultServerHardware, NestedServerHardware, LegacyServerHardware, CiscoCSRServerHardware,
	//SophosUTMServerHardware, F5BigipServerHardware, Q35ServerHardware, Q35NestedServerHardware.
	//HardwareProfile = nil => server hardware is normal type
	HardwareProfile *serverHardwareProfile `json:"hardware_profile,omitempty"`

	//Defines which Availability-Zone the Server is placed. Can be empty
	AvailablityZone string `json:"availability_zone,omitempty"`

	//List of labels. Can be empty.
	Labels []string `json:"labels,omitempty"`

	//Status indicates the status of the object. Can be empty
	Status string `json:"status,omitempty"`

	//If the server should be auto-started in case of a failure (default=true when AutoRecovery=nil).
	AutoRecovery *bool `json:"auto_recovery,omitempty"`

	//The information about other object which are related to this server. the object could be ip, storage, network, and isoimage.
	//*Caution: This field will be deprecated.
	Relations *ServerCreateRequestRelations `json:"relations,omitempty"`
}

//ServerCreateRequestRelations JSOn struct of a list of a server's relations
type ServerCreateRequestRelations struct {
	//Array of objects (ServerCreateRequestIsoimage)
	IsoImages []ServerCreateRequestIsoimage `json:"isoimages"`

	//Array of objects (ServerCreateRequestNetwork)
	Networks []ServerCreateRequestNetwork `json:"networks"`

	//Array of objects (ServerCreateRequestIP)
	PublicIPs []ServerCreateRequestIP `json:"public_ips"`

	//Array of objects (ServerCreateRequestStorage)
	Storages []ServerCreateRequestStorage `json:"storages"`
}

//ServerCreateResponse JSON struct of a response for creating a server
type ServerCreateResponse struct {
	//UUID of object being created. Same as ServerUUID.
	ObjectUUID string `json:"object_uuid"`

	//UUID of the request
	RequestUUID string `json:"request_uuid"`

	//UUID of server being created. Same as ObjectUUID.
	ServerUUID string `json:"server_uuid"`

	//UUIDs of attached networks
	NetworkUUIDs []string `json:"network_uuids"`

	//UUIDs of attached storages
	StorageUUIDs []string `json:"storage_uuids"`

	//UUIDs of attached IP addresses
	IPaddrUUIDs []string `json:"ipaddr_uuids"`
}

//ServerPowerUpdateRequest JSON struct of a request for updating server's power state
type ServerPowerUpdateRequest struct {
	//Power=true => server is on
	//Power=false => server if off
	Power bool `json:"power"`
}

//ServerCreateRequestStorage JSON struct of a relation between a server and a storage
type ServerCreateRequestStorage struct {
	//UUID of the storage being attached to the server
	StorageUUID string `json:"storage_uuid"`

	//Is the storage a boot device?
	BootDevice bool `json:"bootdevice,omitempty"`
}

//ServerCreateRequestNetwork JSON struct of a relation between a server and a network
type ServerCreateRequestNetwork struct {
	//UUID of the networks being attached to the server
	NetworkUUID string `json:"network_uuid"`

	//Is the network a boot device?
	BootDevice bool `json:"bootdevice,omitempty"`
}

//ServerCreateRequestIP JSON struct of a relation between a server and an IP address
type ServerCreateRequestIP struct {
	//UUID of the IP address being attached to the server
	IPaddrUUID string `json:"ipaddr_uuid"`
}

//ServerCreateRequestIsoimage JSON struct of a relation between a server and an ISO-Image
type ServerCreateRequestIsoimage struct {
	//UUID of the ISO-image being attached to the server
	IsoimageUUID string `json:"isoimage_uuid"`
}

//ServerUpdateRequest JSON of a request for updating a server
type ServerUpdateRequest struct {
	//The human-readable name of the object. It supports the full UTF-8 charset, with a maximum of 64 characters.
	//Leave it if you do not want to update the name.
	Name string `json:"name,omitempty"`

	//Defines which Availability-Zone the Server is placed. Leave it if you do not want to update the zone.
	AvailablityZone string `json:"availability_zone,omitempty"`

	//The amount of server memory in GB. Leave it if you do not want to update the memory
	Memory int `json:"memory,omitempty"`

	//The number of server cores. Leave it if you do not want to update the number of the cpu cores.
	Cores int `json:"cores,omitempty"`

	//List of labels. Leave it if you do not want to update the list of labels
	Labels *[]string `json:"labels,omitempty"`

	//If the server should be auto-started in case of a failure (default=true).
	//Leave it if you do not want to update this feature of the server.
	AutoRecovery *bool `json:"auto_recovery,omitempty"`
}

//ServerMetricList JSON struct of a list of a server's metrics
type ServerMetricList struct {
	//Array of a server's metrics
	List []ServerMetricProperties `json:"server_metrics"`
}

//ServerMetric JSON struct of a single metric of a server
type ServerMetric struct {
	//Properties of a server metric
	Properties ServerMetricProperties `json:"server_metric"`
}

//ServerMetricProperties JSON struct
type ServerMetricProperties struct {
	//Defines the begin of the time range.
	BeginTime GSTime `json:"begin_time"`

	//Defines the end of the time range.
	EndTime GSTime `json:"end_time"`

	//The UUID of an object is always unique, and refers to a specific object.
	PaaSServiceUUID string `json:"paas_service_uuid"`

	//Core usage
	CoreUsage struct {
		//Value
		Value float64 `json:"value"`

		//Unit of value
		Unit string `json:"unit"`
	} `json:"core_usage"`

	//Storage usage
	StorageSize struct {
		//Value
		Value float64 `json:"value"`

		//Unit of value
		Unit string `json:"unit"`
	} `json:"storage_size"`
}

//All available server's hardware types
var (
	DefaultServerHardware   = &serverHardwareProfile{"default"}
	NestedServerHardware    = &serverHardwareProfile{"nested"}
	LegacyServerHardware    = &serverHardwareProfile{"legacy"}
	CiscoCSRServerHardware  = &serverHardwareProfile{"cisco_csr"}
	SophosUTMServerHardware = &serverHardwareProfile{"sophos_utm"}
	F5BigipServerHardware   = &serverHardwareProfile{"f5_bigip"}
	Q35ServerHardware       = &serverHardwareProfile{"q35"}
	Q35NestedServerHardware = &serverHardwareProfile{"q35_nested"}
)

//GetServer gets a specific server based on given list
//
//See: https://gridscale.io/en//api-documentation/index.html#operation/getServer
func (c *Client) GetServer(ctx context.Context, id string) (Server, error) {
	if !isValidUUID(id) {
		return Server{}, errors.New("'id' is invalid")
	}
	r := request{
		uri:                 path.Join(apiServerBase, id),
		method:              http.MethodGet,
		skipCheckingRequest: true,
	}
	var response Server
	err := r.execute(ctx, *c, &response)
	return response, err
}

//GetServerList gets a list of available servers
//
//See: https://gridscale.io/en//api-documentation/index.html#operation/getServers
func (c *Client) GetServerList(ctx context.Context) ([]Server, error) {
	r := request{
		uri:                 apiServerBase,
		method:              http.MethodGet,
		skipCheckingRequest: true,
	}
	var response ServerList
	var servers []Server
	err := r.execute(ctx, *c, &response)
	for _, properties := range response.List {
		servers = append(servers, Server{
			Properties: properties,
		})
	}
	return servers, err
}

//CreateServer create a server
//
//NOTE: Allowed values of `HardwareProfile`: nil, DefaultServerHardware, NestedServerHardware, LegacyServerHardware,
//CiscoCSRServerHardware, SophosUTMServerHardware, F5BigipServerHardware, Q35ServerHardware, Q35NestedServerHardware.
//
//See: https://gridscale.io/en//api-documentation/index.html#operation/createServer
func (c *Client) CreateServer(ctx context.Context, body ServerCreateRequest) (ServerCreateResponse, error) {
	//check if these slices are nil
	//make them be empty slice instead of nil
	//so that JSON structure will be valid
	if body.Relations != nil && body.Relations.PublicIPs == nil {
		body.Relations.PublicIPs = make([]ServerCreateRequestIP, 0)
	}
	if body.Relations != nil && body.Relations.Networks == nil {
		body.Relations.Networks = make([]ServerCreateRequestNetwork, 0)
	}
	if body.Relations != nil && body.Relations.IsoImages == nil {
		body.Relations.IsoImages = make([]ServerCreateRequestIsoimage, 0)
	}
	if body.Relations != nil && body.Relations.Storages == nil {
		body.Relations.Storages = make([]ServerCreateRequestStorage, 0)
	}
	r := request{
		uri:    apiServerBase,
		method: http.MethodPost,
		body:   body,
	}
	var response ServerCreateResponse
	err := r.execute(ctx, *c, &response)
	//this fixed the endpoint's bug temporarily when creating server with/without
	//'relations' field
	if response.ServerUUID == "" && response.ObjectUUID != "" {
		response.ServerUUID = response.ObjectUUID
	} else if response.ObjectUUID == "" && response.ServerUUID != "" {
		response.ObjectUUID = response.ServerUUID
	}
	return response, err
}

//DeleteServer deletes a specific server
//
//See: https://gridscale.io/en//api-documentation/index.html#operation/deleteServer
func (c *Client) DeleteServer(ctx context.Context, id string) error {
	if !isValidUUID(id) {
		return errors.New("'id' is invalid")
	}
	r := request{
		uri:    path.Join(apiServerBase, id),
		method: http.MethodDelete,
	}
	return r.execute(ctx, *c, nil)
}

//UpdateServer updates a specific server
//
//See: https://gridscale.io/en//api-documentation/index.html#operation/updateServer
func (c *Client) UpdateServer(ctx context.Context, id string, body ServerUpdateRequest) error {
	if !isValidUUID(id) {
		return errors.New("'id' is invalid")
	}
	r := request{
		uri:    path.Join(apiServerBase, id),
		method: http.MethodPatch,
		body:   body,
	}
	return r.execute(ctx, *c, nil)
}

//GetServerEventList gets a list of a specific server's events
//
//See: https://gridscale.io/en//api-documentation/index.html#operation/getServerEvents
func (c *Client) GetServerEventList(ctx context.Context, id string) ([]Event, error) {
	if !isValidUUID(id) {
		return nil, errors.New("'id' is invalid")
	}
	r := request{
		uri:                 path.Join(apiServerBase, id, "events"),
		method:              http.MethodGet,
		skipCheckingRequest: true,
	}
	var response EventList
	var serverEvents []Event
	err := r.execute(ctx, *c, &response)
	for _, properties := range response.List {
		serverEvents = append(serverEvents, Event{Properties: properties})
	}
	return serverEvents, err
}

//GetServerMetricList gets a list of a specific server's metrics
//
//See: https://gridscale.io/en//api-documentation/index.html#operation/getServerMetrics
func (c *Client) GetServerMetricList(ctx context.Context, id string) ([]ServerMetric, error) {
	if !isValidUUID(id) {
		return nil, errors.New("'id' is invalid")
	}
	r := request{
		uri:                 path.Join(apiServerBase, id, "metrics"),
		method:              http.MethodGet,
		skipCheckingRequest: true,
	}
	var response ServerMetricList
	var serverMetrics []ServerMetric
	err := r.execute(ctx, *c, &response)
	for _, properties := range response.List {
		serverMetrics = append(serverMetrics, ServerMetric{Properties: properties})
	}
	return serverMetrics, err
}

//IsServerOn returns true if the server's power is on, otherwise returns false
func (c *Client) IsServerOn(ctx context.Context, id string) (bool, error) {
	server, err := c.GetServer(ctx, id)
	if err != nil {
		return false, err
	}
	return server.Properties.Power, nil
}

//setServerPowerState turn on/off a specific server.
//turnOn=true to turn on, turnOn=false to turn off
func (c *Client) setServerPowerState(ctx context.Context, id string, powerState bool) error {
	isOn, err := c.IsServerOn(ctx, id)
	if err != nil {
		return err
	}
	if isOn == powerState {
		return nil
	}
	r := request{
		uri:    path.Join(apiServerBase, id, "power"),
		method: http.MethodPatch,
		body: ServerPowerUpdateRequest{
			Power: powerState,
		},
	}
	err = r.execute(ctx, *c, nil)
	if err != nil {
		return err
	}
	if c.Synchronous() {
		return c.waitForServerPowerStatus(ctx, id, powerState)
	}
	return nil
}

//StartServer starts a server
func (c *Client) StartServer(ctx context.Context, id string) error {
	return c.setServerPowerState(ctx, id, true)
}

//StopServer stops a server
func (c *Client) StopServer(ctx context.Context, id string) error {
	return c.setServerPowerState(ctx, id, false)
}

//ShutdownServer shutdowns a specific server
func (c *Client) ShutdownServer(ctx context.Context, id string) error {
	//Make sure the server exists and that it isn't already in the state we need it to be
	server, err := c.GetServer(ctx, id)
	if err != nil {
		return err
	}
	if !server.Properties.Power {
		return nil
	}
	r := request{
		uri:    path.Join(apiServerBase, id, "shutdown"),
		method: http.MethodPatch,
		body:   map[string]string{},
	}

	err = r.execute(ctx, *c, nil)
	if err != nil {
		return err
	}

	if c.Synchronous() {
		//If we get an error, which includes a timeout, power off the server instead
		err = c.waitForServerPowerStatus(ctx, id, false)
		if err != nil {
			return err
		}
	}
	return nil
}

//GetServersByLocation gets a list of servers by location
//
//See: https://gridscale.io/en//api-documentation/index.html#operation/getLocationServers
func (c *Client) GetServersByLocation(ctx context.Context, id string) ([]Server, error) {
	if !isValidUUID(id) {
		return nil, errors.New("'id' is invalid")
	}
	r := request{
		uri:                 path.Join(apiLocationBase, id, "servers"),
		method:              http.MethodGet,
		skipCheckingRequest: true,
	}
	var response ServerList
	var servers []Server
	err := r.execute(ctx, *c, &response)
	for _, properties := range response.List {
		servers = append(servers, Server{Properties: properties})
	}
	return servers, err
}

//GetDeletedServers gets a list of deleted servers
//
//See: https://gridscale.io/en//api-documentation/index.html#operation/getDeletedServers
func (c *Client) GetDeletedServers(ctx context.Context) ([]Server, error) {
	r := request{
		uri:                 path.Join(apiDeletedBase, "servers"),
		method:              http.MethodGet,
		skipCheckingRequest: true,
	}
	var response DeletedServerList
	var servers []Server
	err := r.execute(ctx, *c, &response)
	for _, properties := range response.List {
		servers = append(servers, Server{Properties: properties})
	}
	return servers, err
}

//waitForServerPowerStatus  allows to wait for a server changing its power status.
func (c *Client) waitForServerPowerStatus(ctx context.Context, id string, status bool) error {
	return retryWithContext(ctx, func() (bool, error) {
		server, err := c.GetServer(ctx, id)
		return server.Properties.Power != status, err
	}, c.DelayInterval())
}
