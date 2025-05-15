package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

const baseURL = "http://localhost:10800/api"

var token string

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("== Mercury Backend CLI ==")

	for {
		fmt.Print("\nChoose option [login, timetable, change-password, register-user, add-timetable, add-grade, quit]: ")
		choice, _ := reader.ReadString('\n')
		choice = strings.TrimSpace(choice)

		switch choice {
		case "login":
			login(reader)
		case "timetable":
			getTimetable()
		case "change-password":
			changePassword(reader)
		case "register-user":
			registerUser(reader)
		case "add-timetable":
			addTimetableEntry(reader)
		case "add-grade":
			addGrade(reader)
		case "quit":
			fmt.Println("Goodbye!")
			return
		default:
			fmt.Println("Unknown option.")
		}
	}
}

func login(reader *bufio.Reader) {
	fmt.Print("Email: ")
	email, _ := reader.ReadString('\n')
	fmt.Print("Password: ")
	password, _ := reader.ReadString('\n')

	data := map[string]string{
		"email":    strings.TrimSpace(email),
		"password": strings.TrimSpace(password),
	}
	body, _ := json.Marshal(data)

	resp, err := http.Post(baseURL+"/login", "application/json", bytes.NewBuffer(body))
	if err != nil {
		fmt.Println("Login error:", err)
		return
	}
	defer resp.Body.Close()

	var result map[string]string
	json.NewDecoder(resp.Body).Decode(&result)

	if resp.StatusCode == 200 {
		token = result["token"]
		fmt.Println("Logged in successfully.")
	} else {
		fmt.Println("Error:", result["message"])
	}
}

func getTimetable() {
	resp, err := http.Get(baseURL + "/timetable")
	if err != nil {
		fmt.Println("Failed to fetch timetable:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		fmt.Println("Server error.")
		return
	}

	var timetable []map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&timetable)

	fmt.Println("\n--- Timetable ---")
	for _, entry := range timetable {
		fmt.Printf("%s %s - %s | Room: %s | Subject ID: %v | Class: %v\n",
			entry["day"], entry["start_time"], entry["end_time"],
			entry["room"], entry["subject_id"], entry["class_name"])
	}
}

func changePassword(reader *bufio.Reader) {
	if token == "" {
		fmt.Println("Please login first.")
		return
	}

	fmt.Print("Old password: ")
	oldPass, _ := reader.ReadString('\n')
	fmt.Print("New password: ")
	newPass, _ := reader.ReadString('\n')

	data := map[string]string{
		"old_password": strings.TrimSpace(oldPass),
		"new_password": strings.TrimSpace(newPass),
	}
	body, _ := json.Marshal(data)

	req, _ := http.NewRequest("PUT", baseURL+"/change-password", bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Request error:", err)
		return
	}
	defer resp.Body.Close()

	var result map[string]string
	json.NewDecoder(resp.Body).Decode(&result)

	fmt.Println("Status:", resp.StatusCode)
	fmt.Println("Message:", result["message"])
}

