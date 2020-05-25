package gsclient

import (
	"context"
	"errors"
	"net/http"
	"path"
)

//LocationList JSON struct of a list of locations
type LocationList struct {
	//Array of locations
	List map[string]LocationProperties `json:"locations"`
}

//Location JSON struct of a single location
type Location struct {
	//Properties of a location
	Properties LocationProperties `json:"location"`
}

//LocationProperties JSON struct of properties of a location
type LocationProperties struct {
	//Uses IATA airport code, which works as a location identifier.
	Iata string `json:"iata"`

	//Status indicates the status of the object.
	Status string `json:"status"`

	//List of labels.
	Labels []string `json:"labels"`

	//The human-readable name of the location. It supports the full UTF-8 charset, with a maximum of 64 characters.
	Name string `json:"name"`

	//The UUID of an object is always unique, and refers to a specific object.
	ObjectUUID string `json:"object_uuid"`

	//The human-readable name of the location. It supports the full UTF-8 charset, with a maximum of 64 characters.
	Country string `json:"country"`
}

//GetLocationList gets a list of available locations]
//
//See: https://gridscale.io/en//api-documentation/index.html#operation/getLocations
func (c *Client) GetLocationList(ctx context.Context) ([]Location, error) {
	r := request{
		uri:                 apiLocationBase,
		method:              http.MethodGet,
		skipCheckingRequest: true,
	}
	var response LocationList
	var locations []Location
	err := r.execute(ctx, *c, &response)
	for _, properties := range response.List {
		locations = append(locations, Location{Properties: properties})
	}
	return locations, err
}

//GetLocation gets a specific location
//
//See: https://gridscale.io/en//api-documentation/index.html#operation/getLocation
func (c *Client) GetLocation(ctx context.Context, id string) (Location, error) {
	if !isValidUUID(id) {
		return Location{}, errors.New("'id' is invalid")
	}
	r := request{
		uri:                 path.Join(apiLocationBase, id),
		method:              http.MethodGet,
		skipCheckingRequest: true,
	}
	var location Location
	err := r.execute(ctx, *c, &location)
	return location, err
}
