package controllers

import (
	"context"
	"go_ecommerce/databases"
	"go_ecommerce/models"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type App struct {
	prodCollection *mongo.Collection
	userCollection *mongo.Collection
}

func NewApplication(prodCollection, userCollection *mongo.Collection) *App {
	return &App{
		prodCollection: prodCollection,
		userCollection: userCollection,
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

		err = databases.AddProductToCart(c, app.prodCollection, app.userCollection, productID, userQueryID)
		if err != nil {
			ctx.IndentedJSON(http.StatusInternalServerError, err)
			return
		}
		ctx.IndentedJSON(200, "Success")
	}
}

func (app *App) RemoveItem() gin.HandlerFunc {
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

		err = databases.RemoveCartItem(c, app.prodCollection, app.userCollection, productID, userQueryID)

		if err != nil {
			ctx.IndentedJSON(http.StatusInternalServerError, err)
			return
		}
		ctx.IndentedJSON(http.StatusOK, "Successfull Remove item from the cart")
	}
}

func (app *App) GetItemFromCart() gin.HandlerFunc {
	return func(c *gin.Context) {
		user_id := c.Query("id")

		if user_id == "" {
			c.Header("Content-Type", "application/json")
			c.IndentedJSON(http.StatusNotFound, gin.H{"messsage": "invalid id"})
			c.Abort()
			return
		}

		userID, _ := primitive.ObjectIDFromHex(user_id)

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var filledCart models.User
		err := userCollection.FindOne(ctx, bson.D{primitive.E{Key: "_id", Value: userID}}).Decode(&filledCart)

		if err != nil {
			log.Println(err)
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Not Found"})
			return
		}
		filterMatch := bson.D{{Key: "$match", Value: bson.D{primitive.E{Key: "_id", Value: userID}}}}
		unwind := bson.D{{Key: "$unwind", Value: bson.D{primitive.E{Key: "path", Value: "$usercart"}}}}
		grouping := bson.D{{Key: "$group", Value: bson.D{primitive.E{Key: "_id", Value: "$_id"}, {Key: "total", Value: bson.D{primitive.E{Key: "$sum", Value: "$usercart.price"}}}}}}

		pointCursor, err := userCollection.Aggregate(ctx, mongo.Pipeline{filterMatch, unwind, grouping})

		if err != nil {
			log.Println(err)
		}

		var list []bson.M
		if err = pointCursor.All(ctx, &list); err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
		}
		for _, json := range list {
			c.IndentedJSON(http.StatusOK, json["total"])
			c.IndentedJSON(http.StatusOK, filledCart.UserCart)
		}
		ctx.Done()
	}
}

func (app *App) BuyFromCart() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userQueryID := ctx.Query("id")

		if userQueryID == "" {
			log.Panic("user id is empty")
			_ = ctx.AbortWithError(http.StatusBadRequest, errors.New("User ID is Empty"))
		}
		c, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := databases.BuyItemFromCart(c, app.userCollection, userQueryID); err != nil {
			ctx.IndentedJSON(http.StatusInternalServerError, err)
			return
		}
		ctx.IndentedJSON(http.StatusOK, "Successfully placed the order")
	}
}

func (app *App) InstantBuy() gin.HandlerFunc {
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

		if err = databases.InstantBuyer(c, app.prodCollection, app.userCollection, productID, userQueryID); err != nil {
			ctx.IndentedJSON(http.StatusInternalServerError, err)
		}
		ctx.IndentedJSON(http.StatusOK, "Successfull placed the Order")
	}
}
