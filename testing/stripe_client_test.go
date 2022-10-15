package testing

// import (
// 	"go_ecommerce/service"
// 	"testing"

// 	assert "github.com/stretchr/testify/require"
// 	stripe "github.com/stripe/stripe-go/v73"
// 	_ "github.com/stripe/stripe-go/v73/testing"
// )

// func TestChargeCapture(t *testing.T) {
// 	charge, err := service.Capture("ch_123", &stripe.ChargeCaptureParams{
// 		Amount: stripe.Int64(123),
// 	})
// 	assert.Nil(t, err)
// 	assert.NotNil(t, charge)
// }

// func TestChargeGet(t *testing.T) {
// 	charge, err := service.Get("ch_123", nil)
// 	assert.Nil(t, err)
// 	assert.NotNil(t, charge)
// }

// func TestChargeList(t *testing.T) {
// 	i := service.List(&stripe.ChargeListParams{})

// 	//verify that we do at least one charge
// 	assert.True(t, i.Next())
// 	assert.Nil(t, i.Err())
// 	assert.NotNil(t, i.Charge())
// 	assert.NotNil(t, i.ChargeList())
// }

// func TestChargeSearch(t *testing.T) {
// 	i := service.Search(&stripe.ChargeSearchParams{SearchParams: stripe.SearchParams{
// 		Query: "currency:\"IDR\"",
// 	}})

// 	//Verify that we got some charge
// 	assert.Equal(t, *i.Meta().TotalCount, uint32(1))
// 	assert.True(t, i.Next())
// 	assert.Nil(t, i.Err())
// 	assert.NotNil(t, i.Charge())
// 	assert.False(t, i.Next())
// }

// func TestChargeNew(t *testing.T) {
// 	charge, err := service.New(&stripe.ChargeParams{
// 		Amount:   stripe.Int64(1000000),
// 		Currency: stripe.String(string(stripe.CurrencyIDR)),
// 		Source: &stripe.PaymentSourceSourceParams{
// 			Token: stripe.String("src_123"),
// 		},
// 		Shipping: &stripe.ShippingDetailsParams{
// 			Address: &stripe.AddressParams{
// 				Line1: stripe.String("line1"),
// 				City:  stripe.String("city"),
// 			},
// 			Carrier: stripe.String("carrier"),
// 			Name:    stripe.String("name"),
// 		},
// 	})
// 	assert.Nil(t, err)
// 	assert.NotNil(t, charge)
// }

// func TestChargeUpdate(t *testing.T) {
// 	charge, err := service.Update("ch_123", &stripe.ChargeParams{
// 		Description: stripe.String("Update Description"),
// 	})
// 	assert.Nil(t, err)
// 	assert.NotNil(t, charge)
// }
