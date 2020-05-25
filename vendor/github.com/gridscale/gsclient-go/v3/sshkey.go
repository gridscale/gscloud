package gsclient

import (
	"context"
	"errors"
	"net/http"
	"path"
)

//SshkeyList JSON struct of a list of SSH-keys
type SshkeyList struct {
	//Array of SSH-keys
	List map[string]SshkeyProperties `json:"sshkeys"`
}

//Sshkey JSON struct of a single SSH-key
type Sshkey struct {
	//Properties of a SSH-key
	Properties SshkeyProperties `json:"sshkey"`
}

//SshkeyProperties JSON struct of properties of a single SSH-key
type SshkeyProperties struct {
	//The human-readable name of the object. It supports the full UTF-8 charset, with a maximum of 64 characters.
	Name string `json:"name"`

	//The UUID of an object is always unique, and refers to a specific object.
	ObjectUUID string `json:"object_uuid"`

	//Status indicates the status of the object.
	Status string `json:"status"`

	//Defines the date and time the object was initially created.
	CreateTime GSTime `json:"create_time"`

	//Defines the date and time of the last object change.
	ChangeTime GSTime `json:"change_time"`

	//The OpenSSH public key string (all key types are supported => ed25519, ecdsa, dsa, rsa, rsa1).
	Sshkey string `json:"sshkey"`

	//List of labels.
	Labels []string `json:"labels"`

	//The User-UUID of the account which created this SSH Key.
	UserUUID string `json:"user_uuid"`
}

//SshkeyCreateRequest JSON struct of a request for creating a SSH-key
type SshkeyCreateRequest struct {
	//The human-readable name of the object. It supports the full UTF-8 charset, with a maximum of 64 characters.
	Name string `json:"name"`

	//The OpenSSH public key string (all key types are supported => ed25519, ecdsa, dsa, rsa, rsa1).
	Sshkey string `json:"sshkey"`

	//List of labels. Optional.
	Labels []string `json:"labels,omitempty"`
}

//SshkeyUpdateRequest JSON struct of a request for updating a SSH-key
type SshkeyUpdateRequest struct {
	//The human-readable name of the object. It supports the full UTF-8 charset, with a maximum of 64 characters.
	//Optional.
	Name string `json:"name,omitempty"`

	//The OpenSSH public key string (all key types are supported => ed25519, ecdsa, dsa, rsa, rsa1). Optional.
	Sshkey string `json:"sshkey,omitempty"`

	//List of labels. Optional.
	Labels *[]string `json:"labels,omitempty"`
}

//GetSshkey gets a ssh key
//
//See: https://gridscale.io/en//api-documentation/index.html#operation/getSshKey
func (c *Client) GetSshkey(ctx context.Context, id string) (Sshkey, error) {
	if !isValidUUID(id) {
		return Sshkey{}, errors.New("'id' is invalid")
	}
	r := request{
		uri:                 path.Join(apiSshkeyBase, id),
		method:              http.MethodGet,
		skipCheckingRequest: true,
	}
	var response Sshkey
	err := r.execute(ctx, *c, &response)
	return response, err
}

//GetSshkeyList gets a list of ssh keys
//
//See: https://gridscale.io/en//api-documentation/index.html#operation/getSshKeys
func (c *Client) GetSshkeyList(ctx context.Context) ([]Sshkey, error) {
	r := request{
		uri:                 apiSshkeyBase,
		method:              http.MethodGet,
		skipCheckingRequest: true,
	}

	var response SshkeyList
	var sshKeys []Sshkey
	err := r.execute(ctx, *c, &response)
	for _, properties := range response.List {
		sshKeys = append(sshKeys, Sshkey{Properties: properties})
	}
	return sshKeys, err
}

//CreateSshkey creates a ssh key
//
//See: https://gridscale.io/en//api-documentation/index.html#operation/createSshKey
func (c *Client) CreateSshkey(ctx context.Context, body SshkeyCreateRequest) (CreateResponse, error) {
	r := request{
		uri:    apiSshkeyBase,
		method: "POST",
		body:   body,
	}
	var response CreateResponse
	err := r.execute(ctx, *c, &response)
	return response, err
}

//DeleteSshkey deletes a ssh key
//
//See: https://gridscale.io/en//api-documentation/index.html#operation/deleteSshKey
func (c *Client) DeleteSshkey(ctx context.Context, id string) error {
	if !isValidUUID(id) {
		return errors.New("'id' is invalid")
	}
	r := request{
		uri:    path.Join(apiSshkeyBase, id),
		method: http.MethodDelete,
	}
	return r.execute(ctx, *c, nil)
}

//UpdateSshkey updates a ssh key
//
//See: https://gridscale.io/en//api-documentation/index.html#operation/updateSshKey
func (c *Client) UpdateSshkey(ctx context.Context, id string, body SshkeyUpdateRequest) error {
	if !isValidUUID(id) {
		return errors.New("'id' is invalid")
	}
	r := request{
		uri:    path.Join(apiSshkeyBase, id),
		method: http.MethodPatch,
		body:   body,
	}
	return r.execute(ctx, *c, nil)
}

//GetSshkeyEventList gets a ssh key's events
//
//See: https://gridscale.io/en//api-documentation/index.html#operation/getSshKeyEvents
func (c *Client) GetSshkeyEventList(ctx context.Context, id string) ([]Event, error) {
	if !isValidUUID(id) {
		return nil, errors.New("'id' is invalid")
	}
	r := request{
		uri:                 path.Join(apiSshkeyBase, id, "events"),
		method:              http.MethodGet,
		skipCheckingRequest: true,
	}
	var response EventList
	var sshEvents []Event
	err := r.execute(ctx, *c, &response)
	for _, properties := range response.List {
		sshEvents = append(sshEvents, Event{Properties: properties})
	}
	return sshEvents, err
}
