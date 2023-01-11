package sources

import (
	"golang.org/x/crypto/ssh"
	"net"
)

type MetadataSubnetType string

const (
	MetadataSubnetTypeDHCP   MetadataSubnetType = "dhcp"
	MetadataSubnetTypeStatic MetadataSubnetType = "static"
)

type MetadataSubnet struct {
	IPv4       bool               `json:"ipv4"`
	IPv6       bool               `json:"ipv6"`
	Type       MetadataSubnetType `json:"type"`
	Address    *net.IP            `json:"address,omitempty"`
	Network    *net.IPNet         `json:"network,omitempty"`
	Gateway    *net.IP            `json:"gateway,omitempty"`
	DNSServers []string           `json:"dns_nameservers,omitempty"`
}

type InterfaceType string

const (
	InterfaceTypePhysical InterfaceType = "physical"
	InterfaceTypeBond     InterfaceType = "bond"
	InterfaceTypeBridge   InterfaceType = "bridge"
	InterfaceTypeVLAN     InterfaceType = "vlan"
)

type MetadataInterface struct {
	MacAddress string           `json:"mac_address"`
	Name       string           `json:"name"`
	Type       InterfaceType    `json:"type"`
	Subnets    []MetadataSubnet `json:"subnets"`
	AcceptRA   *bool            `json:"accept-ra,omitempty"` //nolint:tagliatelle // cloud-init metadata requirement
}

type MetadataRoute struct {
	Network string  `json:"network"`
	Netmask string  `json:"netmask"`
	Gateway *string `json:"gateway"`
}

type Metadata struct {
	InstanceID       string                   `json:"id"`
	ProjectID        *string                  `json:"project_id"`
	InstanceType     string                   `json:"instance_type"`
	LocalHostname    string                   `json:"local_hostname"`
	PublicHostname   string                   `json:"public_hostname"`
	AvailabilityZone *string                  `json:"availability_zone"`
	UserData         []byte                   `json:"user_data"`
	VendorData       []byte                   `json:"vendor_data"`
	VendorData2      []byte                   `json:"vendor_data_2"`
	PublicKeys       map[string]ssh.PublicKey `json:"public_keys"`
	Username         *string                  `json:"username,omitempty"`
	Password         *string                  `json:"password,omitempty"`
	Interfaces       []MetadataInterface      `json:"interfaces"`
	Routes           []MetadataRoute          `json:"routes"`
}
