package prize_handlers_test

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"xy.com/mysite/handlers/prize_handlers"
)

func TestExchangePrizeHandler(t *testing.T) {
	// Initialize gin router
	router := gin.Default()
	router.POST("/exchangePrize", prize_handlers.ExchangePrizeHandler)

	// Create a request to pass to our handler
	req, err := http.NewRequest("POST", "/exchangePrize", strings.NewReader(`{"user_id":"1","prize_name":"Prize1"}`))
	if err != nil {
		t.Fatal(err)
	}

	// Initialize a ResponseRecorder to record the response
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code, "Expected response code 200")

}

// Continue similarly for other handlers...
