package main

import (
    "fmt"
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/golang-jwt/jwt/v4"
)

// ValidateToken validates the JWT token from the Authorization header
func ValidateToken(c *gin.Context) (string, string, error) {
    tokenString := c.GetHeader("Authorization")
    if tokenString == "" {
        c.JSON(http.StatusUnauthorized, gin.H{"message": "Missing token"})
        return "", "", fmt.Errorf("missing token")
    }

    // Remove "Bearer " prefix if present
    if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
        tokenString = tokenString[7:]
    }

    // Parse and validate JWT token
    token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
        return jwtKey, nil
    })
    if err != nil || !token.Valid {
        c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid token"})
        return "", "", fmt.Errorf("invalid token")
    }

    // Extract claims
    claims, ok := token.Claims.(*Claims)
    if !ok {
        c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid token claims"})
        return "", "", fmt.Errorf("invalid token claims")
    }

    return claims.Email, claims.Role, nil
}

// TokenAuthMiddleware authenticates requests using JWT
func TokenAuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        email, role, err := ValidateToken(c)
        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
            c.Abort()
            return
        }
        // Set email and role in context for downstream handlers
        c.Set("email", email)
        c.Set("role", role)
        c.Next()
    }
}

// AdminAuthMiddleware restricts access to admin users
func AdminAuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        email, role, err := ValidateToken(c)
        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
            c.Abort()
            return
        }

        // Verify admin role from token claims
        if role != "admin" {
            c.JSON(http.StatusForbidden, gin.H{"message": "Forbidden"})
            c.Abort()
            return
        }

        // Additional database check for consistency
        var storedRole string
        err = db.QueryRow("SELECT role FROM users WHERE email = ?", email).Scan(&storedRole)
        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid credentials"})
            c.Abort()
            return
        }
        if storedRole != "admin" {
            c.JSON(http.StatusForbidden, gin.H{"message": "Forbidden"})
            c.Abort()
            return
        }

        c.Next()
    }
}

// TeacherAuthMiddleware restricts access to teacher users
func TeacherAuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        email, role, err := ValidateToken(c)
        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
            c.Abort()
            return
        }

        // Verify teacher role from token claims
        if role != "teacher" {
            c.JSON(http.StatusForbidden, gin.H{"message": "Forbidden"})
            c.Abort()
            return
        }

        // Additional database check for consistency
        var storedRole string
        err = db.QueryRow("SELECT role FROM users WHERE email = ?", email).Scan(&storedRole)
        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid credentials"})
            c.Abort()
            return
        }
        if storedRole != "teacher" {
            c.JSON(http.StatusForbidden, gin.H{"message": "Forbidden"})
            c.Abort()
            return
        }

        c.Next()
    }
}