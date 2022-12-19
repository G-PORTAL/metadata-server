package gpcloud

import (
	"errors"
	"github.com/g-portal/metadata-server/pkg/sources"
	"net"
)

const Type = "gpcloud"

type Source struct {
	cfg map[string]interface{}
}

func (s *Source) Type() string {
	return Type
}

func (s *Source) Initialize(cfg map[string]interface{}) error {
	s.cfg = cfg
	return nil
}

func (s *Source) GetMetadata(ip net.IP) (*sources.Metadata, error) {
	return nil, errors.New("not implemented")
}

func init() {
	sources.Register(Type, &Source{})
}
