package main

import (
	"github.com/gin-gonic/gin"
)

func (app *application) routes() *gin.Engine {
	mux := gin.New()
	//cors := app.enableCORS(mux)

	mux.Use(gin.Recovery())
	mux.Use(app.enableCORS())

	mux.POST("/auth", app.authenticate)
	mux.POST("/refresh-token", app.refresh)
	users := mux.Group("/users", app.authRequired())
	{
		users.GET("/", app.allUsers)
		users.GET("/:userID", app.getUser)
		users.DELETE("/:userID", app.deleteUser)
		users.PUT("/", app.insertUser)
		users.PATCH("/", app.updateUser)
	}

	return mux
}
