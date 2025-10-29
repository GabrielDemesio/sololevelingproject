package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type ctxKey string

const UserIDKey ctxKey = "uid"

func JWTMiddleware(secret []byte) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			auth := r.Header.Get("Authorization")
			if !strings.HasPrefix(auth, "Bearer ") {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}
			tokenStr := strings.TrimPrefix(auth, "Bearer ")
			claims := jwt.MapClaims{}
			_, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
				return secret, nil
			})
			if err != nil {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}
			uid := claims["sub"].(string)
			if uid == "" {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}
			ctx := context.WithValue(r.Context(), UserIDKey, uid)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
func UserIDFromContext(r *http.Request) (string, bool) {
	v := r.Context().Value(UserIDKey)
	if v == nil {
		return "", false
	}
	s, _ := v.(string)
	return s, s != ""
}
