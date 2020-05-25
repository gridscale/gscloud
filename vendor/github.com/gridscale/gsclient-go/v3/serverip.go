package gsclient

import (
	"context"
	"errors"
	"net/http"
	"path"
)

//ServerIPRelationList JSON struct of a list of relations between a server and IP addresses
type ServerIPRelationList struct {
	//Array of relations between a server and IP addresses
	List []ServerIPRelationProperties `json:"ip_relations"`
}

//ServerIPRelation JSON struct of a single relation between a server and a IP address
type ServerIPRelation struct {
	//Properties of a relation between a server and IP addresses
	Properties ServerIPRelationProperties `json:"ip_relation"`
}

//ServerIPRelationProperties JSON struct of properties of a relation between a server and a IP address
type ServerIPRelationProperties struct {
	//The UUID of the server that this IP is attached to.
	ServerUUID string `json:"server_uuid"`

	//Defines the date and time the object was initially created.
	CreateTime GSTime `json:"create_time"`

	//The prefix of the IP Address.
	Prefix string `json:"prefix"`

	//Either 4 or 6
	Family int `json:"family"`

	//The UUID of an object is always unique, and refers to a specific object.
	ObjectUUID string `json:"object_uuid"`

	//The IP Address (v4 or v6)
	IP string `json:"ip"`
}

//ServerIPRelationCreateRequest JSON struct of request for creating a relation between a server and a IP address
type ServerIPRelationCreateRequest struct {
	//You can only attach 1 IPv4 and/or IPv6 to a server based on the IP address's UUID
	ObjectUUID string `json:"object_uuid"`
}

//GetServerIPList gets a list of a specific server's IPs
//
//See: https://gridscale.io/en//api-documentation/index.html#operation/getServerLinkedIps
func (c *Client) GetServerIPList(ctx context.Context, id string) ([]ServerIPRelationProperties, error) {
	if id == "" {
		return nil, errors.New("'id' is required")
	}
	r := request{
		uri:                 path.Join(apiServerBase, id, "ips"),
		method:              http.MethodGet,
		skipCheckingRequest: true,
	}
	var response ServerIPRelationList
	err := r.execute(ctx, *c, &response)
	return response.List, err
}

//GetServerIP gets an IP of a specific server
//
//See: https://gridscale.io/en//api-documentation/index.html#operation/getServerLinkedIp
func (c *Client) GetServerIP(ctx context.Context, serverID, ipID string) (ServerIPRelationProperties, error) {
	if serverID == "" || ipID == "" {
		return ServerIPRelationProperties{}, errors.New("'serverID' and 'ipID' are required")
	}
	r := request{
		uri:                 path.Join(apiServerBase, serverID, "ips", ipID),
		method:              http.MethodGet,
		skipCheckingRequest: true,
	}
	var response ServerIPRelation
	err := r.execute(ctx, *c, &response)
	return response.Properties, err
}

//CreateServerIP create a link between a server and an IP
//
//See: https://gridscale.io/en//api-documentation/index.html#operation/linkIpToServer
func (c *Client) CreateServerIP(ctx context.Context, id string, body ServerIPRelationCreateRequest) error {
	if id == "" || body.ObjectUUID == "" {
		return errors.New("'server_id' and 'ip_id' are required")
	}
	r := request{
		uri:    path.Join(apiServerBase, id, "ips"),
		method: http.MethodPost,
		body:   body,
	}
	return r.execute(ctx, *c, nil)
}

//DeleteServerIP delete a link between a server and an IP
//
//See: https://gridscale.io/en//api-documentation/index.html#operation/unlinkIpFromServer
func (c *Client) DeleteServerIP(ctx context.Context, serverID, ipID string) error {
	if serverID == "" || ipID == "" {
		return errors.New("'serverID' and 'ipID' are required")
	}
	r := request{
		uri:    path.Join(apiServerBase, serverID, "ips", ipID),
		method: http.MethodDelete,
	}
	return r.execute(ctx, *c, nil)
}

//LinkIP attaches an IP to a server
func (c *Client) LinkIP(ctx context.Context, serverID string, ipID string) error {
	body := ServerIPRelationCreateRequest{
		ObjectUUID: ipID,
	}
	return c.CreateServerIP(ctx, serverID, body)
}

//UnlinkIP removes a link between an IP and a server
func (c *Client) UnlinkIP(ctx context.Context, serverID string, ipID string) error {
	return c.DeleteServerIP(ctx, serverID, ipID)
}
