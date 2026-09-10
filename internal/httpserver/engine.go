package httpserver

import (
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
)

func NewEngine() *gin.Engine {
	r := gin.New()

	r.Use(
		requestid.New(),
		gin.Logger(),
		gin.Recovery(),
	)

	return r
}
