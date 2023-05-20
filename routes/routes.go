// routes/routes.go

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
		orderGroup.GET("/getAllOrders", handlers.GetAllOrdersHandler)
		orderGroup.POST("/", handlers.CreateOrderHandler)
		orderGroup.GET("/:id", handlers.GetOrderByIDHandler)
		orderGroup.GET("/user/:userID", handlers.GetOrdersByUserIDHandler)
		orderGroup.PUT("/:id", handlers.UpdateOrderHandler)
		orderGroup.DELETE("/:id", handlers.DeleteOrderHandler)
		orderGroup.GET("/items/:orderID", handlers.GetOrderItemsByOrderIDHandler)
	}

	// Product routes
	productGroup := router.Group("/products", middleware.AuthMiddleware())
	{
		productGroup.POST("/", handlers.CreateProductHandler)
		productGroup.GET("/:id", handlers.GetProductHandlerByID)
		productGroup.GET("/all", handlers.GetAllProductsHandler)
		productGroup.PUT("/:id", handlers.UpdateProductHandler)
		productGroup.DELETE("/:id", handlers.DeleteProductHandler)
	}

	// Chat routes
	router.GET("/ws", func(c *gin.Context) {
		handlers.HandleConnections(c.Writer, c.Request)
	})

	return router
}
