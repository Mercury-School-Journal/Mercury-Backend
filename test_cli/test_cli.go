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
        fmt.Print("\nChoose option [login, timetable, change-password, register-user, add-timetable, add-grade, delete-account, ping, get-grades, get-user-info, get-subjects, add-attendance, get-lucky-number, get-exams, get-attendance, get-class-members, get-student-grades, get-student-attendance, get-student-info, add-exam, add-class, add-subject, add-class-member, quit]: ")
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
        case "delete-account":
            deleteAccount()
        case "ping":
            ping()
        case "get-grades":
            getGrades()
        case "get-user-info":
            getUserInfo()
        case "get-subjects":
            getSubjects()
        case "add-attendance":
            addAttendance(reader)
        case "get-lucky-number":
            getLuckyNumber()
        case "get-exams":
            getExams()
        case "get-attendance":
            getAttendance()
        case "get-class-members":
            getClassMembers(reader)
        case "get-student-grades":
            getStudentGrades(reader)
        case "get-student-attendance":
            getStudentAttendance(reader)
        case "get-student-info":
            getStudentInfo(reader)
        case "add-exam":
            addExam(reader)
        case "add-class":
            addClass(reader)
        case "add-subject":
            addSubject(reader)
        case "add-class-member":
            addClassMember(reader)
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
    if token == "" {
        fmt.Println("Please login first.")
        return
    }

    req, _ := http.NewRequest("GET", baseURL+"/timetable", nil)
    req.Header.Set("Authorization", "Bearer "+token)

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        fmt.Println("Failed to fetch timetable:", err)
        return
    }
    defer resp.Body.Close()

    if resp.StatusCode != 200 {
        var result map[string]string
        json.NewDecoder(resp.Body).Decode(&result)
        fmt.Println("Error:", result["message"])
        return
    }

    var timetable []map[string]interface{}
    json.NewDecoder(resp.Body).Decode(&timetable)

    fmt.Println("\n--- Timetable ---")
    for _, entry := range timetable {
        fmt.Printf("ID: %v | Day: %s | Start: %s | End: %s | Room: %s | Subject ID: %v | Teacher ID: %v | Class: %s\n",
            entry["id"], entry["day"], entry["time_start"], entry["time_end"],
            entry["room"], entry["subject_id"], entry["teacher_id"], entry["class_name"])
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
        "time_start":  strings.TrimSpace(start),
        "time_end":    strings.TrimSpace(end),
        "room":        strings.TrimSpace(room),
        "teacher_id":  toInt(teacherIDStr),
        "class_name":  strings.TrimSpace(className),
    }

    body, _ := json.Marshal(data)
    req, _ := http.NewRequest("POST", baseURL+"/admin/timetable", bytes.NewBuffer(body))
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
        fmt.Println("Please login as admin or teacher first.")
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

    body, _ := json.Marshal(data)
    url := baseURL + "/teacher/grade"
    if isAdmin() {
        url = baseURL + "/admin/grade"
    }

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

func deleteAccount() {
    if token == "" {
        fmt.Println("Please login first.")
        return
    }

    req, _ := http.NewRequest("DELETE", baseURL+"/delete-account", nil)
    req.Header.Set("Authorization", "Bearer "+token)

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

func ping() {
    resp, err := http.Get(baseURL + "/ping")
    if err != nil {
        fmt.Println("Ping error:", err)
        return
    }
    defer resp.Body.Close()

    if resp.StatusCode == http.StatusNoContent {
        fmt.Println("Ping successful.")
    } else {
        fmt.Println("Ping failed, status:", resp.StatusCode)
    }
}

func getGrades() {
    if token == "" {
        fmt.Println("Please login as student first.")
        return
    }

    req, _ := http.NewRequest("GET", baseURL+"/student/grades", nil)
    req.Header.Set("Authorization", "Bearer "+token)

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        fmt.Println("Request error:", err)
        return
    }
    defer resp.Body.Close()

    if resp.StatusCode != 200 {
        var result map[string]string
        json.NewDecoder(resp.Body).Decode(&result)
        fmt.Println("Error:", result["message"])
        return
    }

    var grades []map[string]interface{}
    json.NewDecoder(resp.Body).Decode(&grades)

    fmt.Println("\n--- Grades ---")
    for _, grade := range grades {
        fmt.Printf("ID: %v | User ID: %v | Subject ID: %v | Grade: %s | Type: %s | Date: %s\n",
            grade["id"], grade["user_id"], grade["subject_id"], grade["grade"], grade["grade_type"], grade["date"])
    }
}

