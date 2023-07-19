// routes/routes.go

package routes

import (
	"github.com/gin-gonic/gin"
	"xy.com/mysite/handlers"
	"xy.com/mysite/handlers/prize_handlers"
	"xy.com/mysite/handlers/shop_handlers"
	"xy.com/mysite/handlers/user_handlers"
	"xy.com/mysite/middleware"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()
	SetupStaticRoutes(router)

	// Auth routes
	authGroup := router.Group("/auth")
	{
		authGroup.POST("/signup", user_handlers.CreateUserHandler)
		authGroup.POST("/login", user_handlers.AuthenticateUserHandlers)
		authGroup.POST("/logout", middleware.AuthMiddleware(), user_handlers.LogoutUserHandler)
	}

	// User routes
	userGroup := router.Group("/users", middleware.AuthMiddleware())
	{
		userGroup.GET("/:id", user_handlers.GetUserHandler)
		userGroup.GET("/", user_handlers.GetUserByEmailHandler)
		userGroup.PUT("/:id", user_handlers.UpdateUserHandler)
		userGroup.DELETE("/:id", user_handlers.DeleteUserHandler)
	}

	// Order routes
	orderGroup := router.Group("/orders", middleware.AuthMiddleware())
	{
		orderGroup.GET("/getAllOrders", shop_handlers.GetAllOrdersHandler)
		orderGroup.POST("/", shop_handlers.CreateOrderHandler)
		orderGroup.GET("/:id", shop_handlers.GetOrderByIDHandler)
		orderGroup.GET("/user/:userID", shop_handlers.GetOrdersByUserIDHandler)
		orderGroup.PUT("/:id", shop_handlers.UpdateOrderHandler)
		orderGroup.DELETE("/:id", shop_handlers.DeleteOrderHandler)
		orderGroup.GET("/items/:orderID", shop_handlers.GetOrderItemsByOrderIDHandler)
	}

	// Product routes
	productGroup := router.Group("/products", middleware.AuthMiddleware())
	{
		productGroup.POST("/", shop_handlers.CreateProductHandler)
		productGroup.GET("/:id", shop_handlers.GetProductHandlerByID)
		productGroup.GET("/all", shop_handlers.GetAllProductsHandler)
		productGroup.PUT("/:id", shop_handlers.UpdateProductHandler)
		productGroup.DELETE("/:id", shop_handlers.DeleteProductHandler)
	}

	// Chat routes
	router.GET("/ws", func(c *gin.Context) {
		handlers.HandleConnections(c.Writer, c.Request)
	})

	pointGroup := router.Group("/point", middleware.AuthMiddleware())
	{
		pointGroup.GET("/points/", prize_handlers.GetPointsSystemHandler)
		pointGroup.POST("/draw/", middleware.CheckRedemptionCode(), prize_handlers.DrawHandler)
		pointGroup.POST("/exchange/", prize_handlers.ExchangeCoinsHandler)
	}

	prizeGroup := router.Group("/prize_handlers", middleware.AuthMiddleware())
	{
		prizeGroup.POST("/addPrize", prize_handlers.AddPrizeHandler)
	}

	exchangeGroup := router.Group("/exchange", middleware.AuthMiddleware())
	{
		exchangeGroup.GET("/checkExchanged/:prizeName", prize_handlers.CheckIfUserExchangedPrizeHandler)
		exchangeGroup.POST("/exchanged", prize_handlers.ExchangePrizeHandler)
		exchangeGroup.GET("/prize_handlers/:prizeName", prize_handlers.GetPrizeByNameHandler)
	}

	adminGroup := router.Group("/admin", middleware.AuthMiddleware())
	{
		adminGroup.POST("/addCode", prize_handlers.AddCodeHandler)
		adminGroup.POST("/addRedemptionCode", prize_handlers.AddRedemptionCodeHandler)
		adminGroup.POST("/addPrize", prize_handlers.AddPrizeHandler)
	}
	return router
}
