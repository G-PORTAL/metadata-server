package keycloak

import (
	"context"
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
	if err := service.refresh(); err != nil {
		return nil, err
	}

	return service, nil
}

// refresh utility function to generate new JWT tokens (normal refresh logic doesn't apply to this authentication mechanism).
func (s *ClientAuthenticationService) refresh() error {
	mutex.Lock()
	defer mutex.Unlock()

	log.Printf("Refreshing Keycloak Token now.")
	jwt, err := s.client.GetToken(context.Background(), s.config.realm, gocloak.TokenOptions{
		ClientID:     &s.config.clientID,
		ClientSecret: &s.config.clientSecret,
		GrantType:    gocloak.StringP("client_credentials"),
		Audience:     gocloak.StringP("dasdsdadasd"),
	})
	if err != nil {
		return err
	}
	s.expireTime = time.Now().Add(time.Second * time.Duration(jwt.ExpiresIn-RefreshGracePeriod))
	s.jwt = jwt
	return nil
}

// GetToken returns the always-up-to-date jwt token.
func (s *ClientAuthenticationService) GetToken() (*gocloak.JWT, error) {
	if s.jwt == nil || s.expireTime.Before(time.Now()) {
		if err := s.refresh(); err != nil {
			return nil, err
		}
	}

	return s.jwt, nil
}

// GetClient returns the gocloak client in case keycloak api requests should be made.
func (s *ClientAuthenticationService) GetClient() *gocloak.GoCloak {
	return s.client
}
