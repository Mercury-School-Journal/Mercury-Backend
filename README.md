# Mercury Backend Application Documentation

## 1. Purpose of the Application
Mercury Backend is a server-side application written in Go that provides a REST API for managing a school system. It enables user registration, login, and management of schedules, grades, attendance, exams, classes, subjects, and personal data of students and teachers. The application uses a SQLite database and JWT-based authentication with user roles (student, teacher, admin).

## 2. Project Structure
The application consists of a single main file, `main.go`, which includes:
- Database and server initialization.
- Definition of API routes and middleware.
- Helper functions, data models, and handlers.

**Files and Their Roles:**
- `main.go`: The main file configuring the Gin server, middleware, API routes, and database initialization.
- `schema.sql`: SQL file defining the database schema (tables: `users`, `persons`, `classes`, `subjects`, `grades`, `timetable`, `attendance`, `exams`, `class_members`).
- **Go Models**: Data structures (e.g., `User`, `Person`, `Grade`, `TimetableEntry`, `Attendance`, `Exam`, `Class`, `Subject`, `ClassMember`) mapping SQL tables.
- **Handlers**: Functions handling HTTP requests (e.g., `RegisterUser`, `Login`, `AddGrade`, `AddAttendance`, `AddExam`).
- **Middleware**: Authentication and authorization functions (`TokenAuthMiddleware`, `AdminAuthMiddleware`, `TeacherAuthMiddleware`, `StudentAuthMiddleware`).

## 3. Database Schema
The SQLite database includes the following tables:
- `users`: Stores user data (`uid`, `email`, `password`, `role`).
- `persons`: User personal data (`user_id`, `first_name`, `last_name`, `birth_date`, `address`, `phone`).
- `classes`: School classes (`id`, `name`).
- `subjects`: School subjects (`id`, `name`, `class_name`, `teacher_id`).
- `grades`: Grades, remarks, and custom values (`id`, `user_id`, `subject_id`, `grade`, `grade_type`, `date`).
- `class_members`: User-class associations (`id`, `user_id`, `class_name`).
- `timetable`: Class schedules (`id`, `day`, `subject_id`, `time_start`, `time_end`, `room`, `teacher_id`, `class_name`).
- `attendance`: Attendance records (`id`, `user_id`, `subject_id`, `status`, `date`).
- `exams`: Exams (`id`, `class_name`, `teacher_id`, `subject_id`, `date`, `type`).

The detailed schema is available in the `schema.sql` file.

## 4. Data Models
Go models map SQL tables and are used in handlers and HTTP requests:
- `User`: { `UID`, `Email`, `Password`, `Role` } – user data.
- `Person`: { `ID`, `UserID`, `FirstName`, `LastName`, `BirthDate`, `Address`, `Phone` } – personal data.
- `Class`: { `ID`, `Name` } – school class.
- `Subject`: { `ID`, `Name`, `ClassName`, `TeacherID` } – subject.
- `Grade`: { `ID`, `UserID`, `SubjectID`, `Grade`, `GradeType`, `Date` } – grade/remark.
- `ClassMember`: { `ID`, `UserID`, `ClassName` } – class association.
- `TimetableEntry`: { `ID`, `Day`, `SubjectID`, `StartTime`, `EndTime`, `Room`, `TeacherID`, `ClassName` } – schedule entry.
- `Attendance`: { `ID`, `UserID`, `SubjectID`, `Status`, `Date` } – attendance.
- `Exam`: { `ID`, `ClassName`, `TeacherID`, `SubjectID`, `Date`, `Type` } – exam.
- `AccessRequest`: { `Email`, `Password`, `Argument` } – login/registration data.
- `Claims`: { `Email`, `Role`, `StandardClaims` } – JWT data.
- `Input`: { `OldPassword`, `NewPassword` } – password change.

## 5. API Endpoints
The API is available at `http://localhost:10800/api` (or HTTPS if certificates are configured). Below is a description of the endpoints:

