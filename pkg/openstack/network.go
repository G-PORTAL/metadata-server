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
	LinkTypeVlan      LinkType = "vlan"
)

type Link struct {
	ID                 string `json:"id"`
	EthernetMacAddress string `json:"ethernet_mac_address"`

	// VLan settings only
	VlanMacAddress *string `json:"vlan_mac_address,omitempty"`
	VlanID         *int    `json:"vlan_id,omitempty"`
	VlanLink       *string `json:"vlan_link,omitempty"`

	Type  LinkType `json:"type"`
	VifID *string  `json:"vif_id,omitempty"`
	Mtu   *int     `json:"mtu,omitempty"`
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
	Metric  *int   `json:"metric,omitempty"`
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

	// Gateway sets the gateway properly even if this flag officially is not
	// supported if you look into the schema file:
	//   https://docs.openstack.org/nova/latest/_downloads/9119ca7ac90aa2990e762c08baea3a36/network_data.json
	// But Cloud-Init does seem to support it, see cloud-init documentation for OpenStack:
	//   https://github.com/canonical/cloud-init/blob/main/cloudinit/sources/helpers/openstack.py#L576
	Gateway *string `json:"gateway,omitempty"`
}
