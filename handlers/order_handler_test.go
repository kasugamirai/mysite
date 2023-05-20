package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"xy.com/mysite/database"
	"xy.com/mysite/handlers"
	"xy.com/mysite/models"
)

func setupTestData() {
	database.InitDB()
}

func setupRouter() *gin.Engine {
	router := gin.Default()
	orderGroup := router.Group("/orders")
	{
		orderGroup.POST("/", handlers.CreateOrderHandler)
		orderGroup.GET("/getall", handlers.GetAllOrdersHandler)
		orderGroup.GET("/:id", handlers.GetOrderByIDHandler)
		orderGroup.GET("/user/:userID", handlers.GetOrdersByUserIDHandler)
		orderGroup.PUT("/:id", handlers.UpdateOrderHandler)
		orderGroup.DELETE("/:id", handlers.DeleteOrderHandler)
		orderGroup.GET("/items/:orderID", handlers.GetOrderItemsByOrderIDHandler)
	}

	return router
}

func TestCreateOrderHandler(t *testing.T) {
	setupTestData()

	newOrder := models.Order{
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
	var createdOrder models.Order
	json.Unmarshal(w.Body.Bytes(), &createdOrder)
	assert.Equal(t, newOrder.UserID, createdOrder.UserID)
	assert.Equal(t, newOrder.TotalCost, createdOrder.TotalCost)

	// Clean up
	database.DB.Delete(&createdOrder)
}

func TestGetAllOrdersHandler(t *testing.T) {
	setupTestData()

	testOrder1 := models.Order{UserID: 1, TotalCost: 100.0}
	testOrder2 := models.Order{UserID: 2, TotalCost: 200.0}
	database.DB.Create(&testOrder1)
	database.DB.Create(&testOrder2)

	req, _ := http.NewRequest("GET", "/orders/getall", nil)
	w := httptest.NewRecorder()

	router := setupRouter()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var orders []models.Order
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

	testOrder := models.Order{UserID: 1, TotalCost: 100.0}
	database.DB.Create(&testOrder)

	req, _ := http.NewRequest("GET", "/orders/"+strconv.Itoa(int(testOrder.ID)), nil)
	w := httptest.NewRecorder()

	router := setupRouter()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var fetchedOrder models.Order
	json.Unmarshal(w.Body.Bytes(), &fetchedOrder)
	assert.Equal(t, testOrder.ID, fetchedOrder.ID)

	// Clean up
	database.DB.Delete(&testOrder)
}

func TestGetOrdersByUserIDHandler(t *testing.T) {
	setupTestData()

	testUser1ID := uint(1)
	testUser2ID := uint(2)
	testOrder1 := models.Order{UserID: testUser1ID, TotalCost: 100.0}
	testOrder2 := models.Order{UserID: testUser1ID, TotalCost: 200.0}
	testOrder3 := models.Order{UserID: testUser2ID, TotalCost: 300.0}

	database.DB.Create(&testOrder1)
	database.DB.Create(&testOrder2)
	database.DB.Create(&testOrder3)

	req, _ := http.NewRequest("GET", "/orders/user/"+strconv.Itoa(int(testUser1ID)), nil)
	w := httptest.NewRecorder()

	router := setupRouter()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var orders []models.Order
	json.Unmarshal(w.Body.Bytes(), &orders)

	assert.Equal(t, 2, len(orders))

	// Clean up
	database.DB.Delete(&testOrder1)
	database.DB.Delete(&testOrder2)
	database.DB.Delete(&testOrder3)
}

func TestUpdateOrderHandler(t *testing.T) {
	setupTestData()

	testOrder := models.Order{UserID: 1, TotalCost: 100.0}
	database.DB.Create(&testOrder)

	updatedOrder := models.Order{
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
	var fetchedOrder models.Order
	json.Unmarshal(w.Body.Bytes(), &fetchedOrder)
	assert.Equal(t, updatedOrder.TotalCost, fetchedOrder.TotalCost)

	// Clean up
	database.DB.Delete(&fetchedOrder)
}

func TestDeleteOrderHandler(t *testing.T) {
	setupTestData()

	testOrder := models.Order{UserID: 1, TotalCost: 100.0}
	database.DB.Create(&testOrder)

	req, _ := http.NewRequest("DELETE", "/orders/"+strconv.Itoa(int(testOrder.ID)), nil)
	w := httptest.NewRecorder()

	router := setupRouter()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var fetchedOrder models.Order
	result := database.DB.First(&fetchedOrder, testOrder.ID)
	assert.Error(t, result.Error)
}

func TestGetOrderItemsByOrderIDHandler(t *testing.T) {
	setupTestData()

	testOrder := models.Order{UserID: 1, TotalCost: 100.0}
	database.DB.Create(&testOrder)

	testProduct := models.Product{Name: "Test Product", Price: 10.0}
	database.DB.Create(&testProduct)

	testOrderItem1 := models.OrderItem{OrderID: testOrder.ID, ProductID: testProduct.ID, Quantity: 2, Price: 10.0}
	testOrderItem2 := models.OrderItem{OrderID: testOrder.ID, ProductID: testProduct.ID, Quantity: 3, Price: 10.0}

	database.DB.Create(&testOrderItem1)
	database.DB.Create(&testOrderItem2)

	req, _ := http.NewRequest("GET", "/orders/items/"+strconv.Itoa(int(testOrder.ID)), nil)
	w := httptest.NewRecorder()

	router := setupRouter()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var orderItems []models.OrderItem
	json.Unmarshal(w.Body.Bytes(), &orderItems)

	assert.Equal(t, 2, len(orderItems))

	// Clean up
	database.DB.Delete(&testOrder)
	database.DB.Delete(&testProduct)
	database.DB.Delete(&testOrderItem1)
	database.DB.Delete(&testOrderItem2)
}
