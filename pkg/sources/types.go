package sources

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
	"sort"
	"strings"
)

type Metadata struct {
	ID               string            `json:"id"`
	AvailabilityZone string            `json:"availability_zone"`
	UserData         []byte            `json:"user_data"`
	PublicKeys       map[string]string `json:"public_keys"`
	Password         string            `json:"password"`
}

type Routes map[string]render.Render

type directoryIndex map[string]int

var openStackVersions = []string{"2021-03-23", "2009-04-04", "2021-03-23", "latest"}

func (r Routes) GetIndex(c *gin.Context) render.Render {
	index := make(directoryIndex)

	requestURL := strings.TrimSuffix(c.Request.URL.Path, "/") + "/"

	for url := range r {
		// Skip all non-matching registered routes
		if !strings.HasPrefix(url, requestURL) {
			continue
		}

		// Trim prefix and split path
		path := strings.TrimPrefix(url, requestURL)
		pathParts := strings.Split(path, "/")

		// next item
		nextItem := pathParts[0]

		itemCount, ok := index[nextItem]
		if !ok {
			index[nextItem] = 0
		}

		if itemCount < len(pathParts) {
			index[nextItem] = len(pathParts)
		}
	}

	fmt.Println(index)

	results := make([]string, 0)
	for url, levelCount := range index {
		item := url
		if levelCount > 1 {
			item += "/"
		}

		results = append(results, item)
	}

	if len(results) > 0 {
		sort.Slice(results, func(i, j int) bool {
			return results[i] > results[j]
		})

		return render.String{Format: strings.Join(results, "\n")}
	}

	return nil
}

// GetRoutes returns a map of all routes that are available for the given metadata
func (m Metadata) GetRoutes() Routes {
	routes := make(Routes)
	if m.ID != "" {
		routes.RegisterVersionedOpenStackMetadataRoute("/instance-id", render.String{Format: m.ID})
	}

	if len(m.PublicKeys) > 0 {
		for id, publicKey := range m.PublicKeys {
			routes.RegisterVersionedOpenStackMetadataRoute(fmt.Sprintf("/public-keys/%s", id), render.String{Format: publicKey})
		}
	}

	if m.UserData != nil {
		routes.RegisterOpenStackRoute("/user_data", render.Data{Data: m.UserData})
	}

	if m.Password != "" {
		routes.RegisterOpenStackRoute("/password", render.String{Format: m.Password})
	}

	return routes
}

// RegisterVersionedOpenStackMetadataRoute registers a OpenStack meta-data route for the given path and render
func (r Routes) RegisterVersionedOpenStackMetadataRoute(suffix string, g render.Render) {
	for _, version := range openStackVersions {
		url := fmt.Sprintf("/%s/meta-data/%s", version, strings.TrimPrefix(suffix, "/"))
		r[url] = g
	}
}

// RegisterOpenStackRoute registers a OpenStack route for the given path and render
func (r Routes) RegisterOpenStackRoute(suffix string, g render.Render) {
	for _, version := range openStackVersions {
		url := fmt.Sprintf("/openstack/%s/%s", version, strings.TrimPrefix(suffix, "/"))
		r[url] = g
	}
}