func getUserInfo() {
    if token == "" {
        fmt.Println("Please login first.")
        return
    }

    req, _ := http.NewRequest("GET", baseURL+"/user", nil)
    req.Header.Set("Authorization", "Bearer "+token)

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        fmt.Println("Request error:", err)
        return
    }
    defer resp.Body.Close()

    if resp.StatusCode != 200 {
        var result map[string]string
        json.NewDecoder(resp.Body).Decode(&result)
        fmt.Println("Error:", result["message"])
        return
    }

    var info map[string]interface{}
    json.NewDecoder(resp.Body).Decode(&info)

    fmt.Println("\n--- User Info ---")
    fmt.Printf("Email: %s\nFirst Name: %s\nLast Name: %s\nBirth Date: %s\nAddress: %s\nPhone: %s\n",
        info["email"], info["first_name"], info["last_name"], info["birth_date"], info["address"], info["phone"])
}

func getSubjects() {
    if token == "" {
        fmt.Println("Please login as student first.")
        return
    }

    req, _ := http.NewRequest("GET", baseURL+"/student/subjects", nil)
    req.Header.Set("Authorization", "Bearer "+token)

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        fmt.Println("Request error:", err)
        return
    }
    defer resp.Body.Close()

    if resp.StatusCode != 200 {
        var result map[string]string
        json.NewDecoder(resp.Body).Decode(&result)
        fmt.Println("Error:", result["message"])
        return
    }

    var subjects []map[string]interface{}
    json.NewDecoder(resp.Body).Decode(&subjects)

    fmt.Println("\n--- Subjects ---")
    for _, subject := range subjects {
        fmt.Printf("ID: %v | Name: %s | Class Name: %s | Teacher ID: %v\n",
            subject["id"], subject["name"], subject["class_name"], subject["teacher_id"])
    }
}

func addAttendance(reader *bufio.Reader) {
    if token == "" {
        fmt.Println("Please login as admin or teacher first.")
        return
    }

    fmt.Println("== Add Attendance ==")
    fmt.Print("Student ID: ")
    userID, _ := reader.ReadString('\n')
    fmt.Print("Subject ID: ")
    subjectID, _ := reader.ReadString('\n')
    fmt.Print("Status (present/absent): ")
    status, _ := reader.ReadString('\n')
    fmt.Print("Date (YYYY-MM-DD): ")
    date, _ := reader.ReadString('\n')

    data := map[string]interface{}{
        "user_id":    toInt(userID),
        "subject_id": toInt(subjectID),
        "status":     strings.TrimSpace(status),
        "date":       strings.TrimSpace(date),
    }

    body, _ := json.Marshal(data)
    url := baseURL + "/teacher/attendance"
    if isAdmin() {
        url = baseURL + "/admin/attendance"
    }

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

func getLuckyNumber() {
    resp, err := http.Get(baseURL + "/lucky-number")
    if err != nil {
        fmt.Println("Request error:", err)
        return
    }
    defer resp.Body.Close()

    if resp.StatusCode != 200 {
        var result map[string]string
        json.NewDecoder(resp.Body).Decode(&result)
        fmt.Println("Error:", result["message"])
        return
    }

    var result map[string]interface{}
    json.NewDecoder(resp.Body).Decode(&result)

    fmt.Println("Lucky Number:", result["lucky_number"])
}

func getExams() {
    if token == "" {
        fmt.Println("Please login first.")
        return
    }

    req, _ := http.NewRequest("GET", baseURL+"/exams", nil)
    req.Header.Set("Authorization", "Bearer "+token)

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        fmt.Println("Request error:", err)
        return
    }
    defer resp.Body.Close()

    if resp.StatusCode != 200 {
        var result map[string]string
        json.NewDecoder(resp.Body).Decode(&result)
        fmt.Println("Error:", result["message"])
        return
    }

    var exams []map[string]interface{}
    json.NewDecoder(resp.Body).Decode(&exams)

    fmt.Println("\n--- Exams ---")
    for _, exam := range exams {
        fmt.Printf("ID: %v | Class Name: %s | Teacher ID: %v | Subject ID: %v | Date: %s | Type: %s\n",
            exam["id"], exam["class_name"], exam["teacher_id"], exam["subject_id"], exam["date"], exam["type"])
    }
}

