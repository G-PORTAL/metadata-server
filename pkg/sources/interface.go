package sources

import (
	"net/http"
)

type Source interface {
	Type() string
	Initialize(cfg SourceConfig) error
	GetMetadata(r *http.Request) (*Metadata, error)
}
