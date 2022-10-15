package testing

import (
	"bytes"
	"encoding/json"
	"go_ecommerce/controllers"
	"go_ecommerce/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func SetUpRouter() *gin.Engine {
	router := gin.Default()
	return router
}
func TestSignUp(t *testing.T) {
	r := SetUpRouter()
	r.POST("/api/v1/users/signup", controllers.SignUp())

	//Error on passing value to struct
	var res = models.User{
		ID:         primitive.NewObjectID(),
		First_Name: new(string),
		Last_Name:  new(string),
		Email:      new(string),
		Password:   new(string),
		Phone:      new(string),
	}
	jsonValue, _ := json.Marshal(res)
	req, _ := http.NewRequest("POST", "/api/v1/users/signup", bytes.NewBuffer(jsonValue))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestLogIn(t *testing.T) {

}
