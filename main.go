package main

import (
	"go_ecommerce/controllers"
	"go_ecommerce/databases"
	"go_ecommerce/middlewares"
	"go_ecommerce/routes"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	app := controllers.NewApplication(databases.ProductData(databases.client, "Products"), databases.UserData(databases.Client, "Users"))

	router := gin.New()
	router.Use(gin.Logger())

	routes.UserRoutes(router)
	router.Use(middlewares.Authentication())

	router.GET("/addtocart", app.AddToCart())
	router.GET("/removeitems", app.RemoveItems())
	router.GET("/cartcheckout", app.CartCheckout())
	router.GET("/instantbuy", app.InstantBuy())

	log.Fatal(router.Run(":" + port))
}
