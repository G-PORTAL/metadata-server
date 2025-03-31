package utils

import (
	"fmt"
	"net"
)

// ConvertCIDRToNetworkAndSubnet takes a CIDR notation string and returns the network address and subnet mask.
func ConvertCIDRToNetworkAndSubnet(cidr string) (string, string, error) {
	_, network, err := net.ParseCIDR(cidr)
	if err != nil {
		return "", "", fmt.Errorf("failed to parse CIDR: %w", err)
	}

	address := network.IP.String()
	netmask := net.IP(network.Mask).String()

	return address, netmask, nil
}

// ConvertCIDRToIPNet takes a CIDR notation string and return the *net.IPNet object.
func ConvertCIDRToIPNet(cidr string) (net.IP, *net.IPNet, error) {
	ip, network, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse CIDR: %w", err)
	}

	return ip, network, nil
}
