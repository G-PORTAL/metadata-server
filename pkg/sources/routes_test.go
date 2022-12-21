package sources_test

import (
	"bytes"
	"github.com/g-portal/metadata-server/pkg/sources"
	"github.com/gin-gonic/gin/render"
	"io"
	"net/http/httptest"
	"testing"
)

func TestGetRoutesWithMockup(t *testing.T) {
	metadata := metadataMockup()
	routes := metadata.GetRoutes()
	if routes == nil {
		t.Fatal("routes should not be nil")
	}

	if len(routes) != 18 {
		t.Fatalf("routes length should not be 18, got %v", len(routes))
	}

	for url, render := range routes {
		switch url {
		case "/2009-04-04/meta-data/instance-id",
			"/2021-03-23/meta-data/instance-id",
			"/latest/meta-data/instance-id":
			expectOutput(url, render, []byte("instance-id"), t)
		case "/2009-04-04/meta-data/local-hostname",
			"/2021-03-23/meta-data/local-hostname",
			"/latest/meta-data/local-hostname":
			expectOutput(url, render, []byte("test"), t)
		case "/2009-04-04/meta-data/public-hostname",
			"/2021-03-23/meta-data/public-hostname",
			"/latest/meta-data/public-hostname":
			expectOutput(url, render, []byte("test.acme.com"), t)
		}
	}

}

func expectOutput(url string, r render.Render, expected []byte, t *testing.T) {
	writer := httptest.NewRecorder()
	r.WriteContentType(writer)
	r.Render(writer)

	b, err := io.ReadAll(writer.Body)
	if err != nil {
		t.Fatal(err)
	}

	if bytes.Compare(b, expected) != 0 {
		t.Fatalf("url %s expected %s, got %s", url, string(expected), string(b))
	}
}

func TestGetRoutesEmpty(t *testing.T) {
	metadata := &sources.Metadata{}
	routes := metadata.GetRoutes()

	if routes == nil {
		t.Fatal("routes should not be nil")
	}
	if len(routes) != 0 {
		t.Fatal("routes length should be zero")
	}
}

func metadataMockup() *sources.Metadata {
	projectID := "project-id"

	return &sources.Metadata{
		InstanceID:       "instance-id",
		ProjectID:        &projectID,
		InstanceType:     "x4.large",
		LocalHostname:    "test",
		PublicHostname:   "test.acme.com",
		AvailabilityZone: nil,
		UserData:         nil,
		VendorData:       nil,
		VendorData2:      nil,
		PublicKeys:       nil,
		Password:         nil,
		Interfaces:       nil,
		Routes:           nil,
	}
}
