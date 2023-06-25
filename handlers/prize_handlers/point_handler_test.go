package prize_handlers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"xy.com/mysite/database"
	"xy.com/mysite/handlers/prize_handlers"
	"xy.com/mysite/models/prize_models"
)

func setupRouter2() *gin.Engine {
	router := gin.Default()
	orderGroup := router.Group("/prize_handlers")
	{
		orderGroup.POST("/draw/:userID", prize_handlers.DrawHandler)
		orderGroup.POST("/exchange/:userID", prize_handlers.ExchangeCoinsHandler)
		orderGroup.GET("/getPointsSystem/:userID", prize_handlers.GetPointsSystemHandler)
	}
	return router
}

func TestDrawHandler(t *testing.T) {
	database.InitDB()

	// Setup data
	userID := "testuser"
	points := 30000 // Set points so that user can draw

	// Add user and points system
	pointsSystem := &prize_models.PointsSystem{UserID: userID, Points: points}
	if err := database.DB.Create(pointsSystem).Error; err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("POST", "/prize_handlers/draw/"+userID, nil)
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()

	router := setupRouter2()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// Check that the points system has been updated
	var updatedPointsSystem prize_models.PointsSystem
	err = database.DB.Where("user_id = ?", userID).First(&updatedPointsSystem).Error
	assert.NoError(t, err)
	assert.Greater(t, updatedPointsSystem.Points, points) // Check that points have increased

	// Clean up
	database.DB.Delete(&pointsSystem)
}

func TestExchangeCoinsHandler(t *testing.T) {
	database.InitDB()

	// Setup data
	userID := "testuser"
	coins := 200 // Set coins so that user can exchange

	// Add user and points system
	pointsSystem := &prize_models.PointsSystem{UserID: userID, Coins: coins}
	if err := database.DB.Create(pointsSystem).Error; err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("POST", "/prize_handlers/exchange/"+userID, nil)
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()

	router := setupRouter2()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// Check that the points system has been updated
	var updatedPointsSystem prize_models.PointsSystem
	err = database.DB.Where("user_id = ?", userID).First(&updatedPointsSystem).Error
	assert.NoError(t, err)
	assert.Equal(t, coins%100, updatedPointsSystem.Coins) // Check that coins have been converted to points

	// Clean up
	database.DB.Delete(&pointsSystem)
}

func TestGetPointsSystemHandler(t *testing.T) {
	database.InitDB()

	// Setup data
	userID := "testuser"

	// Add user and points system
	pointsSystem := &prize_models.PointsSystem{UserID: userID}
	if err := database.DB.Create(pointsSystem).Error; err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("GET", "/prize_handlers/getPointsSystem/"+userID, nil)
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()

	router := setupRouter2()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var retrievedPointsSystem prize_models.PointsSystem
	err = json.Unmarshal(w.Body.Bytes(), &retrievedPointsSystem)
	assert.NoError(t, err)
	assert.Equal(t, userID, retrievedPointsSystem.UserID)

	// Clean up
	database.DB.Delete(&pointsSystem)
}
