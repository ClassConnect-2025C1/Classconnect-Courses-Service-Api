package middleware

import (
	"net/http"
	"os"
	"strings"
	"templateGo/internals/utils"

	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v4"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Obtain the Authorization header
		authHeader := c.GetHeader("Authorization")

		// Verify that the header exists and has the correct format
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			utils.NewErrorResponse(c, http.StatusUnauthorized, "Unauthorized", "Missing or invalid authorization token")
			c.Abort()
			return
		}

		// Extract the token from the header
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Obtain the secret key from the environment variable
		jwtSecret := os.Getenv("JWT_SECRET")
		if jwtSecret == "" {
			jwtSecret = "tu_clave_secreta" // Fallback por si no está configurada
		}

		// Validar el token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
			// Verificar el método de firma
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}

			// Return the secret key to verify the signature
			return []byte(jwtSecret), nil
		})

		if err != nil || !token.Valid {
			utils.NewErrorResponse(c, http.StatusUnauthorized, "Unauthorized", "Invalid token")
			c.Abort()
			return
		}

		// If the token is valid, extract claims and add them to the context
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			utils.NewErrorResponse(c, http.StatusUnauthorized, "Unauthorized", "Invalid token claims")
			c.Abort()
			return
		}

		// Example: store the user ID in the context
		if userID, exists := claims["user_id"]; exists {
			c.Set("user_id", userID)
		}

		c.Next()
	}
}
