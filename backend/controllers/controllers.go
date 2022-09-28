package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"backend/database"
	"backend/models"
	generate "backend/tokens"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var UserCollection *mongo.Collection = database.UserData(database.Client, "Users")
var ProductCollection *mongo.Collection = database.ProductData(database.Client, "Products")
var OrderCollection *mongo.Collection = database.OrderData(database.Client, "Orders")
var Validate = validator.New()

func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err)
	}
	return string(bytes)
}

func VerifyPassword(userPassword string, givenPassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(userPassword), []byte(givenPassword))
	valid := true
	msg := ""
	if err != nil {
		msg = "Username or password is incorrect"
		valid = false
	}
	return valid, msg
}

func SignUp() gin.HandlerFunc {
	return func(c *gin.Context) {
		var response models.Response
		var user models.User
		if err := c.BindJSON(&user); err != nil {
			response.Status = "Failed"
			response.Code = http.StatusBadRequest
			response.Msg = err.Error()
			c.IndentedJSON(http.StatusBadRequest, response)
			return
		}
		validationErr := Validate.Struct(user)
		if validationErr != nil {
			response.Status = "Failed"
			response.Code = http.StatusBadRequest
			response.Msg = validationErr.Error()
			c.IndentedJSON(http.StatusBadRequest, response)
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		count, err := UserCollection.CountDocuments(ctx, bson.M{"phone": user.Phone})
		if err != nil {
			log.Panic(err)
			response.Status = "Failed"
			response.Code = http.StatusInternalServerError
			response.Msg = err.Error()
			c.IndentedJSON(http.StatusInternalServerError, response)
			return
		}
		if count > 0 {
			response.Status = "Failed"
			response.Code = http.StatusBadRequest
			response.Msg = "User already exists"
			c.IndentedJSON(http.StatusBadRequest, response)
			return
		}

		password := HashPassword(user.Password)
		user.Password = password

		user.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.Id = primitive.NewObjectID()
		user.UserId = user.Id.Hex()
		token, refreshToken, _ := generate.TokenGenerator(user.Phone, user.FirstName, user.LastName, user.UserId)
		user.Token = token
		user.RefreshToken = refreshToken
		user.UserCart = make([]models.Product, 0)
		user.AddressDetails = make([]models.Address, 0)
		user.Orders = make([]models.Order, 0)
		_, inserterr := UserCollection.InsertOne(ctx, user)
		if inserterr != nil {
			response.Status = "Failed"
			response.Code = http.StatusInternalServerError
			response.Msg = "Something went wrong. Not created yet"
			c.IndentedJSON(http.StatusInternalServerError, response)
			return
		}

		response.Status = "OK"
		response.Code = http.StatusCreated
		response.Msg = "Successfully signed up"
		c.IndentedJSON(http.StatusCreated, response)
		return
	}
}

func LogIn() gin.HandlerFunc {
	return func(c *gin.Context) {
		var response models.Response
		var user models.User
		var founduser models.User
		if err := c.BindJSON(&user); err != nil {
			response.Status = "Failed"
			response.Code = http.StatusBadRequest
			response.Msg = "Invalid input"
			c.IndentedJSON(http.StatusBadRequest, response)
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		err := UserCollection.FindOne(ctx, bson.M{"phone": user.Phone}).Decode(&founduser)
		if err != nil {
			response.Status = "Failed"
			response.Code = http.StatusInternalServerError
			response.Msg = "Username or password is incorrect"
			c.IndentedJSON(http.StatusInternalServerError, response)
			return
		}
		passwordIsValid, msg := VerifyPassword(founduser.Password, user.Password)
		if !passwordIsValid {
			response.Status = "Failed"
			response.Code = http.StatusInternalServerError
			response.Msg = msg
			c.IndentedJSON(http.StatusInternalServerError, response)
			fmt.Println(msg)
			return
		}
		token, refreshToken, _ := generate.TokenGenerator(founduser.Phone, founduser.FirstName, founduser.LastName, founduser.UserId)

		updateErr := generate.UpdateAllTokens(token, refreshToken, founduser.UserId)
		if updateErr != nil {
			response.Status = "Failed"
			response.Code = http.StatusInternalServerError
			response.Msg = updateErr.Error()
			c.IndentedJSON(http.StatusInternalServerError, response)
			return
		}

		var usert models.User
		UserCollection.FindOne(ctx, bson.M{"phone": user.Phone}).Decode(&usert)
		response.Status = "OK"
		response.Code = http.StatusFound
		response.Msg = "Successfully"
		response.Data = gin.H{"username": usert.Phone,
			"userId":        usert.UserId,
			"access_token":  usert.Token,
			"refresh_token": usert.RefreshToken,
		}
		c.IndentedJSON(http.StatusFound, response)
		return
	}
}

func GetAllProducts() gin.HandlerFunc {
	return func(c *gin.Context) {
		var response models.Response
		var productList []models.Product

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		cursor, err := ProductCollection.Find(ctx, bson.D{{}})
		if err != nil {
			response.Status = "Failed"
			response.Code = http.StatusInternalServerError
			response.Msg = "Something went wrong. Please try again later"
			c.IndentedJSON(http.StatusInternalServerError, response)
			return
		}
		err = cursor.All(ctx, &productList)
		if err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		defer cursor.Close(ctx)
		if err := cursor.Err(); err != nil {
			log.Println(err)
			response.Status = "Failed"
			response.Code = 400
			response.Msg = "Invalid"
			c.IndentedJSON(400, response)
			return
		}

		response.Status = "OK"
		response.Code = 200
		response.Msg = "Successfully"
		response.Data = productList
		c.IndentedJSON(200, response)
		return
	}
}

func SearchProductByQuery() gin.HandlerFunc {
	return func(c *gin.Context) {
		var response models.Response
		var searchedProducts []models.Product
		queryParam := c.Query("name")
		if queryParam == "" {
			log.Println("query is empty")
			c.Header("Content-Type", "application/json")

			response.Status = "Failed"
			response.Code = http.StatusNotFound
			response.Msg = "Invalid searched index"
			c.IndentedJSON(http.StatusNotFound, response)
			c.Abort()
			return
		}
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		searchQueryDB, err := ProductCollection.Find(ctx, bson.M{"product_name": bson.M{"$regex": queryParam}})
		if err != nil {
			response.Status = "Failed"
			response.Code = 404
			response.Msg = "Something went wrong in fetching the db queries"
			c.IndentedJSON(404, response)
			return
		}
		err = searchQueryDB.All(ctx, &searchedProducts)
		if err != nil {
			log.Println(err)

			response.Status = "Failed"
			response.Code = 400
			response.Msg = "Invalid"
			c.IndentedJSON(400, response)
			return
		}
		defer searchQueryDB.Close(ctx)
		if err := searchQueryDB.Err(); err != nil {
			log.Println(err)

			response.Status = "Failed"
			response.Code = 400
			response.Msg = "Invalid request"
			c.IndentedJSON(400, response)
			return
		}

		response.Status = "OK"
		response.Code = 200
		response.Msg = "Successfully"
		response.Data = searchedProducts
		c.IndentedJSON(200, response)
		return
	}
}

func GetAllOrders() gin.HandlerFunc {
	return func(c *gin.Context) {
		var response models.Response
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var orderList []models.Order
		cursor, err := OrderCollection.Find(ctx, bson.D{{}})
		if err != nil {
			response.Status = "Failed"
			response.Code = http.StatusInternalServerError
			response.Msg = "Something went wrong. Please try again later"
			c.IndentedJSON(http.StatusInternalServerError, response)
			return
		}
		err = cursor.All(ctx, &orderList)
		if err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		defer cursor.Close(ctx)
		if err := cursor.Err(); err != nil {
			log.Println(err)
			response.Status = "Failed"
			response.Code = 400
			response.Msg = "Invalid"
			c.IndentedJSON(400, response)
			return
		}

		response.Status = "OK"
		response.Code = 200
		response.Msg = "Successfully"
		response.Data = orderList
		c.IndentedJSON(200, response)
		return
	}
}

func ProductAdderAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		var response models.Response
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var product models.Product
		if err := c.BindJSON(&product); err != nil {
			response.Status = "Failed"
			response.Code = http.StatusBadRequest
			response.Msg = err.Error()
			c.IndentedJSON(http.StatusBadRequest, response)
			return
		}
		product.ProductId = primitive.NewObjectID()
		product.Comments = make([]models.Comment, 0)
		_, anyErr := ProductCollection.InsertOne(ctx, product)
		if anyErr != nil {
			response.Status = "Failed"
			response.Code = http.StatusInternalServerError
			response.Msg = "Not created"
			c.IndentedJSON(http.StatusInternalServerError, response)
			return
		}

		response.Status = "OK"
		response.Code = http.StatusOK
		response.Msg = "New product has been successfully added by an admin"
		c.IndentedJSON(http.StatusOK, response)
		return
	}
}

func ProductUpdaterAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		var response models.Response
		productQueryId := c.Query("productId")
		if productQueryId == "" {
			response.Status = "Failed"
			response.Code = http.StatusBadRequest
			response.Msg = "Missing product id"
			c.IndentedJSON(http.StatusBadRequest, response)
			return
		}

		productId, err := primitive.ObjectIDFromHex(productQueryId)
		if err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		var name string = c.PostForm("name")
		if name == "" {
			response.Status = "Failed"
			response.Code = http.StatusBadRequest
			response.Msg = "Missing name"
			c.IndentedJSON(http.StatusBadRequest, response)
			return
		}
		var priceString string = c.PostForm("price")
		if priceString == "" {
			response.Status = "Failed"
			response.Code = http.StatusBadRequest
			response.Msg = "Missing price"
			c.IndentedJSON(http.StatusBadRequest, response)
			return
		}

		price, err := strconv.Atoi(priceString)
		if err != nil {
			response.Status = "Failed"
			response.Code = http.StatusBadRequest
			response.Msg = "Missing price"
			c.IndentedJSON(http.StatusBadRequest, response)
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		filter := bson.D{primitive.E{Key: "_id", Value: productId}}
		update := bson.D{{Key: "$set", Value: bson.D{primitive.E{Key: "product_name", Value: name}, {Key: "price", Value: price}}}}
		_, err = ProductCollection.UpdateOne(ctx, filter, update)
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
		response.Msg = "Successfully updated the product"
		c.IndentedJSON(200, response)
		return
	}
}
