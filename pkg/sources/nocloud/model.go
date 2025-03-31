package nocloud

type NicConfig struct {
	Renderer    string `yaml:"renderer"`
	Nameservers struct {
		Addresses []string `yaml:"addresses"`
	} `yaml:"nameservers"`
	Match struct {
		MacAddress string `yaml:"macaddress"`
	} `yaml:"match"`
	Dhcp4     bool     `yaml:"dhcp4"`
	Addresses []string `yaml:"addresses"`
	Routes    []struct {
		To  string `yaml:"to"`
		Via string `yaml:"via"`
	} `yaml:"routes"`
}

type VlanConfig struct {
	ID          int      `yaml:"id"`
	Link        string   `yaml:"link"`
	Dhcp4       bool     `yaml:"dhcp4"`
	Addresses   []string `yaml:"addresses"`
	Nameservers struct {
		Addresses []string `yaml:"addresses"`
	} `yaml:"nameservers"`
	Routes []struct {
		To     string `yaml:"to"`
		Via    string `yaml:"via"`
		Metric int    `yaml:"metric,omitempty"`
	} `yaml:"routes"`
}

type NetworkConfigResponse struct {
	Version   int                   `yaml:"version"`
	Ethernets map[string]NicConfig  `yaml:"ethernets"`
	Vlans     map[string]VlanConfig `yaml:"vlans"`
}

type MetadataResponse struct {
	InstanceID    string `yaml:"instance-id"`
	LocalHostname string `yaml:"local-hostname"`
}
