package sources

import (
	"fmt"
	"github.com/gin-gonic/gin/render"
	"strings"
)

var openStackVersions = []string{"2021-03-23", "2009-04-04", "2021-03-23", "latest"}

// RegisterVersionedOpenStackMetadataRoute registers a OpenStack meta-data route for the given path and render
func (r Routes) registerVersionedOpenStackMetadataRoute(suffix string, g render.Render) {
	for _, version := range openStackVersions {
		url := fmt.Sprintf("/%s/meta-data/%s", version, strings.TrimPrefix(suffix, "/"))
		r[url] = g
	}
}

// RegisterOpenStackRoute registers a OpenStack route for the given path and render
func (r Routes) registerOpenStackRoute(suffix string, g render.Render) {
	for _, version := range openStackVersions {
		url := fmt.Sprintf("/openstack/%s/%s", version, strings.TrimPrefix(suffix, "/"))
		r[url] = g
	}
}
