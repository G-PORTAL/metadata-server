package sources

import (
	"net/http"
	"strings"
)

type MetadataClient string

const (
	MetadataClientUnknown   MetadataClient = "unknown"
	MetadataClientCloudInit MetadataClient = "cloud-init"
)

// getMetadataClient extracts the metadata client from the request headers User-Agent.
func getMetadataClient(r *http.Request) MetadataClient {
	if r.Header.Get("User-Agent") == "" {
		return MetadataClientUnknown
	}

	if strings.Contains(r.Header.Get("User-Agent"), "Cloud-Init") {
		return MetadataClientCloudInit
	}

	return MetadataClientUnknown
}
