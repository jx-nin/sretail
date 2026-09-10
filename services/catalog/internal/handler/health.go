package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"service": "catalog",
		"status":  "ok",
	})
}

func Liveness(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

func Readiness(c *gin.Context) {
	c.Status(http.StatusNoContent)
}