func getAttendance() {
    if token == "" {
        fmt.Println("Please login as student first.")
        return
    }

    req, _ := http.NewRequest("GET", baseURL+"/student/attendance", nil)
    req.Header.Set("Authorization", "Bearer "+token)

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        fmt.Println("Request error:", err)
        return
    }
    defer resp.Body.Close()

    if resp.StatusCode != 200 {
        var result map[string]string
        json.NewDecoder(resp.Body).Decode(&result)
        fmt.Println("Error:", result["message"])
        return
    }

    var attendance []map[string]interface{}
    json.NewDecoder(resp.Body).Decode(&attendance)

    fmt.Println("\n--- Attendance ---")
    for _, att := range attendance {
        fmt.Printf("ID: %v | User ID: %v | Subject ID: %v | Status: %s | Date: %s\n",
            att["id"], att["user_id"], att["subject_id"], att["status"], att["date"])
    }
}

func getClassMembers(reader *bufio.Reader) {
    if token == "" {
        fmt.Println("Please login as admin or teacher first.")
        return
    }

    fmt.Print("Class Name: ")
    className, _ := reader.ReadString('\n')

    data := map[string]string{
        "name": strings.TrimSpace(className),
    }
    body, _ := json.Marshal(data)

    url := baseURL + "/teacher/class"
    if isAdmin() {
        url = baseURL + "/admin/class"
    }

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

    if resp.StatusCode != 200 {
        var result map[string]string
        json.NewDecoder(resp.Body).Decode(&result)
        fmt.Println("Error:", result["message"])
        return
    }

    var members []map[string]interface{}
    json.NewDecoder(resp.Body).Decode(&members)

    fmt.Println("\n--- Class Members ---")
    for _, member := range members {
        fmt.Printf("ID: %v | User ID: %v | Class Name: %s\n",
            member["id"], member["user_id"], member["class_name"])
    }
}

func getStudentGrades(reader *bufio.Reader) {
    if token == "" {
        fmt.Println("Please login as admin or teacher first.")
        return
    }

    fmt.Print("Student User ID: ")
    userID, _ := reader.ReadString('\n')

    data := map[string]interface{}{
        "uid": toInt(userID),
    }
    body, _ := json.Marshal(data)

    url := baseURL + "/teacher/student-grades"
    if isAdmin() {
        url = baseURL + "/admin/student-grades"
    }

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

    if resp.StatusCode != 200 {
        var result map[string]string
        json.NewDecoder(resp.Body).Decode(&result)
        fmt.Println("Error:", result["message"])
        return
    }

    var grades []map[string]interface{}
    json.NewDecoder(resp.Body).Decode(&grades)

    fmt.Println("\n--- Student Grades ---")
    for _, grade := range grades {
        fmt.Printf("ID: %v | User ID: %v | Subject ID: %v | Grade: %s | Type: %s | Date: %s\n",
            grade["id"], grade["user_id"], grade["subject_id"], grade["grade"], grade["grade_type"], grade["date"])
    }
}

func getStudentAttendance(reader *bufio.Reader) {
    if token == "" {
        fmt.Println("Please login as admin or teacher first.")
        return
    }

    fmt.Print("Student User ID: ")
    userID, _ := reader.ReadString('\n')

    data := map[string]interface{}{
        "uid": toInt(userID),
    }
    body, _ := json.Marshal(data)

    url := baseURL + "/teacher/student-attendance"
    if isAdmin() {
        url = baseURL + "/admin/student-attendance"
    }

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

    if resp.StatusCode != 200 {
        var result map[string]string
        json.NewDecoder(resp.Body).Decode(&result)
        fmt.Println("Error:", result["message"])
        return
    }

    var attendance []map[string]interface{}
    json.NewDecoder(resp.Body).Decode(&attendance)

    fmt.Println("\n--- Student Attendance ---")
    for _, att := range attendance {
        fmt.Printf("ID: %v | User ID: %v | Subject ID: %v | Status: %s | Date: %s\n",
            att["id"], att["user_id"], att["subject_id"], att["status"], att["date"])
    }
}

