package endpoints

//Credential is JSON struct of credential
type Credential struct {
	//The initial username to authenticate the Service.
	Username string `json:"username"`

	//The initial password to authenticate the Service.
	Password string `json:"password"`

	//The type of Service.
	Type       string `json:"type"`
	KubeConfig string `json:"kubeconfig"`
}

//ResourceLimit is JSON struct of resource limit
type ResourceLimit struct {
	//The name of the resource you would like to cap.
	Resource string `json:"resource"`

	//The maximum number of the specific resource your service can use.
	Limit int `json:"limit"`
}

type PaaSService struct {
	//Properties of a PaaS service
	Properties PaaSServiceProperties `json:"paas_service"`
}

//PaaSServiceProperties is the properties of a single PaaS service
type PaaSServiceProperties struct {
	//The UUID of an object is always unique, and refers to a specific object.
	ObjectUUID string `json:"object_uuid"`

	//List of labels.
	Labels []string `json:"labels"`

	//Contains the initial setup credentials for Service.
	Credentials []Credential `json:"credentials"`

	//Defines the date and time the object was initially created.
	CreateTime GSTime `json:"create_time"`

	//Contains the IPv6 address and port that the Service will listen to,
	//you can use these details to connect internally to a service.
	ListenPorts map[string]map[string]int `json:"listen_ports"`

	//The UUID of the security zone that the service is running in.
	SecurityZoneUUID string `json:"security_zone_uuid"`

	//The template used to create the service, you can find an available list at the /service_templates endpoint.
	ServiceTemplateUUID string `json:"service_template_uuid"`

	//Total minutes the object has been running.
	UsageInMinutes int `json:"usage_in_minutes"`

	//The price for the current period since the last bill.
	CurrentPrice float64 `json:"current_price"`

	//Defines the date and time of the last object change.
	ChangeTime GSTime `json:"change_time"`

	//Status indicates the status of the object.
	Status string `json:"status"`

	//The human-readable name of the object. It supports the full UTF-8 charset, with a maximum of 64 characters.
	Name string `json:"name"`

	//A list of service resource limits.
	ResourceLimits []ResourceLimit `json:"resource_limits"`

	//Contains the service parameters for the service.
	Parameters map[string]interface{} `json:"parameters"`
}

type PaaSKubeCredentialBody struct {
}
