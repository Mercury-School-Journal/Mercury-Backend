package main

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func RegisterUser(c *gin.Context) {
	var request AccessRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input"})
		return
	}

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

	if user.Email == "" || user.Password == "" || user.Role == "" || person.FirstName == "" || person.LastName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Email, password, role, first name, and last name are required"})
		return
	}

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

	hashedPassword, err := HashPassword(user.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error hashing password"})
		return
	}
	user.Password = hashedPassword

	tx, err := db.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error starting transaction"})
		return
	}
	defer tx.Rollback()

	result, err := tx.Exec("INSERT INTO users (email, password, role) VALUES (?, ?, ?)", user.Email, user.Password, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error saving user"})
		return
	}

	userID, err := result.LastInsertId()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error retrieving user ID"})
		return
	}

	_, err = tx.Exec("INSERT INTO persons (user_id, first_name, last_name, birth_date, address, phone) VALUES (?, ?, ?, ?, ?, ?)",
		userID, person.FirstName, person.LastName, person.BirthDate, person.Address, person.Phone)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error saving user details"})
		return
	}

	if err := tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error committing transaction"})
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
	var role string
	err := db.QueryRow("SELECT uid, email, password, role FROM users WHERE email = ?", user.Email).Scan(&storedUser.UID, &storedUser.Email, &storedUser.Password, &role)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid email credentials"})
		return
	}

	if !CheckPasswordHash(user.Password, storedUser.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid password credentials"})
		return
	}

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

func ChangePassword(c *gin.Context) {
	var input Input
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input"})
		return
	}

	email, _ := c.Get("email")

	var user User
	err := db.QueryRow("SELECT uid, email, password FROM users WHERE email = ?", email).Scan(&user.UID, &user.Email, &user.Password)
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
	email, _ := c.Get("email")

	tx, err := db.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error starting transaction"})
		return
	}
	defer tx.Rollback()

	_, err = tx.Exec("DELETE FROM persons WHERE user_id IN (SELECT uid FROM users WHERE email = ?)", email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error deleting user details"})
		return
	}

	_, err = tx.Exec("DELETE FROM users WHERE email = ?", email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error deleting user"})
		return
	}

	if err := tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error committing transaction"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

func Ping(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

func AddTimetableEntry(c *gin.Context) {
	var entry TimetableEntry
	if err := c.ShouldBindJSON(&entry); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input"})
		return
	}

	if entry.SubjectID == 0 || entry.StartTime == "" || entry.EndTime == "" || entry.TeacherID == 0 || entry.ClassName == "" || entry.Day == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Subject ID, start time, end time, teacher ID, class name, and day are required"})
		return
	}

	_, err := db.Exec("INSERT INTO timetable (day, subject_id, time_start, time_end, room, teacher_id, class_name) VALUES (?, ?, ?, ?, ?, ?, ?)",
		entry.Day, entry.SubjectID, entry.StartTime, entry.EndTime, entry.Room, entry.TeacherID, entry.ClassName)
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

	if grade.UserID == 0 || grade.SubjectID == 0 || grade.Grade == "" || grade.GradeType == "" || grade.Date == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "User ID, subject ID, grade, grade type, and date are required"})
		return
	}

	_, err := db.Exec("INSERT INTO grades (user_id, subject_id, grade, grade_type, date) VALUES (?, ?, ?, ?, ?)",
		grade.UserID, grade.SubjectID, grade.Grade, grade.GradeType, grade.Date)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Grade created successfully"})
}

func GetTimetable(c *gin.Context) {
	email, _ := c.Get("email")
	role, _ := c.Get("role")
	var user User
	var classmember ClassMember
	err := db.QueryRow("SELECT uid, email, class_name FROM users INNER JOIN class_members ON user.uid = class_members.user_id WHERE email = ?", email).Scan(&user.UID, &user.Email, &classmember.ClassName)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
		return
	}
	var rows *sql.Rows
	if role == "student" {
		rows, err := db.Query("SELECT id, day, subject_id, time_start, time_end, room, teacher_id, class_name FROM timetable WHERE class_name = ?", classmember.ClassName)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Error retrieving timetable"})
			return
		}
		defer rows.Close()
	} else if role == "teacher" {
		rows, err := db.Query("SELECT id, day, subject_id, time_start, time_end, room, teacher_id, class_name FROM timetable WHERE teacher_id = ?", user.UID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Error retrieving timetable"})
			return
		}
		defer rows.Close()
	}
	var timetable []TimetableEntry
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