### Public Endpoints
#### POST /api/login
- **Description**: Logs in a user and returns a JWT token.
- **Body**: `{ "email": string, "password": string }`
- **Response**:
  - `200`: `{ "token": string }`
  - `400`: `{ "message": "Invalid input" }`
  - `401`: `{ "message": "Invalid email credentials" }` or `{ "message": "Invalid password credentials" }`
  - `500`: `{ "message": "Could not generate token" }`
- **Example**:
  ```json
  POST /api/login
  { "email": "admin@example.com", "password": "secret123" }
  ```

#### GET /api/ping
- **Description**: Returns status 204, confirming server operation.
- **Response**: `204` (No Content)

#### GET /api/lucky-number
- **Description**: Returns a random number.
- **Response**:
  - `200`: `{ "lucky_number": number }`

### Protected Endpoints (Require JWT)
#### PUT /api/change-password (TokenAuthMiddleware)
- **Description**: Changes the user's password.
- **Header**: `Authorization: Bearer <token>`
- **Body**: `{ "old_password": string, "new_password": string }`
- **Response**:
  - `200`: `{ "message": "Password changed successfully" }`
  - `400`: `{ "message": "Invalid input" }`
  - `401`: `{ "message": "Incorrect old password" }`
  - `404`: `{ "message": "User not found" }`
  - `500`: `{ "message": "Error hashing new password" }` or `{ "message": "Error updating password" }`

#### DELETE /api/delete-account (TokenAuthMiddleware)
- **Description**: Deletes the user's account and associated data.
- **Header**: `Authorization: Bearer <token>`
- **Response**:
  - `200`: `{ "message": "User deleted successfully" }`
  - `500`: `{ "message": "Error deleting user details" }` or `{ "message": "Error committing transaction" }`

#### GET /api/timetable (TokenAuthMiddleware)
- **Description**: Retrieves the schedule for the logged-in user (for students: their class; for teachers: their lessons).
- **Header**: `Authorization: Bearer <token>`
- **Response**:
  - `200`: `[{ "id": number, "day": string, "subject_id": number, "time_start": string, "time_end": string, "room": string, "teacher_id": number, "class_name": string }, ...]`
  - `404`: `{ "message": "User not found" }`
  - `500`: `{ "message": "Error retrieving timetable" }` or `{ "message": "Error scanning timetable entry" }`

#### GET /api/user (TokenAuthMiddleware)
- **Description**: Retrieves information about the logged-in user.
- **Header**: `Authorization: Bearer <token>`
- **Response**:
  - `200`: `{ "email": string, "first_name": string, "last_name": string, "birth_date": string, "address": string, "phone": string }`
  - `404`: `{ "message": "User not found" }` or `{ "message": "User details not found" }`

#### GET /api/exams (TokenAuthMiddleware)
- **Description**: Retrieves exams for the logged-in user (for students: their class; for teachers: their exams).
- **Header**: `Authorization: Bearer <token>`
- **Response**:
  - `200`: `[{ "id": number, "class_name": string, "teacher_id": number, "subject_id": number, "date": string, "type": string }, ...]`
  - `404`: `{ "message": "User not found" }`
  - `500`: `{ "message": "Error retrieving exams" }` or `{ "message": "Error scanning exam entry" }`

### Administrative Endpoints (Require admin role)
#### POST /api/register (TokenAuthMiddleware, AdminAuthMiddleware)
- **Description**: Registers a new user and their personal data.
- **Header**: `Authorization: Bearer <token>`
- **Body**:
  ```json
  {
    "email": string,
    "password": string,
    "role": string,
    "argument": {
      "email": string,
      "password": string,
      "role": string,
      "first_name": string,
      "last_name": string,
      "birth_date": string,
      "address": string,
      "phone": string
    }
  }
  ```