func getStudentInfo(reader *bufio.Reader) {
    if token == "" {
        fmt.Println("Please login as admin or teacher first.")
        return
    }

    fmt.Print("Student User ID: ")
    userID, _ := reader.ReadString('\n')

    data := map[string]interface{}{
        "uid": toInt(userID),
    }
    body, _ := json.Marshal(data)

    url := baseURL + "/teacher/student-info"
    if isAdmin() {
        url = baseURL + "/admin/student-info"
    }

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

    if resp.StatusCode != 200 {
        var result map[string]string
        json.NewDecoder(resp.Body).Decode(&result)
        fmt.Println("Error:", result["message"])
        return
    }

    var info map[string]interface{}
    json.NewDecoder(resp.Body).Decode(&info)

    fmt.Println("\n--- Student Info ---")
    fmt.Printf("First Name: %s\nLast Name: %s\nBirth Date: %s\nAddress: %s\nPhone: %s\n",
        info["first_name"], info["last_name"], info["birth_date"], info["address"], info["phone"])
}

func addExam(reader *bufio.Reader) {
    if token == "" {
        fmt.Println("Please login as admin or teacher first.")
        return
    }

    fmt.Println("== Add Exam ==")
    fmt.Print("Class Name: ")
    className, _ := reader.ReadString('\n')
    fmt.Print("Teacher ID: ")
    teacherID, _ := reader.ReadString('\n')
    fmt.Print("Subject ID: ")
    subjectID, _ := reader.ReadString('\n')
    fmt.Print("Date (YYYY-MM-DD): ")
    date, _ := reader.ReadString('\n')
    fmt.Print("Type (e.g. test, quiz): ")
    examType, _ := reader.ReadString('\n')

    data := map[string]interface{}{
        "class_name":  strings.TrimSpace(className),
        "teacher_id":  toInt(teacherID),
        "subject_id":  toInt(subjectID),
        "date":        strings.TrimSpace(date),
        "type":        strings.TrimSpace(examType),
    }

    body, _ := json.Marshal(data)
    url := baseURL + "/teacher/exam"
    if isAdmin() {
        url = baseURL + "/admin/exam"
    }

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

func addClass(reader *bufio.Reader) {
    if token == "" {
        fmt.Println("Please login as admin first.")
        return
    }

    fmt.Println("== Add Class ==")
    fmt.Print("Class Name: ")
    className, _ := reader.ReadString('\n')

    data := map[string]string{
        "name": strings.TrimSpace(className),
    }

    body, _ := json.Marshal(data)
    req, _ := http.NewRequest("POST", baseURL+"/admin/class", bytes.NewBuffer(body))
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

func addSubject(reader *bufio.Reader) {
    if token == "" {
        fmt.Println("Please login as admin first.")
        return
    }

    fmt.Println("== Add Subject ==")
    fmt.Print("Subject Name: ")
    name, _ := reader.ReadString('\n')
    fmt.Print("Class Name: ")
    className, _ := reader.ReadString('\n')
    fmt.Print("Teacher ID: ")
    teacherID, _ := reader.ReadString('\n')

    data := map[string]interface{}{
        "name":       strings.TrimSpace(name),
        "class_name": strings.TrimSpace(className),
        "teacher_id": toInt(teacherID),
    }

    body, _ := json.Marshal(data)
    req, _ := http.NewRequest("POST", baseURL+"/admin/subject", bytes.NewBuffer(body))
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

func addClassMember(reader *bufio.Reader) {
    if token == "" {
        fmt.Println("Please login as admin first.")
        return
    }

    fmt.Println("== Add Class Member ==")
    fmt.Print("User ID: ")
    userID, _ := reader.ReadString('\n')
    fmt.Print("Class Name: ")
    className, _ := reader.ReadString('\n')

    data := map[string]interface{}{
        "user_id":    toInt(userID),
        "class_name": strings.TrimSpace(className),
    }

    body, _ := json.Marshal(data)
    req, _ := http.NewRequest("POST", baseURL+"/admin/class-member", bytes.NewBuffer(body))
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

//# TODO: Implement the isAdmin function to check if the user is an admin
func isAdmin() bool {
    return true
}

func toInt(input string) int {
    value := strings.TrimSpace(input)
    var number int
    fmt.Sscanf(value, "%d", &number)
    return number
}