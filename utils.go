package main
import (
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
func TableExists(tableName string) bool {
    query := `SELECT name FROM sqlite_master WHERE type='table' AND name=?;`
    var name string
    err := db.QueryRow(query, tableName).Scan(&name)
    return err == nil
}