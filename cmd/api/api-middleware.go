package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (app *application) enableCORS(next http.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:8090")
		if c.Request.Method == "OPTIONS" {
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
			c.Writer.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS")
			c.Writer.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, X-CSRF-Token, Authorization")
			return
		} else {
			next.ServeHTTP(c.Writer, c.Request)
		}
	}

}

func (app *application) authRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, _, err := app.getTokenFromHeaderAndVerify(c.Writer, c.Request)
		if err != nil {
			c.Writer.WriteHeader(http.StatusUnauthorized)
			return
		}

		c.Next()
		return
	}
}
