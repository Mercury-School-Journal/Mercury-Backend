# Mercury Backend Documentation

## 1. Purpose
Mercury Backend is a Go-based server application providing a REST API for managing a school system. It supports user registration, authentication, timetable management, grades, and personal data for students and teachers. The application uses a SQLite database and JWT-based authentication with user roles (`student`, `teacher`, `admin`).

## 2. Project Structure
The application is primarily contained in a single `main.go` file, which includes:
- Database and server initialization.
- API route definitions and middleware.
- Helper functions, data models, and request handlers.

### Files and Roles:
- **main.go**: Main file configuring the Gin server, middleware, API routes, and database initialization.
- **schema.sql**: SQL file defining the database schema (tables `users`, `persons`, `classes`, `subjects`, etc.).
- **Go Models**: Data structures (e.g., `User`, `Person`, `Grade`, `TimetableEntry`) mapping to SQL tables.
- **Handlers**: Functions handling HTTP requests (e.g., `RegisterUser`, `Login`, `AddGrade`).
- **Middleware**: Authentication and authorization functions (`TokenAuthMiddleware`, `AdminAuthMiddleware`, `TeacherAuthMiddleware`).

## 3. Database Schema
The SQLite database includes the following tables:

- **users**: Stores user data (uid, email, password, role).
- **persons**: Stores personal information (user_id, first_name, last_name, birth_date, address, phone).
- **classes**: School classes (id, name).
- **subjects**: School subjects (id, name, class_name, teacher_id).
- **students_subjects**: Student-subject assignments (user_id, subject_id).
- **teachers_subjects**: Teacher-subject assignments (user_id, subject_id).
- **grades**: Grades, comments, or custom values (user_id, subject_id, grade, grade_type, date).
- **class_members**: User-class assignments (user_id, class_name).
- **timetable**: Timetable entries (day, subject_id, time_start, time_end, room, teacher_id, class_name).

The full schema is defined in `schema.sql`.

## 4. Data Models
Go models map to SQL tables and are used in handlers and HTTP requests:

- **User**: `{ UID, Email, Password, Role }` – User data.
- **Person**: `{ ID, UserID, FirstName, LastName, BirthDate, Address, Phone }` – Personal information.
- **Class**: `{ ID, Name }` – School class.
- **Subject**: `{ ID, Name, ClassName, TeacherID }` – School subject.
- **StudentSubject**: `{ ID, UserID, SubjectID }` – Student-subject assignment.
- **TeacherSubject**: `{ ID, UserID, SubjectID }` – Teacher-subject assignment.
- **Grade**: `{ ID, UserID, SubjectID, Grade, GradeType, Date }` – Grade or comment.
- **ClassMember**: `{ ID, UserID, ClassName }` – User-class assignment.
- **TimetableEntry**: `{ ID, Day, SubjectID, StartTime, EndTime, Room, TeacherID, ClassName }` – Timetable entry.
- **AccessRequest**: `{ Email, Password, Argument }` – Login/registration data.
- **Claims**: `{ Email, Role, StandardClaims }` – JWT claims.
- **Input**: `{ OldPassword, NewPassword }` – Password change request.

## 5. API Endpoints
The API is accessible at `http://localhost:10800/api` (or HTTPS if certificates are configured). Below are the endpoint details:

### Public Endpoints
- **POST /api/login**
  - **Description**: Authenticates a user and returns a JWT token.
  - **Body**: `{ "email": string, "password": string }`
  - **Responses**:
    - 200: `{ "token": string }`
    - 400: `{ "message": "Invalid input" }`
    - 401: `{ "message": "Invalid credentials" }`
  - **Example**: `POST /api/login { "email": "admin@example.com", "password": "secret" }`

- **GET /api/ping**
  - **Description**: Returns a 204 status to confirm server availability.
  - **Response**: 204 (No Content)

- **GET /api/timetable**
  - **Description**: Retrieves the entire timetable.
  - **Responses**:
    - 200: `[{ "id": number, "day": string, "subject_id": number, "start_time": string, "end_time": string, "room": string, "teacher_id": number, "class_name": string }, ...]`
    - 500: `{ "message": "Error retrieving timetable" }`

### Protected Endpoints (Require JWT)
- **PUT /api/change-password** (TokenAuthMiddleware)
  - **Description**: Updates the user’s password.
  - **Header**: `Authorization: Bearer <token>`
  - **Body**: `{ "old_password": string, "new_password": string }`
  - **Responses**:
    - 200: `{ "message": "Password changed successfully" }`
    - 400: `{ "message": "Invalid input" }`
    - 401: `{ "message": "Incorrect old password" }`
    - 404: `{ "message": "User not found" }`

- **DELETE /api/delete-account** (TokenAuthMiddleware)
  - **Description**: Deletes the user’s account.
  - **Header**: `Authorization: Bearer <token>`
  - **Responses**:
    - 200: `{ "message": "User deleted successfully" }`
    - 500: `{ "message": "Error deleting user" }`

