package auth

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

// @title           Authentication API
// @version         1.0
// @description     API para autenticação JWT
// @termsOfService  http://swagger.io/terms/

// @license.name  MIT
// @license.url   http://opensource.org/licenses/MIT

// @host      localhost:6769
// @BasePath  /

var secret = []byte(mustGetEnv("JWT_TOKEN"))

func validate(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")

	if authHeader == "" {
		http.Error(w, "Missing token", http.StatusUnauthorized)
		return
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})

	if err != nil || !token.Valid {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	claims := token.Claims.(jwt.MapClaims)

	w.Header().Set("X-User-Id", fmt.Sprintf("%v", claims["sub"]))
	w.Header().Set("X-User-Role", fmt.Sprintf("%v", claims["role"]))

	w.WriteHeader(http.StatusOK)
}

func main() {
	http.HandleFunc("/validate", validate)
	http.ListenAndServe(":6769", nil)
}
