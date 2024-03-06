package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (app *application) enableCORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:8090")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, PATCH, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}

}

func (app *application) authRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, _, err := app.getTokenFromHeaderAndVerify(c)
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
		}

		c.Next()
	}
}
