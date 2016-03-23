package main

import (
	"encoding/json"
	"github.com/appleboy/gin-jwt-server/tests"
	"github.com/appleboy/gin-status-api"
	"github.com/gin-gonic/gin"
	"github.com/icrowley/fake"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

var (
	username string = fake.FullName()
	password string = "1234"
	token    string
)

func TestRegisterHandler(t *testing.T) {
	initDB()

	// Missing usename or password
	data := `{"username":"` + username + `"}`
	tests.RunSimplePost("/register", data,
		RegisterHandler,
		func(r *httptest.ResponseRecorder) {
			var rd map[string]interface{}
			err := json.NewDecoder(r.Body).Decode(&rd)

			if err != nil {
				log.Fatalf("JSON Decode fail: %v", err)
			}

			assert.Equal(t, rd["message"], "Missing usename or password")
			assert.Equal(t, r.Code, http.StatusBadRequest)
		})

	// Register success.
	data = `{"username":"` + username + `","password":"` + password + `"}`
	tests.RunSimplePost("/register", data,
		RegisterHandler,
		func(r *httptest.ResponseRecorder) {
			var rd map[string]interface{}
			err := json.NewDecoder(r.Body).Decode(&rd)

			if err != nil {
				log.Fatalf("JSON Decode fail: %v", err)
			}

			assert.Equal(t, rd["message"], "ok")
			assert.Equal(t, r.Code, http.StatusOK)
		})

	// Username is already exist.
	data = `{"username":"` + username + `","password":"` + password + `"}`
	tests.RunSimplePost("/register", data,
		RegisterHandler,
		func(r *httptest.ResponseRecorder) {
			var rd map[string]interface{}
			err := json.NewDecoder(r.Body).Decode(&rd)

			if err != nil {
				log.Fatalf("JSON Decode fail: %v", err)
			}

			assert.Equal(t, rd["message"], "Username is already exist")
			assert.Equal(t, r.Code, http.StatusBadRequest)
		})
}

func TestLoginHandler(t *testing.T) {
	initDB()

	// Missing usename or password
	data := `{"username":"` + username + `"}`
	tests.RunSimplePost("/login", data,
		LoginHandler,
		func(r *httptest.ResponseRecorder) {
			var rd map[string]interface{}
			err := json.NewDecoder(r.Body).Decode(&rd)

			if err != nil {
				log.Fatalf("JSON Decode fail: %v", err)
			}

			assert.Equal(t, rd["message"], "Missing usename or password")
			assert.Equal(t, r.Code, http.StatusBadRequest)
		})

	// incorrect password
	data = `{"username":"` + username + `","password":"test"}`
	tests.RunSimplePost("/login", data,
		LoginHandler,
		func(r *httptest.ResponseRecorder) {
			var rd map[string]interface{}
			err := json.NewDecoder(r.Body).Decode(&rd)

			if err != nil {
				log.Fatalf("JSON Decode fail: %v", err)
			}

			assert.Equal(t, rd["message"], "Incorrect Username / Password")
			assert.Equal(t, r.Code, http.StatusUnauthorized)
		})

	// login success
	data = `{"username":"` + username + `","password":"` + password + `"}`
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
			assert.Equal(t, r.Code, http.StatusOK)

			token = rd["token"].(string)
		})
}

func performRequest(r http.Handler, method, path string, token string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	return w
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

	w := performRequest(r, "GET", "/v1/hello", token)
	assert.Equal(t, w.Code, http.StatusOK)

	w = performRequest(r, "GET", "/v1/refresh_token", token)
	assert.Equal(t, w.Code, http.StatusOK)

	w = performRequest(r, "GET", "/v1/refresh_token", "1234")
	assert.Equal(t, w.Code, http.StatusUnauthorized)
}

func TestStatusHandler(t *testing.T) {
	initDB()

	gin.SetMode(gin.TestMode)
	r := gin.New()
	api := r.Group("/api")
	api.Use(Auth())
	{
		api.GET("/status", status.StatusHandler)
	}

	w := performRequest(r, "GET", "/api/status", token)
	assert.Equal(t, w.Code, http.StatusOK)
}
