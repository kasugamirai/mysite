package prize_handlers_test

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"xy.com/mysite/database"
	"xy.com/mysite/handlers/prize_handlers"
	"xy.com/mysite/models/prize_models"
)

func setupRouter() *gin.Engine {
	router := gin.Default()
	orderGroup := router.Group("/prize_handlers")
	{
		orderGroup.POST("/addPrize", prize_handlers.AddPrizeHandler)
	}
	return router
}

func TestAddPrizeHandler(t *testing.T) {
	database.InitDB()

	// Setup data
	prizeName := "testprize"
	cost := 100

	reqBody := map[string]interface{}{
		"prize_name": prizeName,
		"cost":       cost,
	}
	reqBodyJson, _ := json.Marshal(reqBody)
	body := bytes.NewReader(reqBodyJson)

	req, err := http.NewRequest("POST", "/prize_handlers/addPrize", body)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	router := setupRouter()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	// Check that the new prize_handlers has been added
	var prize prize_models.Prize
	err = database.DB.Where("prize_name = ?", prizeName).First(&prize).Error
	assert.NoError(t, err)
	assert.Equal(t, cost, prize.Cost)

	// Clean up
	database.DB.Delete(&prize)
}
