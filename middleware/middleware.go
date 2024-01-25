package middleware

import (
	"estiam/logger"
	"net/http"
	"time"
)

func LoggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        logger.Logger.Printf("Time: %v, Method: %s, Path: %s\n", time.Now().Format(time.RFC3339), r.Method, r.URL.Path)
        next.ServeHTTP(w, r)
    })
}
const AuthToken = "YouCanTryButYouNeedAToken"

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")

		if token != AuthToken {
			http.Error(w, "Accès non autorisé", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}