package gsclient

import (
	"context"
	"errors"
	"net/http"
	"path"
)

//StorageList JSON struct of a list of storages
type StorageList struct {
	//Array of storages
	List map[string]StorageProperties `json:"storages"`
}

//DeletedStorageList JSON struct of a list of storages
type DeletedStorageList struct {
	//Array of deleted storages
	List map[string]StorageProperties `json:"deleted_storages"`
}

//Storage JSON struct of a single storage
type Storage struct {
	//Properties of a storage
	Properties StorageProperties `json:"storage"`
}

//StorageProperties JSON struct of properties of a storage
type StorageProperties struct {
	//Defines the date and time of the last object change.
	ChangeTime GSTime `json:"change_time"`

	//Uses IATA airport code, which works as a location identifier.
	LocationIata string `json:"location_iata"`

	//Status indicates the status of the object.
	Status string `json:"status"`

	//If a template has been used that requires a license key (e.g. Windows Servers)
	//this shows the product_no of the license (see the /prices endpoint for more details).
	LicenseProductNo int `json:"license_product_no"`

	//The human-readable name of the location. It supports the full UTF-8 charset, with a maximum of 64 characters.
	LocationCountry string `json:"location_country"`

	//Total minutes the object has been running.
	UsageInMinutes int `json:"usage_in_minutes"`

	//Indicates the UUID of the last used template on this storage.
	LastUsedTemplate string `json:"last_used_template"`

	//The price for the current period since the last bill.
	CurrentPrice float64 `json:"current_price"`

	//The capacity of a storage/ISO-Image/template/snapshot in GB.
	Capacity int `json:"capacity"`

	//Helps to identify which datacenter an object belongs to.
	LocationUUID string `json:"location_uuid"`

	//(one of storage, storage_high, storage_insane).
	StorageType string `json:"storage_type"`

	//The UUID of the Storage used to create this Snapshot.
	ParentUUID string `json:"parent_uuid"`

	//The human-readable name of the object. It supports the full UTF-8 charset, with a maximum of 64 characters.
	Name string `json:"name"`

	//The human-readable name of the location. It supports the full UTF-8 charset, with a maximum of 64 characters.
	LocationName string `json:"location_name"`

	//The UUID of an object is always unique, and refers to a specific object.
	ObjectUUID string `json:"object_uuid"`

	//Snapshots list in this storage.
	Snapshots []StorageSnapshotRelation `json:"snapshots"`

	//The information about other object which are related to this storage.
	//The object could be servers and/or snapshot schedules.
	Relations StorageRelations `json:"relations"`

	//List of labels.
	Labels []string `json:"labels"`

	//Defines the date and time the object was initially created.
	CreateTime GSTime `json:"create_time"`
}

//StorageRelations JSON struct of a list of a storage's relations
type StorageRelations struct {
	//Array of related servers
	Servers []StorageServerRelation `json:"servers"`

	//Array if related snapshot schedules
	SnapshotSchedules []StorageAndSnapshotScheduleRelation `json:"snapshot_schedules"`
}

//StorageServerRelation JSON struct of a relation between a storage and a server
type StorageServerRelation struct {
	//Whether the server boots from this iso image or not.
	Bootdevice bool `json:"bootdevice"`

	//Defines the SCSI target ID. The SCSI defines transmission routes like Serial Attached SCSI (SAS),
	//Fibre Channel and iSCSI. The target ID is a device (e.g. disk).
	Target int `json:"target"`

	//Defines the SCSI controller id. The SCSI defines transmission routes such as Serial Attached SCSI (SAS), Fibre Channel and iSCSI.
	Controller int `json:"controller"`

	//The SCSI bus id. The SCSI defines transmission routes like Serial Attached SCSI (SAS), Fibre Channel and iSCSI.
	//Each SCSI device is addressed via a specific number. Each SCSI bus can have multiple SCSI devices connected to it.
	Bus int `json:"bus"`

	//The UUID of an object is always unique, and refers to a specific object.
	ObjectUUID string `json:"object_uuid"`

	//Is the common SCSI abbreviation of the Logical Unit Number. A lun is a unique identifier for a single disk or a composite of disks.
	Lun int `json:"lun"`

	//Defines the date and time the object was initially created.
	CreateTime GSTime `json:"create_time"`

	//The human-readable name of the object. It supports the full UTF-8 charset, with a maximum of 64 characters.
	ObjectName string `json:"object_name"`
}

