package keycloak

import (
	"context"
	"fmt"
	"github.com/Nerzal/gocloak/v12"
	"log"
	"sync"
	"time"
)

const RefreshGracePeriod = 10

var mutex sync.RWMutex

// ClientAuthenticationService is a utility tool for service tokens and handles refreshing the token.
type ClientAuthenticationService struct {
	client *gocloak.GoCloak
	jwt    *gocloak.JWT
	config struct {
		clientID     string
		clientSecret string
		realm        string
	}
	expireTime time.Time
}

// NewClientAuthenticationService returns utility tool for service tokens and handles refreshing the token.
func NewClientAuthenticationService(url, realm, clientID, clientSecret string) (*ClientAuthenticationService, error) {
	service := &ClientAuthenticationService{
		client: gocloak.NewClient(url),
		config: struct {
			clientID     string
			clientSecret string
			realm        string
		}{clientID: clientID, clientSecret: clientSecret, realm: realm},
	}

	// Initial refresh
	if err := service.refresh(context.Background()); err != nil {
		return nil, err
	}

	return service, nil
}

// refresh utility function to generate new JWT tokens.
func (s *ClientAuthenticationService) refresh(ctx context.Context) error {
	mutex.Lock()
	defer mutex.Unlock()

	log.Printf("Refreshing Keycloak Token now.")
	var err error
	if s.jwt, err = s.client.GetToken(ctx, s.config.realm, gocloak.TokenOptions{
		ClientID:     &s.config.clientID,
		ClientSecret: &s.config.clientSecret,
		GrantType:    gocloak.StringP("client_credentials"),
	}); err != nil {
		return fmt.Errorf("failed to get token: %w", err)
	}
	s.expireTime = time.Now().Add(time.Second * time.Duration(s.jwt.ExpiresIn-RefreshGracePeriod))

	return nil
}

// GetToken returns the always-up-to-date jwt token.
func (s *ClientAuthenticationService) GetToken(ctx context.Context) (*gocloak.JWT, error) {
	if s.jwt == nil || s.expireTime.Before(time.Now()) {
		if err := s.refresh(ctx); err != nil {
			return nil, err
		}
	}

	return s.jwt, nil
}

// GetClient returns the gocloak client in case keycloak api requests should be made.
func (s *ClientAuthenticationService) GetClient() *gocloak.GoCloak {
	return s.client
}
