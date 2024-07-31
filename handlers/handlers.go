package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/iadityanath8/gopi/middleware"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	v1 := router.Group("/api/v1")
	{

		//users
		v1.POST("/login", Login)
		v1.POST("/register", Register)
		// userend

		// products
		v1.GET("/products", GetProducts)
		v1.POST("/products", middleware.VerifyToken, CreateProduct)
		v1.GET("/products/:id", GetProduct)
		v1.PUT("/products/:id", middleware.VerifyToken, UpdataProduct)
		v1.DELETE("/products/:id", middleware.VerifyToken, DeleteProduct)
		// productend

		// cartendpoint
		v1.POST("/addToCart", middleware.VerifyToken, AddItemToCart)
		v1.GET("/getCart", middleware.VerifyToken, GetCartItems)
		v1.DELETE("/deleteCart", middleware.VerifyToken, RemoveItemCart)
		// cartendpointend
	}

	return router
}
