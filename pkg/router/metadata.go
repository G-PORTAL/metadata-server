package router

import (
	"errors"
	"github.com/g-portal/metadata-server/pkg/sources"
	"github.com/gin-gonic/gin"
	"net/http"
)

func MetadataRequest(c *gin.Context) {
	// Get first matching metadata response from registered sources.
	metadata, err := sources.GetMetadata(c.Request)
	if err != nil {
		if errors.Is(err, sources.ErrNoMatchingMetadata) {
			NotFoundRequest(c)

			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})

		return
	}

	// Check if we found a matching metadata response by URL, if yes we render that.
	routes := metadata.GetRoutes()
	if routes[c.Request.URL.Path] != nil {
		c.Render(http.StatusOK, routes[c.Request.URL.Path])

		return
	}

	// If nothing found, we check if we find any directory listing for the requested URL.
	if res := routes.GetIndex(c); res != nil {
		c.Render(http.StatusOK, res)

		return
	}

	NotFoundRequest(c)
}
