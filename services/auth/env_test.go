package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetEnv(t *testing.T) {
	key := "TEST_VAR"
	expectedValue := "test_value"
	os.Setenv(key, expectedValue)
	defer os.Unsetenv(key)

	t.Run("should return the correct environment variable value", func(t *testing.T) {
		value := GetEnv(key)
		assert.Equal(t, expectedValue, value)
	})
}
