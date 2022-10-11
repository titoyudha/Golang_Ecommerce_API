package main

import (
	"go_ecommerce/databases"
	"go_ecommerce/routes"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// app := controllers.NewApplication(databases.ProductData(databases.Client, "Products"), databases.UserData(databases.Client, "Users"))
	client, ctx, cancel, err := databases.Connect("mongodb://localhost:27017")
	if err != nil {
		panic(err)
	}
	defer databases.Close(client, ctx, cancel)
	router := gin.New()
	router.Use(gin.Logger())
	routes.UserRoutes(router)
	router.Use(gin.Logger())

	// routes.UserRoutes(router)
	// router.Use(middlewares.Authentication())

	// router.GET("/addtocart", app.AddToCart())
	// router.GET("/removeitem", app.RemoveItem())
	// router.POST("/addaddress", controllers.AddAddress())
	// router.PUT("/edithomeaddress", controllers.EditHomeAddress())
	// router.PUT("/editworkaddress", controllers.EditWorkAddress())
	// router.GET("/deleteaddresses", controllers.DeleteAddress())
	// router.GET("/cartcheckout", app.BuyFromCart())
	// router.GET("/instantbuy", app.InstantBuy())

	router.Run()

}
