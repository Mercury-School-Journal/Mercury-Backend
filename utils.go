package main
import (
	"golang.org/x/crypto/bcrypt"
    "math/rand"
    "sync"
    "time"
)
var (
	randomNum        int
	lastGeneratedDate time.Time
	mu               sync.Mutex
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

func generateRandomNumber() int {
    return rand.Intn(35) + 1
}

func getRandomNumber() int {
    mu.Lock()
    defer mu.Unlock()

    today := time.Now().Truncate(24 * time.Hour)

    if lastGeneratedDate.Equal(today) {
        return randomNum
    }

    randomNum = generateRandomNumber()
    lastGeneratedDate = today
    return randomNum
}