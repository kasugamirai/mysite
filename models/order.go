package models

import (
	"gorm.io/gorm"
)

// Order represents an order entity in the system.
type Order struct {
	gorm.Model
	UserID     uint        `json:"user_id" gorm:"index:idx_user_OrderItems"`
	OrderItems []OrderItem `json:"order_items" gorm:"foreignKey:OrderID"`
	TotalCost  float64     `json:"total_cost"`
}

type OrderItem struct {
	gorm.Model
	OrderID   uint `json:"-"`
	ProductID uint
	Product   Product `json:"-"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
}

// CreateOrder creates a new order in the database.
func CreateOrder(db *gorm.DB, order *Order) error {
	return db.Create(order).Error
}

// GetOrderByID retrieves an order from the database by ID.
func GetOrderByID(db *gorm.DB, id uint) (*Order, error) {
	var order Order
	err := db.Preload("OrderItems.Product").Where("id = ?", id).First(&order).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

// GetOrdersByUserID retrieves all orders for a specific user from the database.
func GetOrdersByUserID(db *gorm.DB, userID uint) ([]Order, error) {
	var orders []Order
	err := db.Preload("OrderItems.Product").Where("user_id = ?", userID).Find(&orders).Error
	if err != nil {
		return nil, err
	}
	return orders, nil
}

// UpdateOrder updates the order data in the database.
func UpdateOrder(db *gorm.DB, order *Order) error {
	return db.Save(order).Error
}

// DeleteOrder deletes an order from the database.
func DeleteOrder(db *gorm.DB, id uint) error {
	return db.Delete(&Order{}, id).Error
}

// GetOrderItemsByOrderID retrieves all order items for a specific order from the database.
func GetOrderItemsByOrderID(db *gorm.DB, orderID uint) ([]OrderItem, error) {
	var orderItems []OrderItem
	err := db.Preload("Product").Where("order_id = ?", orderID).Find(&orderItems).Error
	if err != nil {
		return nil, err
	}
	return orderItems, nil
}
