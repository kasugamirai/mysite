package routes

import (
	"github.com/gin-gonic/gin"
	"xy.com/mysite/handlers"
	"xy.com/mysite/middleware"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	// Auth routes
	authGroup := router.Group("/auth")
	{
		authGroup.POST("/signup", handlers.CreateUserHandler)
		authGroup.POST("/login", handlers.AuthenticateUserHandlers)
		authGroup.POST("/logout", middleware.AuthMiddleware(), handlers.LogoutUserHandler)
	}

	// User routes
	userGroup := router.Group("/users", middleware.AuthMiddleware())
	{
		userGroup.GET("/:id", handlers.GetUserHandler)
		userGroup.GET("/", handlers.GetUserByEmailHandler)
		userGroup.PUT("/:id", handlers.UpdateUserHandler)
		userGroup.DELETE("/:id", handlers.DeleteUserHandler)
	}

	// Order routes
	orderGroup := router.Group("/orders", middleware.AuthMiddleware())
	{
		orderGroup.POST("/", handlers.CreateOrderHandler)
		orderGroup.GET("/:id", handlers.GetOrderHandler)
		orderGroup.GET("/user/:userID", handlers.GetOrdersByUserIDHandler)
		orderGroup.PUT("/:id", handlers.UpdateOrderHandler)
		orderGroup.DELETE("/:id", handlers.DeleteOrderHandler)
		orderGroup.GET("/items/:orderID", handlers.GetOrderItemsByOrderIDHandler)
	}

	return router
}
