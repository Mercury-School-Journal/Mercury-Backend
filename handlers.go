package main

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

// RegisterUser handles user registration
func RegisterUser(c *gin.Context) {
	var request AccessRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input"})
		return
	}

	// Expect a map containing user and person details
	userData, ok := request.Argument.(map[string]interface{})
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid argument format"})
		return
	}

	user := User{
		Email:    userData["email"].(string),
		Password: userData["password"].(string),
		Role:     userData["role"].(string),
	}
	person := Person{
		FirstName: userData["first_name"].(string),
		LastName:  userData["last_name"].(string),
		BirthDate: userData["birth_date"].(string),
		Address:   userData["address"].(string),
		Phone:     userData["phone"].(string),
	}

	// Validate required fields
	if user.Email == "" || user.Password == "" || user.Role == "" || person.FirstName == "" || person.LastName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Email, password, role, first name, and last name are required"})
		return
	}

	// Check if email is already taken
	var existingUser User
	err := db.QueryRow("SELECT uid, email, password FROM users WHERE email = ?", user.Email).Scan(&existingUser.UID, &existingUser.Email, &existingUser.Password)
	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"message": "Email already taken"})
		return
	}
	if err != sql.ErrNoRows {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error checking email"})
		return
	}

	// Hash the password
	hashedPassword, err := HashPassword(user.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error hashing password"})
		return
	}
	user.Password = hashedPassword

	// Begin transaction to insert user and person
	tx, err := db.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error starting transaction"})
		return
	}
	defer tx.Rollback()

	// Insert into users table
	result, err := tx.Exec("INSERT INTO users (email, password, role) VALUES (?, ?, ?)", user.Email, user.Password, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error saving user"})
		return
	}

	// Get user ID
	userID, err := result.LastInsertId()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error retrieving user ID"})
		return
	}

	// Insert into persons table
	_, err = tx.Exec("INSERT INTO persons (user_id, first_name, last_name, birth_date, address, phone) VALUES (?, ?, ?, ?, ?, ?)",
		userID, person.FirstName, person.LastName, person.BirthDate, person.Address, person.Phone)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error saving user details"})
		return
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error committing transaction"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully"})
}

// Login handles user authentication
func Login(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input"})
		return
	}

	// Query user by email
	var storedUser User
	var role string
	err := db.QueryRow("SELECT uid, email, password, role FROM users WHERE email = ?", user.Email).Scan(&storedUser.UID, &storedUser.Email, &storedUser.Password, &role)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid email credentials"})
		return
	}
	// Verify password
	if !CheckPasswordHash(user.Password, storedUser.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid password credentials"})
		return
	}

	// Generate JWT
	expirationTime := time.Now().Add(7 * 24 * time.Hour)
	claims := &Claims{
		Email: storedUser.Email,
		Role:  role,
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

// ChangePassword handles password updates
func ChangePassword(c *gin.Context) {
	var input Input
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input"})
		return
	}

	// Get email from JWT claims
	email, _ := c.Get("email")

	// Query user by email
	var user User
	err := db.QueryRow("SELECT uid, email, password FROM users WHERE email = ?", email).Scan(&user.UID, &user.Email, &user.Password)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
		return
	}

	// Verify old password
	if !CheckPasswordHash(input.OldPassword, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Incorrect old password"})
		return
	}

	// Hash new password
	hashedPassword, err := HashPassword(input.NewPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error hashing new password"})
		return
	}

	// Update password
	_, err = db.Exec("UPDATE users SET password = ? WHERE email = ?", hashedPassword, user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error updating password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password changed successfully"})
}

// DeleteAccount handles user account deletion
func DeleteAccount(c *gin.Context) {
	// Get email from JWT claims
	email, _ := c.Get("email")

	// Begin transaction
	tx, err := db.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error starting transaction"})
		return
	}
	defer tx.Rollback()

	// Delete from persons table
	_, err = tx.Exec("DELETE FROM persons WHERE user_id IN (SELECT uid FROM users WHERE email = ?)", email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error deleting user details"})
		return
	}

	// Delete from users table
	_, err = tx.Exec("DELETE FROM users WHERE email = ?", email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error deleting user"})
		return
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error committing transaction"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

// Ping responds with a no-content status
func Ping(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

// AddTimetableEntry adds a new timetable entry
func AddTimetableEntry(c *gin.Context) {
	var entry TimetableEntry
	if err := c.ShouldBindJSON(&entry); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input"})
		return
	}

	// Validate required fields
	if entry.SubjectID == 0 || entry.StartTime == "" || entry.EndTime == "" || entry.TeacherID == 0 || entry.ClassName == "" || entry.Day == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Subject ID, start time, end time, teacher ID, class name, and day are required"})
		return
	}

	// Insert into timetable table
	_, err := db.Exec("INSERT INTO timetable (day, subject_id, time_start, time_end, room, teacher_id, class_name) VALUES (?, ?, ?, ?, ?, ?, ?)",
		entry.Day, entry.SubjectID, entry.StartTime, entry.EndTime, entry.Room, entry.TeacherID, entry.ClassName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error saving timetable entry"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Timetable entry created successfully"})
}

// AddGrade adds a new grade, comment, or custom value
func AddGrade(c *gin.Context) {
	var grade Grade
	if err := c.ShouldBindJSON(&grade); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input"})
		return
	}

	// Validate required fields
	if grade.UserID == 0 || grade.SubjectID == 0 || grade.Grade == "" || grade.GradeType == "" || grade.Date == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "User ID, subject ID, grade, grade type, and date are required"})
		return
	}

	// Insert into grades table
	_, err := db.Exec("INSERT INTO grades (user_id, subject_id, grade, grade_type, date) VALUES (?, ?, ?, ?, ?)",
		grade.UserID, grade.SubjectID, grade.Grade, grade.GradeType, grade.Date)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Grade created successfully"})
}

// GetTimetable retrieves the timetable
func GetTimetable(c *gin.Context) {
	var timetable []TimetableEntry
	rows, err := db.Query("SELECT id, day, subject_id, time_start, time_end, room, teacher_id, class_name FROM timetable")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error retrieving timetable"})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var entry TimetableEntry
		if err := rows.Scan(&entry.ID, &entry.Day, &entry.SubjectID, &entry.StartTime, &entry.EndTime, &entry.Room, &entry.TeacherID, &entry.ClassName); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Error scanning timetable entry"})
			return
		}
		timetable = append(timetable, entry)
	}

	c.JSON(http.StatusOK, timetable)
}

// GetGrades retrieves grades for a specific user
func GetGrades(c *gin.Context) {
	userID := c.Param("user_id")
	var grades []Grade
	rows, err := db.Query("SELECT id, user_id, subject_id, grade, grade_type, date FROM grades WHERE user_id = ?", userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error retrieving grades"})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var grade Grade
		if err := rows.Scan(&grade.ID, &grade.UserID, &grade.SubjectID, &grade.Grade, &grade.GradeType, &grade.Date); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Error scanning grade"})
			return
		}
		grades = append(grades, grade)
	}

	c.JSON(http.StatusOK, grades)
}
