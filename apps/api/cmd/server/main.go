package main

import (
	"log"

	"github.com/azizcodeproject/ghaura-go/apps/api/internal/config"
	"github.com/azizcodeproject/ghaura-go/apps/api/internal/handler"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()
	if cfg.AppEnv != "local" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery(), handler.CORSMiddleware())
	handler.RegisterRoutes(r)

	log.Printf("ghaura api listening on %s (env=%s)", cfg.HTTPAddr, cfg.AppEnv)
	if err := r.Run(cfg.HTTPAddr); err != nil {
		log.Fatal(err)
	}
}
