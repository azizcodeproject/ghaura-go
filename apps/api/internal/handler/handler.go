package handler

import (
	"github.com/azizcodeproject/ghaura-go/apps/api/internal/cache"
	"github.com/azizcodeproject/ghaura-go/apps/api/internal/repository"
)

type Handler struct {
	customers *repository.CustomerRepository
	shipments *repository.ShipmentRepository
	cache     *cache.ShipmentCache
}

func New(
	customers *repository.CustomerRepository,
	shipments *repository.ShipmentRepository,
	shipmentCache *cache.ShipmentCache,
) *Handler {
	return &Handler{
		customers: customers,
		shipments: shipments,
		cache:     shipmentCache,
	}
}

func actorFromRequest(createdBy string) string {
	if createdBy != "" {
		return createdBy
	}
	return "system"
}
