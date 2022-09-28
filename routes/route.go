package routes

import (
	"go_ecommerce/controllers"

	"github.com/gin-gonic/gin"
)

func UserRoutes(route *gin.Engine) {
	route.POST("/users/signup", controllers.SignUp())
	route.POST("/users/login", controllers.LogIn())
	route.POST("/admin/addproduct", controllers.ProductViewerAdmin())
	route.GET("/users/productview", controllers.SearchProduct())
	route.GET("/users/search", SearchProductbyQuery())
}
