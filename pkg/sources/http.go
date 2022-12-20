package sources

import (
	"fmt"
	"github.com/gin-gonic/gin/render"
	"net/http"
)

type Routes map[string]render.Render

// GetMetadata Try to find a valid metadata response by the registered sources
func GetMetadata(r *http.Request) (*Metadata, error) {
	// TODO: Sort sources by priority
	for _, source := range registration {
		if result, err := source.GetMetadata(r); err == nil && result != nil {
			return result, nil
		}
	}

	return nil, fmt.Errorf("no matching metadata found")
}
