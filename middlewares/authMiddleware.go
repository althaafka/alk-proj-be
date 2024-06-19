package middlewares

import (
	"net/http"
	"strings"
	"context"

	"github.com/althaafka/alk-proj-be.git/helpers"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bearerToken := r.Header.Get("Authorization")
		if bearerToken == "" {
			helpers.RespondWithJSON(w, http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
			return
		}

		tokenString := strings.TrimPrefix(bearerToken, "Bearer ")

		claims, err := helpers.ValidateToken(tokenString)
		if err != nil {
			helpers.RespondWithJSON(w, http.StatusUnauthorized, map[string]string{"error": "Invalid token"})
			return
		}

		ctx := context.WithValue(r.Context(), "userID", claims.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
