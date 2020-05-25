package gsclient

import (
	"context"
	"errors"
	"net/http"
	"path"
)

//ServerIsoImageRelationList JSON struct of a list of relations between a server and ISO-Images
type ServerIsoImageRelationList struct {
	//Array of relations between a server and ISO-Images
	List []ServerIsoImageRelationProperties `json:"isoimage_relations"`
}

//ServerIsoImageRelation JSON struct of a single relation between a server and an ISO-Image
type ServerIsoImageRelation struct {
	//Properties of a relation between a server and an ISO-Image
	Properties ServerIsoImageRelationProperties `json:"isoimage_relation"`
}

//ServerIsoImageRelationProperties JSON struct of properties of a relation between a server and an ISO-Image
type ServerIsoImageRelationProperties struct {
	//The UUID of an object is always unique, and refers to a specific object.
	ObjectUUID string `json:"object_uuid"`

	//The human-readable name of the object. It supports the full UTF-8 charset, with a maximum of 64 characters.
	ObjectName string `json:"object_name"`

	//Whether the ISO-Image is private or not.
	Private bool `json:"private"`

	//Defines the date and time the object was initially created.
	CreateTime GSTime `json:"create_time"`

	//Whether the server boots from this iso image or not.
	Bootdevice bool `json:"bootdevice"`
}

//ServerIsoImageRelationCreateRequest JSON struct of a request for creating a relation between a server and an ISO-Image
type ServerIsoImageRelationCreateRequest struct {
	//The UUID of the ISO-image you are requesting.
	ObjectUUID string `json:"object_uuid"`
}

//ServerIsoImageRelationUpdateRequest JSON struct of a request for updating a relation between a server and an ISO-Image
type ServerIsoImageRelationUpdateRequest struct {
	//Whether the server boots from this ISO-image or not.
	BootDevice bool   `json:"bootdevice"`
	Name       string `json:"name"`
}

//GetServerIsoImageList gets a list of a specific server's ISO images
//
//See: https://gridscale.io/en//api-documentation/index.html#operation/getServerLinkedIsoimages
func (c *Client) GetServerIsoImageList(ctx context.Context, id string) ([]ServerIsoImageRelationProperties, error) {
	if !isValidUUID(id) {
		return nil, errors.New("'id' is invalid")
	}
	r := request{
		uri:                 path.Join(apiServerBase, id, "isoimages"),
		method:              http.MethodGet,
		skipCheckingRequest: true,
	}
	var response ServerIsoImageRelationList
	err := r.execute(ctx, *c, &response)
	return response.List, err
}

//GetServerIsoImage gets an ISO image of a specific server
//
//See: https://gridscale.io/en//api-documentation/index.html#operation/getServerLinkedIsoimage
func (c *Client) GetServerIsoImage(ctx context.Context, serverID, isoImageID string) (ServerIsoImageRelationProperties, error) {
	if !isValidUUID(serverID) || !isValidUUID(isoImageID) {
		return ServerIsoImageRelationProperties{}, errors.New("'id' is invalid")
	}
	r := request{
		uri:                 path.Join(apiServerBase, serverID, "isoimages", isoImageID),
		method:              http.MethodGet,
		skipCheckingRequest: true,
	}
	var response ServerIsoImageRelation
	err := r.execute(ctx, *c, &response)
	return response.Properties, err
}

//UpdateServerIsoImage updates a link between a storage and an ISO image
//
//See: https://gridscale.io/en//api-documentation/index.html#operation/updateServerLinkedIsoimage
func (c *Client) UpdateServerIsoImage(ctx context.Context, serverID, isoImageID string, body ServerIsoImageRelationUpdateRequest) error {
	if !isValidUUID(serverID) || !isValidUUID(isoImageID) {
		return errors.New("'serverID' or 'isoImageID' is invalid")
	}
	r := request{
		uri:    path.Join(apiServerBase, serverID, "isoimages", isoImageID),
		method: http.MethodPatch,
		body:   body,
	}
	return r.execute(ctx, *c, nil)
}

//CreateServerIsoImage creates a link between a server and an ISO image
//
//See: https://gridscale.io/en//api-documentation/index.html#operation/linkIsoimageToServer
func (c *Client) CreateServerIsoImage(ctx context.Context, id string, body ServerIsoImageRelationCreateRequest) error {
	if !isValidUUID(id) || !isValidUUID(body.ObjectUUID) {
		return errors.New("'serverID' or 'isoImageID' is invalid")
	}
	r := request{
		uri:    path.Join(apiServerBase, id, "isoimages"),
		method: http.MethodPost,
		body:   body,
	}
	return r.execute(ctx, *c, nil)
}

//DeleteServerIsoImage deletes a link between an ISO image and a server
//
//See: https://gridscale.io/en//api-documentation/index.html#operation/unlinkIsoimageFromServer
func (c *Client) DeleteServerIsoImage(ctx context.Context, serverID, isoImageID string) error {
	if !isValidUUID(serverID) || !isValidUUID(isoImageID) {
		return errors.New("'serverID' or 'isoImageID' is invalid")
	}
	r := request{
		uri:    path.Join(apiServerBase, serverID, "isoimages", isoImageID),
		method: http.MethodDelete,
	}
	return r.execute(ctx, *c, nil)
}

//LinkIsoImage attaches an ISO image to a server
func (c *Client) LinkIsoImage(ctx context.Context, serverID string, isoimageID string) error {
	body := ServerIsoImageRelationCreateRequest{
		ObjectUUID: isoimageID,
	}
	return c.CreateServerIsoImage(ctx, serverID, body)
}

//UnlinkIsoImage removes the link between an ISO image and a server
func (c *Client) UnlinkIsoImage(ctx context.Context, serverID string, isoimageID string) error {
	return c.DeleteServerIsoImage(ctx, serverID, isoimageID)
}
