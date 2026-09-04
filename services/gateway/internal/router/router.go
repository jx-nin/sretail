package router

import (
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"github.com/jx-nin/sretail/services/gateway/internal/handler"
)

func New() *gin.Engine {
	r := gin.New()

	r.Use(
		requestid.New(),
		gin.Logger(),
		gin.Recovery(),
	)

	r.GET("/health", handler.Health)
	r.GET("/health/live", handler.Liveness)
	r.GET("/health/ready", handler.Readiness)

	return r
}
