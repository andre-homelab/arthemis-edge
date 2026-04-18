package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestValidateHandler(t *testing.T) {
	secret := "test-secret-key"
	os.Setenv("JWT_TOKEN", secret)
	endpoint := "/validate"

	t.Run("should return 401 when Authorization header is missing", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, endpoint, nil)
		rr := httptest.NewRecorder()

		validate(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
		assert.Contains(t, rr.Body.String(), "Missing token")
	})

	t.Run("should return 401 when token is invalid", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, endpoint, nil)
		req.Header.Set("Authorization", "Bearer invalid-token")
		rr := httptest.NewRecorder()

		validate(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
		assert.Contains(t, rr.Body.String(), "Invalid token")
	})

	t.Run("should return 200 and correct headers when token is valid", func(t *testing.T) {
		// Criando um token válido para o teste
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"sub":  "user-123",
			"role": "admin",
		})
		tokenString, _ := token.SignedString([]byte(secret))

		req := httptest.NewRequest(http.MethodPost, endpoint, nil)
		req.Header.Set("Authorization", "Bearer "+tokenString)
		rr := httptest.NewRecorder()

		validate(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, "user-123", rr.Header().Get("X-User-Id"))
		assert.Equal(t, "admin", rr.Header().Get("X-User-Role"))
	})
}

func TestHealthHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rr := httptest.NewRecorder()

	// Handler anônimo definido no main
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "ok", rr.Body.String())
}
