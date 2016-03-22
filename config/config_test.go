package config

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"os"
	"testing"
)

// Test file is missing
func TestMissingFile(t *testing.T) {
	filename := "test"
	_, err := ReadConfig(filename)

	assert.NotNil(t, err)
}

// Test wrong json format
func TestWrongJSONFormat(t *testing.T) {
	content := []byte(`{"DB_HOST": "127.0.0.1""DB_USERNAME": "root","DB_PASSWORD": "","DB_PORT": 3306,"DB_NAME": "test"}`)
	filename := "tempfile"

	if err := ioutil.WriteFile(filename, content, 0644); err != nil {
		log.Fatalf("WriteFile %s: %v", filename, err)
	}

	// clean up
	defer os.Remove(filename)

	// parse JSON format error
	_, err := ReadConfig(filename)

	assert.NotNil(t, err)
}

// Test config file.
func TestReadConfig(t *testing.T) {
	content := []byte(`{"DB_HOST": "127.0.0.1","DB_USERNAME": "root","DB_PASSWORD": "","DB_PORT": 3306,"DB_NAME": "test"}`)
	filename := "tempfile"

	if err := ioutil.WriteFile(filename, content, 0644); err != nil {
		log.Fatalf("WriteFile %s: %v", filename, err)
	}

	// clean up
	defer os.Remove(filename)

	configs, err := ReadConfig(filename)

	assert.Nil(t, err)
	assert.Equal(t, configs.DB_HOST, "127.0.0.1")
	assert.Equal(t, configs.DB_USERNAME, "root")
	assert.Empty(t, configs.DB_PASSWORD)
	assert.Equal(t, configs.DB_PORT, 3306)
	assert.Equal(t, configs.DB_NAME, "test")
}
