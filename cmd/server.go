package main

import (
	"github.com/g-portal/metadata-server/pkg/config"
	"github.com/g-portal/metadata-server/pkg/router"
	"github.com/g-portal/metadata-server/pkg/sources"
	_ "github.com/g-portal/metadata-server/pkg/sources/gpcore"
	_ "github.com/g-portal/metadata-server/pkg/sources/nocloud"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	if err := config.ReloadConfig(); err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	cfg := config.GetConfig()
	datasourceList, err := sources.Load()
	if err != nil {
		log.Fatalf("Failed to load datasource: %v", err)
	}

	if !cfg.Debug {
		gin.SetMode(gin.ReleaseMode)
	}

	engine := gin.Default()
	engine.Use(gin.Recovery())
	engine.Use(func(ctx *gin.Context) {
		ctx.Set("datasources", datasourceList)
	})
	engine.Use(func(ctx *gin.Context) {
		if len(cfg.Whitelist) == 0 ||
			cfg.Whitelist.Contains(sources.GetServer(ctx.Request)) {
			ctx.Next()
		} else {
			router.ForbiddenRequest(ctx)
			ctx.Abort()
		}
	})

	router.LoadRoutes(engine)
	if err = engine.Run(cfg.Listen); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
