package gsclient

import (
	"context"
	"errors"
	"net/http"
	"path"
)

//LoadBalancers is the JSON struct of a list of loadbalancers
type LoadBalancers struct {
	//Array of loadbalancers
	List map[string]LoadBalancerProperties `json:"loadbalancers"`
}

//LoadBalancer is the JSON struct of a loadbalancer
type LoadBalancer struct {
	//Properties of a loadbalancer
	Properties LoadBalancerProperties `json:"loadbalancer"`
}

//LoadBalancerProperties is the properties of a loadbalancer
type LoadBalancerProperties struct {
	//The UUID of an object is always unique, and refers to a specific object.
	ObjectUUID string `json:"object_uuid"`

	//Defines the numbering of the Data Centers on a given IATA location (e.g. where fra is the location_iata, the site is then 1, 2, 3, ...).
	LocationSite string `json:"location_site"`

	//The human-readable name of the object. It supports the full UTF-8 charset, with a maximum of 64 characters.
	Name string `json:"name"`

	//Forwarding rules of a loadbalancer
	ForwardingRules []ForwardingRule `json:"forwarding_rules"`

	//Uses IATA airport code, which works as a location identifier.
	LocationIata string `json:"location_iata"`

	//Helps to identify which datacenter an object belongs to.
	LocationUUID string `json:"location_uuid"`

	//The servers that this Load balancer can communicate with.
	BackendServers []BackendServer `json:"backend_servers"`

	//Defines the date and time of the last object change.
	ChangeTime GSTime `json:"change_time"`

	//Status indicates the status of the object.
	Status string `json:"status"`

	//The price for the current period since the last bill.
	CurrentPrice float64 `json:"current_price"`

	//The human-readable name of the location. It supports the full UTF-8 charset, with a maximum of 64 characters.
	LocationCountry string `json:"location_country"`

	//Whether the Load balancer is forced to redirect requests from HTTP to HTTPS.
	RedirectHTTPToHTTPS bool `json:"redirect_http_to_https"`

	//List of labels.
	Labels []string `json:"labels"`

	//The human-readable name of the location. It supports the full UTF-8 charset, with a maximum of 64 characters.
	LocationName string `json:"location_name"`

	//Total minutes of cores used
	UsageInMinutes int `json:"usage_in_minutes"`

	//The algorithm used to process requests. Accepted values: roundrobin / leastconn.
	Algorithm string `json:"algorithm"`

	//Defines the date and time the object was initially created.
	CreateTime GSTime `json:"create_time"`

	//The UUID of the IPv6 address the Load balancer will listen to for incoming requests.
	ListenIPv6UUID string `json:"listen_ipv6_uuid"`

	//The UUID of the IPv4 address the Load balancer will listen to for incoming requests.
	ListenIPv4UUID string `json:"listen_ipv4_uuid"`
}

//BackendServer is the JSON struct of backend server
type BackendServer struct {
	//Weight of the server
	Weight int `json:"weight"`

	//Host of the server. Can be URL or IP address.
	Host string `json:"host"`
}

//ForwardingRule is the JSON struct of forwarding rule
type ForwardingRule struct {
	//SSL from Letsencrypt
	LetsencryptSSL interface{} `json:"letsencrypt_ssl"`

	//Listen port
	ListenPort int `json:"listen_port"`

	//Mode of forwarding
	Mode string `json:"mode"`

	//Target port
	TargetPort int `json:"target_port"`
}

//LoadBalancerCreateRequest is the JSON struct for creating a loadbalancer request
type LoadBalancerCreateRequest struct {
	//The human-readable name of the object. It supports the full UTF-8 charset, with a maximum of 64 characters.
	Name string `json:"name"`

	//The human-readable name of the object. It supports the full UTF-8 charset, with a maximum of 64 characters.
	ListenIPv6UUID string `json:"listen_ipv6_uuid"`

	//The UUID of the IPv4 address the loadbalancer will listen to for incoming requests.
	ListenIPv4UUID string `json:"listen_ipv4_uuid"`

	//The algorithm used to process requests. Allowed values: `LoadbalancerRoundrobinAlg`, `LoadbalancerLeastConnAlg`
	Algorithm loadbalancerAlgorithm `json:"algorithm"`

	//An array of ForwardingRule objects containing the forwarding rules for the loadbalancer
	ForwardingRules []ForwardingRule `json:"forwarding_rules"`

	//The servers that this loadbalancer can communicate with
	BackendServers []BackendServer `json:"backend_servers"`

	//List of labels.
	Labels []string `json:"labels"`

	//Whether the Load balancer is forced to redirect requests from HTTP to HTTPS
	RedirectHTTPToHTTPS bool `json:"redirect_http_to_https"`

	//Status indicates the status of the object.
	Status string `json:"status,omitempty"`
}

