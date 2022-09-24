package main

import (
	"log"
	"os"

	"backend/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	router := gin.New()

	router.Use(gin.Logger())
	routes.Routes(router)

	log.Fatal(router.Run(":" + port))
}