- **Response**:
  - `201`: `{ "message": "User created successfully" }`
  - `400`: `{ "message": "Invalid input" }` or `{ "message": "Invalid argument format" }`
  - `409`: `{ "message": "Email already taken" }`
  - `500`: `{ "message": "Error checking email" }`, `{ "message": "Error hashing password" }`, `{ "message": "Error saving user" }`, `{ "message": "Error retrieving user ID" }`, `{ "message": "Error saving user details" }`, or `{ "message": "Error committing transaction" }`

#### POST /api/timetable (TokenAuthMiddleware, AdminAuthMiddleware)
- **Description**: Adds a new schedule entry.
- **Header**: `Authorization: Bearer <token>`
- **Body**:
  ```json
  {
    "day": string,
    "subject_id": number,
    "time_start": string,
    "time_end": string,
    "room": string,
    "teacher_id": number,
    "class_name": string
  }
  ```
- **Response**:
  - `201`: `{ "message": "Timetable entry created successfully" }`
  - `400`: `{ "message": "Invalid input" }` or `{ "message": "Subject ID, start time, end time, teacher ID, class name, and day are required" }`
  - `500`: `{ "message": "Error saving timetable entry" }`

#### POST /api/admin/class (TokenAuthMiddleware, AdminAuthMiddleware)
- **Description**: Adds a new class.
- **Header**: `Authorization: Bearer <token>`
- **Body**: `{ "name": string }`
- **Response**:
  - `201`: `{ "message": "Class created successfully" }`
  - `400`: `{ "message": "Invalid input" }` or `{ "message": "Class name is required" }`
  - `500`: `{ "message": "Error saving class" }`

#### POST /api/admin/subject (TokenAuthMiddleware, AdminAuthMiddleware)
- **Description**: Adds a new subject.
- **Header**: `Authorization: Bearer <token>`
- **Body**: `{ "name": string, "class_name": string, "teacher_id": number }`
- **Response**:
  - `201`: `{ "message": "Subject created successfully" }`
  - `400`: `{ "message": "Invalid input" }` or `{ "message": "Subject name, class name, and teacher ID are required" }`
  - `500`: `{ "message": "Error saving subject" }`

#### POST /api/admin/class-member (TokenAuthMiddleware, AdminAuthMiddleware)
- **Description**: Adds a user to a class.
- **Header**: `Authorization: Bearer <token>`
- **Body**: `{ "user_id": number, "class_name": string }`
- **Response**:
  - `201`: `{ "message": "Class member added successfully" }`
  - `400`: `{ "message": "Invalid input" }` or `{ "message": "User ID and class name are required" }`
  - `500`: `{ "message": "Error saving class member" }`

#### POST /api/admin/grade (TokenAuthMiddleware, AdminAuthMiddleware)
- **Description**: Adds a grade, remark, or custom value for a student.
- **Header**: `Authorization: Bearer <token>`
- **Body**:
  ```json
  {
    "user_id": number,
    "subject_id": number,
    "grade": string,
    "grade_type": string,
    "date": string
  }
  ```
- **Response**:
  - `201`: `{ "message": "Grade created successfully" }`
  - `400`: `{ "message": "Invalid input" }` or `{ "message": "User ID, subject ID, grade, grade type, and date are required" }`
  - `500`: `{ "message": "Error saving grade" }`

#### POST /api/admin/attendance (TokenAuthMiddleware, AdminAuthMiddleware)
- **Description**: Adds attendance for a student.
- **Header**: `Authorization: Bearer <token>`
- **Body**:
  ```json
  {
    "user_id": number,
    "subject_id": number,
    "status": string,
    "date": string
  }
  ```
- **Response**:
  - `201`: `{ "message": "Attendance added successfully" }`
  - `400`: `{ "message": "Invalid input" }` or `{ "message": "User ID, subject ID, status, and date are required" }`
  - `500`: `{ "message": "Error saving attendance" }`

#### POST /api/admin/exam (TokenAuthMiddleware, AdminAuthMiddleware)
- **Description**: Adds a new exam.
- **Header**: `Authorization: Bearer <token>`
- **Body**:
  ```json
  {
    "class_name": string,
    "teacher_id": number,
    "subject_id": number,
    "date": string,
    "type": string
  }
  ```
