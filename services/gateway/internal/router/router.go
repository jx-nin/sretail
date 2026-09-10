package router

import (
	"github.com/gin-gonic/gin"
	"github.com/jx-nin/sretail/internal/httpserver"
	"github.com/jx-nin/sretail/services/gateway/internal/handler"
)

func New() *gin.Engine {
	r := httpserver.NewEngine()

	r.GET("/health", handler.Health)
	r.GET("/health/live", handler.Liveness)
	r.GET("/health/ready", handler.Readiness)

	return r
}