### Admin Endpoints (Require admin role)
- **POST /api/register** (TokenAuthMiddleware, AdminAuthMiddleware)
  - **Description**: Registers a new user and their personal details.
  - **Header**: `Authorization: Bearer <token>`
  - **Body**: `{ "email": string, "password": string, "role": string, "argument": { "email": string, "password": string, "role": string, "first_name": string, "last_name": string, "birth_date": string, "address": string, "phone": string } }`
  - **Responses**:
    - 201: `{ "message": "User created successfully" }`
    - 400: `{ "message": "Invalid input" }`
    - 409: `{ "message": "Email already taken" }`
    - 500: `{ "message": "Error saving user" }`

- **POST /api/timetable** (TokenAuthMiddleware, AdminAuthMiddleware)
  - **Description**: Adds a new timetable entry.
  - **Header**: `Authorization: Bearer <token>`
  - **Body**: `{ "day": string, "subject_id": number, "start_time": string, "end_time": string, "room": string, "teacher_id": number, "class_name": string }`
  - **Responses**:
    - 201: `{ "message": "Timetable entry created successfully" }`
    - 400: `{ "message": "Invalid input" }`
    - 500: `{ "message": "Error saving timetable entry" }`

### Teacher Endpoints (Require teacher role)
- **POST /api/grades/:user_id** (TokenAuthMiddleware, TeacherAuthMiddleware)
  - **Description**: Adds a grade, comment, or custom value for a student.
  - **Header**: `Authorization: Bearer <token>`
  - **Parameter**: `user_id` (Student ID)
  - **Body**: `{ "user_id": number, "subject_id": number, "grade": string, "grade_type": string, "date": string }`
  - **Responses**:
    - 201: `{ "message": "Grade created successfully" }`
    - 400: `{ "message": "Invalid input" }`
    - 500: `{ "message": "Error saving grade" }`

## 6. Middleware
The application uses three middleware for authentication and authorization:

- **TokenAuthMiddleware**:
  - Validates the JWT token in the `Authorization` header (format: `Bearer <token>`).
  - Sets `email` and `role` in the request context.
  - Applied to all protected routes.

- **AdminAuthMiddleware**:
  - Verifies the user has the `admin` role (via JWT claims and database).
  - Applied to `/api/register` and `/api/timetable` (POST).

- **TeacherAuthMiddleware**:
  - Verifies the user has the `teacher` role (via JWT claims and database).
  - Applied to `/api/grades/:user_id` (POST).

Additional middleware:
- **LoggerMiddleware**: Logs HTTP request details (method, path, status, latency).
- **CORS**: Allows cross-origin requests with `Authorization` and `Content-Type` headers.

## 7. Configuration
The application requires the following environment variables:

- **JWT_KEY**: Secret key for signing JWT tokens.
- **ADMIN_EMAIL**: Email for the admin user (created during initialization).
- **ADMIN_PASSWORD**: Password for the admin user.
- **DB_PATH** (optional): Path to the SQLite database file (default: `./database.db`).
- **PORT** (optional): Server port (default: `:10800`).
- **CERT_PATH** (optional): Path to the SSL certificate (default: `cert.pem`).
- **KEY_PATH** (optional): Path to the SSL key (default: `key.pem`).

### Example `.env` file:
```plaintext
JWT_KEY=your-secret-key
ADMIN_EMAIL=admin@example.com
ADMIN_PASSWORD=secret123
DB_PATH=./database.db
PORT=:10800
CERT_PATH=cert.pem
KEY_PATH=key.pem
```

## 8. Setup and Running
1. **Install Dependencies:**
   ```bash
   go mod tidy
   ```
   Required packages: 
   - `github.com/gin-gonic/gin`
   - `github.com/gin-contrib/cors`
   - `github.com/golang-jwt/jwt/v4`
   - `golang.org/x/crypto/bcrypt`
   - `modernc.org/sqlite`.

2. **Create `schema.sql`:**
   Copy the updated SQL schema (including `persons`, `subject_id`, `grade_type`) to `schema.sql` in the project directory.

3. **Set Environment Variables:**
   Use a `.env` file with `godotenv` or set variables in the system:
   ```bash
   export JWT_KEY=your-secret-key
   export ADMIN_EMAIL=admin@example.com
   export ADMIN_PASSWORD=secret123
   ```

4. **Run the Application:**
   ```bash
   go run .
   ```
   The server starts on `http://localhost:10800` (or HTTPS if certificates are provided).

### Test the API:
Use tools like Postman or curl, e.g.:
```bash
curl -X POST http://localhost:10800/api/login -H "Content-Type: application/json" -d '{"email":"admin@example.com","password":"secret123"}'
```

## 9. Security
- **Passwords**: Passwords are hashed using bcrypt before storage.
- **JWT**: Tokens are signed with `JWT_KEY` and valid for 7 days.
- **Roles**: `AdminAuthMiddleware` and `TeacherAuthMiddleware` restrict access to appropriate roles.
- **CORS**: Configured to allow requests from any origin, which may need tightening in production.
- **HTTPS**: Optional support for HTTPS (requires certificates).
