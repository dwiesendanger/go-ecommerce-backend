package handlers

import (
	"ecommerce-platform/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	Service *services.CheckoutService
}

func NewOrderHandler(service *services.CheckoutService) *OrderHandler {
	return &OrderHandler{Service: service}
}

// CreateOrder
func (h *OrderHandler) CreateOrder(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	order, err := h.Service.PlaceOrder(userID.(uint))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":      "Order placed successfully",
		"order_number": order.OrderNumber,
		"total":        order.TotalAmount,
	})
}

// GetOrders (History)
func (h *OrderHandler) GetOrders(c *gin.Context) {
	userID, _ := c.Get("userID")

	orders, err := h.Service.GetOrdersByUser(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch orders"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"orders": orders})
}
