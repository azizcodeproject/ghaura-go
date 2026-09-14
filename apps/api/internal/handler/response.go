package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func writeError(c *gin.Context, status int, message string) {
	c.JSON(status, gin.H{"error": message})
}

func writeInternalError(c *gin.Context, err error) {
	c.Error(err)
	writeError(c, http.StatusInternalServerError, "terjadi kesalahan pada server")
}
