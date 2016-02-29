package main

import (
	"testing"
	"net/http/httptest"
	"encoding/json"
	"log"
	"github.com/appleboy/gin-jwt-server/tests"
	"github.com/stretchr/testify/assert"
	"github.com/icrowley/fake"
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
