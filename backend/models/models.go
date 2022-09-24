package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	Id             primitive.ObjectID `json:"_id" bson:"_id"`
	FirstName      string             `json:"first_name" validate:"required,min=2,max=30"`
	LastName       string             `json:"last_name"  validate:"required,min=2,max=30"`
	Password       string             `json:"password"   validate:"required,min=6"`
	Phone          string             `json:"phone"      validate:"required"`
	Token          string             `json:"token"`
	RefreshToken   string             `josn:"refresh_token"`
	CreatedAt      time.Time          `json:"created_at"`
	UpdatedAt      time.Time          `json:"updated_at"`
	UserId         string             `json:"user_id"`
	UserCart       []Product          `json:"user_cart" bson:"user_cart"`
	AddressDetails []Address          `json:"addresses" bson:"addresses"`
	Orders         []Order            `json:"orders" bson:"orders"`
}

type Product struct {
	ProductId   primitive.ObjectID `bson:"_id"`
	ProductName string             `json:"product_name"`
	Price       uint64             `json:"price"`
	Rating      float32            `json:"rating"`
	Image       string             `json:"image"`
	Comments    []Comment          `json:"comments" bson:"comments"`
}

type Address struct {
	AddressId primitive.ObjectID `bson:"_id"`
	House     string             `json:"house" bson:"house"`
	Street    string             `json:"street" bson:"street"`
	Ward      string             `json:"ward" bson:"ward"`
	District  string             `json:"district" bson:"district"`
	City      string             `json:"city" bson:"city"`
}

type Order struct {
	OrderId       primitive.ObjectID `bson:"_id"`
	OrderCart     []Product          `json:"order_list"  bson:"order_list"`
	OrdereredAt   time.Time          `json:"ordered_at"  bson:"ordered_at"`
	Price         uint64             `json:"total_price" bson:"total_price"`
	Discount      int                `json:"discount"    bson:"discount"`
	PaymentMethod Payment            `json:"payment_method" bson:"payment_method"`
}

type Payment struct {
	Digital bool `json:"digital" bson:"digital"`
	COD     bool `json:"cod"     bson:"cod"`
}

type Comment struct {
	CommentId primitive.ObjectID `bson:"_id"`
	UserId    string             `json:"user_id" bson:"user_id"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	Content   string             `json:"content" bson:"content"`
}

type Response struct {
	Status string      `json:"status"`
	Code   uint        `json:"code"`
	Msg    interface{} `json:"message"`
	Data   interface{} `json:"data"`
}
