package gsclient

import (
	"context"
	"net/http"
)

//LabelList JSON struct of a list of labels
type LabelList struct {
	//List of labels
	List map[string]LabelProperties `json:"labels"`
}

//Label JSON struct of a single label
type Label struct {
	//Properties of a label
	Properties LabelProperties `json:"label"`
}

//LabelProperties JSON struct of properties of a label
type LabelProperties struct {
	//Label's name
	Label string `json:"label"`

	//Create time of a label
	CreateTime GSTime `json:"create_time"`

	//Time of the last change of a label
	ChangeTime GSTime `json:"change_time"`

	//Relations of a label
	Relations []interface{} `json:"relations"`

	//Status indicates the status of a label.
	Status string `json:"status"`
}

//LabelCreateRequest JSON struct of a request for creating a label
type LabelCreateRequest struct {
	//Name of the new label
	Label string `json:"label"`
}

//GetLabelList gets a list of available labels
//
//See: https://gridscale.io/en//api-documentation/index.html#operation/GetLabels
func (c *Client) GetLabelList(ctx context.Context) ([]Label, error) {
	r := request{
		uri:                 apiLabelBase,
		method:              http.MethodGet,
		skipCheckingRequest: true,
	}
	var response LabelList
	var labels []Label
	err := r.execute(ctx, *c, &response)
	for _, properties := range response.List {
		labels = append(labels, Label{Properties: properties})
	}
	return labels, err
}
