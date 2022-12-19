package sources

import "net"

type Source interface {
	Type() string
	Initialize(cfg map[string]interface{}) error
	GetMetadata(ip net.IP) (*Metadata, error)
}
