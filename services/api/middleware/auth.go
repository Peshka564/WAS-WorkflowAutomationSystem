package middleware

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/Peshka564/WAS-WorkflowAutomationSystem/services/api/utils"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			utils.SendError(w, http.StatusUnauthorized, "Authorization header required")
			return
		}

		bearer, tokenString, found := strings.Cut(authHeader, " ")
		if !found || bearer != "Bearer" {
			utils.SendError(w, http.StatusUnauthorized, "Invalid authorization format")
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		fmt.Println(token)
		if err != nil || !token.Valid {
			fmt.Println(token.Valid)
			utils.SendError(w, http.StatusUnauthorized, "Invalid or expired token")
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			if floatID, ok := claims["user_id"].(float64); ok {
				ctx := context.WithValue(r.Context(), "user_id", int64(floatID))
				next.ServeHTTP(w, r.WithContext(ctx))
			} else {
				utils.SendError(w, http.StatusUnauthorized, "Invalid token claims")
				return
			}
		} else {
			utils.SendError(w, http.StatusUnauthorized, "Invalid token claims")
			return
		}
	})
}