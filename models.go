package main
import (
	"github.com/golang-jwt/jwt/v4"
)

type User struct {
	UID       uint   `json:"uid"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}
type AccessRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Argument       interface{}   `json:"argument"`
}

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}
type Input struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}
type TimetableEntry struct {
	Subject string `json:"subject"`
	Teacher string `json:"teacher"`
	StartTime string `json:"start_time"`
	EndTime string `json:"end_time"`
	Room string `json:"room"`
	Day string `json:"day"`
	Class string `json:"class"`
}
type Grade struct {
	UserID uint `json:"user_id"`
	Subject string `json:"subject"`
	Grade string `json:"grade"`
	
}