// handlers/chat_handler_test.go

package handlers

import (
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

func TestHandleConnectionsAndMessages(t *testing.T) {
	router := gin.Default()
	router.GET("/ws", func(c *gin.Context) {
		HandleConnections(c.Writer, c.Request)
	})

	server := httptest.NewServer(router)
	defer server.Close()

	u := url.URL{Scheme: "ws", Host: server.Listener.Addr().String(), Path: "/ws"}

	// Create WebSocket Connect
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		t.Fatalf("dial: %v", err)
	}
	defer c.Close()

	// Send A Message
	message := Message{
		Email:    "test@example.com",
		Username: "test_user",
		Message:  "Hello, world!",
	}

	err = c.WriteJSON(message)
	if err != nil {
		t.Fatalf("write: %v", err)
	}

	go HandleMessages()

	// Read Message and Verify
	var receivedMessage Message
	err = c.ReadJSON(&receivedMessage)
	if err != nil {
		t.Fatalf("read: %v", err)
	}

	assert.Equal(t, message.Email, receivedMessage.Email)
	assert.Equal(t, message.Username, receivedMessage.Username)
	assert.Equal(t, message.Message, receivedMessage.Message)
}
