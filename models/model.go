package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID              primitive.ObjectID `json:"_id" bson:"_id"`
	First_Name      *string            `json:"first_name" validate:"required,min=2,max=30"`
	Last_Name       *string            `json:"last_name" validate:"required,min=2,max=30"`
	Email           *string            `json:"email" validate:"email, required"`
	Password        *string            `json:"password" validate:"required,min=8"`
	Phone           *string            `json:"phone" validate:"required"`
	Token           *string            `json:"token"`
	Refresh_Token   *string            `json:"refresh_token"`
	Created_At      time.Time          `json:"created_at"`
	Update_At       time.Time          `json:"update_at"`
	User_ID         string             `json:"user_id"`
	UserCart        []ProductUser      `json:"user_cart" bson:"usercart"`
	Address_Details []Address          `json:"address_details" bson:"address"`
	Order_Status    []Order            `json:"order_status" bson:"orders"`
}

type Product struct {
	ProductID    primitive.ObjectID `bson:"_id"`
	Product_Name *string            `json:"product_name"`
	Price        *uint64            `json:"price"`
	Rating       *uint8             `json:"rating"`
	Image        *string            `json:"image"`
}

type ProductUser struct {
	ProductID    primitive.ObjectID `bson:"_id"`
	Product_Name *string            `json:"product_name" bson:"product_name"`
	Price        int                `json:"price" bson:"price"`
	Rating       *uint8             `json:"rating" bson:"rating"`
	Image        *string            `json:"image" bson:"image"`
}

type Address struct {
	Address_id  primitive.ObjectID `bson:"_id"`
	Description *string            `json:"description" bson:"description"`
	Street      *string            `json:"street" bson:"street_name"`
	City        *string            `json:"city,omitempty" bson:"city_name"`
	Country     *string            `json:"country,omitempty" bson:"country_name"`
	Pincode     *string            `json:"pincode,omitempty" bson:"pincode"`
}

type Order struct {
	Order_ID       primitive.ObjectID `bson:"_id"`
	Order_Cart     []ProductUser      `json:"order_cart,omitempty" bson:"order_list"`
	Ordered_At     time.Time          `json:"ordered_at,omitempty" bson:"ordered_at"`
	Price          int                `json:"price,omitempty" bson:"total_price"`
	Discount       *int               `json:"discount,omitempty" bson:"discount"`
	Payment_Method Payment            `json:"payment_method,omitempty" bson:"payment_method"`
}

type Payment struct {
	Digital bool
	COD     bool
}
