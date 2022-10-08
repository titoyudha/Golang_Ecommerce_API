package controllers

import (
	"context"
	"go_ecommerce/models"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func AddAddress() gin.HandlerFunc {
	return func(c *gin.Context) {
		user_id := c.Query("id")

		if user_id == "" {
			c.Header("Content-Type", "application/json")
			c.IndentedJSON(http.StatusNotFound, gin.H{"Message": "Invalid search Index"})
			c.Abort()
			return
		}
		address, err := primitive.ObjectIDFromHex(user_id)

		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"Message": "Internal Server Error"})
		}

		var addresses models.Address
		addresses.Address_id = primitive.NewObjectID()

		if err = c.BindJSON(&addresses); err != nil {
			c.IndentedJSON(http.StatusUnprocessableEntity, err.Error())
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		matchFilter := bson.D{{Key: "$match", Value: bson.D{primitive.E{Key: "_id", Value: address}}}}
		unwind := bson.D{{Key: "$unwind", Value: bson.D{primitive.E{Key: "path", Value: "$address"}}}}
		group := bson.D{{Key: "$group", Value: bson.D{primitive.E{Key: "_id", Value: "$address_id"}, {Key: "count", Value: bson.D{primitive.E{Key: "$sum", Value: 1}}}}}}

		pointCursor, err := userCollection.Aggregate(ctx, mongo.Pipeline{matchFilter, unwind, group})
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Internal Server Error"})
		}

		var addressInfo []bson.M
		if err = pointCursor.All(ctx, &addressInfo); err != nil {
			panic(err.Error())
		}

		var size int32

		for _, address_no := range addressInfo {
			count := address_no["count"]
			size = count.(int32)
		}

		if size < 2 {
			filter := bson.D{primitive.E{Key: "_id", Value: address}}
			update := bson.D{{Key: "$push", Value: bson.D{primitive.E{Key: "address", Value: addresses}}}}
			_, err := userCollection.UpdateOne(ctx, filter, update)

			if err != nil {
				log.Println(err)
			}

		} else {
			c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Method not Allowed"})
		}
		defer cancel()
		ctx.Done()
	}
}

func EditHomeAddress() gin.HandlerFunc {
	return func(c *gin.Context) {
		user_id := c.Query("id")

		if user_id == "" {
			c.Header("Content-Type", "application/json")
			c.IndentedJSON(http.StatusNotFound, gin.H{"Message": "Invalid ID"})
			c.Abort()
			return
		}
		userID, err := primitive.ObjectIDFromHex(user_id)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"Message": "Internal Server Error"})
		}
		var editAddress models.Address

		if err := c.BindJSON(&editAddress); err != nil {
			c.IndentedJSON(http.StatusUnprocessableEntity, gin.H{"Message": "Please enter valid entity"})
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		filter := bson.D{primitive.E{Key: "i_d", Value: userID}}
		update := bson.D{{Key: "$set", Value: bson.D{primitive.E{Key: "address.0.description", Value: editAddress.Description}, {Key: "address.0.street_name", Value: editAddress.Street}, {Key: "address.0.city_name", Value: editAddress.City}, {Key: "address.0.pin_code", Value: editAddress.Pincode}}}}
		_, err = userCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"Message": "Something Went Wrong"})
			return
		}
		ctx.Done()
		c.IndentedJSON(http.StatusOK, "Successful update")
	}
}

func EditWorkAddress() gin.HandlerFunc {
	return func(c *gin.Context) {
		user_id := c.Query("id")

		if user_id == "" {
			c.Header("Content-Type", "application/json")
			c.IndentedJSON(http.StatusNotFound, gin.H{"Message": "Invalid ID"})
			c.Abort()
			return
		}
		userID, err := primitive.ObjectIDFromHex(user_id)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"Message": "Internal Server Error"})
		}
		var editAddress models.Address
		if err := c.BindJSON(editAddress); err != nil {
			c.IndentedJSON(http.StatusBadRequest, err.Error())
		}
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		filter := bson.D{primitive.E{Key: "_id", Value: userID}}
		update := bson.D{{Key: "$set", Value: bson.D{primitive.E{Key: "address.1.description", Value: editAddress.Description}, {Key: "address.1.street_name", Value: editAddress.Street}, {Key: "address.1.city_name", Value: editAddress.City}, {Key: "address.1.pin_code", Value: editAddress.Pincode}}}}
		_, err = userCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"Message": "Something went wrong"})
			return
		}
		defer cancel()
		ctx.Done()
		c.IndentedJSON(http.StatusOK, gin.H{"Message": "Successfull Update Working Address"})
	}
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
