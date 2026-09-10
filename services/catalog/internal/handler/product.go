package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jx-nin/sretail/services/catalog/internal/product"
)

type Product struct {
	store product.Store
}

func NewProduct(store product.Store) *Product {
	return &Product{store: store}
}

func (h *Product) List(c *gin.Context) {
	products, err := h.store.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, products)
}

func (h *Product) GetByID(c *gin.Context) {
	found, err := h.store.GetByID(c.Request.Context(), c.Param("id"))
	switch {
	case err == nil:
		c.JSON(http.StatusOK, found)

	case errors.Is(err, product.ErrNotFound):
		c.JSON(http.StatusNotFound, gin.H{
			"error": "product not found",
		})

	default:
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
	}
}
