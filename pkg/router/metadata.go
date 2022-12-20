package router

import (
	"github.com/g-portal/metadata-server/pkg/sources"
	"github.com/gin-gonic/gin"
)

func MetadataRequest(c *gin.Context) {
	// Get first matching metadata response from registered sources.
	metadata, err := sources.GetMetadata(c.Request)
	if err != nil {
		c.Error(err)
		return
	}

	// Check if we found a matching metadata response by URL, if yes we render that.
	routes := metadata.GetRoutes()
	if routes[c.Request.URL.Path] != nil {
		c.Render(200, routes[c.Request.URL.Path])
		return
	}

	// If nothing found, we check if we find any directory listing for the requested URL.
	if res := routes.GetIndex(c); res != nil {
		c.Render(200, res)
		return
	}

	NotFoundRequest(c)
}