func registerUser(reader *bufio.Reader) {
	if token == "" {
		fmt.Println("Please login as admin first.")
		return
	}

	fmt.Println("== Register New User ==")
	fmt.Print("New user email: ")
	email, _ := reader.ReadString('\n')
	fmt.Print("Password: ")
	password, _ := reader.ReadString('\n')
	fmt.Print("Role [student, teacher, admin]: ")
	role, _ := reader.ReadString('\n')
	fmt.Print("First name: ")
	firstName, _ := reader.ReadString('\n')
	fmt.Print("Last name: ")
	lastName, _ := reader.ReadString('\n')
	fmt.Print("Birth date (YYYY-MM-DD): ")
	birthDate, _ := reader.ReadString('\n')
	fmt.Print("Address: ")
	address, _ := reader.ReadString('\n')
	fmt.Print("Phone: ")
	phone, _ := reader.ReadString('\n')

	user := map[string]interface{}{
		"email":    strings.TrimSpace(email),
		"password": strings.TrimSpace(password),
		"role":     strings.TrimSpace(role),
		"argument": map[string]interface{}{
			"email":      strings.TrimSpace(email),
			"password":   strings.TrimSpace(password),
			"role":       strings.TrimSpace(role),
			"first_name": strings.TrimSpace(firstName),
			"last_name":  strings.TrimSpace(lastName),
			"birth_date": strings.TrimSpace(birthDate),
			"address":    strings.TrimSpace(address),
			"phone":      strings.TrimSpace(phone),
		},
	}

	body, _ := json.Marshal(user)
	req, _ := http.NewRequest("POST", baseURL+"/admin/register", bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Request error:", err)
		return
	}
	defer resp.Body.Close()

	var result map[string]string
	json.NewDecoder(resp.Body).Decode(&result)

	fmt.Println("Status:", resp.StatusCode)
	fmt.Println("Message:", result["message"])
}
func addTimetableEntry(reader *bufio.Reader) {
	if token == "" {
		fmt.Println("Please login as admin first.")
		return
	}

	fmt.Println("== Add Timetable Entry ==")
	fmt.Print("Day (e.g. Monday): ")
	day, _ := reader.ReadString('\n')
	fmt.Print("Subject ID (number): ")
	subjectIDStr, _ := reader.ReadString('\n')
	fmt.Print("Start Time (e.g. 08:00): ")
	start, _ := reader.ReadString('\n')
	fmt.Print("End Time (e.g. 09:30): ")
	end, _ := reader.ReadString('\n')
	fmt.Print("Room: ")
	room, _ := reader.ReadString('\n')
	fmt.Print("Teacher ID (number): ")
	teacherIDStr, _ := reader.ReadString('\n')
	fmt.Print("Class Name: ")
	className, _ := reader.ReadString('\n')

	data := map[string]interface{}{
		"day":         strings.TrimSpace(day),
		"subject_id":  toInt(subjectIDStr),
		"start_time":  strings.TrimSpace(start),
		"end_time":    strings.TrimSpace(end),
		"room":        strings.TrimSpace(room),
		"teacher_id":  toInt(teacherIDStr),
		"class_name":  strings.TrimSpace(className),
	}

	body, _ := json.Marshal(data)
	req, _ := http.NewRequest("POST", baseURL+"/timetable", bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Request error:", err)
		return
	}
	defer resp.Body.Close()

	var result map[string]string
	json.NewDecoder(resp.Body).Decode(&result)

	fmt.Println("Status:", resp.StatusCode)
	fmt.Println("Message:", result["message"])
}
func addGrade(reader *bufio.Reader) {
	if token == "" {
		fmt.Println("Please login as teacher first.")
		return
	}

	fmt.Println("== Add Grade for Student ==")
	fmt.Print("Student ID: ")
	studentID, _ := reader.ReadString('\n')
	fmt.Print("Subject ID: ")
	subjectID, _ := reader.ReadString('\n')
	fmt.Print("Grade (e.g. 5, A, etc): ")
	grade, _ := reader.ReadString('\n')
	fmt.Print("Grade Type ('numeric', 'comment', 'custom'): ")
	gradeType, _ := reader.ReadString('\n')
	fmt.Print("Date (YYYY-MM-DD): ")
	date, _ := reader.ReadString('\n')

	data := map[string]interface{}{
		"user_id":    toInt(studentID),
		"subject_id": toInt(subjectID),
		"grade":      strings.TrimSpace(grade),
		"grade_type": strings.TrimSpace(gradeType),
		"date":       strings.TrimSpace(date),
	}

	url := fmt.Sprintf("%s/grades/%d", baseURL, toInt(studentID))
	body, _ := json.Marshal(data)

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Request error:", err)
		return
	}
	defer resp.Body.Close()

	var result map[string]string
	json.NewDecoder(resp.Body).Decode(&result)

	fmt.Println("Status:", resp.StatusCode)
	fmt.Println("Message:", result["message"])
}
func toInt(input string) int {
	value := strings.TrimSpace(input)
	var number int
	fmt.Sscanf(value, "%d", &number)
	return number
}
