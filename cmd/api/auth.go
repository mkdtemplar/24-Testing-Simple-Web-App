package main

import (
	"24-Testing-Simple-Web-App/pkg/data"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const (
	jwtTokenExpiry     = time.Minute * 15
	refreshTokenExpiry = time.Hour * 24
)

type TokenPairs struct {
	Token        string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type Claims struct {
	UserName string `json:"name"`
	jwt.RegisteredClaims
}

func (app *application) getTokenFromHeaderAndVerify(c *gin.Context) (string, *Claims, error) {
	c.Writer.Header().Set("Vary", "Authorization")
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return "", nil, errors.New("no authorization header")
	}

	headerParts := strings.Split(authHeader, " ")
	if len(headerParts) != 2 {
		return "", nil, errors.New("invalid authorization header")
	}

	if headerParts[0] != "Bearer" {
		return "", nil, errors.New("no bearer in  authorization header")
	}

	token := headerParts[1]

	claims := &Claims{}
	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method %v", token.Header["alg"])
		}
		return []byte(app.JWTSecret), nil
	})

	if err != nil {
		if strings.HasPrefix(err.Error(), "token is expired by") {
			return "", nil, errors.New("token expired")
		}
		return "", nil, err
	}

	if claims.Issuer != app.Domain {
		return "", nil, errors.New("incorrect issuer")
	}

	return token, claims, err
}

func (app *application) generateTokenPairs(user *data.User) (TokenPairs, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["name"] = fmt.Sprintf("%s %s", user.FirstName, user.LastName)
	claims["sub"] = fmt.Sprint(user.ID)
	claims["aud"] = app.Domain
	claims["iss"] = app.Domain

	if user.IsAdmin == 1 {
		claims["admin"] = true
	} else {
		claims["admin"] = false
	}

	claims["exp"] = time.Now().Add(jwtTokenExpiry).Unix()

	signedAccessToken, err := token.SignedString([]byte(app.JWTSecret))
	if err != nil {
		return TokenPairs{}, err
	}

	// create the refresh token
	refreshToken := jwt.New(jwt.SigningMethodHS256)
	refreshTokenClaims := refreshToken.Claims.(jwt.MapClaims)
	refreshTokenClaims["sub"] = fmt.Sprint(user.ID)

	// set expiry; must be longer than jwt expiry
	refreshTokenClaims["exp"] = time.Now().Add(refreshTokenExpiry).Unix()

	// create signed refresh token
	signedRefreshToken, err := refreshToken.SignedString([]byte(app.JWTSecret))
	if err != nil {
		return TokenPairs{}, err
	}

	var tokenPairs = TokenPairs{
		Token:        signedAccessToken,
		RefreshToken: signedRefreshToken,
	}

	return tokenPairs, nil
}
