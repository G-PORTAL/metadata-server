package gpcore

import (
	"buf.build/gen/go/gportal/gpcore/grpc/go/gpcore/api/metadata/v1/metadatav1grpc"
	cloudv1 "buf.build/gen/go/gportal/gpcore/protocolbuffers/go/gpcore/api/cloud/v1"
	metadatav1 "buf.build/gen/go/gportal/gpcore/protocolbuffers/go/gpcore/api/metadata/v1"
	"context"
	"crypto/tls"
	"fmt"
	grpcclient "github.com/g-portal/metadata-server/pkg/grpc"
	"github.com/g-portal/metadata-server/pkg/keycloak"
	"github.com/g-portal/metadata-server/pkg/sources"
	"golang.org/x/crypto/ssh"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"log"
	"net"
)

const Type = "gpcore"

type Source struct {
	cfg        sources.SourceConfig
	grpcClient *grpc.ClientConn
}

func (s *Source) Type() string {
	return Type
}

func (s *Source) Initialize(cfg sources.SourceConfig) error {
	s.cfg = cfg

	clientAuth, err := keycloak.NewClientAuthenticationService(
		cfg.GetString("auth_url"),
		cfg.GetString("realm"),
		cfg.GetString("client_id"),
		cfg.GetString("client_secret"),
	)
	if err != nil {
		return fmt.Errorf("failed to create client authentication service: %w", err)
	}

	if s.grpcClient, err = grpc.NewClient(cfg.GetString("grpc_host"),
		grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{MinVersion: tls.VersionTLS13})),
		grpc.WithPerRPCCredentials(grpcclient.KeycloakClientAuthenticationAuth{
			Service: clientAuth,
		})); err != nil {
		return fmt.Errorf("failed to connect to gRPC server %s: %w", cfg.GetString("grpc_host"), err)
	}

	return nil
}

func (s *Source) GetMetadataClient() metadatav1grpc.MetadataServiceClient {
	return metadatav1grpc.NewMetadataServiceClient(s.grpcClient)
}

func (s *Source) GetMetadata(ip net.IP) (*sources.Metadata, error) {
	resp, err := s.GetMetadataClient().GetMetadata(context.Background(), &metadatav1.GetMetadataRequest{
		IpAddress: ip.String(),
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get metadata: %w", err)
	}

	sshKeys := make(map[string]ssh.PublicKey)
	for _, key := range resp.GetMetadata().GetSshKeys() {
		sshPublicKey, _, _, _, err := ssh.ParseAuthorizedKey([]byte(key.GetPublicKey()))
		if err != nil {
			log.Printf("Failed to parse SSH key: %s", err)

			continue
		}

		sshKeys[key.GetId()] = sshPublicKey
	}

	nicList := make([]sources.MetadataInterface, 0)
	routeList := make([]sources.MetadataRoute, 0)
	for i, networkInterface := range resp.GetMetadata().GetInterfaces() {
		var gateway *net.IP
		if networkInterface.Ipv4.Gateway != nil {
			gw := net.ParseIP(networkInterface.GetIpv4().GetGateway())
			gateway = &gw
		}

		ip, network, err := net.ParseCIDR(fmt.Sprintf("%s/%v",
			networkInterface.GetIpv4().GetIpAddress(), networkInterface.GetIpv4().GetPrefix()))
		if err != nil {
			return nil, fmt.Errorf("failed to parse CIDR: %w", err)
		}

		subnets := make([]sources.MetadataSubnet, 0)
		subnets = append(subnets, sources.MetadataSubnet{
			IPv4:       true,
			IPv6:       false,
			Type:       sources.MetadataSubnetTypeStatic,
			Address:    &ip,
			Network:    network,
			Gateway:    gateway,
			DNSServers: resp.GetMetadata().GetDns().GetNameservers(),
		})

		nicList = append(nicList, sources.MetadataInterface{
			MacAddress: networkInterface.GetMacAddress(),
			Name:       fmt.Sprintf("eth%d", i),
			Type:       sources.InterfaceTypePhysical,
			Subnets:    subnets,
			AcceptRA:   nil,
		})
	}

	var username *string
	var password *string
	if resp.GetMetadata().GetUsername() != "" {
		newUsername := resp.GetMetadata().GetUsername()
		username = &newUsername
	}
	if resp.GetMetadata().GetPassword() != "" {
		newPassword := resp.GetMetadata().GetPassword()
		password = &newPassword
	}

	return &sources.Metadata{
		ProjectID:        &resp.Metadata.ProjectId,
		InstanceID:       resp.GetMetadata().GetInstanceId(),
		InstanceType:     resp.GetMetadata().GetFlavour(),
		PublicHostname:   resp.GetMetadata().GetHostname(),
		LocalHostname:    resp.GetMetadata().GetHostname(),
		AvailabilityZone: &resp.Metadata.AvailabilityZone,
		UserData:         resp.GetMetadata().GetUserData(),
		VendorData:       resp.GetMetadata().GetVendorData(),
		VendorData2:      resp.GetMetadata().GetVendorData_2(),
		Username:         username,
		Password:         password,
		PublicKeys:       sshKeys,
		Interfaces:       nicList,
		Routes:           routeList,
	}, nil
}

func (s *Source) ReportLog(message sources.ReportMessage) error {
	level := cloudv1.ServerLogLevelType_SERVER_LOG_LEVEL_TYPE_INFO
	switch message.Level {
	case sources.ReportMessageLevelTypeError:
		level = cloudv1.ServerLogLevelType_SERVER_LOG_LEVEL_TYPE_ERROR
	case sources.ReportMessageLevelTypeWarning:
		level = cloudv1.ServerLogLevelType_SERVER_LOG_LEVEL_TYPE_WARNING
	case sources.ReportMessageLevelTypeInfo:
		break
	}

	report := &cloudv1.MetadataReport{
		IpAddress: message.IP.String(),
		Message:   message.Message,
		Timestamp: grpcclient.TimeToTimestamp(message.Timestamp),
		Level:     level,
	}
	_, err := s.GetMetadataClient().Report(context.Background(), &metadatav1.ReportRequest{
		Report: report,
	})

	if err != nil {
		return fmt.Errorf("failed to report metadata: %w", err)
	}

	return nil
}

func init() {
	sources.Register(Type, &Source{
		cfg:        nil,
		grpcClient: nil,
	})
}
