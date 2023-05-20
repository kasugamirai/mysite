package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"xy.com/mysite/database"
	"xy.com/mysite/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"xy.com/mysite/handlers"
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