//StorageSnapshotRelation JSON struct of a relation between a storage and a snapshot
type StorageSnapshotRelation struct {
	//Indicates the UUID of the last used template on this storage.
	LastUsedTemplate string `json:"last_used_template"`

	//The UUID of an object is always unique, and refers to a specific object.
	ObjectUUID string `json:"object_uuid"`

	//The UUID of an object is always unique, and refers to a specific object.
	StorageUUID string `json:"storage_uuid"`

	//The human-readable name of the object. It supports the full UTF-8 charset, with a maximum of 64 characters.
	SchedulesSnapshotName string `json:"schedules_snapshot_name"`

	//The UUID of an object is always unique, and refers to a specific object.
	SchedulesSnapshotUUID string `json:"schedules_snapshot_uuid"`

	//Capacity of the snapshot (in GB)
	ObjectCapacity int `json:"object_capacity"`

	//Defines the date and time the object was initially created.
	CreateTime GSTime `json:"create_time"`

	//The human-readable name of the object. It supports the full UTF-8 charset, with a maximum of 64 characters.
	ObjectName string `json:"object_name"`
}

//StorageAndSnapshotScheduleRelation JSON struct of a relation between a storage and a snapshot schedule
type StorageAndSnapshotScheduleRelation struct {
	//The interval at which the schedule will run (in minutes)
	RunInterval int `json:"run_interval"`

	//The amount of Snapshots to keep before overwriting the last created Snapshot.
	KeepSnapshots int `json:"keep_snapshots"`

	//The human-readable name of the object. It supports the full UTF-8 charset, with a maximum of 64 characters.
	ObjectName string `json:"object_name"`

	//The date and time that the snapshot schedule will be run.
	NextRuntime GSTime `json:"next_runtime"`

	//The UUID of an object is always unique, and refers to a specific object.
	ObjectUUID string `json:"object_uuid"`

	//The human-readable name of the object. It supports the full UTF-8 charset, with a maximum of 64 characters.
	Name string `json:"name"`

	//Defines the date and time the object was initially created.
	CreateTime GSTime `json:"create_time"`
}

//StorageTemplate JSON struct of a storage template
type StorageTemplate struct {
	//List of SSH Keys uuid. Optional
	Sshkeys []string `json:"sshkeys,omitempty"`

	//The UUID of a template (public or private).
	TemplateUUID string `json:"template_uuid"`

	//The root (Linux) or Administrator (Windows) password to set for the installed storage. Valid only for public templates.
	//The password has to be either plaintext or a crypt string (modular crypt format - MCF). Optional.
	Password string `json:"password,omitempty"`

	//Password type. Allowed values: nil, PlainPasswordType, CryptPasswordType. Optional.
	PasswordType *passwordType `json:"password_type,omitempty"`

	//Hostname to set for the installed storage. The running server will use this as its hostname.
	//Valid only for public Linux and Windows templates. Optional.
	Hostname string `json:"hostname,omitempty"`
}

//StorageCreateRequest JSON struct of a request for creating a storage
type StorageCreateRequest struct {
	//Required (integer - minimum: 1 - maximum: 4096).
	Capacity int `json:"capacity"`

	//The human-readable name of the object. It supports the full UTF-8 charset, with a maximum of 64 characters.
	Name string `json:"name"`

	//Storage type. Allowed values: nil, DefaultStorageType, HighStorageType, InsaneStorageType. Optional.
	StorageType *storageType `json:"storage_type,omitempty"`

	//An object holding important values such as hostnames, passwords, and SSH keys.
	//Creating a storage with a template is required either sshkey or password.
	//Optional
	Template *StorageTemplate `json:"template,omitempty"`

	//List of labels. Optional.
	Labels []string `json:"labels,omitempty"`
}

//StorageUpdateRequest JSON struct of a request for updating a storage
type StorageUpdateRequest struct {
	//The human-readable name of the object. It supports the full UTF-8 charset, with a maximum of 64 characters. Optional.
	Name string `json:"name,omitempty"`

	//List of labels. Optional.
	Labels *[]string `json:"labels,omitempty"`

	//The Capacity of the Storage in GB. Optional.
	Capacity int `json:"capacity,omitempty"`
}

//All allowed storage type's values
var (
	DefaultStorageType = &storageType{"storage"}
	HighStorageType    = &storageType{"storage_high"}
	InsaneStorageType  = &storageType{"storage_insane"}
)

//All allowed password type's values
var (
	PlainPasswordType = &passwordType{"plain"}
	CryptPasswordType = &passwordType{"crypt"}
)

//GetStorage get a storage
//
//See: https://gridscale.io/en//api-documentation/index.html#operation/getStorage
func (c *Client) GetStorage(ctx context.Context, id string) (Storage, error) {
	if !isValidUUID(id) {
		return Storage{}, errors.New("'id' is invalid")
	}
	r := request{
		uri:                 path.Join(apiStorageBase, id),
		method:              http.MethodGet,
		skipCheckingRequest: true,
	}
	var response Storage
	err := r.execute(ctx, *c, &response)
	return response, err
}

