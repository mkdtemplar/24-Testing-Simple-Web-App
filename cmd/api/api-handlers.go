package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type Credentials struct {
	UserName string `json:"email"`
	Password string `json:"password"`
}

func (app *application) authenticate(c *gin.Context) {
	var creds Credentials
	if err := c.BindJSON(&creds); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"Message": err})
		return
	}

	user, err := app.DB.GetUserByEmail(creds.UserName)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"Message": err})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"Message": err})
		return
	}

	tokenPairs, err := app.generateTokenPairs(user)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"Message": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"Token Pairs": tokenPairs})
}

func (app *application) refresh(c *gin.Context) {

}

func (app *application) allUsers(c *gin.Context) {

}

func (app *application) getUser(c *gin.Context) {

}

func (app *application) updateUser(c *gin.Context) {

}

func (app *application) deleteUser(c *gin.Context) {

}

func (app *application) insertUser(c *gin.Context) {

}
