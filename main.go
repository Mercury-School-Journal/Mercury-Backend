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

	db, err = sql.Open("sqlite", "./users.db")
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		email TEXT UNIQUE,
		password TEXT
	);`)
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
	r.POST("/api/register", RegisterUser)
	r.POST("/api/login", Login)
	r.GET("/api/ping",Ping)

	// err := r.RunTLS(":10800", "cert.pem", "key.pem") //  for HTTPS
	err := r.Run(":10800")
	if err != nil {
		log.Fatal(err)
		return
	}
}