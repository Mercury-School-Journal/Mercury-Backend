package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "modernc.org/sqlite"
)

var jwtKey []byte
var db *sql.DB

func init() {
	var err error

	jwtKeyStr, exists := os.LookupEnv("JWT_KEY")
	if !exists {
		log.Fatal("JWT_KEY environment variable is not set")
	}
	jwtKey = []byte(jwtKeyStr)

	adminEmail, exists := os.LookupEnv("ADMIN_EMAIL")
	if !exists {
		log.Fatal("ADMIN_EMAIL environment variable is not set")
	}
	adminPassword, exists := os.LookupEnv("ADMIN_PASSWORD")
	if !exists {
		log.Fatal("ADMIN_PASSWORD environment variable is not set")
	}
	dbPath, exists := os.LookupEnv("DB_PATH")
	if !exists {
		dbPath = "./database.db"
	}

	db, err = sql.Open("sqlite", dbPath)
	if err != nil {
		log.Fatal(err)
	}

	if !TableExists("users") {
		schema, err := os.ReadFile("schema.sql")
		if err != nil {
			log.Fatal(err)
		}

		_, err = db.Exec("PRAGMA foreign_keys = ON")
		if err != nil {
			log.Fatal(err)
		}

		_, err = db.Exec(string(schema))
		if err != nil {
			log.Fatal(err)
		}
		hashedPassword, err := HashPassword(adminPassword)
		if err != nil {
			log.Fatal(err)
		}
		tx, err := db.Begin()
		if err != nil {
			log.Fatal(err)
		}
		result, err := tx.Exec("INSERT INTO users (email, password, role) VALUES (?, ?, 'admin')", adminEmail, hashedPassword)
		if err != nil {
			log.Fatal(err)
		}
		adminUserID, err := result.LastInsertId()
		if err != nil {
			log.Fatal(err)
		}
		_, err = tx.Exec("INSERT INTO persons (user_id, first_name, last_name) VALUES (?, ?, ?)", adminUserID, "Admin", "User")
		if err != nil {
			log.Fatal(err)
		}
		if err := tx.Commit(); err != nil {
			log.Fatal(err)
		}
	}
}

func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()
		log.Printf("%s %s %d %s", method, path, status, latency)
	}
}

// main sets up and runs the web server
func main() {
	r := gin.Default()

	r.Use(LoggerMiddleware())

	r.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Public routes
	r.POST("/api/login", Login)
	r.GET("/api/ping", Ping)
	r.GET("/api/lucky-number", GetLuckyNumber)

	// Authenticated routes
	auth := r.Group("/api").Use(TokenAuthMiddleware())
	{
		auth.PUT("/change-password", ChangePassword)
		auth.DELETE("/delete-account", DeleteAccount)
		auth.GET("/timetable", GetTimetable)
		auth.GET("/user", GetUserInfo)
		auth.GET("/exams", GetExams)
	}

	// Admin routes
	admin := r.Group("/api/admin").Use(TokenAuthMiddleware(), AdminAuthMiddleware())
	{
		admin.POST("/register", RegisterUser)
		admin.POST("/timetable", AddTimetableEntry)
		admin.POST("/class", AddClass)
		admin.POST("/subject", AddSubject)
		admin.POST("/class-member", AddClassMember)

		admin.POST("/grade", AddGrade)
		admin.POST("/attendance", AddAttendance)
		admin.POST("/exam", AddExam)
		admin.GET("/class", GetClassMembers)
		admin.GET("/student-grades", GetStudentGrades)
		admin.GET("/student-attendance", GetStudentAttendance)
		admin.GET("/student-info", GetStudentInfo)
	}

	// Teacher routes
	teacher := r.Group("/api/teacher").Use(TokenAuthMiddleware(), TeacherAuthMiddleware())
	{
		teacher.POST("/grade", AddGrade)
		teacher.POST("/attendance", AddAttendance)
		teacher.POST("/exam", AddExam)
		teacher.GET("/class", GetClassMembers)
		teacher.GET("/student-grades", GetStudentGrades)
		teacher.GET("/student-attendance", GetStudentAttendance)
		teacher.GET("student-info", GetStudentInfo)
	}
	// Student routes
	student := r.Group("/api/student").Use(TokenAuthMiddleware(), StudentAuthMiddleware())
	{
		student.GET("/grades", GetGrades)
		student.GET("/subjects", GetSubjects)
		student.GET("/attendance", GetAttendance)
	}

	port, exists := os.LookupEnv("PORT")
	if !exists {
		port = ":10800"
	}
	certPath, exists := os.LookupEnv("CERT_PATH")
	if !exists {
		certPath = "cert.pem"
	}
	keyPath, exists := os.LookupEnv("KEY_PATH")
	if !exists {
		keyPath = "key.pem"
	}

	server := &http.Server{
		Addr:    port,
		Handler: r,
	}
	go func() {
		var err error
		if _, certErr := os.Stat(certPath); certErr == nil {
			if _, keyErr := os.Stat(keyPath); keyErr == nil {
				log.Printf("Starting HTTPS server on %s", port)
				err = server.ListenAndServeTLS(certPath, keyPath)
			} else {
				log.Printf("Key file not found, starting HTTP server on %s", port)
				err = server.ListenAndServe()
			}
		} else {
			log.Printf("Cert file not found, starting HTTP server on %s", port)
			err = server.ListenAndServe()
		}
		if err != nil && !strings.Contains(err.Error(), "Server closed") {
			log.Fatal(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}

	if err := db.Close(); err != nil {
		log.Fatal("Error closing database: ", err)
	}
	log.Println("Server exited")
}
