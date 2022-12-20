package gpcloud

import (
	"context"
	"crypto/tls"
	"fmt"
	grpcclient "github.com/g-portal/metadata-server/pkg/grpc"
	"github.com/g-portal/metadata-server/pkg/keycloak"
	metadatav1 "github.com/g-portal/metadata-server/pkg/proto/gpcloud/api/metadata/v1"
	"github.com/g-portal/metadata-server/pkg/sources"
	"golang.org/x/crypto/ssh"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"log"
	"net"
	"net/http"
)

const Type = "gpcloud"

type Source struct {
	cfg        sources.SourceConfig
	grpcClient *grpc.ClientConn
}

func (s *Source) Type() string {
	return Type
}

func (s *Source) Initialize(cfg sources.SourceConfig) error {
	s.cfg = cfg

	clientAuth, err := keycloak.NewClientAuthenticationService(cfg.GetString("auth_url"), cfg.GetString("realm"), cfg.GetString("client_id"), cfg.GetString("client_secret"))
	if err != nil {
		return err
	}

	if s.grpcClient, err = grpc.Dial(cfg.GetString("grpc_host"),
		grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{})),
		grpc.WithPerRPCCredentials(grpcclient.KeycloakClientAuthenticationAuth{
			Service: clientAuth,
		})); err != nil {
		return err
	}

	return nil
}

func (s *Source) GetMetadataClient() metadatav1.MetadataServiceClient {
	return metadatav1.NewMetadataServiceClient(s.grpcClient)
}

func (s *Source) GetMetadata(r *http.Request) (*sources.Metadata, error) {
	remoteIP, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return nil, err
	}

	if r.Header.Get("X-Forwarded-For") != "" {
		remoteIP = r.Header.Get("X-Forwarded-For")
	}

	log.Printf("Remote IP: %s", remoteIP)
	resp, err := s.GetMetadataClient().GetMetadata(context.Background(), &metadatav1.GetMetadataRequest{
		IpAddress: remoteIP,
	})

	if err != nil {
		return nil, err
	}

	sshKeys := make(map[string]ssh.PublicKey)
	for _, key := range resp.Metadata.SshKeys {
		sshPublicKey, _, _, _, err := ssh.ParseAuthorizedKey([]byte(key.PublicKey))
		if err != nil {
			log.Printf("Failed to parse SSH key: %s", err)
			continue
		}

		sshKeys[key.Id] = sshPublicKey
	}

	nicList := make([]sources.MetadataInterface, 0)
	routeList := make([]sources.MetadataRoute, 0)
	for i, networkInterface := range resp.Metadata.Interfaces {
		var gateway *net.IP
		if networkInterface.Ipv4.Gateway != nil {
			gw := net.ParseIP(*networkInterface.Ipv4.Gateway)
			gateway = &gw
		}

		_, net, err := net.ParseCIDR(fmt.Sprintf("%s/%v", networkInterface.Ipv4.IpAddress, networkInterface.Ipv4.Prefix))
		if err != nil {
			return nil, err
		}

		subnets := make([]sources.MetadataSubnet, 0)
		subnets = append(subnets, sources.MetadataSubnet{
			IPv4:       true,
			Type:       sources.MetadataSubnetTypeStatic,
			Address:    net,
			Gateway:    gateway,
			DnsServers: resp.Metadata.Dns.Nameservers,
		})

		routeList = append(routeList, sources.MetadataRoute{})
		nicList = append(nicList, sources.MetadataInterface{
			MacAddress: networkInterface.MacAddress,
			Name:       fmt.Sprintf("eth%d", i),
			Type:       sources.InterfaceTypePhysical,
			Subnets:    subnets,
		})
	}

	return &sources.Metadata{
		InstanceID:       resp.Metadata.InstanceId,
		InstanceType:     resp.Metadata.Flavour,
		PublicHostname:   resp.Metadata.Hostname,
		LocalHostname:    resp.Metadata.Hostname,
		AvailabilityZone: &resp.Metadata.AvailabilityZone,
		UserData:         resp.Metadata.UserData,
		VendorData:       resp.Metadata.VendorData,
		VendorData2:      resp.Metadata.VendorData_2,
		Password:         resp.Metadata.Password,
		PublicKeys:       sshKeys,
		Interfaces:       nicList,
		Routes:           routeList,
	}, nil
}

func init() {
	sources.Register(Type, &Source{})
}
