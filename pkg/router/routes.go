package router

import (
	"errors"
	"github.com/g-portal/metadata-server/pkg/sources"
	"github.com/gin-gonic/gin"
	"net/http"
)

func LoadRoutes(r *gin.Engine) {
	r.GET("/healthz", func(c *gin.Context) {
		if ds, ok := c.MustGet("datasource").(sources.Source); ok {
			c.JSON(http.StatusOK, gin.H{
				"datasource": ds.Type(),
			})
		} else {
			_ = c.Error(errors.New("no datasource found"))
		}
	})
}
