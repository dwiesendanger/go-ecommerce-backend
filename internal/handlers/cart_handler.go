package handlers

import (
	"ecommerce-platform/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CartHandler struct {
	Service *services.CartService
}

func NewCartHandler(service *services.CartService) *CartHandler {
	return &CartHandler{Service: service}
}

type AddItemInput struct {
	ProductID uint `json:"product_id" binding:"required"`
	Quantity  int  `json:"quantity" binding:"required,min=1"`
}

func (h *CartHandler) AddItem(c *gin.Context) {
	var input AddItemInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	err := h.Service.AddItem(userID.(uint), input.ProductID, input.Quantity)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Item added to cart"})
}

func (h *CartHandler) GetCart(c *gin.Context) {
	userID, _ := c.Get("userID")

	cart, err := h.Service.GetCart(userID.(uint))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"items": []interface{}{}})
		return
	}

	c.JSON(http.StatusOK, gin.H{"cart": cart})
}
