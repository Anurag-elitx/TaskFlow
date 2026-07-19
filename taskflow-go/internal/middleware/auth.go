package middleware

import (
	"net/http"
	"strings"
	"taskflow/internal/handlers"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			handlers.RespondError(c, handlers.NewAppError(http.StatusUnauthorized, "Missing Authorization header"))
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			handlers.RespondError(c, handlers.NewAppError(http.StatusUnauthorized, "Invalid Authorization header format"))
			c.Abort()
			return
		}

		tokenString := parts[1]
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtSecret), nil
		})

		if err != nil || !token.Valid {
			handlers.RespondError(c, handlers.NewAppError(http.StatusUnauthorized, "Invalid or expired token"))
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			handlers.RespondError(c, handlers.NewAppError(http.StatusUnauthorized, "Invalid token claims"))
			c.Abort()
			return
		}

		c.Set("user_id", claims["user_id"])
		c.Next()
	}
}
