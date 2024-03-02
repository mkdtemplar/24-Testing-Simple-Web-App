package main

import (
	"24-Testing-Simple-Web-App/pkg/repository/dbrepo"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
)

var app application
var expiredToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZG1pbiI6dHJ1ZSwiYXVkIjoiZXhhbXBsZS5jb20iLCJleHAiOjE3MDk1NTMxMjUsImlzcyI6ImV4YW1wbGUuY29tIiwibmFtZSI6IkpvaG4gRG9lIiwic3ViIjoiMSJ9.Xa4qLlOUWf6R-ttL9_eZlCK8Qwm4tiG13nHtUUVZyuI"

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	app.DB = &dbrepo.TestDBRepo{}
	app.Domain = "example.com"
	app.JWTSecret = "2dce505d96a53c5768052ee90f3df2055657518dad489160df9913f66042e160"
	os.Exit(m.Run())

}
