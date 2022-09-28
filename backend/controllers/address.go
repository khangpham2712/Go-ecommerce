package controllers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"backend/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func AddAddress() gin.HandlerFunc {
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
		address, err := primitive.ObjectIDFromHex(userId)
		if err != nil {
			response.Status = "Failed"
			response.Code = 500
			response.Msg = "Internal Server Error"
			c.IndentedJSON(500, response)
			return
		}
		var addresses models.Address
		addresses.AddressId = primitive.NewObjectID()
		if err = c.BindJSON(&addresses); err != nil {
			response.Status = "Failed"
			response.Code = http.StatusNotAcceptable
			response.Msg = "Invalid address"
			c.IndentedJSON(http.StatusNotAcceptable, response)
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		matchFilter := bson.D{{Key: "$match", Value: bson.D{primitive.E{Key: "_id", Value: address}}}}
		unwind := bson.D{{Key: "$unwind", Value: bson.D{primitive.E{Key: "path", Value: "$addresses"}}}}
		group := bson.D{{Key: "$group", Value: bson.D{primitive.E{Key: "_id", Value: "$address_id"}, {Key: "count", Value: bson.D{primitive.E{Key: "$sum", Value: 1}}}}}}

		pointCursor, err := UserCollection.Aggregate(ctx, mongo.Pipeline{matchFilter, unwind, group})
		if err != nil {
			response.Status = "Failed"
			response.Code = 500
			response.Msg = "Internal Server Error"
			c.IndentedJSON(500, response)
			return
		}

		var addressInfo []bson.M
		if err = pointCursor.All(ctx, &addressInfo); err != nil {
			panic(err)
		}

		var size int32
		for _, address_no := range addressInfo {
			count := address_no["count"]
			size = count.(int32)
		}
		if size < 2 {
			filter := bson.D{primitive.E{Key: "_id", Value: address}}
			update := bson.D{{Key: "$push", Value: bson.D{primitive.E{Key: "addresses", Value: addresses}}}}
			_, err := UserCollection.UpdateOne(ctx, filter, update)
			if err != nil {
				fmt.Println(err)
			}
		} else {
			response.Status = "Failed"
			response.Code = 400
			response.Msg = "Not allowed"
			c.IndentedJSON(400, response)
			return
		}

		ctx.Done()

		response.Status = "OK"
		response.Code = http.StatusOK
		response.Msg = "Successfully"
		c.IndentedJSON(http.StatusOK, response)
		return
	}
}

func EditHomeAddress() gin.HandlerFunc {
	return func(c *gin.Context) {
		var response models.Response
		userId := c.Query("userId")
		if userId == "" {
			c.Header("Content-Type", "application/json")

			response.Status = "Failed"
			response.Code = http.StatusNotFound
			response.Msg = "Invalid"
			c.IndentedJSON(http.StatusNotFound, response)
			c.Abort()
			return
		}
		usertId, err := primitive.ObjectIDFromHex(userId)
		if err != nil {
			response.Status = "Failed"
			response.Code = 500
			response.Msg = err.Error()
			c.IndentedJSON(500, response)
			return
		}
		var editAddress models.Address
		if err := c.BindJSON(&editAddress); err != nil {
			response.Status = "Failed"
			response.Code = http.StatusBadRequest
			response.Msg = err.Error()
			c.IndentedJSON(http.StatusBadRequest, response)
			return
		}
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		filter := bson.D{primitive.E{Key: "_id", Value: usertId}}
		update := bson.D{{Key: "$set", Value: bson.D{primitive.E{Key: "addresses.0.house", Value: editAddress.House}, {Key: "addresses.0.street", Value: editAddress.Street}, {Key: "addresses.0.ward", Value: editAddress.Ward}, {Key: "addresses.0.district", Value: editAddress.District}, {Key: "addresses.0.city", Value: editAddress.City}}}}
		_, err = UserCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			response.Status = "Failed"
			response.Code = 500
			response.Msg = "Something went wrong"
			c.IndentedJSON(500, response)
			return
		}

		ctx.Done()

		response.Status = "OK"
		response.Code = 200
		response.Msg = "Successfully updated home address"
		c.IndentedJSON(200, response)
		return
	}
}

func EditWorkAddress() gin.HandlerFunc {
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
		usertId, err := primitive.ObjectIDFromHex(userId)
		if err != nil {
			response.Status = "Failed"
			response.Code = 500
			response.Msg = err.Error()
			c.IndentedJSON(500, response)
			return
		}
		var editAddress models.Address
		if err := c.BindJSON(&editAddress); err != nil {
			response.Status = "Failed"
			response.Code = http.StatusBadRequest
			response.Msg = err.Error()
			c.IndentedJSON(http.StatusBadRequest, response)
			return
		}
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		filter := bson.D{primitive.E{Key: "_id", Value: usertId}}
		update := bson.D{{Key: "$set", Value: bson.D{primitive.E{Key: "addresses.1.house", Value: editAddress.House}, {Key: "addresses.1.street", Value: editAddress.Street}, {Key: "addresses.1.ward", Value: editAddress.Ward}, {Key: "addresses.1.district", Value: editAddress.District}, {Key: "addresses.1.city", Value: editAddress.City}}}}
		_, err = UserCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			response.Status = "Failed"
			response.Code = 500
			response.Msg = "Something went wrong"
			c.IndentedJSON(500, response)
			return
		}

		ctx.Done()

		response.Status = "OK"
		response.Code = 200
		response.Msg = "Successfully updated work address"
		c.IndentedJSON(200, response)
		return
	}
}

func DeleteAddress() gin.HandlerFunc {
	return func(c *gin.Context) {
		var response models.Response
		userId := c.Query("userId")
		if userId == "" {
			c.Header("Content-Type", "application/json")

			response.Status = "Failed"
			response.Code = http.StatusNotFound
			response.Msg = "Invalid searched index"
			c.IndentedJSON(http.StatusNotFound, response)
			c.Abort()
			return
		}
		addresses := make([]models.Address, 0)
		usertId, err := primitive.ObjectIDFromHex(userId)
		if err != nil {
			response.Status = "Failed"
			response.Code = 500
			response.Msg = "Internal Server Error"
			c.IndentedJSON(500, response)
			return
		}
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		filter := bson.D{primitive.E{Key: "_id", Value: usertId}}
		update := bson.D{{Key: "$set", Value: bson.D{primitive.E{Key: "addresses", Value: addresses}}}}
		_, err = UserCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			response.Status = "Failed"
			response.Code = 404
			response.Msg = "Wrong"
			c.IndentedJSON(404, response)
			return
		}

		ctx.Done()

		response.Status = "OK"
		response.Code = 200
		response.Msg = "Successfully deleted the address(es)"
		c.IndentedJSON(200, response)
		return
	}
}
