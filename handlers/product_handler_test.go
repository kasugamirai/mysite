package handlers_test

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
	"xy.com/mysite/handlers"
	"xy.com/mysite/models"
)

func setupProductRouter() *gin.Engine {
	router := gin.Default()
	productGroup := router.Group("/products")
	{
		productGroup.POST("/", handlers.CreateProductHandler)
		productGroup.GET("/:id", handlers.GetProductHandlerByID)
		productGroup.GET("/all", handlers.GetAllProductsHandler)
		productGroup.PUT("/:id", handlers.UpdateProductHandler)
		productGroup.DELETE("/:id", handlers.DeleteProductHandler)
	}
	return router
}

func TestCreateProductHandler(t *testing.T) {
	setupTestData()
	product := models.Product{Name: "Test Product", Price: 10.0}
	productJSON, _ := json.Marshal(product)

	req, err := http.NewRequest("POST", "/products/", bytes.NewReader(productJSON))
	req.Header.Set("Content-Type", "application/json")
	assert.NoError(t, err)

	resp := httptest.NewRecorder()
	router := setupProductRouter()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusCreated, resp.Code)

	var returnedProduct models.Product
	json.Unmarshal(resp.Body.Bytes(), &returnedProduct)
	assert.Equal(t, product.Name, returnedProduct.Name)
	assert.Equal(t, product.Price, returnedProduct.Price)
	database.DB.Delete(&returnedProduct)
}

func TestGetAllProductsHandler(t *testing.T) {
	setupTestData()

	product1 := &models.Product{Name: "test1", Price: 10.0}
	product2 := &models.Product{Name: "test2", Price: 20.0}
	database.DB.Create(product1)
	database.DB.Create(product2)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/products/all", nil)
	router := setupProductRouter()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var products []models.Product
	json.Unmarshal(w.Body.Bytes(), &products)

	foundProduct1 := false
	foundProduct2 := false

	for _, product := range products {
		if product.Name == product1.Name {
			foundProduct1 = true
		}
		if product.Name == product2.Name {
			foundProduct2 = true
		}
	}

	assert.True(t, foundProduct1)
	assert.True(t, foundProduct2)

	database.DB.Delete(product1)
	database.DB.Delete(product2)
}

func TestGetProductHandlerByID(t *testing.T) {
	setupTestData()
	product := &models.Product{Name: "test", Price: 123.0}

	database.DB.Create(product)
	req, _ := http.NewRequest("GET", "/products/"+strconv.Itoa(int(product.ID)), nil)
	w := httptest.NewRecorder()
	router := setupProductRouter()
	router.ServeHTTP(w, req)

	var products models.Product
	json.Unmarshal(w.Body.Bytes(), &products)
	assert.Equal(t, product.ID, products.ID)

	database.DB.Delete(product)
}

func TestUpdateProductHandler(t *testing.T) {
	setupTestData()
	product := &models.Product{Name: "origin", Price: 1.0}
	updatedProduct := &models.Product{Name: "updated", Price: 2.0}

	database.DB.Create(product)
	productHJson, _ := json.Marshal(updatedProduct)
	body := bytes.NewReader(productHJson)
	req, _ := http.NewRequest("PUT", "/products/"+strconv.Itoa(int(product.ID)), body)
	w := httptest.NewRecorder()
	router := setupProductRouter()
	router.ServeHTTP(w, req)

	var newproduct models.Product
	json.Unmarshal(w.Body.Bytes(), &newproduct)
	assert.Equal(t, updatedProduct.Name, newproduct.Name)

	database.DB.Delete(updatedProduct)
}

func TestDeleteProductHandler(t *testing.T) {
	setupTestData()
	product1 := &models.Product{Name: "test", Price: 1.0}

	database.DB.Create(product1)
	req, _ := http.NewRequest("DELETE", "/products/1", nil)
	w := httptest.NewRecorder()
	route := setupProductRouter()
	route.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
