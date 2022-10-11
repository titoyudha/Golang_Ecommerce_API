package controllers

import (
	"context"
	"fmt"
	"go_ecommerce/cache"
	"go_ecommerce/databases"
	"go_ecommerce/models"
	jwt "go_ecommerce/tokens"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var userCollection *mongo.Collection = databases.UserData(databases.Client, "Users")
var prodCollection *mongo.Collection = databases.ProductData(databases.Client, "Products")
var Validate = validator.New()

var (
	redisCache cache.RedisCache
	userCache  cache.UserCache
)

func HashPassword(password string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 15)
	if err != nil {
		log.Panic(err)
	}
	return string(hash)
}

func VerifyPassword(userPassword string, givenPassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(givenPassword), []byte(userPassword))
	valid := true
	msg := ""

	if err != nil {
		msg = "Incorrect Password"
		valid = false
	}
	return valid, msg
}

func SignUp() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 50*time.Second)
		defer cancel()

		var user models.User

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err})
			return
		}

		validationErr := Validate.Struct(user)
		if validationErr != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": validationErr})
			return
		}

		count, err := userCollection.CountDocuments(ctx, bson.M{"email": user.Email})
		if err != nil {
			log.Fatal(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}

		if count > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User already exist"})
		}

		count, err = userCollection.CountDocuments(ctx, bson.M{"phone": user.Phone})
		defer cancel()

		if err != nil {
			log.Fatal(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}

		if count > 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "this phone number already used"})
		}

		password := HashPassword(*user.Password)
		user.Password = &password

		user.Created_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.Update_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.ID = primitive.NewObjectID()
		user.User_ID = user.ID.Hex()
		token, refreshToken, _ := jwt.TokenGenerator(*user.Email, *user.First_Name, *user.Last_Name, *&user.User_ID)
		user.Token = &token
		user.Refresh_Token = &refreshToken
		user.UserCart = make([]models.ProductUser, 0)
		user.Address_Details = make([]models.Address, 0)
		user.Order_Status = make([]models.Order, 0)

		var cacheID = user.ID.Hex()
		var userCache *models.User = userCache.Get(cacheID)
		if userCache == nil {
			return
		}
		redisCache.Set(cacheID, userCache)
		_, insertErr := userCollection.InsertOne(ctx, user)
		if insertErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "The User didn't get created"})
			return
		}
		defer cancel()
		c.JSON(http.StatusOK, gin.H{"message": "Successfull Signed In"})
	}
}

func LogIn() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 50*time.Second)
		defer cancel()

		var user models.User
		var foundUser models.User
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}

		err := userCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&foundUser)
		defer cancel()

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "email or password incorrect"})
			return
		}

		passwordValid, msg := VerifyPassword(*user.Password, *foundUser.Password)
		defer cancel()

		if !passwordValid {
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			fmt.Println(msg)
			return
		}
		token, refreshToken, _ := jwt.TokenGenerator(*foundUser.Email, *foundUser.Password, *foundUser.First_Name, *&foundUser.User_ID)

		jwt.UpdateAllToken(token, refreshToken, foundUser.User_ID)

		c.JSON(http.StatusFound, foundUser)

	}
}

func ProductViewerAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var products models.Product
		if err := c.BindJSON(&products); err != nil {
			c.IndentedJSON(http.StatusUnprocessableEntity, gin.H{"Message": err})
			return
		}
		products.ProductID = primitive.NewObjectID()
		_, collection := prodCollection.InsertOne(ctx, &products)
		if collection != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Message": "Not Inserted"})
			return
		}
		defer cancel()
		c.IndentedJSON(http.StatusOK, gin.H{"Message": "Successful added"})
	}
}

func SearchProduct() gin.HandlerFunc {
	return func(c *gin.Context) {
		var productList []models.Product
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		cursor, err := prodCollection.Find(ctx, bson.D{{}})
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, "Something went wrong, please try again")
			return
		}

		err = cursor.All(ctx, productList)

		if err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		defer cursor.Close(ctx)

		if err := cursor.Err(); err != nil {
			log.Println(err)
			c.IndentedJSON(400, "Invalid")
			return
		}
		defer cancel()
		c.IndentedJSON(http.StatusOK, productList)
	}
}

func SearchProductbyQuery() gin.HandlerFunc {
	return func(c *gin.Context) {
		var searchProduct []models.Product
		queryParam := c.Query("name")

		// if params is empty
		if queryParam == "" {
			log.Println("Query is empty")
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"Error": "Invalid search index"})
			c.Abort()
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		queryDB, err := prodCollection.Find(ctx, bson.M{"product_name": bson.M{"$regex": queryParam}})

		if err != nil {
			c.IndentedJSON(http.StatusNotFound, "Something wrong when fetching the data")
			return
		}
		err = queryDB.All(ctx, &searchProduct)

		if err != nil {
			log.Println(err)
			c.IndentedJSON(400, "Invalid")
			return
		}

		defer queryDB.Close(ctx)

		if err := queryDB.Err(); err != nil {
			log.Println(err)
			c.IndentedJSON(400, "Invalid Request")
			return
		}
		defer cancel()

		c.IndentedJSON(http.StatusOK, searchProduct)
	}
}
