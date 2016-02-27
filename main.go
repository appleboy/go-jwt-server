package main

import (
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"time"
	"fmt"
)

const JWTSigningKey string = "appleboy"

func Auth(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		_, err := jwt.ParseFromRequest(c.Request, func(token *jwt.Token) (interface{}, error) {
			b := ([]byte(JWTSigningKey))

			fmt.Println(b)
			return b, nil
		})

		if err != nil {
			c.AbortWithError(401, err)
		}
	}
}

func loginHandler(c *gin.Context) {
	username := c.DefaultPostForm("username", "test")
	password := c.DefaultPostForm("password", "test")

	// Create the token
	token := jwt.New(jwt.SigningMethodHS256)
	// Set some claims
	token.Claims["id"] = username
	token.Claims["exp"] = time.Now().Add(time.Minute * 60 * 24 * 30).Unix()
	// Sign and get the complete encoded token as a string
	tokenString, err := token.SignedString([]byte(JWTSigningKey))

	if err != nil {
		c.AbortWithError(401, err)
	}

	c.JSON(http.StatusOK, gin.H{
		"username": username,
		"password": password,
		"token": tokenString,
		"expire": time.Now().Add(time.Minute * 60 * 24 * 30).Format(time.RFC3339),
	})
}

func helloHandler(c *gin.Context) {
	currentTime := time.Now()
	currentTime.Format(time.RFC3339)
	c.JSON(200, gin.H{
		"current_time": currentTime,
		"text":         "You are login now.",
	})
}

func main() {
	port := os.Getenv("PORT")
	r := gin.Default()
	if port == "" {
		port = "8000"
	}
	r.POST("/login", loginHandler)

	auth := r.Group("/auth")
	auth.Use(Auth("test"))
	{
		auth.GET("/hello", helloHandler)
	}

	r.Run(":" + port)
}
