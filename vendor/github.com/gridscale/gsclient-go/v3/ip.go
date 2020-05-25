package gsclient

import (
	"context"
	"errors"
	"net/http"
	"path"
)

//IPList is JSON struct of a list of IPs
type IPList struct {
	//Array of IP addresses
	List map[string]IPProperties `json:"ips"`
}

//DeletedIPList is JSON struct of a list of deleted IPs
type DeletedIPList struct {
	//Array of deleted IP addresses
	List map[string]IPProperties `json:"deleted_ips"`
}

//IP is JSON struct if a single IP
type IP struct {
	//Properties of an IP address
	Properties IPProperties `json:"ip"`
}

//IPProperties is JSON struct of an IP's properties
type IPProperties struct {
	//The human-readable name of the object. It supports the full UTF-8 charset, with a maximum of 64 characters.
	Name string `json:"name"`

	//The human-readable name of the location. It supports the full UTF-8 charset, with a maximum of 64 characters.
	LocationCountry string `json:"location_country"`

	//Helps to identify which datacenter an object belongs to.
	LocationUUID string `json:"location_uuid"`

	//The UUID of an object is always unique, and refers to a specific object.
	ObjectUUID string `json:"object_uuid"`

	//Defines the reverse DNS entry for the IP Address (PTR Resource Record).
	ReverseDNS string `json:"reverse_dns"`

	//Enum:4 6. The IP Address family (v4 or v6)
	Family int `json:"family"`

	//Status indicates the status of the object.
	Status string `json:"status"`

	//Defines the date and time the object was initially created.
	CreateTime GSTime `json:"create_time"`

	//Sets failover mode for this IP. If true, then this IP is no longer available for DHCP and can no longer be related to any server.
	Failover bool `json:"failover"`

	//Defines the date and time of the last object change.
	ChangeTime GSTime `json:"change_time"`

	//Uses IATA airport code, which works as a location identifier.
	LocationIata string `json:"location_iata"`

	//The human-readable name of the location. It supports the full UTF-8 charset, with a maximum of 64 characters.
	LocationName string `json:"location_name"`

	//The IP prefix.
	Prefix string `json:"prefix"`

	//Defines the IP Address (v4 or v6).
	IP string `json:"ip"`

	//Defines if the object is administratively blocked. If true, it can not be deleted by the user.
	DeleteBlock bool `json:"delete_block"`

	//Total minutes the object has been running.
	UsagesInMinutes float64 `json:"usage_in_minutes"`

	//The price for the current period since the last bill.
	CurrentPrice float64 `json:"current_price"`

	//List of labels.
	Labels []string `json:"labels"`

	//The information about other object which are related to this IP. the object could be servers and/or loadbalancer.
	Relations IPRelations `json:"relations"`
}

//IPRelations is JSON struct of a list of an IP's relations
type IPRelations struct {
	//Array of object (IPLoadbalancer)
	Loadbalancers []IPLoadbalancer `json:"loadbalancers"`

	//Array of object (IPServer)
	Servers []IPServer `json:"servers"`

	//Array of object (ServerIPRelationProperties)
	PublicIPs []ServerIPRelationProperties `json:"public_ips"`

	//Array of object (ServerStorageRelationProperties)
	Storages []ServerStorageRelationProperties `json:"storages"`
}

//IPLoadbalancer is JSON struct of the relation between an IP and a Load Balancer
type IPLoadbalancer struct {
	//Defines the date and time the object was initially created.
	CreateTime GSTime `json:"create_time"`

	//The human-readable name of the object. It supports the full UTF-8 charset, with a maximum of 64 characters.
	LoadbalancerName string `json:"loadbalancer_name"`

	//The UUID of loadbalancer.
	LoadbalancerUUID string `json:"loadbalancer_uuid"`
}

//IPServer is JSON struct of the relation between an IP and a Server
type IPServer struct {
	//Defines the date and time the object was initially created.
	CreateTime GSTime `json:"create_time"`

	//The human-readable name of the object. It supports the full UTF-8 charset, with a maximum of 64 characters.
	ServerName string `json:"server_name"`

	//The UUID of the server.
	ServerUUID string `json:"server_uuid"`
}

//IPCreateResponse is JSON struct of a response for creating an IP
type IPCreateResponse struct {
	//Request's UUID
	RequestUUID string `json:"request_uuid"`

	//UUID of the IP address being created
	ObjectUUID string `json:"object_uuid"`

	//The IP prefix.
	Prefix string `json:"prefix"`

	//The IP Address (v4 or v6).
	IP string `json:"ip"`
}

//IPCreateRequest is JSON struct of a request for creating an IP
type IPCreateRequest struct {
	//Name of an IP address being created. Can be an empty string
	Name string `json:"name,omitempty"`

	//IP address family. Can only be either `IPv4Type` or `IPv6Type`
	Family ipAddressType `json:"family"`

	//Sets failover mode for this IP. If true, then this IP is no longer available for DHCP and can no longer be related to any server.
	Failover bool `json:"failover,omitempty"`

	//Defines the reverse DNS entry for the IP Address (PTR Resource Record).
	ReverseDNS string `json:"reverse_dns,omitempty"`

	//List of labels.
	Labels []string `json:"labels,omitempty"`
}

