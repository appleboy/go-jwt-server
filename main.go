package main

import (
	"fmt"
	"github.com/appleboy/gin-jwt-server/config"
	"github.com/appleboy/gin-jwt-server/model"
	"github.com/appleboy/gin-jwt-server/input"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"
)

const (
	JWTSigningKey string        = "appleboy"
	ExpireTime    time.Duration = time.Minute * 60 * 24 * 30
	Realm         string        = "jwt auth"
)

var (
	orm *xorm.Engine
)

func AbortWithError(c *gin.Context, code int, message string) {
	c.Header("WWW-Authenticate", "JWT realm="+Realm)
	c.JSON(code, gin.H{
		"code":    code,
		"message": message,
	})
	c.Abort()
}

func Auth(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		_, err := jwt.ParseFromRequest(c.Request, func(token *jwt.Token) (interface{}, error) {
			b := ([]byte(JWTSigningKey))

			fmt.Println(b)
			return b, nil
		})

		if err != nil {
			AbortWithError(c, http.StatusUnauthorized, "Invaild User Token")
			return
		}
	}
}

func LoginHandler(c *gin.Context) {
	var form input.Login
	var user model.User

	if c.BindJSON(&form) != nil {
		AbortWithError(c, http.StatusBadRequest, "Missing usename or password")
		return
	}

	found, err := orm.Where("username = ?", form.Username).Get(&user)

	if err != nil {
		AbortWithError(c, http.StatusInternalServerError, "DB Query Error")
		return
	}

	if found == false || bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(form.Password)) != nil {
		AbortWithError(c, http.StatusUnauthorized, "Incorrect Username / Password")
		return
	}

	expire := time.Now().Add(ExpireTime)

	// Create the token
	token := jwt.New(jwt.SigningMethodHS256)
	// Set some claims
	token.Claims["id"] = form.Username
	token.Claims["exp"] = expire.Unix()
	// Sign and get the complete encoded token as a string
	tokenString, err := token.SignedString([]byte(JWTSigningKey))

	if err != nil {
		AbortWithError(c, http.StatusUnauthorized, "Create JWT Token faild")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token":    tokenString,
		"expire":   expire.Format(time.RFC3339),
	})
}

func RegisterHandler(c *gin.Context) {
	var form input.Login
	var user model.User

	if c.BindJSON(&form) != nil {
		AbortWithError(c, http.StatusBadRequest, "Missing usename or password")
		return
	}

	has, err := orm.Where("username = ?", form.Username).Get(&user)

	if has {
		AbortWithError(c, http.StatusBadRequest, "username is already exist.")
		return
	}

	userId := uuid.NewV4().String()

	if digest, err := bcrypt.GenerateFromPassword([]byte(form.Password), bcrypt.DefaultCost); err != nil {
		AbortWithError(c, http.StatusInternalServerError, err.Error())
		return
	} else {
		form.Password = string(digest)
	}

	_, err = orm.Insert(&model.User{
		Id:       userId,
		Username: form.Username,
		Password: form.Password,
	})

	if err != nil {
		AbortWithError(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    "200",
		"message": "ok",
	})
}

func RefreshHandler(c *gin.Context) {
	expire := time.Now().Add(ExpireTime)

	// Create the token
	token := jwt.New(jwt.SigningMethodHS256)
	// Set some claims
	token.Claims["exp"] = expire.Unix()
	// Sign and get the complete encoded token as a string
	tokenString, err := token.SignedString([]byte(JWTSigningKey))

	if err != nil {
		AbortWithError(c, http.StatusUnauthorized, "Create JWT Token faild")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token":  tokenString,
		"expire": expire.Format(time.RFC3339),
	})
}

func HelloHandler(c *gin.Context) {
	currentTime := time.Now()
	currentTime.Format(time.RFC3339)
	c.JSON(200, gin.H{
		"current_time": currentTime,
		"text":         "You are login now.",
	})
}

func initDB() {
	configs := config.ReadConfig("config.json")

	connectStr := &mysql.Config{
		User:   configs.DB_USERNAME,
		Passwd: configs.DB_PASSWORD,
		Net:    "tcp",
		Addr:   net.JoinHostPort(configs.DB_HOST, strconv.Itoa(configs.DB_PORT)),
		DBName: configs.DB_NAME,
		Params: map[string]string{
			"charset": "utf8",
		},
	}

	db, err := xorm.NewEngine("mysql", connectStr.FormatDSN())

	if err != nil {
		log.Panic("DB connection initialization failed", err)
	}

	orm = db

	_, err = orm.Insert(&model.User{
		Id: uuid.NewV4().String(),
	})

	if err != nil {
		fmt.Println(err)
		return
	}
}

func main() {
	port := os.Getenv("PORT")
	r := gin.Default()

	if port == "" {
		port = "8000"
	}

	// initial DB setting
	initDB()

	r.POST("/login", LoginHandler)
	r.POST("/register", RegisterHandler)

	auth := r.Group("/auth")
	auth.Use(Auth("test"))
	{
		auth.GET("/hello", HelloHandler)
		auth.GET("/refresh_token", RefreshHandler)
	}

	r.Run(":" + port)
}
