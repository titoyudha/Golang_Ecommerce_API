package main

import (
	"go_ecommerce/controllers"
	"go_ecommerce/databases"
	"go_ecommerce/middlewares"
	"go_ecommerce/routes"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	app := controllers.NewApplication(databases.ProductData(databases.Client, "Products"), databases.UserData(databases.Client, "Users"))

	client, ctx, cancel, err := databases.Connect("mongodb://localhost:27018")
	if err != nil {
		panic(err)
	}
	defer databases.Close(client, ctx, cancel)
	router := gin.New()
	router.Use(gin.Logger())
	routes.UserRoutes(router)
	router.Use(gin.Logger())

	router.Use(middlewares.Authentication())

	router.GET("/api/v1/addtocart", app.AddToCart())
	router.GET("/api/v1/removeitem", app.RemoveItem())
	router.POST("/api/v1/addaddress", controllers.AddAddress())
	router.PUT("/api/v1/edithomeaddress", controllers.EditHomeAddress())
	router.PUT("/api/v1/editworkaddress", controllers.EditWorkAddress())
	router.GET("/api/v1/deleteaddresses", controllers.DeleteAddress())
	router.GET("/api/v1/cartcheckout", app.BuyFromCart())
	router.GET("/api/v1/instantbuy", app.InstantBuy())

	router.Run()

}
