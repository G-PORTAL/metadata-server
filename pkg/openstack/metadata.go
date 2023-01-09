package openstack

type MetadataKeyDefinition struct {
	Name string `json:"name"`
	Type string `json:"type"`
	Data string `json:"data"`
}

type Metadata struct {
	UUID             string                  `json:"uuid"`
	PublicKeys       map[string]string       `json:"public_keys"`
	Keys             []MetadataKeyDefinition `json:"keys"`
	Hostname         string                  `json:"hostname"`
	Name             string                  `json:"name"`
	LaunchIndex      int                     `json:"launch_index"`
	AvailabilityZone string                  `json:"availability_zone"`
	RandomSeed       string                  `json:"random_seed"`
	ProjectID        string                  `json:"project_id"`
	Devices          []interface{}           `json:"devices"`
	AdminUsername    *string                 `json:"admin_username,omitempty"`
	AdminPassword    *string                 `json:"admin_pass,omitempty"`
}

type VendorData struct {
	CloudInit string `json:"cloud-init"` //nolint:tagliatelle // cloud-init metadata requirement
}
