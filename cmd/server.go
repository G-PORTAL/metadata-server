package main

import (
	"github.com/g-portal/metadata-server/pkg/config"
	"github.com/g-portal/metadata-server/pkg/router"
	"github.com/g-portal/metadata-server/pkg/sources"
	_ "github.com/g-portal/metadata-server/pkg/sources/gpcloud"
	"github.com/gin-gonic/gin"

	"log"
)

func main() {
	if err := config.ReloadConfig(); err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	datasourceList, err := sources.Load()
	if err != nil {
		log.Fatalf("Failed to load datasource: %v", err)
	}

	engine := gin.Default()
	engine.Use(gin.Recovery())
	engine.Use(func(ctx *gin.Context) {
		ctx.Set("datasources", datasourceList)
	})

	router.LoadRoutes(engine)
	if err = engine.Run(config.GetConfig().Listen); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
