package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, h *Handler) {
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"service": "ghaura-api",
			"time":    time.Now().UTC(),
		})
	})

	api := r.Group("/api/v1")
	{
		api.GET("/ping", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "pong"})
		})

		api.POST("/customers", h.CreateCustomer)
		api.GET("/customers", h.ListCustomers)

		api.POST("/shipments", h.CreateShipment)
		api.GET("/shipments", h.ListShipments)
		api.GET("/shipments/by-resi/:resiNumber", h.GetShipmentByResi)
		api.GET("/shipments/:id", h.GetShipmentByID)
		api.PATCH("/shipments/:id/status", h.UpdateShipmentStatus)
	}
}
