package router

import (
	"github.com/g-portal/metadata-server/pkg/sources"
	"github.com/gin-gonic/gin"
	"net/http"
)

// LoadRoutes Registers all routes on *gin.Engine inclusive the metadata catch-all rute.
func LoadRoutes(r *gin.Engine) {
	r.GET("/metrics")
	r.GET("/healthz", func(c *gin.Context) {
		if sourceList, ok := c.MustGet("datasources").([]sources.Source); ok {
			sources := make([]string, 0)
			for _, source := range sourceList {
				sources = append(sources, source.Type())
			}

			c.JSON(http.StatusOK, gin.H{
				"datasources": sources,
			})
		} else {
			_ = c.Error(sources.ErrNoDatasourceFound)
		}
	})

	r.NoRoute(MetadataRequest)
}
