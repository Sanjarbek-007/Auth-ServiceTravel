package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

var secretKey = []byte("your-secret-key")

// AuthMiddleware is a Gin middleware function for JWT token authentication
func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Get token from Authorization header
		tokenString := ctx.GetHeader("Authorization")
		

		// Check if token is missing
		if tokenString == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Authorization header is missing"})
			ctx.Abort()
			return
		}

		// Parse JWT token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Validate signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			// Provide the secret key used to sign the token
			return secretKey, nil
		})

		// Handle token parsing errors
		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid token signature"})
				ctx.Abort()
				return
			}
			ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid token"})
			ctx.Abort()
			return
		}

		// Check token validity
		if !token.Valid {
			ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Token is invalid"})
			ctx.Abort()
			return
		}

		// Token is valid, proceed to the next handler
		ctx.Next()
	}
}
