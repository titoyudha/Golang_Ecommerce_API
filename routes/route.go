package routes

import (
	"go_ecommerce/controllers"

	"github.com/gin-gonic/gin"
)

func UserRoutes(route *gin.Engine) {
	route.POST("/api/v1/users/signup", controllers.SignUp())
	route.POST("/api/v1/users/login", controllers.LogIn())
	route.POST("/api/v1/admin/addproduct", controllers.ProductViewerAdmin())
	route.GET("/api/v1/users/productview", controllers.SearchProduct())
	route.GET("/api/v1/users/search", controllers.SearchProductbyQuery())
	route.GET("/api/v1/payment-stripe", controllers.Charge())
}
