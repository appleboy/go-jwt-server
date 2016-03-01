package main

import (
	"testing"
	"net/http/httptest"
	"encoding/json"
	"log"
	"net/http"
	"github.com/appleboy/gin-jwt-server/tests"
	"github.com/stretchr/testify/assert"
	"github.com/icrowley/fake"
	"github.com/gin-gonic/gin"
)

var (
	username string = fake.FullName()
	password string = "1234"
	token string
)

func TestRegisterHandler(t *testing.T) {
	initDB()

	// Missing usename or password
	data := `{"username":"`+username+`"}`
	tests.RunSimplePost("/register", data,
		RegisterHandler,
		func(r *httptest.ResponseRecorder) {
			var rd map[string]interface{}
			err := json.NewDecoder(r.Body).Decode(&rd)

			if err != nil {
				log.Fatalf("JSON Decode fail: %v", err)
			}

			assert.Equal(t, rd["message"], "Missing usename or password")
			assert.Equal(t, r.Code, 400)
		})

	// Register success.
	data = `{"username":"` + username + `","password":"`+password+`"}`
	tests.RunSimplePost("/register", data,
		RegisterHandler,
		func(r *httptest.ResponseRecorder) {
			var rd map[string]interface{}
			err := json.NewDecoder(r.Body).Decode(&rd)

			if err != nil {
				log.Fatalf("JSON Decode fail: %v", err)
			}

			assert.Equal(t, rd["message"], "ok")
			assert.Equal(t, r.Code, 200)
		})

	// Username is already exist.
	data = `{"username":"`+username+`","password":"`+password+`"}`
	tests.RunSimplePost("/register", data,
		RegisterHandler,
		func(r *httptest.ResponseRecorder) {
			var rd map[string]interface{}
			err := json.NewDecoder(r.Body).Decode(&rd)

			if err != nil {
				log.Fatalf("JSON Decode fail: %v", err)
			}

			assert.Equal(t, rd["message"], "Username is already exist")
			assert.Equal(t, r.Code, 400)
		})
}

func TestLoginHandler(t *testing.T) {
	initDB()

	// Missing usename or password
	data := `{"username":"`+username+`"}`
	tests.RunSimplePost("/login", data,
		LoginHandler,
		func(r *httptest.ResponseRecorder) {
			var rd map[string]interface{}
			err := json.NewDecoder(r.Body).Decode(&rd)

			if err != nil {
				log.Fatalf("JSON Decode fail: %v", err)
			}

			assert.Equal(t, rd["message"], "Missing usename or password")
			assert.Equal(t, r.Code, 400)
		})

	// incorrect password
	data = `{"username":"`+username+`","password":"test"}`
	tests.RunSimplePost("/login", data,
		LoginHandler,
		func(r *httptest.ResponseRecorder) {
			var rd map[string]interface{}
			err := json.NewDecoder(r.Body).Decode(&rd)

			if err != nil {
				log.Fatalf("JSON Decode fail: %v", err)
			}

			assert.Equal(t, rd["message"], "Incorrect Username / Password")
			assert.Equal(t, r.Code, 401)
		})

	// login success
	data = `{"username":"`+username+`","password":"`+password+`"}`
	tests.RunSimplePost("/login", data,
		LoginHandler,
		func(r *httptest.ResponseRecorder) {
			var rd map[string]interface{}
			err := json.NewDecoder(r.Body).Decode(&rd)

			if err != nil {
				log.Fatalf("JSON Decode fail: %v", err)
			}

			assert.Contains(t, "token", r.Body.String())
			assert.Contains(t, "expire", r.Body.String())
			assert.Equal(t, r.Code, 200)

			token = rd["token"].(string)
		})
}

func Result(t *testing.T, router *gin.Engine, path string, token string, code int) {
	// RUN
	req, err := http.NewRequest("GET", path, nil)

	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Authorization", "Bearer " + token)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// TEST
	assert.Equal(t, w.Code, code)
}

func TestHelloHandler(t *testing.T) {
	initDB()

	gin.SetMode(gin.TestMode)
	r := gin.New()
	v1 := r.Group("/v1")
	v1.Use(Auth())
	{
		v1.GET("/hello", HelloHandler)
		v1.GET("/refresh_token", RefreshHandler)
	}

	Result(t, r, "/v1/hello", token, 200)
	Result(t, r, "/v1/refresh_token", token, 200)
	Result(t, r, "/v1/refresh_token", "1234", 401)
}