//IPUpdateRequest is JSON struct of a request for updating an IP
type IPUpdateRequest struct {
	//New name. Leave it if you do not want to update the name
	Name string `json:"name,omitempty"`

	//Sets failover mode for this IP. If true, then this IP is no longer available for DHCP and can no longer be related to any server.
	Failover bool `json:"failover"`

	//Defines the reverse DNS entry for the IP Address (PTR Resource Record). Leave it if you do not want to update the reverse DNS.
	ReverseDNS string `json:"reverse_dns,omitempty"`

	//List of labels. Leave it if you do not want to update the labels.
	Labels *[]string `json:"labels,omitempty"`
}

//Allowed IP address versions
var (
	IPv4Type = ipAddressType{4}
	IPv6Type = ipAddressType{6}
)

//GetIP get a specific IP based on given id
//
//See: https://gridscale.io/en//api-documentation/index.html#operation/getIp
func (c *Client) GetIP(ctx context.Context, id string) (IP, error) {
	if !isValidUUID(id) {
		return IP{}, errors.New("'id' is invalid")
	}
	r := request{
		uri:                 path.Join(apiIPBase, id),
		method:              http.MethodGet,
		skipCheckingRequest: true,
	}

	var response IP
	err := r.execute(ctx, *c, &response)

	return response, err
}

//GetIPList gets a list of available IPs
//
//https://gridscale.io/en//api-documentation/index.html#operation/getIps
func (c *Client) GetIPList(ctx context.Context) ([]IP, error) {
	r := request{
		uri:                 apiIPBase,
		method:              http.MethodGet,
		skipCheckingRequest: true,
	}

	var response IPList
	var IPs []IP
	err := r.execute(ctx, *c, &response)
	for _, properties := range response.List {
		IPs = append(IPs, IP{Properties: properties})
	}

	return IPs, err
}

//CreateIP creates an IP
//
//Note: IP address family can only be either `IPv4Type` or `IPv6Type`
//
//See: https://gridscale.io/en//api-documentation/index.html#operation/createIp
func (c *Client) CreateIP(ctx context.Context, body IPCreateRequest) (IPCreateResponse, error) {
	r := request{
		uri:    apiIPBase,
		method: http.MethodPost,
		body:   body,
	}

	var response IPCreateResponse
	err := r.execute(ctx, *c, &response)
	return response, err
}

//DeleteIP deletes a specific IP based on given id
//
//See: https://gridscale.io/en//api-documentation/index.html#operation/deleteIp
func (c *Client) DeleteIP(ctx context.Context, id string) error {
	if !isValidUUID(id) {
		return errors.New("'id' is invalid")
	}
	r := request{
		uri:    path.Join(apiIPBase, id),
		method: http.MethodDelete,
	}
	return r.execute(ctx, *c, nil)
}

//UpdateIP updates a specific IP based on given id
//
//See: https://gridscale.io/en//api-documentation/index.html#operation/updateIp
func (c *Client) UpdateIP(ctx context.Context, id string, body IPUpdateRequest) error {
	if !isValidUUID(id) {
		return errors.New("'id' is invalid")
	}
	r := request{
		uri:    path.Join(apiIPBase, id),
		method: http.MethodPatch,
		body:   body,
	}
	return r.execute(ctx, *c, nil)
}

//GetIPEventList gets a list of an IP's events
//
//See: https://gridscale.io/en//api-documentation/index.html#operation/getIpEvents
func (c *Client) GetIPEventList(ctx context.Context, id string) ([]Event, error) {
	if !isValidUUID(id) {
		return nil, errors.New("'id' is invalid")
	}
	r := request{
		uri:                 path.Join(apiIPBase, id, "events"),
		method:              http.MethodGet,
		skipCheckingRequest: true,
	}
	var response EventList
	var IPEvents []Event
	err := r.execute(ctx, *c, &response)
	for _, properties := range response.List {
		IPEvents = append(IPEvents, Event{Properties: properties})
	}
	return IPEvents, err
}

//GetIPVersion gets IP's version, returns 0 if an error was encountered
func (c *Client) GetIPVersion(ctx context.Context, id string) int {
	ip, err := c.GetIP(ctx, id)
	if err != nil {
		return 0
	}
	return ip.Properties.Family
}

//GetIPsByLocation gets a list of IPs by location
//
//See: https://gridscale.io/en//api-documentation/index.html#operation/getLocationIps
func (c *Client) GetIPsByLocation(ctx context.Context, id string) ([]IP, error) {
	if !isValidUUID(id) {
		return nil, errors.New("'id' is invalid")
	}
	r := request{
		uri:                 path.Join(apiLocationBase, id, "ips"),
		method:              http.MethodGet,
		skipCheckingRequest: true,
	}
	var response IPList
	var IPs []IP
	err := r.execute(ctx, *c, &response)
	for _, properties := range response.List {
		IPs = append(IPs, IP{Properties: properties})
	}
	return IPs, err
}

//GetDeletedIPs gets a list of deleted IPs
//
//See: https://gridscale.io/en//api-documentation/index.html#operation/getDeletedIps
func (c *Client) GetDeletedIPs(ctx context.Context) ([]IP, error) {
	r := request{
		uri:                 path.Join(apiDeletedBase, "ips"),
		method:              http.MethodGet,
		skipCheckingRequest: true,
	}
	var response DeletedIPList
	var IPs []IP
	err := r.execute(ctx, *c, &response)
	for _, properties := range response.List {
		IPs = append(IPs, IP{Properties: properties})
	}
	return IPs, err
}
