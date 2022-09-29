package controllers

import (
	"context"
	"go_ecommerce/databases"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type App struct {
	prodCollection *mongo.Collection
	userCollection *mongo.Collection
}

func NewApplication(prodCollection, userCollection *mongo.Client) *App {
	return &App{
		prodCollection: (*mongo.Collection)(prodCollection),
		userCollection: (*mongo.Collection)(userCollection),
	}
}

func (app *App) AddToCart() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		productQueryID := ctx.Query("id")
		if productQueryID == "" {
			log.Println("Product is empty")
			_ = ctx.AbortWithError(http.StatusBadRequest, errors.New("Product is empty"))
			return
		}

		userQueryID := ctx.Query("userID")
		if userQueryID == "" {
			log.Println("user id is empty")
			_ = ctx.AbortWithError(http.StatusBadRequest, errors.New("User ID is empty"))
			return
		}

		productID, err := primitive.ObjectIDFromHex(productQueryID)
		if err != nil {
			log.Println(err)
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		var c, cancel = context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err = databases.AddProductToChart(c, app.prodCollection, app.userCollection, productID, userQueryID)
		if err != nil {
			ctx.IndentedJSON(http.StatusInternalServerError, err)
			return
		}
		ctx.IndentedJSON(200, "Success")
	}
}

func RemoteItem() gin.HandlerFunc {

}

func GetItemFromCart() gin.HandlerFunc {

}

func BuyFromCart() gin.HandlerFunc {

}

func InstantBuy() gin.HandlerFunc {

}
