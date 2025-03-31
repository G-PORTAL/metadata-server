package sources

import (
	"github.com/g-portal/metadata-server/pkg/openstack"
	"github.com/gin-gonic/gin/render"
	"golang.org/x/crypto/ssh"
	"net"
	"strings"
)

func (m Metadata) OpenStackNetworkData() render.JSON {
	config := openstack.Interfaces{
		Links:    make([]openstack.Link, 0),
		Networks: nil,
		Services: make([]openstack.Service, 0),
	}

	for _, metadataInterface := range m.Interfaces {
		linkType := openstack.LinkTypePhysical
		if metadataInterface.Type == InterfaceTypeVLAN {
			linkType = openstack.LinkTypeVlan
		}

		link := openstack.Link{
			ID:   metadataInterface.Name,
			Type: linkType,
		}

		if linkType == openstack.LinkTypePhysical {
			link.EthernetMacAddress = metadataInterface.MacAddress
		}

		if linkType == openstack.LinkTypeVlan {
			link.VlanMacAddress = &metadataInterface.MacAddress
			link.VlanID = metadataInterface.VlanID
			link.VlanLink = metadataInterface.VlanLink
		}

		config.Links = append(config.Links, link)

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

			// Add matching routes for the network interface where the gateway is in the subnet
			for _, route := range m.Routes {
				if route.Gateway != nil &&
					subnet.Network.Contains(net.ParseIP(*route.Gateway)) {
					network.Routes = append(network.Routes, openstack.Route{
						Network: route.Network,
						Netmask: route.Netmask,
						Gateway: *route.Gateway,
						Metric:  route.Metric,
					})
				}
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
		//nolint:godox // implement later
		RandomSeed: m.InstanceID, // TODO: do something
	}

	for keyID, key := range m.PublicKeys {
		sshKey := strings.TrimSpace(string(ssh.MarshalAuthorizedKey(key)))
		metadata.Keys = append(metadata.Keys, openstack.MetadataKeyDefinition{
			Name: keyID,
			Type: "ssh",
			Data: sshKey,
		})

		metadata.PublicKeys[keyID] = sshKey
	}

	if m.Username != nil {
		metadata.AdminUsername = m.Username
	}

	if m.Password != nil {
		metadata.AdminPassword = m.Password
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
