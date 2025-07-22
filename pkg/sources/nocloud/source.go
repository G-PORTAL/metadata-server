package nocloud

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/g-portal/metadata-server/pkg/sources"
	"github.com/g-portal/metadata-server/pkg/utils"
	"gopkg.in/yaml.v3"
	"io"
	"net"
	"net/http"
	"strings"
)

const Type = "nocloud"

type Source struct {
	cfg      sources.SourceConfig
	endpoint string
}

func (s *Source) Type() string {
	return Type
}

func (s *Source) Initialize(cfg sources.SourceConfig) error {
	s.cfg = cfg

	endpoint := s.cfg.GetString("endpoint")
	if endpoint == "" {
		return errors.New("failed to initialize nocloud because endpoint config is missing")
	}

	s.endpoint = endpoint

	return nil
}

func (s *Source) GetMetadata(ip net.IP, client sources.MetadataClient) (*sources.Metadata, error) {
	// meta-data response
	metadata, err := s.fetchMetadata(ip)
	if err != nil {
		return nil, fmt.Errorf("failed to receive metadata: %w", err)
	}

	// optional user-data
	userData, err := s.fetchUserData(ip)
	if err != nil {
		userData = nil
	}

	// network config
	networkConfig, err := s.fetchNetworkConfig(ip)
	if err != nil {
		return nil, fmt.Errorf("failed to receive network-config: %w", err)
	}

	interfaceList := make([]sources.MetadataInterface, 0)
	routeList := make([]sources.MetadataRoute, 0)

	for nicName, networkInterface := range networkConfig.Ethernets {
		for _, route := range networkInterface.Routes {
			network, subnet, err := utils.ConvertCIDRToNetworkAndSubnet(route.To)
			if err != nil {
				continue
			}

			routeList = append(routeList, sources.MetadataRoute{
				Network: network,
				Netmask: subnet,
				Gateway: &route.Via,
			})
		}

		subnets := make([]sources.MetadataSubnet, 0)
		for _, address := range networkInterface.Addresses {
			ipAddress, network, err := utils.ConvertCIDRToIPNet(address)
			if err != nil {
				continue
			}

			subnets = append(subnets, sources.MetadataSubnet{
				IPv4:       true,
				IPv6:       false,
				Type:       sources.MetadataSubnetTypeStatic,
				Address:    &ipAddress,
				Network:    network,
				DNSServers: networkInterface.Nameservers.Addresses,
			})
		}

		macAddress := ""
		if networkInterface.Match.MacAddress != "" {
			macAddress = networkInterface.Match.MacAddress
		}

		interfaceList = append(interfaceList, sources.MetadataInterface{
			Name:       nicName,
			Type:       sources.InterfaceTypePhysical,
			MacAddress: macAddress,
			Subnets:    subnets,
		})
	}

	for vlanName, vlanConfig := range networkConfig.Vlans {
		for _, route := range vlanConfig.Routes {
			network, subnet, err := utils.ConvertCIDRToNetworkAndSubnet(route.To)
			if err != nil {
				continue
			}

			routeList = append(routeList, sources.MetadataRoute{
				Network: network,
				Netmask: subnet,
				Gateway: &route.Via,
				Metric:  &route.Metric,
			})
		}

		subnets := make([]sources.MetadataSubnet, 0)
		for _, address := range vlanConfig.Addresses {
			ipAddress, network, err := utils.ConvertCIDRToIPNet(address)
			if err != nil {
				continue
			}

			subnets = append(subnets, sources.MetadataSubnet{
				IPv4:       true,
				IPv6:       false,
				Type:       sources.MetadataSubnetTypeStatic,
				Address:    &ipAddress,
				Network:    network,
				DNSServers: vlanConfig.Nameservers.Addresses,
			})
		}

		var nicConfig *NicConfig
		var linkName string
		for nicName, config := range networkConfig.Ethernets {
			if vlanConfig.Link == nicName {
				nicConfig = &config
				linkName = nicName
				break
			}
		}

		macAddress := ""
		if nicConfig != nil && nicConfig.Match.MacAddress != "" {
			macAddress = nicConfig.Match.MacAddress
		}

		interfaceList = append(interfaceList, sources.MetadataInterface{
			Name:       vlanName,
			Type:       sources.InterfaceTypeVLAN,
			VlanID:     &vlanConfig.ID,
			VlanLink:   &linkName,
			MacAddress: macAddress,
			Subnets:    subnets,
		})
	}

	return &sources.Metadata{
		InstanceID:     metadata.InstanceID,
		LocalHostname:  metadata.LocalHostname,
		UserData:       userData,
		Routes:         routeList,
		Interfaces:     interfaceList,
		MetadataClient: client,
	}, nil
}

func (s *Source) ReportLog(message sources.ReportMessage) error {
	payload, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal report message: %w", err)
	}

	req, err := http.NewRequest(
		http.MethodPost, fmt.Sprintf("%s/report", strings.TrimSuffix(s.endpoint, "/")), bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Add("X-Real-IP", message.IP.String())
	req.Header.Add("X-Forwarded-For", message.IP.String())

	httpClient := http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("invalid http status code returned: %d", resp.StatusCode)
	}

	return nil
}

func (s *Source) fetchMetadata(ip net.IP) (*MetadataResponse, error) {
	response, err := s.sendRequest(ip, "meta-data")
	if err != nil {
		return nil, err
	}

	var metadata MetadataResponse
	if err = yaml.Unmarshal(response, &metadata); err != nil {
		return nil, fmt.Errorf("failed to parse yaml from response: %w", err)
	}

	return &metadata, nil
}

func (s *Source) fetchNetworkConfig(ip net.IP) (*NetworkConfigResponse, error) {
	response, err := s.sendRequest(ip, "network-config")
	if err != nil {
		return nil, err
	}

	var networkConfig NetworkConfigResponse
	if err = yaml.Unmarshal(response, &networkConfig); err != nil {
		return nil, fmt.Errorf("failed to parse yaml from response: %w", err)
	}

	return &networkConfig, nil
}

func (s *Source) fetchUserData(ip net.IP) ([]byte, error) {
	userData, err := s.sendRequest(ip, "user-data")
	if err != nil {
		return nil, err
	}

	return userData, nil
}

func (s *Source) sendRequest(ip net.IP, url string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/%s", strings.TrimSuffix(s.endpoint, "/"), strings.TrimPrefix(url, "/")), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Add("X-Real-IP", ip.String())
	req.Header.Add("X-Forwarded-For", ip.String())

	httpClient := http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("invalid http status code returned: %d", resp.StatusCode)
	}

	responseBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return responseBytes, nil
}

func init() {
	sources.Register(Type, &Source{
		cfg: nil,
	})
}