- **Response**:
  - `201`: `{ "message": "Exam created successfully" }`
  - `400`: `{ "message": "Invalid input" }` or `{ "message": "Class name, teacher ID, subject ID, date, and type are required" }`
  - `500`: `{ "message": "Error saving exam" }`

#### POST /api/admin/class (TokenAuthMiddleware, AdminAuthMiddleware)
- **Description**: Retrieves the list of class members.
- **Header**: `Authorization: Bearer <token>`
- **Body**: `{ "name": string }`
- **Response**:
  - `200`: `[{ "id": number, "user_id": number, "class_name": string }, ...]`
  - `400`: `{ "message": "Invalid input" }`
  - `500`: `{ "message": "Error retrieving class members" }` or `{ "message": "Error scanning class member" }`

#### POST /api/admin/student-grades (TokenAuthMiddleware, AdminAuthMiddleware)
- **Description**: Retrieves grades for a specific student.
- **Header**: `Authorization: Bearer <token>`
- **Body**: `{ "uid": number }`
- **Response**:
  - `200`: `[{ "id": number, "user_id": number, "subject_id": number, "grade": string, "grade_type": string, "date": string }, ...]`
  - `400`: `{ "message": "Invalid input" }`
  - `500`: `{ "message": "Error retrieving grades" }` or `{ "message": "Error scanning grade" }`

#### POST /api/admin/student-attendance (TokenAuthMiddleware, AdminAuthMiddleware)
- **Description**: Retrieves attendance for a specific student.
- **Header**: `Authorization: Bearer <token>`
- **Body**: `{ "uid": number }`
- **Response**:
  - `200`: `[{ "id": number, "user_id": number, "subject_id": number, "status": string, "date": string }, ...]`
  - `400`: `{ "message": "Invalid input" }`
  - `500`: `{ "message": "Error retrieving attendance" }` or `{ "message": "Error scanning attendance" }`

#### POST /api/admin/student-info (TokenAuthMiddleware, AdminAuthMiddleware)
- **Description**: Retrieves personal data for a specific student.
- **Header**: `Authorization: Bearer <token>`
- **Body**: `{ "uid": number }`
- **Response**:
  - `200`: `{ "first_name": string, "last_name": string, "birth_date": string, "address": string, "phone": string }`
  - `400`: `{ "message": "Invalid input" }`
  - `404`: `{ "message": "User details not found" }`

### Teacher Endpoints (Require teacher role)
#### POST /api/grades/:user_id (TokenAuthMiddleware, TeacherAuthMiddleware)
- **Description**: Adds a grade, remark, or custom value for a student.
- **Header**: `Authorization: Bearer <token>`
- **Parameter**: `user_id` (student ID)
- **Body**:
  ```json
  {
    "user_id": number,
    "subject_id": number,
    "grade": string,
    "grade_type": string,
    "date": string
  }
  ```
- **Response**:
  - `201`: `{ "message": "Grade created successfully" }`
  - `400`: `{ "message": "Invalid input" }` or `{ "message": "User ID, subject ID, grade, grade type, and date are required" }`
  - `500`: `{ "message": "Error saving grade" }`

#### POST /api/teacher/attendance (TokenAuthMiddleware, TeacherAuthMiddleware)
- **Description**: Adds attendance for a student.
- **Header**: `Authorization: Bearer <token>`
- **Body**:
  ```json
  {
    "user_id": number,
    "subject_id": number,
    "status": string,
    "date": string
  }
  ```
- **Response**:
  - `201`: `{ "message": "Attendance added successfully" }`
  - `400`: `{ "message": "Invalid input" }` or `{ "message": "User ID, subject ID, status, and date are required" }`
  - `500`: `{ "message": "Error saving attendance" }`

