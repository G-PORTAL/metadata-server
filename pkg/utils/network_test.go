package utils_test

import (
	"github.com/g-portal/metadata-server/pkg/utils"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConvertCIDRToNetworkAndSubnet(t *testing.T) {
	tests := []struct {
		cidr          string
		expectedAddr  string
		expectedMask  string
		expectedError bool
	}{
		{"192.168.1.1/24", "192.168.1.0", "255.255.255.0", false},
		{"10.0.0.1/8", "10.0.0.0", "255.0.0.0", false},
		{"invalid", "", "", true},
	}

	for _, tt := range tests {
		addr, mask, err := utils.ConvertCIDRToNetworkAndSubnet(tt.cidr)
		if tt.expectedError {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedAddr, addr)
			assert.Equal(t, tt.expectedMask, mask)
		}
	}
}

func TestConvertCIDRToIPNet(t *testing.T) {
	tests := []struct {
		cidr          string
		ip            net.IP
		expectedNet   *net.IPNet
		expectedError bool
	}{
		{"192.168.1.1/24", net.IPv4(192, 168, 1, 1), &net.IPNet{
			IP:   net.IPv4(192, 168, 1, 0).To4(),
			Mask: net.CIDRMask(24, 32),
		}, false},
		{"10.0.0.1/8", net.IPv4(10, 0, 0, 1), &net.IPNet{
			IP:   net.IPv4(10, 0, 0, 0).To4(),
			Mask: net.CIDRMask(8, 32),
		}, false},
		{"invalid", nil, nil, true},
	}

	for _, tt := range tests {
		ipAddr, ipNet, err := utils.ConvertCIDRToIPNet(tt.cidr)
		if tt.expectedError {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedNet, ipNet)
			assert.Equal(t, tt.ip, ipAddr)
		}
	}
}