func GetGrades(c *gin.Context) {
	email, _ := c.Get("email")

	var user User
	err := db.QueryRow("SELECT uid, email FROM users WHERE email = ?", email).Scan(&user.UID, &user.Email)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
		return
	}

	var grades []Grade
	rows, err := db.Query("SELECT id, user_id, subject_id, grade, grade_type, date FROM grades WHERE user_id = ?", user.UID)
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

func GetUserInfo(c *gin.Context) {
	email, _ := c.Get("email")
	var user User
	err := db.QueryRow("SELECT uid, email FROM users WHERE email = ?", email).Scan(&user.UID, &user.Email)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
		return
	}
	var person Person
	err = db.QueryRow("SELECT first_name, last_name, birth_date, address, phone FROM persons WHERE user_id = ?", user.UID).Scan(&person.FirstName, &person.LastName, &person.BirthDate, &person.Address, &person.Phone)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "User details not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"email":      user.Email,
		"first_name": person.FirstName,
		"last_name":  person.LastName,
		"birth_date": person.BirthDate,
		"address":    person.Address,
		"phone":      person.Phone,
	})
}

func GetSubjects(c *gin.Context) {

	email, _ := c.Get("email")
	var user User
	err := db.QueryRow("SELECT uid, email FROM users WHERE email = ?", email).Scan(&user.UID, &user.Email)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
		return
	}

	var classmember ClassMember
	err = db.QueryRow("SELECT class_name FROM class_members WHERE user_id = ?", user.UID).Scan(&classmember.ClassName)

	var subjects []Subject
	rows, err := db.Query("SELECT id, name, class_name, teacher_id FROM subjects WHERE class_name = ?", classmember.ClassName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error retrieving subjects"})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var subject Subject
		if err := rows.Scan(&subject.ID, &subject.Name, &subject.ClassName, &subject.TeacherID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Error scanning subject"})
			return
		}
		subjects = append(subjects, subject)
	}

	c.JSON(http.StatusOK, subjects)
}
func AddAttendance(c *gin.Context) {
	var attendance Attendance
	if err := c.ShouldBindJSON(&attendance); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input"})
		return
	}

	if attendance.UserID == 0 || attendance.SubjectID == 0 || attendance.Status == "" || attendance.Date == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "User ID, subject ID, subjectid, status, and date are required"})
		return
	}
	_, err := db.Exec("INSERT INTO attendance (user_id, subject_id, status, date) VALUES (?, ?, ?, ?)", attendance.UserID, attendance.SubjectID, attendance.Status, attendance.Date)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Attendance added successfully"})
}

