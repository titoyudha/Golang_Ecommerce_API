package controllers

import (
	"go_ecommerce/service"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/charge"
)

func Charge() gin.HandlerFunc {
	return func(c *gin.Context) {
		var json service.ChargeJSON
		err := c.BindJSON(&json)
		if err != nil {
			c.IndentedJSON(http.StatusNotAcceptable, gin.H{"Message": err})
			log.Panic(err)
			return
		}
		apiKey := os.Getenv("SK_TEST_KEY")
		stripe.Key = apiKey

		_, err = charge.New(&stripe.ChargeParams{
			Amount:       stripe.Int64(json.Amount),
			Currency:     stripe.String(string(stripe.CurrencyIDR)),
			Source:       &stripe.SourceParams{Token: stripe.String("token_visa")},
			ReceiptEmail: stripe.String(json.ReceiptEmail),
		})
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"Message": "Request Charged Failed"})
			return
		}
		c.IndentedJSON(http.StatusOK, gin.H{"Message": "Charged Successfull"})
	}
}
