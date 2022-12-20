package sources

import (
	"github.com/g-portal/metadata-server/pkg/openstack"
	"github.com/gin-gonic/gin/render"
	"golang.org/x/crypto/ssh"
	"strings"
)

func (m Metadata) OpenStackNetworkData() render.JSON {
	config := openstack.Interfaces{
		Links:    make([]openstack.Link, 0),
		Networks: nil,
		Services: make([]openstack.Service, 0),
	}

	for _, metadataInterface := range m.Interfaces {
		config.Links = append(config.Links, openstack.Link{
			ID:                 metadataInterface.Name,
			Type:               openstack.LinkTypePhysical,
			EthernetMacAddress: metadataInterface.MacAddress,
		})

		for _, subnet := range metadataInterface.Subnets {
			for _, dnsServer := range subnet.DNSServers {
				config.Services = append(config.Services, openstack.Service{
					Type:    openstack.ServiceTypeDNS,
					Address: dnsServer,
				})
			}

			network := openstack.Network{
				Link: metadataInterface.Name,
				Type: openstack.NetworkTypeIPv4,
			}

			if subnet.Address != nil && subnet.Network != nil {
				ip := subnet.Address.String()
				cidr := strings.Split(subnet.Network.String(), "/")
				ip += "/" + cidr[1]
				network.IPAddress = &ip
			}

			if subnet.Gateway != nil {
				ip := subnet.Gateway.String()
				network.Routes = append(network.Routes, openstack.Route{
					Network: "0.0.0.0",
					Netmask: "0.0.0.0",
					Gateway: ip,
				})
			}

			config.Networks = append(config.Networks, network)
		}
	}

	return render.JSON{
		Data: config,
	}
}

func (m Metadata) OpenStackVendorData(data []byte) render.JSON {
	return render.JSON{
		Data: openstack.VendorData{
			CloudInit: string(data),
		},
	}
}

func (m Metadata) OpenStackMetaData() render.JSON {
	metadata := openstack.Metadata{
		UUID:        m.InstanceID,
		Keys:        make([]openstack.MetadataKeyDefinition, 0),
		PublicKeys:  map[string]string{},
		Hostname:    m.PublicHostname,
		Name:        m.LocalHostname,
		LaunchIndex: 0,
		RandomSeed:  m.InstanceID, // TODO: do something
	}

	for keyID, key := range m.PublicKeys {
		metadata.Keys = append(metadata.Keys, openstack.MetadataKeyDefinition{
			Name: keyID,
			Type: "ssh",
			Data: string(ssh.MarshalAuthorizedKey(key)),
		})

		metadata.PublicKeys[keyID] = string(ssh.MarshalAuthorizedKey(key))
	}

	if m.ProjectID != nil {
		metadata.ProjectID = *m.ProjectID
	}

	if m.AvailabilityZone != nil {
		metadata.AvailabilityZone = *m.AvailabilityZone
	}

	return render.JSON{
		Data: metadata,
	}
}