#### POST /api/teacher/exam (TokenAuthMiddleware, TeacherAuthMiddleware)
- **Description**: Adds a new exam.
- **Header**: `Authorization: Bearer <token>`
- **Body**:
  ```json
  {
    "class_name": string,
    "teacher_id": number,
    "subject_id": number,
    "date": string,
    "type": string
  }
  ```
- **Response**:
  - `201`: `{ "message": "Exam created successfully" }`
  - `400`: `{ "message": "Invalid input" }` or `{ "message": "Class name, teacher ID, subject ID, date, and type are required" }`
  - `500`: `{ "message": "Error saving exam" }`

#### POST /api/teacher/class (TokenAuthMiddleware, TeacherAuthMiddleware)
- **Description**: Retrieves the list of class members.
- **Header**: `Authorization: Bearer <token>`
- **Body**: `{ "name": string }`
- **Response**:
  - `200`: `[{ "id": number, "user_id": number, "class_name": string }, ...]`
  - `400`: `{ "message": "Invalid input" }`
  - `500`: `{ "message": "Error retrieving class members" }` or `{ "message": "Error scanning class member" }`

#### POST /api/teacher/student-grades (TokenAuthMiddleware, TeacherAuthMiddleware)
- **Description**: Retrieves grades for a specific student.
- **Header**: `Authorization: Bearer <token>`
- **Body**: `{ "uid": number }`
- **Response**:
  - `200`: `[{ "id": number, "user_id": number, "subject_id": number, "grade": string, "grade_type": string, "date": string }, ...]`
  - `400`: `{ "message": "Invalid input" }`
  - `500`: `{ "message": "Error retrieving grades" }` or `{ "message": "Error scanning grade" }`

#### POST /api/teacher/student-attendance (TokenAuthMiddleware, TeacherAuthMiddleware)
- **Description**: Retrieves attendance for a specific student.
- **Header**: `Authorization: Bearer <token>`
- **Body**: `{ "uid": number }`
- **Response**:
  - `200`: `[{ "id": number, "user_id": number, "subject_id": number, "status": string, "date": string }, ...]`
  - `400`: `{ "message": "Invalid input" }`
  - `500`: `{ "message": "Error retrieving attendance" }` or `{ "message": "Error scanning attendance" }`

#### POST /api/teacher/student-info (TokenAuthMiddleware, TeacherAuthMiddleware)
- **Description**: Retrieves personal data for a specific student.
- **Header**: `Authorization: Bearer <token>`
- **Body**: `{ "uid": number }`
- **Response**:
  - `200`: `{ "first_name": string, "last_name": string, "birth_date": string, "address": string, "phone": string }`
  - `400`: `{ "message": "Invalid input" }`
  - `404`: `{ "message": "User details not found" }`

### Student Endpoints (Require student role)
#### GET /api/student/grades (TokenAuthMiddleware, StudentAuthMiddleware)
- **Description**: Retrieves grades for the logged-in student.
- **Header**: `Authorization: Bearer <token>`
- **Response**:
  - `200`: `[{ "id": number, "user_id": number, "subject_id": number, "grade": string, "grade_type": string, "date": string }, ...]`
  - `404`: `{ "message": "User not found" }`
  - `500`: `{ "message": "Error retrieving grades" }` or `{ "message": "Error scanning grade" }`

#### GET /api/student/subjects (TokenAuthMiddleware, StudentAuthMiddleware)
- **Description**: Retrieves subjects for the logged-in student's abrasion resistant coating.
- **Header**: `Authorization: Bearer <token>`
- **Response**:
  - `200`: `[{ "id": number, "name": string, "class_name": string, "teacher_id": number }, ...]`
  - `404`: `{ "message": "User not found" }`
  - `500`: `{ "message": "Error retrieving subjects" }` or `{ "message": "Error scanning subject" }`

