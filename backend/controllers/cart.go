package controllers

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"backend/database"
	"backend/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Application struct {
	productCollection *mongo.Collection
	userCollection    *mongo.Collection
	orderCollection   *mongo.Collection
}

func NewApplication(productCollection *mongo.Collection, userCollection *mongo.Collection, orderCollection *mongo.Collection) *Application {
	return &Application{
		productCollection: productCollection,
		userCollection:    userCollection,
		orderCollection:   orderCollection,
	}
}

func (app *Application) AddToCart() gin.HandlerFunc {
	return func(c *gin.Context) {
		var response models.Response
		productQueryId := c.Query("productId")
		if productQueryId == "" {
			log.Println("product id is empty")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("product id is empty"))
			return
		}
		userQueryId := c.Query("userId")
		if userQueryId == "" {
			log.Println("user id is empty")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("user id is empty"))
			return
		}
		productId, err := primitive.ObjectIDFromHex(productQueryId)
		if err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err = database.AddProductToCart(ctx, app.productCollection, app.userCollection, productId, userQueryId)
		if err != nil {
			response.Status = "Failed"
			response.Code = http.StatusInternalServerError
			response.Msg = err.Error()
			c.IndentedJSON(http.StatusInternalServerError, response)
			return
		}

		response.Status = "OK"
		response.Code = 200
		response.Msg = "Successfully added to the cart"
		c.IndentedJSON(200, response)
		return
	}
}

func (app *Application) RemoveItem() gin.HandlerFunc {
	return func(c *gin.Context) {
		var response models.Response
		productQueryId := c.Query("productId")
		if productQueryId == "" {
			log.Println("Missing product id")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("Missing product id"))
			return
		}

		userQueryId := c.Query("userId")
		if userQueryId == "" {
			log.Println("Missing user id")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("Missing user id"))
		}

		productId, err := primitive.ObjectIDFromHex(productQueryId)
		if err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		err = database.RemoveCartItem(ctx, app.productCollection, app.userCollection, productId, userQueryId)
		if err != nil {
			response.Status = "Failed"
			response.Code = http.StatusInternalServerError
			response.Msg = err.Error()
			c.IndentedJSON(http.StatusInternalServerError, response)
			return
		}

		response.Status = "OK"
		response.Code = 200
		response.Msg = "Successfully removed from cart"
		c.IndentedJSON(200, response)
		return
	}
}

func GetItemsFromCart() gin.HandlerFunc {
	return func(c *gin.Context) {
		var response models.Response
		userId := c.Query("userId")
		if userId == "" {
			c.Header("Content-Type", "application/json")

			response.Status = "Failed"
			response.Code = http.StatusNotFound
			response.Msg = "Missing userId"
			c.IndentedJSON(http.StatusNotFound, response)
			c.Abort()
			return
		}

		usertId, _ := primitive.ObjectIDFromHex(userId)

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var filledCart models.User
		err := UserCollection.FindOne(ctx, bson.D{primitive.E{Key: "_id", Value: usertId}}).Decode(&filledCart)
		if err != nil {
			log.Println(err)

			response.Status = "Failed"
			response.Code = 500
			response.Msg = "User id not found"
			c.IndentedJSON(500, response)
			return
		}

		filterMatch := bson.D{{Key: "$match", Value: bson.D{primitive.E{Key: "_id", Value: usertId}}}}
		unwind := bson.D{{Key: "$unwind", Value: bson.D{primitive.E{Key: "path", Value: "$user_cart"}}}}
		grouping := bson.D{{Key: "$group", Value: bson.D{primitive.E{Key: "_id", Value: "$_id"}, {Key: "total", Value: bson.D{primitive.E{Key: "$sum", Value: "$user_cart.price"}}}}}}
		pointCursor, err := UserCollection.Aggregate(ctx, mongo.Pipeline{filterMatch, unwind, grouping})
		if err != nil {
			log.Println(err)
		}
		var listing []bson.M
		if err = pointCursor.All(ctx, &listing); err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
		}

		for _, json := range listing {
			response.Status = "OK"
			response.Code = 200
			response.Msg = json["total"]
			response.Data = filledCart.UserCart
			c.IndentedJSON(200, response)
			return
		}

		ctx.Done()
	}
}

func (app *Application) BuyFromCart() gin.HandlerFunc {
	return func(c *gin.Context) {
		var response models.Response
		userQueryId := c.Query("userId")
		if userQueryId == "" {
			log.Panicln("user id is empty")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("UserID is empty"))
		}
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		err := database.BuyItemFromCart(ctx, app.userCollection, userQueryId, app.orderCollection)
		if err != nil {
			response.Status = "Failed"
			response.Code = http.StatusInternalServerError
			response.Msg = err.Error()
			c.IndentedJSON(http.StatusInternalServerError, response)
			return
		}

		response.Status = "OK"
		response.Code = 200
		response.Msg = "Successfully placed the order"
		c.IndentedJSON(200, response)
		return
	}
}

// func (app *Application) InstantBuy() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		var response models.Response
// 		userQueryId := c.Query("userId")
// 		if userQueryId == "" {
// 			log.Println("UserID is empty")
// 			_ = c.AbortWithError(http.StatusBadRequest, errors.New("UserID is empty"))
// 		}
// 		productQueryId := c.Query("productId")
// 		if productQueryId == "" {
// 			log.Println("Product_ID id is empty")
// 			_ = c.AbortWithError(http.StatusBadRequest, errors.New("product_id is empty"))
// 		}
// 		productId, err := primitive.ObjectIDFromHex(productQueryId)
// 		if err != nil {
// 			log.Println(err)
// 			c.AbortWithStatus(http.StatusInternalServerError)
// 			return
// 		}

// 		var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
// 		defer cancel()
// 		err = database.InstantBuyer(ctx, app.productCollection, app.userCollection, productId, userQueryId)
// 		if err != nil {
// 			response.Status = "Failed"
// 			response.Code = http.StatusInternalServerError
// 			response.Msg = err.Error()
// 			c.IndentedJSON(http.StatusInternalServerError, response)
// 			return
// 		}

// 		response.Status = "OK"
// 		response.Code = 200
// 		response.Msg = "Successfully placed the order"
// 		c.IndentedJSON(200, response)
// 		return
// 	}
// }
