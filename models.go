package main

import (
	"github.com/golang-jwt/jwt/v4"
)

// User represents a user in the system (students, teachers, admins)
type User struct {
	UID      uint   `json:"uid"`
	Email    string `json:"email"`
	Password string `json:"password"` // User password, excluded from JSON
	Role     string `json:"role"`     // User role: "student", "teacher", or "admin"
}

// Person represents personal information for a user
type Person struct {
	ID        uint   `json:"id"`
	UserID    uint   `json:"user_id"`              // Reference to users(uid)
	FirstName string `json:"first_name"`           // First name
	LastName  string `json:"last_name"`            // Last name
	BirthDate string `json:"birth_date,omitempty"` // Birth date in YYYY-MM-DD format
	Address   string `json:"address,omitempty"`    // Address
	Phone     string `json:"phone,omitempty"`      // Phone number
}

// Class represents a school class (group of students)
type Class struct {
	ID   uint   `json:"id"`
	Name string `json:"name"` // Unique class name (e.g., "1A", "2B")
}

// Subject represents a school subject
type Subject struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`       // Unique subject name (e.g., "Mathematics")
	ClassName string `json:"class_name"` // Reference to classes(name)
	TeacherID uint   `json:"teacher_id"` // Reference to users(uid)
}

// StudentSubject represents a student-subject assignment
type StudentSubject struct {
	ID        uint `json:"id"`
	UserID    uint `json:"user_id"`    // Reference to users(uid)
	SubjectID uint `json:"subject_id"` // Reference to subjects(id)
}

// TeacherSubject represents a teacher-subject assignment
type TeacherSubject struct {
	ID        uint `json:"id"`
	UserID    uint `json:"user_id"`    // Reference to users(uid)
	SubjectID uint `json:"subject_id"` // Reference to subjects(id)
}

// Grade represents a grade, comment, or custom value for a student
type Grade struct {
	ID        uint   `json:"id"`
	UserID    uint   `json:"user_id"`    // Reference to users(uid)
	SubjectID uint   `json:"subject_id"` // Reference to subjects(id)
	Grade     string `json:"grade"`      // Numeric grade, comment, or custom value
	GradeType string `json:"grade_type"` // Type: "numeric", "comment", or "custom"
	Date      string `json:"date"`       // Date of entry in YYYY-MM-DD format
}

// ClassMember represents a user (student or teacher) assigned to a class
type ClassMember struct {
	ID        uint   `json:"id"`
	UserID    uint   `json:"user_id"`    // Reference to users(uid)
	ClassName string `json:"class_name"` // Reference to classes(name)
}

// TimetableEntry represents a single timetable entry
type TimetableEntry struct {
	ID        uint   `json:"id"`
	Day       string `json:"day"`        // Day of the week (e.g., "Monday")
	SubjectID uint   `json:"subject_id"` // Reference to subjects(id)
	StartTime string `json:"start_time"` // Start time in HH:MM format
	EndTime   string `json:"end_time"`   // End time in HH:MM format
	Room      string `json:"room"`       // Room number or name
	TeacherID uint   `json:"teacher_id"` // Reference to users(uid)
	ClassName string `json:"class_name"` // Reference to classes(name)
}

// AccessRequest represents a login request
type AccessRequest struct {
	Email    string      `json:"email"`    // User email
	Password string      `json:"password"` // User password
	Argument interface{} `json:"argument"` // Additional data for the request
}

// Claims represents JWT claims for authentication
type Claims struct {
	Email string `json:"email"` // User email
	Role  string `json:"role"`  // User role
	jwt.StandardClaims
}

// Input represents a password change request
type Input struct {
	OldPassword string `json:"old_password"` // Current password
	NewPassword string `json:"new_password"` // New password
}

type Attendance struct {
	ID        uint   `json:"id"`
	UserID    uint   `json:"user_id"`    // Reference to users(uid)
	SubjectID uint   `json:"subject_id"` // Reference to subjects(id)
	Date      string `json:"date"`       // Date of attendance in YYYY-MM-DD format
	Status    string `json:"status"`     // Attendance status: "present", "absent", or "late"
}

// Exam represents a exam or test
type Exam struct {
	ID        uint   `json:"id"`
	ClassName string `json:"class_name"` // Reference to classes(name)
	SubjectID uint   `json:"subject_id"` // Reference to subjects(id)
	TeacherID uint   `json:"teacher_id"` // Reference to users(uid)
	Date      string `json:"date"`       // Date of the exam in YYYY-MM-DD format
	Type      string `json:"type"`      // Type of exam (e.g., "exam", "test", "quiz")
	Description string `json:"description"` // Description of the exam
}