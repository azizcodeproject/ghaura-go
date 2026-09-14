package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type createCustomerRequest struct {
	Name      string `json:"name"`
	NPWP      string `json:"npwp"`
	CreatedBy string `json:"createdBy"`
}

func (h *Handler) CreateCustomer(c *gin.Context) {
	var req createCustomerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeError(c, http.StatusBadRequest, "body JSON tidak valid")
		return
	}

	name := strings.TrimSpace(req.Name)
	if name == "" {
		writeError(c, http.StatusBadRequest, "nama pelanggan wajib diisi")
		return
	}

	customer, err := h.customers.Create(c.Request.Context(), name, strings.TrimSpace(req.NPWP), actorFromRequest(strings.TrimSpace(req.CreatedBy)))
	if err != nil {
		writeInternalError(c, err)
		return
	}

	c.JSON(http.StatusCreated, customer)
}

func (h *Handler) ListCustomers(c *gin.Context) {
	customers, err := h.customers.List(c.Request.Context())
	if err != nil {
		writeInternalError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"items": customers})
}
