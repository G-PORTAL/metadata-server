package grpcclient

import (
	"context"
	"github.com/g-portal/metadata-server/pkg/keycloak"
)

// KeycloakClientAuthenticationAuth is a gRPC client authentication method.
type KeycloakClientAuthenticationAuth struct {
	Service *keycloak.ClientAuthenticationService
}

// GetRequestMetadata Append access token to the request metadata.
func (t KeycloakClientAuthenticationAuth) GetRequestMetadata(ctx context.Context, in ...string) (map[string]string, error) {
	token, err := t.Service.GetToken()
	if err != nil {
		return nil, err
	}

	return map[string]string{
		"authorization": "Bearer " + token.AccessToken,
	}, nil
}

// RequireTransportSecurity Indicates whether the credentials requires transport security.
func (KeycloakClientAuthenticationAuth) RequireTransportSecurity() bool {
	return true
}
