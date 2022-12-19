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

	datasource, err := sources.Load()
	if err != nil {
		log.Fatalf("Failed to load datasource: %v", err)
	}

	r := gin.Default()
	r.Use(gin.Recovery())
	r.Use(func(ctx *gin.Context) {
		ctx.Set("datasource", datasource)
	})

	router.LoadRoutes(r)
	err = r.Run(config.GetConfig().Listen)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