#### GET /api/student/attendance (TokenAuthMiddleware, StudentAuthMiddleware)
- **Description**: Retrieves attendance for the logged-in student.
- **Header**: `Authorization: Bearer <token>`
- **Response**:
  - `200`: `[{ "id": number, "user_id": number, "subject_id": number, "status": string, "date": string }, ...]`
  - `404`: `{ "message": "User not found" }`
  - `500`: `{ "message": "Error retrieving attendance" }` or `{ "message": "Error scanning attendance" }`

## 6. Middleware
The application uses four middleware for authentication and authorization:
- **TokenAuthMiddleware**:
  - Verifies the JWT token in the `Authorization` header (format: `Bearer <token>`).
  - Sets email and role in the request context.
  - Used for all protected routes.
- **AdminAuthMiddleware**:
  - Checks if the user has the `admin` role (based on JWT and database).
  - Used for routes in the `/api/admin` group.
- **TeacherAuthMiddleware**:
  - Checks if the user has the `teacher` role (based on JWT and database).
  - Used for routes in the `/api/teacher` group.
- **StudentAuthMiddleware**:
  - Checks if the user has the `student` role (based on JWT and database).
  - Used for routes in the `/api/student` group.

**Additionally**:
- **LoggerMiddleware**: Logs HTTP request details (method, path, status, response time).
- **CORS**: Allows requests from any origin with `Authorization` and `Content-Type` headers.

## 7. Configuration
The application requires the following environment variables:
- `JWT_KEY`: Key for signing JWT tokens.
- `ADMIN_EMAIL`: Administrator's email (created during initialization).
- `ADMIN_PASSWORD`: Administrator's password.
- `DB_PATH` (optional): Path to the database file (default: `./database.db`).
- `PORT` (optional): Server port (default: `:10800`).
- `CERT_PATH` (optional): Path to the SSL certificate (default: `cert.pem`).
- `KEY_PATH` (optional): Path to the SSL key (default: `key.pem`).

**Example `.env` file**:
```
JWT_KEY=your-secret-key
ADMIN_EMAIL=admin@example.com
ADMIN_PASSWORD=secret123
DB_PATH=./database.db
PORT=:10800
CERT_PATH=cert.pem
KEY_PATH=key.pem
```

## 8. Running the Application
1. **Install dependencies**:
   ```bash
   go mod tidy
   ```
   Required packages:
   - `github.com/gin Stuart: `github.com/gin-gonic/gin`
   - `github.com/gin-contrib/cors`
   - `github.com/golang-jwt/jwt/v4`
   - `golang.org/x/crypto/bcrypt`
   - `modernc.org/sqlite`

2. **Create the `schema.sql` file**:
   - Copy the SQL schema (with tables `users`, `persons`, `classes`, `subjects`, `grades`, `timetable`, `attendance`, `exams`, `class_members`) to the `schema.sql` file in the project directory.

3. **Set environment variables**:
   - Use a `.env` file with the `godotenv` package or set variables in the system:
     ```bash
     export JWT_KEY=your-secret-key
     export ADMIN_EMAIL=admin@example.com
     export ADMIN_PASSWORD=secret123
     ```

4. **Run the application**:
   ```bash''
   go run .
   ```
   The server will start on `http://localhost:10800` (or HTTPS if certificates are provided).

5. **Testing**:
   - Use tools like Postman or curl, e.g.:
     ```bash
     curl -X POST http://localhost:10800/api/login -H "Content-Type: application/json" -d '{"email":"admin@example.com","password":"secret123"}'
     ```

## 10. Implementation Notes
- The `/api/grades/:user_id` endpoint requires a `user_id` parameter in the URL path, which is handled in the CLI client by dynamically building the address.
- Functions like `/api/admin/class`, `/api/teacher/class`, `/api/admin/student-grades`, `/api/teacher/student-grades`, etc., require sending data in JSON format with the `name` or `uid` field.
- The CLI code assumes the user knows their role (`admin`, `teacher`, `student`) since it does not decode the JWT token locally. If needed, token decoding can be added to automatically select appropriate paths (`/api/admin/*` or `/api/teacher/*`).

