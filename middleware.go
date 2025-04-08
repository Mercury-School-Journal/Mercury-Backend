package main
import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func ValidateToken(c *gin.Context) (string, error) {
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Missing token"})
		return "", fmt.Errorf("missing token")
	}

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid token"})
		return "", fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid token"})
		return "", fmt.Errorf("invalid token")
	}

	return claims.Username, nil
}
func TokenAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		username, err := ValidateToken(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
			c.Abort()
			return
		}
		c.Set("username", username)
		c.Next()
	}
}
func AdminAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		username, err := ValidateToken(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
			c.Abort()
			return
		}
		var storedUser User
		err = db.QueryRow("SELECT role FROM users WHERE email = ?", username).Scan(&storedUser.Role)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid credentials"})
			c.Abort()
			return
		}
		if storedUser.Role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"message": "Forbidden"})
			c.Abort()
			return
		}
		c.Next()
	}
}
func TeacherAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		username, err := ValidateToken(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
			c.Abort()
			return
		}
		var storedUser User
		err = db.QueryRow("SELECT role FROM users WHERE email = ?", username).Scan(&storedUser.Role)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid credentials"})
			c.Abort()
			return
		}
		if storedUser.Role != "teacher" {
			c.JSON(http.StatusForbidden, gin.H{"message": "Forbidden"})
			c.Abort()
			return
		}
		c.Next()
	}
}