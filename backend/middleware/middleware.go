package middleware

import (
	"net/http"

	"backend/models"
	token "backend/tokens"

	"github.com/gin-gonic/gin"
)

func Authorization() gin.HandlerFunc {
	return func(c *gin.Context) {
		var response models.Response
		clientToken := c.Request.Header.Get("token")
		if clientToken == "" {
			response.Status = "Failed"
			response.Code = http.StatusInternalServerError
			response.Msg = "No access token provided"
			c.IndentedJSON(http.StatusInternalServerError, response)
			c.Abort()
			return
		}
		claims, err := token.ValidateToken(clientToken)
		if err != "" {
			response.Status = "Failed"
			response.Code = http.StatusInternalServerError
			response.Msg = err
			c.IndentedJSON(http.StatusInternalServerError, response)
			c.Abort()
			return
		}

		c.Set("phone", claims.Phone)
		c.Set("uid", claims.Uid)
		c.Next()
	}
}
