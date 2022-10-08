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
	//controllers.NewApplication(databases.(*mongo.Client)(ProductData()(databases.client, "Products"), databases.UserData(databases.Client, "Users")))

	app := controllers.NewApplication(databases.ProductData(databases.Client, "Products"), databases.UserData(databases.Client, "Users"))
	router := gin.New()
	router.Use(gin.Logger())

	routes.UserRoutes(router)
	router.Use(middlewares.Authentication())

	router.GET("/addtocart", app.AddToCart())
	router.GET("/removeitems", app.RemoveItem())
	router.GET("/cartcheckout", app.GetItemFromCart())
	router.GET("/instantbuy", app.InstantBuy())

	log.Fatal(router.Run(":" + port))
}
