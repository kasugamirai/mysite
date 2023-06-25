package shop_models_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"testing"
	"xy.com/mysite/models/shop_models"
)

func setupDatabase() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&shop_models.Order{}, &shop_models.OrderItem{}, &shop_models.Product{})
	return db, err
}

func TestCreateOrder(t *testing.T) {
	db, err := setupDatabase()
	require.NoError(t, err)

	order := &shop_models.Order{
		UserID: 1,
		OrderItems: []shop_models.OrderItem{
			{
				ProductID: 1,
				Quantity:  2,
				Price:     10.0,
			},
			{
				ProductID: 2,
				Quantity:  1,
				Price:     5.0,
			},
		},
		TotalCost: 25.0,
	}

	err = shop_models.CreateOrder(db, order)
	require.NoError(t, err)
	assert.NotZero(t, order.ID)
}

func TestGetOrderByID(t *testing.T) {
	db, err := setupDatabase()
	require.NoError(t, err)

	// Create a sample order
	order := &shop_models.Order{
		UserID: 1,
		OrderItems: []shop_models.OrderItem{
			{
				ProductID: 1,
				Quantity:  2,
				Price:     10.0,
			},
			{
				ProductID: 2,
				Quantity:  1,
				Price:     5.0,
			},
		},
		TotalCost: 25.0,
	}

	err = shop_models.CreateOrder(db, order)
	require.NoError(t, err)

	retrievedOrder, err := shop_models.GetOrderByID(db, order.ID)
	require.NoError(t, err)
	assert.NotNil(t, retrievedOrder)
	assert.Equal(t, order.ID, retrievedOrder.ID)
	assert.Equal(t, order.UserID, retrievedOrder.UserID)
	assert.Equal(t, len(order.OrderItems), len(retrievedOrder.OrderItems))
}

func TestGetOrdersByUserID(t *testing.T) {
	db, err := setupDatabase()
	require.NoError(t, err)

	// Create a sample order
	order := &shop_models.Order{
		UserID: 1,
		OrderItems: []shop_models.OrderItem{
			{
				ProductID: 1,
				Quantity:  2,
				Price:     10.0,
			},
			{
				ProductID: 2,
				Quantity:  1,
				Price:     5.0,
			},
		},
		TotalCost: 25.0,
	}

	err = shop_models.CreateOrder(db, order)
	require.NoError(t, err)

	orders, err := shop_models.GetOrdersByUserID(db, order.UserID)
	require.NoError(t, err)
	assert.NotNil(t, orders)
	assert.Len(t, orders, 1)
	assert.Equal(t, order.ID, orders[0].ID)
	assert.Equal(t, order.UserID, orders[0].UserID)
}

func TestUpdateOrder(t *testing.T) {
	db, err := setupDatabase()
	require.NoError(t, err)

	// Create a sample order
	order := &shop_models.Order{
		UserID: 1,
		OrderItems: []shop_models.OrderItem{
			{
				ProductID: 1,
				Quantity:  2,
				Price:     10.0,
			},
			{
				ProductID: 2,
				Quantity:  1,
				Price:     5.0,
			},
		},
		TotalCost: 25.0,
	}

	err = shop_models.CreateOrder(db, order)
	require.NoError(t, err)

	// Update the order
	order.TotalCost = 30.0
	err = shop_models.UpdateOrder(db, order)
	require.NoError(t, err)

	// Retrieve the updated order
	updatedOrder, err := shop_models.GetOrderByID(db, order.ID)
	require.NoError(t, err)

	assert.Equal(t, order.TotalCost, updatedOrder.TotalCost)
}

func TestDeleteOrder(t *testing.T) {
	db, err := setupDatabase()
	require.NoError(t, err)

	// Create a sample order
	order := &shop_models.Order{
		UserID: 1,
		OrderItems: []shop_models.OrderItem{
			{
				ProductID: 1,
				Quantity:  2,
				Price:     10.0,
			},
			{
				ProductID: 2,
				Quantity:  1,
				Price:     5.0,
			},
		},
		TotalCost: 25.0,
	}

	err = shop_models.CreateOrder(db, order)
	require.NoError(t, err)

	// Delete the order
	err = shop_models.DeleteOrder(db, order.ID)
	require.NoError(t, err)

	// Try to retrieve the deleted order
	deletedOrder, err := shop_models.GetOrderByID(db, order.ID)
	assert.Error(t, err)
	assert.Nil(t, deletedOrder)
}

func TestGetOrderItemsByOrderID(t *testing.T) {
	db, err := setupDatabase()
	require.NoError(t, err)

	// Create a sample order
	order := &shop_models.Order{
		UserID: 1,
		OrderItems: []shop_models.OrderItem{
			{
				ProductID: 1,
				Quantity:  2,
				Price:     10.0,
			},
			{
				ProductID: 2,
				Quantity:  1,
				Price:     5.0,
			},
		},
		TotalCost: 25.0,
	}

	err = shop_models.CreateOrder(db, order)
	require.NoError(t, err)

	// Retrieve order items by order ID
	orderItems, err := shop_models.GetOrderItemsByOrderID(db, order.ID)
	require.NoError(t, err)
	assert.NotNil(t, orderItems)
	assert.Len(t, orderItems, 2)

	for i, orderItem := range orderItems {
		assert.Equal(t, order.OrderItems[i].ProductID, orderItem.ProductID)
		assert.Equal(t, order.OrderItems[i].Quantity, orderItem.Quantity)
		assert.Equal(t, order.OrderItems[i].Price, orderItem.Price)
	}
}
