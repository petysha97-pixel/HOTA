package service

import (
	"context"
	"log"
	"net/http"
	"strings"

	_ "modernc.org/sqlite"
)

type contextKey string

const userIDKey contextKey = "userID"

// CORSMiddleware — единственный, что должен остаться из «методных» мидлварей.
// Оборачивает весь мукс один раз, всё остальное делает роутер Go 1.22.
func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenHeader := r.Header.Get("Authorization")
		if tokenHeader == "" {
			http.Error(w, "Отсутствует токен", http.StatusUnauthorized)
			return
		}

		token := strings.TrimPrefix(tokenHeader, "Bearer ")

		userID, err := ValidateToken(token)
		if err != nil {
			log.Printf("ошибка валидации токена: %v", err)
			http.Error(w, "ошибка валидации токена", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), userIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
