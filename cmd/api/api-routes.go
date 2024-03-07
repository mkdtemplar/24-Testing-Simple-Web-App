package main

import (
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

func (app *application) routes() *gin.Engine {
	mux := gin.New()

	mux.Use(gin.Recovery())
	mux.Use(app.enableCORS())
	mux.Use(static.Serve("/", static.LocalFile("./html/", false)))
	web := mux.Group("/web")
	{
		web.Use(app.authRequired())
		web.POST("/auth", app.authenticate)
		// /refresh-token
		// /logout
	}
	mux.POST("/auth", app.authenticate)
	mux.POST("/refresh-token", app.refresh)
	users := mux.Group("/users")
	{
		users.Use(app.authRequired())
		users.GET("/", app.allUsers)
		users.GET("/:userID", app.getUser)
		users.DELETE("/:userID", app.deleteUser)
		users.PUT("/", app.insertUser)
		users.PATCH("/", app.updateUser)

	}

	return mux
}
