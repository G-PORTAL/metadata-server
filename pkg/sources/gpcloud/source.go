package gpcloud

import (
	"context"
	"crypto/tls"
	grpcclient "github.com/g-portal/metadata-server/pkg/grpc"
	"github.com/g-portal/metadata-server/pkg/keycloak"
	metadatav1 "github.com/g-portal/metadata-server/pkg/proto/gpcloud/api/metadata/v1"
	"github.com/g-portal/metadata-server/pkg/sources"
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

	log.Printf("Remote IP: %s", remoteIP)
	resp, err := s.GetMetadataClient().GetMetadata(context.Background(), &metadatav1.GetMetadataRequest{
		IpAddress: "176.57.191.140",
	})

	if err != nil {
		return nil, err
	}

	sshKeys := make(map[string]string)
	for _, key := range resp.Metadata.SshKeys {
		sshKeys[key.Id] = key.PublicKey
	}

	return &sources.Metadata{
		ID:               resp.Metadata.InstanceId,
		AvailabilityZone: resp.Metadata.Region,
		UserData:         []byte(resp.Metadata.UserData),
		PublicKeys:       sshKeys,
	}, nil
}

func init() {
	sources.Register(Type, &Source{})
}
