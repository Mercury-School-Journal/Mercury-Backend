package main
import (
	"github.com/golang-jwt/jwt/v4"
)

type User struct {
	UID       uint   `json:"uid"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}