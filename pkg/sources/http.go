package sources

import (
	"github.com/g-portal/metadata-server/pkg/config"
	"github.com/g-portal/metadata-server/pkg/metrics"
	"github.com/gin-gonic/gin/render"
	"net"
	"net/http"
)

type Routes map[string]render.Render

// GetMetadata Try to find a valid metadata response by the registered sources.
func GetMetadata(r *http.Request) (*Metadata, error) {
	ip := GetServer(r)
	if ip == nil {
		return nil, ErrFailedGetRemoteAddress
	}

	metrics.MetadataRequests.With(map[string]string{
		"url":    r.URL.Path,
		"source": ip.String(),
	}).Inc()

	//nolint:godox // implement later
	// TODO: Sort sources by priority
	for _, source := range registration {
		if result, err := source.GetMetadata(ip, getMetadataClient(r)); err == nil && result != nil {
			return result, nil
		}
	}

	return nil, ErrNoMatchingMetadata
}

func GetServer(r *http.Request) net.IP {
	ipAddress, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return nil
	}

	if r.Header.Get("X-Forwarded-For") != "" && isLegallyForwarded(ipAddress) {
		ipAddress = r.Header.Get("X-Forwarded-For")
	}

	return net.ParseIP(ipAddress)
}

// isLegallyForwarded Checks if the given IP is allowed to send
// X-Forwarded-For headers.
func isLegallyForwarded(i string) bool {
	c := config.GetConfig()
	ip := net.ParseIP(i)

	return c.ForwardedForWhitelist.Contains(ip)
}
