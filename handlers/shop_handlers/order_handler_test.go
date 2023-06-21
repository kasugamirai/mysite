package shop_handlers_test

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"xy.com/mysite/database"
	"xy.com/mysite/handlers/shop_handlers"
	"xy.com/mysite/models/shop_models"
)

func setupTestData() {
	database.InitDB()
}

func setupRouter() *gin.Engine {
	router := gin.Default()
	orderGroup := router.Group("/orders")
	{
		orderGroup.POST("/", shop_handlers.CreateOrderHandler)
		orderGroup.GET("/getall", shop_handlers.GetAllOrdersHandler)
		orderGroup.GET("/:id", shop_handlers.GetOrderByIDHandler)
		orderGroup.GET("/user/:userID", shop_handlers.GetOrdersByUserIDHandler)
		orderGroup.PUT("/:id", shop_handlers.UpdateOrderHandler)
		orderGroup.DELETE("/:id", shop_handlers.DeleteOrderHandler)
		orderGroup.GET("/items/:orderID", shop_handlers.GetOrderItemsByOrderIDHandler)
	}

	return router
}

func TestCreateOrderHandler(t *testing.T) {
	setupTestData()

	newOrder := shop_models.Order{
		UserID:    1,
		TotalCost: 100.0,
	}

	orderJson, _ := json.Marshal(newOrder)
	body := bytes.NewReader(orderJson)

	req, _ := http.NewRequest("POST", "/orders/", body)
	w := httptest.NewRecorder()

	router := setupRouter()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	var createdOrder shop_models.Order
	json.Unmarshal(w.Body.Bytes(), &createdOrder)
	assert.Equal(t, newOrder.UserID, createdOrder.UserID)
	assert.Equal(t, newOrder.TotalCost, createdOrder.TotalCost)

	// Clean up
	database.DB.Delete(&createdOrder)
}

func TestGetAllOrdersHandler(t *testing.T) {
	setupTestData()

	testOrder1 := shop_models.Order{UserID: 1, TotalCost: 100.0}
	testOrder2 := shop_models.Order{UserID: 2, TotalCost: 200.0}
	database.DB.Create(&testOrder1)
	database.DB.Create(&testOrder2)

	req, _ := http.NewRequest("GET", "/orders/getall", nil)
	w := httptest.NewRecorder()

	router := setupRouter()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var orders []shop_models.Order
	json.Unmarshal(w.Body.Bytes(), &orders)

	foundTestOrder1 := false
	foundTestOrder2 := false

	for _, order := range orders {
		if order.ID == testOrder1.ID {
			foundTestOrder1 = true
		}
		if order.ID == testOrder2.ID {
			foundTestOrder2 = true
		}
	}

	assert.True(t, foundTestOrder1)
	assert.True(t, foundTestOrder2)

	// Clean up
	database.DB.Delete(&testOrder1)
	database.DB.Delete(&testOrder2)
}

func TestGetOrderByIDHandler(t *testing.T) {
	setupTestData()

	testOrder := shop_models.Order{UserID: 1, TotalCost: 100.0}
	database.DB.Create(&testOrder)

	req, _ := http.NewRequest("GET", "/orders/"+strconv.Itoa(int(testOrder.ID)), nil)
	w := httptest.NewRecorder()

	router := setupRouter()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var fetchedOrder shop_models.Order
	json.Unmarshal(w.Body.Bytes(), &fetchedOrder)
	assert.Equal(t, testOrder.ID, fetchedOrder.ID)

	// Clean up
	database.DB.Delete(&testOrder)
}

func TestGetOrdersByUserIDHandler(t *testing.T) {
	setupTestData()

	testUser1ID := uint(1)
	testUser2ID := uint(2)
	testOrder1 := shop_models.Order{UserID: testUser1ID, TotalCost: 100.0}
	testOrder2 := shop_models.Order{UserID: testUser1ID, TotalCost: 200.0}
	testOrder3 := shop_models.Order{UserID: testUser2ID, TotalCost: 300.0}

	database.DB.Create(&testOrder1)
	database.DB.Create(&testOrder2)
	database.DB.Create(&testOrder3)

	req, _ := http.NewRequest("GET", "/orders/user/"+strconv.Itoa(int(testUser1ID)), nil)
	w := httptest.NewRecorder()

	router := setupRouter()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var orders []shop_models.Order
	json.Unmarshal(w.Body.Bytes(), &orders)

	assert.Equal(t, 2, len(orders))

	// Clean up
	database.DB.Delete(&testOrder1)
	database.DB.Delete(&testOrder2)
	database.DB.Delete(&testOrder3)
}

func TestUpdateOrderHandler(t *testing.T) {
	setupTestData()

	testOrder := shop_models.Order{UserID: 1, TotalCost: 100.0}
	database.DB.Create(&testOrder)

	updatedOrder := shop_models.Order{
		UserID:    testOrder.UserID,
		TotalCost: 200.0,
	}

	orderJson, _ := json.Marshal(updatedOrder)
	body := bytes.NewReader(orderJson)

	req, _ := http.NewRequest("PUT", "/orders/"+strconv.Itoa(int(testOrder.ID)), body)
	w := httptest.NewRecorder()

	router := setupRouter()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var fetchedOrder shop_models.Order
	json.Unmarshal(w.Body.Bytes(), &fetchedOrder)
	assert.Equal(t, updatedOrder.TotalCost, fetchedOrder.TotalCost)

	// Clean up
	database.DB.Delete(&fetchedOrder)
}

func TestDeleteOrderHandler(t *testing.T) {
	setupTestData()

	testOrder := shop_models.Order{UserID: 1, TotalCost: 100.0}
	database.DB.Create(&testOrder)

	req, _ := http.NewRequest("DELETE", "/orders/"+strconv.Itoa(int(testOrder.ID)), nil)
	w := httptest.NewRecorder()

	router := setupRouter()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var fetchedOrder shop_models.Order
	result := database.DB.First(&fetchedOrder, testOrder.ID)
	assert.Error(t, result.Error)
}

func TestGetOrderItemsByOrderIDHandler(t *testing.T) {
	setupTestData()

	testOrder := shop_models.Order{UserID: 1, TotalCost: 100.0}
	database.DB.Create(&testOrder)

	testProduct := shop_models.Product{Name: "Test Product", Price: 10.0}
	database.DB.Create(&testProduct)

	testOrderItem1 := shop_models.OrderItem{OrderID: testOrder.ID, ProductID: testProduct.ID, Quantity: 2, Price: 10.0}
	testOrderItem2 := shop_models.OrderItem{OrderID: testOrder.ID, ProductID: testProduct.ID, Quantity: 3, Price: 10.0}

	database.DB.Create(&testOrderItem1)
	database.DB.Create(&testOrderItem2)

	req, _ := http.NewRequest("GET", "/orders/items/"+strconv.Itoa(int(testOrder.ID)), nil)
	w := httptest.NewRecorder()

	router := setupRouter()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var orderItems []shop_models.OrderItem
	json.Unmarshal(w.Body.Bytes(), &orderItems)

	assert.Equal(t, 2, len(orderItems))

	// Clean up
	database.DB.Delete(&testOrder)
	database.DB.Delete(&testProduct)
	database.DB.Delete(&testOrderItem1)
	database.DB.Delete(&testOrderItem2)
}