//LoadBalancerUpdateRequest is the JSON struct for updating a loadbalancer request
type LoadBalancerUpdateRequest struct {
	//The human-readable name of the object. It supports the full UTF-8 charset, with a maximum of 64 characters.
	Name string `json:"name"`

	//The human-readable name of the object. It supports the full UTF-8 charset, with a maximum of 64 characters.
	ListenIPv6UUID string `json:"listen_ipv6_uuid"`

	//The UUID of the IPv4 address the loadbalancer will listen to for incoming requests.
	ListenIPv4UUID string `json:"listen_ipv4_uuid"`

	//The algorithm used to process requests. Allowed values: `LoadbalancerRoundrobinAlg`, `LoadbalancerLeastConnAlg`
	Algorithm loadbalancerAlgorithm `json:"algorithm"`

	//An array of ForwardingRule objects containing the forwarding rules for the loadbalancer
	ForwardingRules []ForwardingRule `json:"forwarding_rules"`

	//The servers that this loadbalancer can communicate with
	BackendServers []BackendServer `json:"backend_servers"`

	//List of labels.
	Labels []string `json:"labels"`

	//Whether the Load balancer is forced to redirect requests from HTTP to HTTPS
	RedirectHTTPToHTTPS bool `json:"redirect_http_to_https"`

	//Status indicates the status of the object.
	Status string `json:"status,omitempty"`
}

//LoadBalancerCreateResponse is the JSON struct for a loadbalancer response
type LoadBalancerCreateResponse struct {
	//Request's UUID
	RequestUUID string `json:"request_uuid"`

	//UUID of the loadbalancer being created
	ObjectUUID string `json:"object_uuid"`
}

//All available loadbalancer algorithms
var (
	LoadbalancerRoundrobinAlg = loadbalancerAlgorithm{"roundrobin"}
	LoadbalancerLeastConnAlg  = loadbalancerAlgorithm{"leastconn"}
)

//GetLoadBalancerList returns a list of loadbalancers
//
//See: https://gridscale.io/en//api-documentation/index.html#operation/getLoadbalancers
func (c *Client) GetLoadBalancerList(ctx context.Context) ([]LoadBalancer, error) {
	r := request{
		uri:                 apiLoadBalancerBase,
		method:              http.MethodGet,
		skipCheckingRequest: true,
	}
	var response LoadBalancers
	var loadBalancers []LoadBalancer
	err := r.execute(ctx, *c, &response)
	for _, properties := range response.List {
		loadBalancers = append(loadBalancers, LoadBalancer{Properties: properties})
	}
	return loadBalancers, err
}

//GetLoadBalancer returns a loadbalancer of a given uuid
//
//See: https://gridscale.io/en//api-documentation/index.html#operation/getLoadbalancer
func (c *Client) GetLoadBalancer(ctx context.Context, id string) (LoadBalancer, error) {
	if !isValidUUID(id) {
		return LoadBalancer{}, errors.New("'id' is invalid")
	}
	r := request{
		uri:                 path.Join(apiLoadBalancerBase, id),
		method:              http.MethodGet,
		skipCheckingRequest: true,
	}
	var response LoadBalancer
	err := r.execute(ctx, *c, &response)
	return response, err
}

//CreateLoadBalancer creates a new loadbalancer
//
//Note: loadbalancer's algorithm can only be either `LoadbalancerRoundrobinAlg` or `LoadbalancerLeastConnAlg`
//
//See: https://gridscale.io/en//api-documentation/index.html#operation/createLoadbalancer
func (c *Client) CreateLoadBalancer(ctx context.Context, body LoadBalancerCreateRequest) (LoadBalancerCreateResponse, error) {
	if body.Labels == nil {
		body.Labels = make([]string, 0)
	}
	r := request{
		uri:    apiLoadBalancerBase,
		method: http.MethodPost,
		body:   body,
	}
	var response LoadBalancerCreateResponse
	err := r.execute(ctx, *c, &response)
	return response, err
}

//UpdateLoadBalancer update configuration of a loadbalancer
//
//Note: loadbalancer's algorithm can only be either `LoadbalancerRoundrobinAlg` or `LoadbalancerLeastConnAlg`
//
//See: https://gridscale.io/en//api-documentation/index.html#operation/updateLoadbalancer
func (c *Client) UpdateLoadBalancer(ctx context.Context, id string, body LoadBalancerUpdateRequest) error {
	if !isValidUUID(id) {
		return errors.New("'id' is invalid")
	}
	if body.Labels == nil {
		body.Labels = make([]string, 0)
	}
	r := request{
		uri:    path.Join(apiLoadBalancerBase, id),
		method: http.MethodPatch,
		body:   body,
	}
	return r.execute(ctx, *c, nil)
}

//GetLoadBalancerEventList retrieves events of a given uuid
//
//See: https://gridscale.io/en//api-documentation/index.html#operation/getLoadbalancerEvents
func (c *Client) GetLoadBalancerEventList(ctx context.Context, id string) ([]Event, error) {
	if !isValidUUID(id) {
		return nil, errors.New("'id' is invalid")
	}
	r := request{
		uri:                 path.Join(apiLoadBalancerBase, id, "events"),
		method:              http.MethodGet,
		skipCheckingRequest: true,
	}
	var response EventList
	var loadBalancerEvents []Event
	err := r.execute(ctx, *c, &response)
	for _, properties := range response.List {
		loadBalancerEvents = append(loadBalancerEvents, Event{Properties: properties})
	}
	return loadBalancerEvents, err
}

//DeleteLoadBalancer deletes a loadbalancer
//
//See: https://gridscale.io/en//api-documentation/index.html#operation/deleteLoadbalancer
func (c *Client) DeleteLoadBalancer(ctx context.Context, id string) error {
	if !isValidUUID(id) {
		return errors.New("'id' is invalid")
	}
	r := request{
		uri:    path.Join(apiLoadBalancerBase, id),
		method: http.MethodDelete,
	}
	return r.execute(ctx, *c, nil)
}
