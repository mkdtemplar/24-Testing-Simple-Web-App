package main

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
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
	err := c.Request.ParseForm()
	if err != nil {
		c.Writer.WriteHeader(http.StatusBadRequest)
		return
	}

	refreshToken := c.Request.Form.Get("refresh_token")
	claims := &Claims{}

	_, err = jwt.ParseWithClaims(refreshToken, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(app.JWTSecret), nil
	})

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	if time.Unix(claims.ExpiresAt.Unix(), 0).Sub(time.Now()) > 30*time.Second {
		c.JSON(http.StatusTooEarly, gin.Error{
			Err: errors.New("refresh token not required to be renewed"),
		})
		return
	}

	userId, err := strconv.Atoi(claims.Subject)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	user, err := app.DB.GetUser(userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.Error{
			Err: errors.New("unknown user"),
		})
		return
	}

	tokenPairs, err := app.generateTokenPairs(user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	c.SetCookie("__Host-refresh_token", tokenPairs.RefreshToken, int(refreshTokenExpiry.Seconds()),
		"/", "localhost", true, true)

	c.JSON(http.StatusOK, tokenPairs)
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
