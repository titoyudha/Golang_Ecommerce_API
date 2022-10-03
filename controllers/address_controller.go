package controllers

import (
	"context"
	"go_ecommerce/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func AddAddress() gin.HandlerFunc {

}

func EditHomeAddress() gin.HandlerFunc {

}

func EditWorkAddress() gin.HandlerFunc {

}

func DeleteAddress() gin.HandlerFunc {
	return func(c *gin.Context) {
		user_id := c.Query("id")

		if user_id == "" {
			c.Header("Content-Type", "application/json")
			c.IndentedJSON(http.StatusNotFound, gin.H{"Message": "Invalid Search Index"})
			c.Abort()
			return
		}

		addresses := make([]models.Address, 0)
		userID, err := primitive.ObjectIDFromHex(user_id)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Internal Server Eror"})
		}
		ctx, cancel := context.WithTimeout(context.Background(), 50*time.Second)
		defer cancel()

		filter := bson.D{primitive.E{Key: "_id", Value: userID}}
		update := bson.D{{Key: "$set", Value: bson.D{primitive.E{Key: "address", Value: addresses}}}}
		_, err = userCollection.UpdateOne(ctx, filter, update)

		if err != nil {
			c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Wrong Command"})
			return
		}
		defer cancel()
		ctx.Done()
		c.IndentedJSON(http.StatusOK, gin.H{"message": "Successfully Deleted"})
	}
}
