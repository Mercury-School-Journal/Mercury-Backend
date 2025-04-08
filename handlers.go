package main
import (
	"time"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func RegisterUser(c *gin.Context) {

	var request AccessRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input"})
		return
	}
	var user = request.Argument.(User)

	if user.Email == "" || user.Password == "" || user.Role == "" || user.FirstName == "" || user.LastName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Email, password, first name, last name and role are required"})
		return
	}

	var existingUser User
	err := db.QueryRow("SELECT uid, email, password FROM users WHERE email = ?", user.Email).Scan(&existingUser.UID, &existingUser.Email, &existingUser.Password)
	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"message": "Email already taken"})
		return
	}

	hashedPassword, err := HashPassword(user.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error hashing password"})
		return
	}
	user.Password = hashedPassword

	_, err = db.Exec("INSERT INTO users (email, password, role) VALUES (?, ?, ?)", user.Email, user.Password, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error saving user"})
		return
	}
	
	var userID int64
	err = db.QueryRow("SELECT uid FROM users WHERE email = ?", user.Email).Scan(&userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error retrieving user ID"})
		return
	}
	
	if user.Role == "student" {
		_, err = db.Exec("INSERT INTO students (user_id, first_name, last_name) VALUES (?, ?, ?)", userID, user.FirstName, user.LastName)
	} else if user.Role == "teacher" {
		_, err = db.Exec("INSERT INTO teachers (user_id, first_name, last_name) VALUES (?, ?, ?)", userID, user.FirstName, user.LastName)
	}
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error saving user details"})
		return
	}
	
	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully"})
}

func Login(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input"})
		return
	}

	var storedUser User
	err := db.QueryRow("SELECT uid, email, password FROM users WHERE email = ?", user.Email).Scan(&storedUser.UID, &storedUser.Email, &storedUser.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid credentials sql"})
		return
	}

	if !CheckPasswordHash(user.Password, storedUser.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid credentials password"})
		return
	}

	expirationTime := time.Now().Add(7 * 24 * time.Hour)
	claims := &Claims{
		Username: storedUser.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Could not generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}

func ChangePassword(c *gin.Context) {
	var input Input
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input"})
		return
	}

	username, _ := c.Get("username")

	var user User
	err := db.QueryRow("SELECT uid, email, password FROM users WHERE email = ?", username).Scan(&user.UID, &user.Email, &user.Password)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
		return
	}

	if !CheckPasswordHash(input.OldPassword, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Incorrect old password"})
		return
	}

	hashedPassword, err := HashPassword(input.NewPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error hashing new password"})
		return
	}

	_, err = db.Exec("UPDATE users SET password = ? WHERE email = ?", hashedPassword, user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error updating password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password changed successfully"})
}

func DeleteAccount(c *gin.Context) {
	username, _ := c.Get("username")
	_, err := db.Exec("DELETE FROM users WHERE email = ?", username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error deleting user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

func Ping(c *gin.Context){
	c.Status(http.StatusNoContent)
}
func AddTimetableEntry(c *gin.Context) {
	var timetableEntry TimetableEntry
	if err := c.ShouldBindJSON(&timetableEntry); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input"})
		return
	}

	if timetableEntry.Subject == "" || timetableEntry.StartTime == "" || timetableEntry.EndTime == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Subject, start time and end time are required"})
		return
	}

	_, err := db.Exec("INSERT INTO timetable (subject, start_time, end_time) VALUES (?, ?, ?)", timetableEntry.Subject, timetableEntry.StartTime, timetableEntry.EndTime)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error saving timetable entry"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Timetable entry created successfully"})
}
func AddGrade(c *gin.Context) {
	var grade Grade
	if err := c.ShouldBindJSON(&grade); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input"})
		return
	}

	if grade.UserID == 0 || grade.Subject == "" || grade.Grade == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "User ID, subject and grade are required"})
		return
	}

	_, err := db.Exec("INSERT INTO grades (user_id, subject, grade) VALUES (?, ?, ?)", grade.UserID, grade.Subject, grade.Grade)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error saving grade"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Grade created successfully"})
}
func GetTimetable (c *gin.Context) {
	var timetable []TimetableEntry
	rows, err := db.Query("SELECT subject, start_time, end_time FROM timetable")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error retrieving timetable"})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var entry TimetableEntry
		if err := rows.Scan(&entry.Subject, &entry.StartTime, &entry.EndTime); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Error scanning timetable entry"})
			return
		}
		timetable = append(timetable, entry)
	}

	c.JSON(http.StatusOK, timetable)
}
func GetGrades(c *gin.Context) {
	userID := c.Param("user_id")
	var grades []Grade
	rows, err := db.Query("SELECT user_id, subject, grade FROM grades WHERE user_id = ?", userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error retrieving grades"})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var grade Grade
		if err := rows.Scan(&grade.UserID, &grade.Subject, &grade.Grade); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Error scanning grade"})
			return
		}
		grades = append(grades, grade)
	}

	c.JSON(http.StatusOK, grades)
}