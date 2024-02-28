package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (app *application) routes() http.Handler {
	mux := gin.New()

	mux.Use(gin.Recovery())

	mux.POST("/auth", app.authenticate)
	mux.POST("/refresh-token", app.refresh)
	mux.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"Message": "Hello from Gin",
		})
	})
	users := mux.Group("/users")
	{
		users.GET("/", app.allUsers)
		users.GET("/{userID}", app.getUser)
		users.DELETE("/{userID}", app.deleteUser)
		users.PUT("/", app.insertUser)
		users.PATCH("/", app.updateUser)
	}
	err := mux.Run(":8090")
	if err != nil {
		return nil
	}

	return mux
}
