package prize_handlers_test

import (
	"bytes"
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

func setupRouter1() *gin.Engine {
	router := gin.Default()
	orderGroup := router.Group("/prize_handlers")
	{
		orderGroup.POST("/exchangePrize", prize_handlers.ExchangePrizeHandler)
		orderGroup.GET("/checkIfUserExchangedPrize/:userID/:prizeName", prize_handlers.CheckIfUserExchangedPrizeHandler)
		orderGroup.GET("/getPrizeByName/:prizeName", prize_handlers.GetPrizeByNameHandler)
	}
	return router
}

func TestExchangePrizeHandler(t *testing.T) {
	database.InitDB()

	// Setup data
	userID := "testuser"
	prizeName := "testprize"
	cost := 100
	initialPoints := 200

	// Add user and points system
	pointsSystem := &prize_models.PointsSystem{UserID: userID, Points: initialPoints}
	if err := database.DB.Create(pointsSystem).Error; err != nil {
		t.Fatal(err)
	}

	// Add prize
	prize := &prize_models.Prize{PrizeName: prizeName, Cost: cost}
	if err := database.DB.Create(prize).Error; err != nil {
		t.Fatal(err)
	}

	// Add code to the database
	code := "ABC123"
	if err := prize_models.AddCode(database.DB, code); err != nil {
		t.Fatal(err)
	}

	reqBody := map[string]interface{}{
		"user_id":    userID,
		"prize_name": prizeName,
	}
	reqBodyJson, _ := json.Marshal(reqBody)
	body := bytes.NewReader(reqBodyJson)

	req, err := http.NewRequest("POST", "/prize_handlers/exchangePrize", body)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	router := setupRouter1()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// Check that the points system has not been updated (points remain unchanged)
	var updatedPointsSystem prize_models.PointsSystem
	err = database.DB.Where("user_id = ?", userID).First(&updatedPointsSystem).Error
	assert.NoError(t, err)
	assert.Equal(t, initialPoints, updatedPointsSystem.Points)

	// Check that the prize has been exchanged
	var exchangedPrize prize_models.ExchangedPrize
	err = database.DB.Where("user_id = ? AND prize_name = ?", userID, prizeName).First(&exchangedPrize).Error
	assert.NoError(t, err)
	assert.Equal(t, code, exchangedPrize.RedemptionCode)

	// Clean up
	database.DB.Delete(&exchangedPrize)
	database.DB.Delete(&updatedPointsSystem)
	database.DB.Delete(&prize)
}

// Test for CheckIfUserExchangedPrizeHandler
func TestCheckIfUserExchangedPrizeHandler(t *testing.T) {
	database.InitDB()

	// Setup data
	userID := "testuser"
	prizeName := "testprize"

	// Add user and prize
	exchangedPrize := &prize_models.ExchangedPrize{UserID: userID, PrizeName: prizeName}
	if err := database.DB.Create(exchangedPrize).Error; err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("GET", "/prize_handlers/checkIfUserExchangedPrize/"+userID+"/"+prizeName, nil)
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()

	router := setupRouter1()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var result map[string]bool
	err = json.Unmarshal(w.Body.Bytes(), &result)
	assert.NoError(t, err)
	assert.True(t, result["hasExchanged"])

	// Clean up
	database.DB.Delete(&exchangedPrize)
}

// Test for GetPrizeByNameHandler
func TestGetPrizeByNameHandler(t *testing.T) {
	database.InitDB()

	// Setup data
	prizeName := "testprize"
	cost := 100

	// Add prize
	prize := &prize_models.Prize{PrizeName: prizeName, Cost: cost}
	if err := database.DB.Create(prize).Error; err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("GET", "/prize_handlers/getPrizeByName/"+prizeName, nil)
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()

	router := setupRouter1()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var retrievedPrize prize_models.Prize
	err = json.Unmarshal(w.Body.Bytes(), &retrievedPrize)
	assert.NoError(t, err)
	assert.Equal(t, prizeName, retrievedPrize.PrizeName)
	assert.Equal(t, cost, retrievedPrize.Cost)

	// Clean up
	database.DB.Delete(&prize)
}
