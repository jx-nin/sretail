package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jx-nin/sretail/services/gateway/internal/response"
)

func Health(c *gin.Context) {
	res := response.HealthResponse{
		Service: "gateway",
		Status:  "ok",
	}

	c.JSON(http.StatusOK, res)
}

func Liveness(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

func Readiness(c *gin.Context) {
	c.Status(http.StatusNoContent)
}
