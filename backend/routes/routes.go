package routes

import (
	"backend/controllers"
	"backend/database"
	"backend/middleware"

	"github.com/gin-gonic/gin"
)

func Routes(router *gin.Engine) {
	app := controllers.NewApplication(database.ProductData(database.Client, "Products"), database.UserData(database.Client, "Users"), database.UserData(database.Client, "Orders"))

	router.POST("/user/sign-up", controllers.SignUp())
	router.POST("/user/log-in", controllers.LogIn())
	router.GET("/user/view-products", controllers.GetAllProducts())
	router.GET("/user/search", controllers.SearchProductByQuery())

	router.GET("/admin/view-orders", controllers.GetAllOrders())
	router.POST("/admin/add-product", controllers.ProductAdderAdmin())
	router.PATCH("/admin/update-product", controllers.ProductUpdaterAdmin())

	router.Use(middleware.Authorization())

	router.GET("/user/list-cart", controllers.GetItemsFromCart())
	router.POST("/user/add-address", controllers.AddAddress())
	router.PATCH("/user/edit-home-address", controllers.EditHomeAddress())
	router.PATCH("/user/edit-work-address", controllers.EditWorkAddress())
	router.DELETE("/user/delete-addresses", controllers.DeleteAddress())

	router.PATCH("/user/add-to-cart", app.AddToCart())
	router.PATCH("/user/remove-item", app.RemoveItem())
	router.GET("/user/cart-checkout", app.BuyFromCart())
	// router.GET("/user/instant-buy", app.InstantBuy())
}
