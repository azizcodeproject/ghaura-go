package main

import (
	"context"
	"log"
	"time"

	"github.com/azizcodeproject/ghaura-go/apps/api/internal/cache"
	"github.com/azizcodeproject/ghaura-go/apps/api/internal/config"
	"github.com/azizcodeproject/ghaura-go/apps/api/internal/db"
	"github.com/azizcodeproject/ghaura-go/apps/api/internal/handler"
	"github.com/azizcodeproject/ghaura-go/apps/api/internal/repository"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()
	if cfg.AppEnv != "local" {
		gin.SetMode(gin.ReleaseMode)
	}

	database, err := db.Open(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("database: %v", err)
	}
	defer database.Close()

	if err := db.RunMigrations(database); err != nil {
		log.Fatalf("migrations: %v", err)
	}

	shipmentCache := cache.NewShipmentCache(cfg.RedisAddr)
	defer shipmentCache.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := shipmentCache.Ping(ctx); err != nil {
		log.Printf("warning: redis belum siap di %s: %v", cfg.RedisAddr, err)
	}

	apiHandler := handler.New(
		repository.NewCustomerRepository(database),
		repository.NewShipmentRepository(database),
		shipmentCache,
	)

	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery(), handler.CORSMiddleware())
	handler.RegisterRoutes(router, apiHandler)

	log.Printf("ghaura api listening on %s (env=%s)", cfg.HTTPAddr, cfg.AppEnv)
	if err := router.Run(cfg.HTTPAddr); err != nil {
		log.Fatal(err)
	}
}
