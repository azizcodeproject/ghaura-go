package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/azizcodeproject/ghaura-go/apps/api/internal/domain"
	"github.com/redis/go-redis/v9"
)

const (
	shipmentResiTTL = 5 * time.Minute
	resiKeyPrefix   = "shipment:resi:"
)

type ShipmentCache struct {
	client *redis.Client
}

func NewShipmentCache(redisAddr string) *ShipmentCache {
	client := redis.NewClient(&redis.Options{
		Addr:         redisAddr,
		DialTimeout:  2 * time.Second,
		ReadTimeout:  2 * time.Second,
		WriteTimeout: 2 * time.Second,
	})
	return &ShipmentCache{client: client}
}

func (c *ShipmentCache) Ping(ctx context.Context) error {
	return c.client.Ping(ctx).Err()
}

func (c *ShipmentCache) GetByResi(ctx context.Context, resiNumber string) (domain.Shipment, bool) {
	payload, err := c.client.Get(ctx, resiCacheKey(resiNumber)).Bytes()
	if err != nil {
		if err != redis.Nil {
			log.Printf("redis get %s: %v", resiCacheKey(resiNumber), err)
		}
		return domain.Shipment{}, false
	}

	var shipment domain.Shipment
	if err := json.Unmarshal(payload, &shipment); err != nil {
		log.Printf("redis unmarshal %s: %v", resiCacheKey(resiNumber), err)
		return domain.Shipment{}, false
	}
	return shipment, true
}

func (c *ShipmentCache) SetByResi(ctx context.Context, shipment domain.Shipment) {
	payload, err := json.Marshal(shipment)
	if err != nil {
		log.Printf("redis marshal shipment %s: %v", shipment.ResiNumber, err)
		return
	}
	if err := c.client.Set(ctx, resiCacheKey(shipment.ResiNumber), payload, shipmentResiTTL).Err(); err != nil {
		log.Printf("redis set %s: %v", resiCacheKey(shipment.ResiNumber), err)
	}
}

func (c *ShipmentCache) InvalidateByResi(ctx context.Context, resiNumber string) {
	if err := c.client.Del(ctx, resiCacheKey(resiNumber)).Err(); err != nil {
		log.Printf("redis del %s: %v", resiCacheKey(resiNumber), err)
	}
}

func (c *ShipmentCache) Close() error {
	return c.client.Close()
}

func resiCacheKey(resiNumber string) string {
	return fmt.Sprintf("%s%s", resiKeyPrefix, resiNumber)
}