func GetLuckyNumber(c *gin.Context){
	c.JSON(http.StatusOK, gin.H{"lucky_number": getRandomNumber()})
}
func GetExams(c *gin.Context){
	email, _ := c.Get("email")
	role, _ := c.Get("role")
	var user User
	var classmember ClassMember
	err := db.QueryRow("SELECT uid, email, class_name FROM users INNER JOIN class_members ON user.uid = class_members.user_id WHERE email = ?", email).Scan(&user.UID, &user.Email, &classmember.ClassName)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
		return
	}
	var rows *sql.Rows
	if role == "student" {
		rows, err := db.Query("SELECT id, class_name, teacher_id, subject_id, date, type, FROM exams WHERE class_name = ?", user.UID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Error retrieving exams"})
			return
		}
		defer rows.Close()
	} else if role == "teacher" {
		rows, err := db.Query("SELECT id, class_name, teacher_id, subject_id, date, type, FROM exams WHERE teacher_id = ?", user.UID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Error retrieving exams"})
			return
		}
		defer rows.Close()
	}
	var exams []Exam
	for rows.Next() {
		var exam Exam
		if err := rows.Scan(&exam.ID, &exam.ClassName, &exam.TeacherID, &exam.SubjectID, &exam.Date, &exam.Type); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Error scanning exam entry"})
			return
		}
		exams = append(exams, exam)
	}
	c.JSON(http.StatusOK, exams)
}
func GetAttendance(c *gin.Context){
	email, _ := c.Get("email")
	var user User
	err := db.QueryRow("SELECT uid, email FROM users WHERE email = ?", email).Scan(&user.UID, &user.Email)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
		return
	}
	var attendance []Attendance
	rows, err := db.Query("SELECT id, user_id, subject_id, status, date FROM attendance WHERE user_id = ?", user.UID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error retrieving attendance"})
		return
	}
	defer rows.Close()
	for rows.Next() {
		var att Attendance
		if err := rows.Scan(&att.ID, &att.UserID, &att.SubjectID, &att.Status, &att.Date); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Error scanning attendance"})
			return
		}
		attendance = append(attendance, att)
	}
	c.JSON(http.StatusOK, attendance)
}
func GetClassMembers(c *gin.Context){
	var class Class
	if err := c.ShouldBindJSON(&class); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input"})
		return
	}
	var classmembers []ClassMember
	rows, err := db.Query("SELECT id, user_id, class_name FROM class_members WHERE class_name = ?", class.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error retrieving class members"})
		return
	}
	defer rows.Close()
	for rows.Next() {
		var classmember ClassMember
		if err := rows.Scan(&classmember.ID, &classmember.UserID, &classmember.ClassName); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Error scanning class member"})
			return
		}
		classmembers = append(classmembers, classmember)
	}
	c.JSON(http.StatusOK, classmembers)
}
func GetStudentGrades(c *gin.Context){
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input"})
		return
	}
	var grades []Grade
	rows, err := db.Query("SELECT id, user_id, subject_id, grade, grade_type, date FROM grades WHERE user_id = ?", user.UID)
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
func GetStudentAttendance(c *gin.Context){
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input"})
		return
	}
	var attendance []Attendance
	rows, err := db.Query("SELECT id, user_id, subject_id, status, date FROM attendance WHERE user_id = ?", user.UID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error retrieving attendance"})
		return
	}
	defer rows.Close()
	for rows.Next() {
		var att Attendance
		if err := rows.Scan(&att.ID, &att.UserID, &att.SubjectID, &att.Status, &att.Date); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Error scanning attendance"})
			return
		}
		attendance = append(attendance, att)
	}
	c.JSON(http.StatusOK, attendance)
}
func GetStudentInfo(c *gin.Context){
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input"})
		return
	}
	var person Person
	err := db.QueryRow("SELECT first_name, last_name, birth_date, address, phone FROM persons WHERE user_id = ?", user.UID).Scan(&person.FirstName, &person.LastName, &person.BirthDate, &person.Address, &person.Phone)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "User details not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"first_name": person.FirstName,
		"last_name":  person.LastName,
		"birth_date": person.BirthDate,
		"address":    person.Address,
		"phone":      person.Phone,
	})
}
func AddExam(c *gin.Context){
	var exam Exam
	if err := c.ShouldBindJSON(&exam); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input"})
		return
	}
	if exam.ClassName == "" || exam.TeacherID == 0 || exam.SubjectID == 0 || exam.Date == "" || exam.Type == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Class name, teacher ID, subject ID, date, and type are required"})
		return
	}
	_, err := db.Exec("INSERT INTO exams (class_name, teacher_id, subject_id, date, type) VALUES (?, ?, ?, ?, ?)", exam.ClassName, exam.TeacherID, exam.SubjectID, exam.Date, exam.Type)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Exam created successfully"})
}
func AddClass(c *gin.Context){
	var class Class
	if err := c.ShouldBindJSON(&class); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input"})
		return
	}
	if class.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Class name is required"})
		return
	}
	_, err := db.Exec("INSERT INTO classes (name) VALUES (?)", class.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Class created successfully"})
}
func AddSubject(c *gin.Context){
	var subject Subject
	if err := c.ShouldBindJSON(&subject); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input"})
		return
	}
	if subject.Name == "" || subject.ClassName == "" || subject.TeacherID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Subject name, class name, and teacher ID are required"})
		return
	}
	_, err := db.Exec("INSERT INTO subjects (name, class_name, teacher_id) VALUES (?, ?, ?)", subject.Name, subject.ClassName, subject.TeacherID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Subject created successfully"})
}
func AddClassMember(c *gin.Context){
	var classmember ClassMember
	if err := c.ShouldBindJSON(&classmember); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input"})
		return
	}
	if classmember.UserID == 0 || classmember.ClassName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "User ID and class name are required"})
		return
	}
	_, err := db.Exec("INSERT INTO class_members (user_id, class_name) VALUES (?, ?)", classmember.UserID, classmember.ClassName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Class member added successfully"})
}