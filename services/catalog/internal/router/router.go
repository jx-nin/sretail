package router

import (
	"github.com/gin-gonic/gin"
	"github.com/jx-nin/sretail/internal/httpserver"
	"github.com/jx-nin/sretail/services/catalog/internal/handler"
	"github.com/jx-nin/sretail/services/catalog/internal/product"
)

func New(store product.Store) *gin.Engine {
	r := httpserver.NewEngine()
	products := handler.NewProduct(store)

	r.GET("/products", products.List)
	r.GET("/products/:id", products.GetByID)

	return r
}
