package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	_ "github.com/andre-homelab/arthemis-edge/docs"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/golang-jwt/jwt/v5"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title           Authentication API
// @version         1.0
// @description     API service for JWT token validation and authentication.
// @termsOfService  http://swagger.io/terms/

// @license.name  MIT
// @license.url   http://opensource.org/licenses/MIT

// @host      localhost:6769
// @BasePath  /

// @securityDefinitions.apikey BearerAuth
// @in                         header
// @name                       Authorization
// @description                Type "Bearer " followed by your JWT token.

// validate checks the JWT token validity and returns user claims in headers.
// @Summary      Validate JWT Token
// @Description  Parses the Authorization header, validates the HMAC signature, and returns user identity in response headers.
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string  true  "Bearer <token>"
// @Success      200  {string}  string    "Token is valid. User info returned in X-User-Id and X-User-Role headers."
// @Failure      401  {string}  string    "Unauthorized: Missing or invalid token"
// @Router       /validate [post]
func validate(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")

	if authHeader == "" {
		http.Error(w, "Missing token", http.StatusUnauthorized)
		return
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	var secret = []byte(GetEnv("JWT_TOKEN"))

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secret, nil
	})

	if err != nil || !token.Valid {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		http.Error(w, "Invalid claims", http.StatusUnauthorized)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	w.Header().Set("X-User-Id", fmt.Sprintf("%v", claims["sub"]))
	w.Header().Set("X-User-Role", fmt.Sprintf("%v", claims["role"]))

	w.WriteHeader(http.StatusOK)
}

func main() {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RealIP)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	r.Route("/validate", func(r chi.Router) {
		r.Post("/", validate)
	})

	// Rota do Swagger UI
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("doc.json"),
	))

	slog.Info("Server running on", "port", 6769)
	slog.Info("API documentation on", "URL", "http://localhost:6769/swagger/index.html")
	if err := http.ListenAndServe(":6769", r); err != nil {
		slog.Error("Failed to start server", "err", err)
	}
}