//GetStorageList gets a list of available storages
//
//See: https://gridscale.io/en//api-documentation/index.html#operation/getStorages
func (c *Client) GetStorageList(ctx context.Context) ([]Storage, error) {
	r := request{
		uri:                 apiStorageBase,
		method:              http.MethodGet,
		skipCheckingRequest: true,
	}
	var response StorageList
	var storages []Storage
	err := r.execute(ctx, *c, &response)
	for _, properties := range response.List {
		storages = append(storages, Storage{
			Properties: properties,
		})
	}
	return storages, err
}

//CreateStorage create a storage
//
//NOTE:
//
// - Allowed value for `StorageType`: nil, DefaultStorageType, HighStorageType, InsaneStorageType.
//
// - Allowed value for `PasswordType`: nil, PlainPasswordType, CryptPasswordType.
//
//See: https://gridscale.io/en//api-documentation/index.html#operation/createStorage
func (c *Client) CreateStorage(ctx context.Context, body StorageCreateRequest) (CreateResponse, error) {
	r := request{
		uri:    apiStorageBase,
		method: http.MethodPost,
		body:   body,
	}
	var response CreateResponse
	err := r.execute(ctx, *c, &response)
	return response, err
}

//DeleteStorage delete a storage
//
//See: https://gridscale.io/en//api-documentation/index.html#operation/deleteStorage
func (c *Client) DeleteStorage(ctx context.Context, id string) error {
	if !isValidUUID(id) {
		return errors.New("'id' is invalid")
	}
	r := request{
		uri:    path.Join(apiStorageBase, id),
		method: http.MethodDelete,
	}
	return r.execute(ctx, *c, nil)
}

//UpdateStorage update a storage
//
//See: https://gridscale.io/en//api-documentation/index.html#operation/updateStorage
func (c *Client) UpdateStorage(ctx context.Context, id string, body StorageUpdateRequest) error {
	if !isValidUUID(id) {
		return errors.New("'id' is invalid")
	}
	r := request{
		uri:    path.Join(apiStorageBase, id),
		method: http.MethodPatch,
		body:   body,
	}
	return r.execute(ctx, *c, nil)
}

//GetStorageEventList get list of a storage's event
//
//See: https://gridscale.io/en//api-documentation/index.html#operation/getStorageEvents
func (c *Client) GetStorageEventList(ctx context.Context, id string) ([]Event, error) {
	if !isValidUUID(id) {
		return nil, errors.New("'id' is invalid")
	}
	r := request{
		uri:                 path.Join(apiStorageBase, id, "events"),
		method:              http.MethodGet,
		skipCheckingRequest: true,
	}
	var response EventList
	var storageEvents []Event
	err := r.execute(ctx, *c, &response)
	for _, properties := range response.List {
		storageEvents = append(storageEvents, Event{Properties: properties})
	}
	return storageEvents, err
}

//GetStoragesByLocation gets a list of storages by location
//
//See: https://gridscale.io/en//api-documentation/index.html#operation/getLocationStorages
func (c *Client) GetStoragesByLocation(ctx context.Context, id string) ([]Storage, error) {
	if !isValidUUID(id) {
		return nil, errors.New("'id' is invalid")
	}
	r := request{
		uri:                 path.Join(apiLocationBase, id, "storages"),
		method:              http.MethodGet,
		skipCheckingRequest: true,
	}
	var response StorageList
	var storages []Storage
	err := r.execute(ctx, *c, &response)
	for _, properties := range response.List {
		storages = append(storages, Storage{Properties: properties})
	}
	return storages, err
}

//GetDeletedStorages gets a list of deleted storages
//
//See: https://gridscale.io/en//api-documentation/index.html#operation/getDeletedStorages
func (c *Client) GetDeletedStorages(ctx context.Context) ([]Storage, error) {
	r := request{
		uri:                 path.Join(apiDeletedBase, "storages"),
		method:              http.MethodGet,
		skipCheckingRequest: true,
	}
	var response DeletedStorageList
	var storages []Storage
	err := r.execute(ctx, *c, &response)
	for _, properties := range response.List {
		storages = append(storages, Storage{Properties: properties})
	}
	return storages, err
}

//CloneStorage clones a specific storage
//
func (c *Client) CloneStorage(ctx context.Context, id string) (CreateResponse, error) {
	var response CreateResponse
	if !isValidUUID(id) {
		return response, errors.New("'id' is invalid")
	}
	r := request{
		uri:    path.Join(apiStorageBase, id, "clone"),
		method: http.MethodPost,
	}
	err := r.execute(ctx, *c, &response)
	return response, err
}
