package openstack

type Interfaces struct {
	Links    []Link    `json:"links"`
	Networks []Network `json:"networks"`
	Services []Service `json:"services"`
}

type LinkType string

const (
	LinkTypeBridge    LinkType = "bridge"
	LinkTypeDVS       LinkType = "dvs"
	LinkTypeHwVeb     LinkType = "hw_veb"
	LinkTypeHyperV    LinkType = "hyperv"
	LinkTypeOVS       LinkType = "ovs"
	LinkTypeTyp       LinkType = "tap"
	LinkTypeVhostUser LinkType = "vhostuser"
	LinkTypeVif       LinkType = "vif"
	LinkTypePhysical  LinkType = "phy"
)

type Link struct {
	ID                 string   `json:"id"`
	EthernetMacAddress string   `json:"ethernet_mac_address"`
	Type               LinkType `json:"type"`
	VifID              *string  `json:"vif_id,omitempty"`
	Mtu                *int     `json:"mtu,omitempty"`
}

type ServiceType string

const (
	ServiceTypeDNS ServiceType = "dns"
)

type Service struct {
	Type    ServiceType `json:"type"`
	Address string      `json:"address"`
}

type Route struct {
	Network string `json:"network"`
	Netmask string `json:"netmask"`
	Gateway string `json:"gateway"`
}

type NetworkType string

const (
	NetworkTypeIPv4     NetworkType = "ipv4"
	NetworkTypeIPv4DHCP NetworkType = "ipv4_dhcp"
)

type Network struct {
	ID        string      `json:"id"`
	Link      string      `json:"link"`
	NetworkID string      `json:"network_id"`
	Type      NetworkType `json:"type"`
	IPAddress *string     `json:"ip_address,omitempty"`
	Netmask   *string     `json:"netmask,omitempty"`
	Routes    []Route     `json:"routes,omitempty"`
}
