package main

import (
	"time"
	"log"
	"database/sql"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_"modernc.org/sqlite"
)

var jwtKey []byte
var db *sql.DB

func init() {
	var err error
	key, exists := os.LookupEnv("jwtKey")
	if !exists {
		log.Fatal("jwtKey environment variable is not set")
	}
	jwtKey = []byte(key)

	db, err = sql.Open("sqlite", "./database.db")
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
}

func main() {
	r := gin.Default()
    r.Use(cors.New(cors.Config{
        AllowAllOrigins: true,
        AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
        ExposeHeaders:    []string{"Content-Length"},
        AllowCredentials: true,
        MaxAge: 12 * time.Hour,
    }))
	r.Use(TokenAuthMiddleware())
	{
		r.PUT("/api/change-password", ChangePassword)
		r.DELETE("/api/delete-account", DeleteAccount)
	}
	r.Use(AdminAuthMiddleware())
	{
		r.POST("/api/register", RegisterUser)
		r.POST("/api/timetable", AddTimetableEntry)
	}
	r.Use(TeacherAuthMiddleware())
	{
		r.POST("/api/grades/:user_id", AddGrade)
		r.GET("/api/grades/:user_id", GetGrades)
	}	
	r.POST("/api/login", Login)
	r.GET("/api/ping",Ping)
	
	r.POST("/api/grades", AddGrade)
	r.GET("/api/grades/:user_id", GetGrades)
	
	r.GET("/api/timetable", GetTimetable)
	

	// err := r.RunTLS(":10800", "cert.pem", "key.pem") //  for HTTPS
	err := r.Run(":10800")
	if err != nil {
		log.Fatal(err)
		return
	}
}