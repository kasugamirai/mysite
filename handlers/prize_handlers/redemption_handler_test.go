package prize_handlers_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"xy.com/mysite/handlers/prize_handlers"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"xy.com/mysite/database"
	"xy.com/mysite/models/prize_models"
)

func setupRouter5() *gin.Engine {
	router := gin.Default()
	prizeGroup := router.Group("/prize_handlers")
	{
		prizeGroup.POST("/addRedemptionCode", prize_handlers.AddRedemptionCodeHandler)
	}
	return router
}

func TestAddRedemptionCodeHandler(t *testing.T) {
	database.InitDB()

	// Setup data
	code := `{"code":"testcode"}`

	req, err := http.NewRequest("POST", "/prize_handlers/addRedemptionCode", strings.NewReader(code))
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()

	router := setupRouter5()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// Check that the code has been added
	var addedCode prize_models.RedemptionCode
	err = database.DB.Where("code = ?", "testcode").First(&addedCode).Error
	assert.NoError(t, err)

	// Clean up
	database.DB.Delete(&addedCode)
}
