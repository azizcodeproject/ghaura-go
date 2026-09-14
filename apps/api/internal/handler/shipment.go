package handler

import (
	"database/sql"
	"errors"
	"net/http"
	"strings"

	"github.com/azizcodeproject/ghaura-go/apps/api/internal/domain"
	"github.com/gin-gonic/gin"
)

type createShipmentRequest struct {
	CustomerID  string `json:"customerId"`
	Origin      string `json:"origin"`
	Destination string `json:"destination"`
	CreatedBy   string `json:"createdBy"`
}

type updateShipmentStatusRequest struct {
	Status string `json:"status"`
}

type shipmentByResiResponse struct {
	domain.Shipment
	Cache string `json:"cache"`
}

func (h *Handler) CreateShipment(c *gin.Context) {
	var req createShipmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeError(c, http.StatusBadRequest, "body JSON tidak valid")
		return
	}

	customerID := strings.TrimSpace(req.CustomerID)
	origin := strings.TrimSpace(req.Origin)
	destination := strings.TrimSpace(req.Destination)
	if customerID == "" || origin == "" || destination == "" {
		writeError(c, http.StatusBadRequest, "customerId, origin, dan destination wajib diisi")
		return
	}

	exists, err := h.customers.Exists(c.Request.Context(), customerID)
	if err != nil {
		writeInternalError(c, err)
		return
	}
	if !exists {
		writeError(c, http.StatusBadRequest, "pelanggan tidak ditemukan")
		return
	}

	shipment, err := h.shipments.Create(
		c.Request.Context(),
		customerID,
		origin,
		destination,
		actorFromRequest(strings.TrimSpace(req.CreatedBy)),
	)
	if err != nil {
		writeInternalError(c, err)
		return
	}

	c.JSON(http.StatusCreated, shipment)
}

func (h *Handler) ListShipments(c *gin.Context) {
	statusFilter := strings.TrimSpace(c.Query("status"))
	if statusFilter != "" {
		if _, ok := domain.ParseShipmentStatus(statusFilter); !ok {
			writeError(c, http.StatusBadRequest, "status tidak valid")
			return
		}
	}

	shipments, err := h.shipments.List(c.Request.Context(), statusFilter)
	if err != nil {
		writeInternalError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"items": shipments})
}

func (h *Handler) GetShipmentByID(c *gin.Context) {
	shipment, err := h.shipments.GetByID(c.Request.Context(), c.Param("id"))
	if errors.Is(err, sql.ErrNoRows) {
		writeError(c, http.StatusNotFound, "pengiriman tidak ditemukan")
		return
	}
	if err != nil {
		writeInternalError(c, err)
		return
	}

	c.JSON(http.StatusOK, shipment)
}

func (h *Handler) GetShipmentByResi(c *gin.Context) {
	resiNumber := strings.TrimSpace(c.Param("resiNumber"))
	if resiNumber == "" {
		writeError(c, http.StatusBadRequest, "nomor resi wajib diisi")
		return
	}

	ctx := c.Request.Context()
	if cached, ok := h.cache.GetByResi(ctx, resiNumber); ok {
		c.Header("X-Cache", "HIT")
		c.JSON(http.StatusOK, shipmentByResiResponse{Shipment: cached, Cache: "HIT"})
		return
	}

	shipment, err := h.shipments.GetByResi(ctx, resiNumber)
	if errors.Is(err, sql.ErrNoRows) {
		writeError(c, http.StatusNotFound, "resi tidak ditemukan")
		return
	}
	if err != nil {
		writeInternalError(c, err)
		return
	}

	h.cache.SetByResi(ctx, shipment)
	c.Header("X-Cache", "MISS")
	c.JSON(http.StatusOK, shipmentByResiResponse{Shipment: shipment, Cache: "MISS"})
}

func (h *Handler) UpdateShipmentStatus(c *gin.Context) {
	var req updateShipmentStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeError(c, http.StatusBadRequest, "body JSON tidak valid")
		return
	}

	status, ok := domain.ParseShipmentStatus(strings.TrimSpace(req.Status))
	if !ok {
		writeError(c, http.StatusBadRequest, "status harus salah satu dari DRAFT, BOOKED, IN_TRANSIT, DELIVERED, CANCELLED")
		return
	}

	ctx := c.Request.Context()
	shipment, err := h.shipments.UpdateStatus(ctx, c.Param("id"), status)
	if errors.Is(err, sql.ErrNoRows) {
		writeError(c, http.StatusNotFound, "pengiriman tidak ditemukan")
		return
	}
	if err != nil {
		writeInternalError(c, err)
		return
	}

	h.cache.InvalidateByResi(ctx, shipment.ResiNumber)
	c.JSON(http.StatusOK, shipment)
}
