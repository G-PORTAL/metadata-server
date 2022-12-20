package sources

import (
	"net/http"
)

// Source defines the interface for a metadata source, each source needs to implement this.
type Source interface {
	Type() string
	Initialize(cfg SourceConfig) error
	GetMetadata(r *http.Request) (*Metadata, error)
}
