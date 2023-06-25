package prize_handlers_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"xy.com/mysite/database"
	"xy.com/mysite/handlers/prize_handlers"
	"xy.com/mysite/models/prize_models"
)

func setupRouter4() *gin.Engine {
	router := gin.Default()
	orderGroup := router.Group("/prize_handlers")
	{
		orderGroup.POST("/addCode", prize_handlers.AddCodeHandler)
		orderGroup.GET("/getCode", prize_handlers.GetCodeHandler)
	}
	return router
}

func TestAddCodeHandler(t *testing.T) {
	database.InitDB()

	// Setup data
	code := `{"code":"testcode"}`

	req, err := http.NewRequest("POST", "/prize_handlers/addCode", strings.NewReader(code))
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()

	router := setupRouter4()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// Check that the code has been added
	var addedCode prize_models.Code
	err = database.DB.Where("code = ?", "testcode").First(&addedCode).Error
	assert.NoError(t, err)

	// Clean up
	database.DB.Delete(&addedCode)
}

func TestGetCodeHandler(t *testing.T) {
	database.InitDB()

	// Setup data
	code := prize_models.Code{Code: "testcode", IsUsed: false}
	if err := database.DB.Create(&code).Error; err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("GET", "/prize_handlers/getCode", nil)
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()

	router := setupRouter4()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// Check that the code has been used
	var usedCode prize_models.Code
	err = database.DB.Where("code = ?", "testcode").First(&usedCode).Error
	assert.NoError(t, err)
	assert.Equal(t, true, usedCode.IsUsed)

	// Clean up
	database.DB.Delete(&usedCode)
}
