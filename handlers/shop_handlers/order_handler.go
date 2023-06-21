package shop_handlers

import (
	"net/http"
	"strconv"
	"xy.com/mysite/models/shop_models"

	"github.com/gin-gonic/gin"
	"xy.com/mysite/database"
)

// CreateOrderHandler handles the creation of a new order.
func CreateOrderHandler(c *gin.Context) {
	var order shop_models.Order
	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := shop_models.CreateOrder(database.DB, &order); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, order)
}

// GetAllOrdersHandler handles fetching all orders.
func GetAllOrdersHandler(c *gin.Context) {
	orders, err := shop_models.GetAllOrders(database.DB)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"GetAllOrdersHandler": err.Error()})
		return
	}

	c.JSON(http.StatusOK, orders)
}

// GetOrderByIDHandler handles fetching an order by ID.
func GetOrderByIDHandler(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	order, err := shop_models.GetOrderByID(database.DB, uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, order)
}

// GetOrdersByUserIDHandler handles fetching all orders for a specific user.
func GetOrdersByUserIDHandler(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("userID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	orders, err := shop_models.GetOrdersByUserID(database.DB, uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, orders)
}

// UpdateOrderHandler handles updating an order.
func UpdateOrderHandler(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var order shop_models.Order
	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	order.ID = uint(id)
	if err := shop_models.UpdateOrder(database.DB, &order); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, order)
}

// DeleteOrderHandler handles deleting an order.
func DeleteOrderHandler(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := shop_models.DeleteOrder(database.DB, uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

// GetOrderItemsByOrderIDHandler handles fetching all order items for a specific order.
func GetOrderItemsByOrderIDHandler(c *gin.Context) {
	orderID, err := strconv.Atoi(c.Param("orderID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	orderItems, err := shop_models.GetOrderItemsByOrderID(database.DB, uint(orderID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, orderItems)
}
