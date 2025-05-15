
# Dokumentacja aplikacji Mercury Backend

## 1. Cel aplikacji
Mercury Backend to aplikacja serwerowa napisana w Go, która dostarcza REST API do zarządzania systemem szkolnym. Umożliwia rejestrację użytkowników, logowanie, zarządzanie planem lekcji, ocenami, obecnościami, egzaminami, klasami, przedmiotami oraz danymi osobowymi uczniów i nauczycieli. Aplikacja używa bazy danych SQLite oraz uwierzytelniania opartego na JWT z rolami użytkowników (student, teacher, admin).

## 2. Struktura projektu
Aplikacja składa się z jednego głównego pliku `main.go`, który zawiera:
- Inicjalizację bazy danych i serwera.
- Definicję tras API i middleware.
- Funkcje pomocnicze, modele danych oraz handlery.

**Pliki i ich role:**
- `main.go`: Główny plik konfigurujący serwer Gin, middleware, trasy API oraz inicjalizację bazy danych.
- `schema.sql`: Plik SQL definiujący schemat bazy danych (tabele `users`, `persons`, `classes`, `subjects`, `grades`, `timetable`, `attendance`, `exams`, `class_members`).
- **Modele Go**: Struktury danych (np. `User`, `Person`, `Grade`, `TimetableEntry`, `Attendance`, `Exam`, `Class`, `Subject`, `ClassMember`) mapujące tabele SQL.
- **Handlery**: Funkcje obsługujące żądania HTTP (np. `RegisterUser`, `Login`, `AddGrade`, `AddAttendance`, `AddExam`).
- **Middleware**: Funkcje uwierzytelniania i autoryzacji (`TokenAuthMiddleware`, `AdminAuthMiddleware`, `TeacherAuthMiddleware`, `StudentAuthMiddleware`).

## 3. Schemat bazy danych
Baza danych SQLite zawiera następujące tabele:
- `users`: Przechowuje dane użytkowników (`uid`, `email`, `password`, `role`).
- `persons`: Dane osobowe użytkowników (`user_id`, `first_name`, `last_name`, `birth_date`, `address`, `phone`).
- `classes`: Klasy szkolne (`id`, `name`).
- `subjects`: Przedmioty szkolne (`id`, `name`, `class_name`, `teacher_id`).
- `grades`: Oceny, uwagi i wartości niestandardowe (`id`, `user_id`, `subject_id`, `grade`, `grade_type`, `date`).
- `class_members`: Powiązania użytkowników z klasami (`id`, `user_id`, `class_name`).
- `timetable`: Plan lekcji (`id`, `day`, `subject_id`, `time_start`, `time_end`, `room`, `teacher_id`, `class_name`).
- `attendance`: Obecności (`id`, `user_id`, `subject_id`, `status`, `date`).
- `exams`: Egzaminy (`id`, `class_name`, `teacher_id`, `subject_id`, `date`, `type`).

Szczegółowy schemat znajduje się w pliku `schema.sql`.

## 4. Modele danych
Modele Go mapują tabele SQL i są używane w handlerach oraz żądaniach HTTP:
- `User`: { `UID`, `Email`, `Password`, `Role` } – dane użytkownika.
- `Person`: { `ID`, `UserID`, `FirstName`, `LastName`, `BirthDate`, `Address`, `Phone` } – dane osobowe.
- `Class`: { `ID`, `Name` } – klasa szkolna.
- `Subject`: { `ID`, `Name`, `ClassName`, `TeacherID` } – przedmiot.
- `Grade`: { `ID`, `UserID`, `SubjectID`, `Grade`, `GradeType`, `Date` } – ocena/uwaga.
- `ClassMember`: { `ID`, `UserID`, `ClassName` } – powiązanie z klasą.
- `TimetableEntry`: { `ID`, `Day`, `SubjectID`, `StartTime`, `EndTime`, `Room`, `TeacherID`, `ClassName` } – wpis w planie lekcji.
- `Attendance`: { `ID`, `UserID`, `SubjectID`, `Status`, `Date` } – obecność.
- `Exam`: { `ID`, `ClassName`, `TeacherID`, `SubjectID`, `Date`, `Type` } – egzamin.
- `AccessRequest`: { `Email`, `Password`, `Argument` } – dane logowania/rejestracji.
- `Claims`: { `Email`, `Role`, `StandardClaims` } – dane JWT.
- `Input`: { `OldPassword`, `NewPassword` } – zmiana hasła.

## 5. Endpointy API
API jest dostępne pod adresem `http://localhost:10800/api` (lub HTTPS, jeśli skonfigurowano certyfikaty). Poniżej opis endpointów:

### Publiczne endpointy
#### POST /api/login
- **Opis**: Loguje użytkownika i zwraca token JWT.
- **Body**: `{ "email": string, "password": string }`
- **Odpowiedź**:
  - `200`: `{ "token": string }`
  - `400`: `{ "message": "Invalid input" }`
  - `401`: `{ "message": "Invalid email credentials" }` lub `{ "message": "Invalid password credentials" }`
  - `500`: `{ "message": "Could not generate token" }`
- **Przykład**:
  ```json
  POST /api/login
  { "email": "admin@example.com", "password": "secret123" }
  ```

#### GET /api/ping
- **Opis**: Zwraca status 204, potwierdzając działanie serwera.
- **Odpowiedź**: `204` (No Content)

#### GET /api/lucky-number
- **Opis**: Zwraca losowy numer.
- **Odpowiedź**:
  - `200`: `{ "lucky_number": number }`

### Endpointy chronione (wymagają JWT)
#### PUT /api/change-password (TokenAuthMiddleware)
- **Opis**: Zmienia hasło użytkownika.
- **Nagłówek**: `Authorization: Bearer <token>`
- **Body**: `{ "old_password": string, "new_password": string }`
- **Odpowiedź**:
  - `200`: `{ "message": "Password changed successfully" }`
  - `400`: `{ "message": "Invalid input" }`
  - `401`: `{ "message": "Incorrect old password" }`
  - `404`: `{ "message": "User not found" }`
  - `500`: `{ "message": "Error hashing new password" }` lub `{ "message": "Error updating password" }`

#### DELETE /api/delete-account (TokenAuthMiddleware)
- **Opis**: Usuwa konto użytkownika i powiązane dane.
- **Nagłówek**: `Authorization: Bearer <token>`
- **Odpowiedź**:
  - `200`: `{ "message": "User deleted successfully" }`
  - `500`: `{ "message": "Error deleting user details" }` lub `{ "message": "Error committing transaction" }`

#### GET /api/timetable (TokenAuthMiddleware)
- **Opis**: Pobiera plan lekcji dla zalogowanego użytkownika (dla studenta: dla jego klasy, dla nauczyciela: dla jego lekcji).
- **Nagłówek**: `Authorization: Bearer <token>`
- **Odpowiedź**:
  - `200`: `[{ "id": number, "day": string, "subject_id": number, "time_start": string, "time_end": string, "room": string, "teacher_id": number, "class_name": string }, ...]`
  - `404`: `{ "message": "User not found" }`
  - `500`: `{ "message": "Error retrieving timetable" }` lub `{ "message": "Error scanning timetable entry" }`

#### GET /api/user (TokenAuthMiddleware)
- **Opis**: Pobiera informacje o zalogowanym użytkowniku.
- **Nagłówek**: `Authorization: Bearer <token>`
- **Odpowiedź**:
  - `200`: `{ "email": string, "first_name": string, "last_name": string, "birth_date": string, "address": string, "phone": string }`
  - `404`: `{ "message": "User not found" }` lub `{ "message": "User details not found" }`

#### GET /api/exams (TokenAuthMiddleware)
- **Opis**: Pobiera egzaminy dla zalogowanego użytkownika (dla studenta: dla jego klasy, dla nauczyciela: jego egzaminy).
- **Nagłówek**: `Authorization: Bearer <token>`
- **Odpowiedź**:
  - `200`: `[{ "id": number, "class_name": string, "teacher_id": number, "subject_id": number, "date": string, "type": string }, ...]`
  - `404`: `{ "message": "User not found" }`
  - `500`: `{ "message": "Error retrieving exams" }` lub `{ "message": "Error scanning exam entry" }`

### Endpointy administracyjne (wymagają roli admin)
#### POST /api/register (TokenAuthMiddleware, AdminAuthMiddleware)
- **Opis**: Rejestruje nowego użytkownika i jego dane osobowe.
- **Nagłówek**: `Authorization: Bearer <token>`
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
- **Odpowiedź**:
  - `201`: `{ "message": "User created successfully" }`
  - `400`: `{ "message": "Invalid input" }` lub `{ "message": "Invalid argument format" }`
  - `409`: `{ "message": "Email already taken" }`
  - `500`: `{ "message": "Error checking email" }`, `{ "message": "Error hashing password" }`, `{ "message": "Error saving user" }`, `{ "message": "Error retrieving user ID" }`, `{ "message": "Error saving user details" }`, lub `{ "message": "Error committing transaction" }`

#### POST /api/timetable (TokenAuthMiddleware, AdminAuthMiddleware)
- **Opis**: Dodaje nowy wpis do planu lekcji.
- **Nagłówek**: `Authorization: Bearer <token>`
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
- **Odpowiedź**:
  - `201`: `{ "message": "Timetable entry created successfully" }`
  - `400`: `{ "message": "Invalid input" }` lub `{ "message": "Subject ID, start time, end time, teacher ID, class name, and day are required" }`
  - `500`: `{ "message": "Error saving timetable entry" }`

#### POST /api/admin/class (TokenAuthMiddleware, AdminAuthMiddleware)
- **Opis**: Dodaje nową klasę.
- **Nagłówek**: `Authorization: Bearer <token>`
- **Body**: `{ "name": string }`
- **Odpowiedź**:
  - `201`: `{ "message": "Class created successfully" }`
  - `400`: `{ "message": "Invalid input" }` lub `{ "message": "Class name is required" }`
  - `500`: `{ "message": "Error saving class" }`

#### POST /api/admin/subject (TokenAuthMiddleware, AdminAuthMiddleware)
- **Opis**: Dodaje nowy przedmiot.
- **Nagłówek**: `Authorization: Bearer <token>`
- **Body**: `{ "name": string, "class_name": string, "teacher_id": number }`
- **Odpowiedź**:
  - `201`: `{ "message": "Subject created successfully" }`
  - `400`: `{ "message": "Invalid input" }` lub `{ "message": "Subject name, class name, and teacher ID are required" }`
  - `500`: `{ "message": "Error saving subject" }`

#### POST /api/admin/class-member (TokenAuthMiddleware, AdminAuthMiddleware)
- **Opis**: Dodaje użytkownika do klasy.
- **Nagłówek**: `Authorization: Bearer <token>`
- **Body**: `{ "user_id": number, "class_name": string }`
- **Odpowiedź**:
  - `201`: `{ "message": "Class member added successfully" }`
  - `400`: `{ "message": "Invalid input" }` lub `{ "message": "User ID and class name are required" }`
  - `500`: `{ "message": "Error saving class member" }`

#### POST /api/admin/grade (TokenAuthMiddleware, AdminAuthMiddleware)
- **Opis**: Dodaje ocenę, uwagę lub wartość niestandardową dla ucznia.
- **Nagłówek**: `Authorization: Bearer <token>`
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
- **Odpowiedź**:
  - `201`: `{ "message": "Grade created successfully" }`
  - `400`: `{ "message": "Invalid input" }` lub `{ "message": "User ID, subject ID, grade, grade type, and date are required" }`
  - `500`: `{ "message": "Error saving grade" }`

#### POST /api/admin/attendance (TokenAuthMiddleware, AdminAuthMiddleware)
- **Opis**: Dodaje obecność dla ucznia.
- **Nagłówek**: `Authorization: Bearer <token>`
- **Body**:
  ```json
  {
    "user_id": number,
    "subject_id": number,
    "status": string,
    "date": string
  }
  ```
- **Odpowiedź**:
  - `201`: `{ "message": "Attendance added successfully" }`
  - `400`: `{ "message": "Invalid input" }` lub `{ "message": "User ID, subject ID, status, and date are required" }`
  - `500`: `{ "message": "Error saving attendance" }`

#### POST /api/admin/exam (TokenAuthMiddleware, AdminAuthMiddleware)
- **Opis**: Dodaje nowy egzamin.
- **Nagłówek**: `Authorization: Bearer <token>`
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
- **Odpowiedź**:
  - `201`: `{ "message": "Exam created successfully" }`
  - `400`: `{ "message": "Invalid input" }` lub `{ "message": "Class name, teacher ID, subject ID, date, and type are required" }`
  - `500`: `{ "message": "Error saving exam" }`

#### POST /api/admin/class (TokenAuthMiddleware, AdminAuthMiddleware)
- **Opis**: Pobiera listę członków klasy.
- **Nagłówek**: `Authorization: Bearer <token>`
- **Body**: `{ "name": string }`
- **Odpowiedź**:
  - `200`: `[{ "id": number, "user_id": number, "class_name": string }, ...]`
  - `400`: `{ "message": "Invalid input" }`
  - `500`: `{ "message": "Error retrieving class members" }` lub `{ "message": "Error scanning class member" }`

#### POST /api/admin/student-grades (TokenAuthMiddleware, AdminAuthMiddleware)
- **Opis**: Pobiera oceny konkretnego ucznia.
- **Nagłówek**: `Authorization: Bearer <token>`
- **Body**: `{ "uid": number }`
- **Odpowiedź**:
  - `200`: `[{ "id": number, "user_id": number, "subject_id": number, "grade": string, "grade_type": string, "date": string }, ...]`
  - `400`: `{ "message": "Invalid input" }`
  - `500`: `{ "message": "Error retrieving grades" }` lub `{ "message": "Error scanning grade" }`

#### POST /api/admin/student-attendance (TokenAuthMiddleware, AdminAuthMiddleware)
- **Opis**: Pobiera obecności konkretnego ucznia.
- **Nagłówek**: `Authorization: Bearer <token>`
- **Body**: `{ "uid": number }`
- **Odpowiedź**:
  - `200`: `[{ "id": number, "user_id": number, "subject_id": number, "status": string, "date": string }, ...]`
  - `400`: `{ "message": "Invalid input" }`
  - `500`: `{ "message": "Error retrieving attendance" }` lub `{ "message": "Error scanning attendance" }`

#### POST /api/admin/student-info (TokenAuthMiddleware, AdminAuthMiddleware)
- **Opis**: Pobiera dane osobowe konkretnego ucznia.
- **Nagłówek**: `Authorization: Bearer <token>`
- **Body**: `{ "uid": number }`
- **Odpowiedź**:
  - `200`: `{ "first_name": string, "last_name": string, "birth_date": string, "address": string, "phone": string }`
  - `400`: `{ "message": "Invalid input" }`
  - `404`: `{ "message": "User details not found" }`

### Endpointy nauczycielskie (wymagają roli teacher)
#### POST /api/grades/:user_id (TokenAuthMiddleware, TeacherAuthMiddleware)
- **Opis**: Dodaje ocenę, uwagę lub wartość niestandardową dla ucznia.
- **Nagłówek**: `Authorization: Bearer <token>`
- **Parametr**: `user_id` (ID ucznia)
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
- **Odpowiedź**:
  - `201`: `{ "message": "Grade created successfully" }`
  - `400`: `{ "message": "Invalid input" }` lub `{ "message": "User ID, subject ID, grade, grade type, and date are required" }`
  - `500`: `{ "message": "Error saving grade" }`

#### POST /api/teacher/attendance (TokenAuthMiddleware, TeacherAuthMiddleware)
- **Opis**: Dodaje obecność dla ucznia.
- **Nagłówek**: `Authorization: Bearer <token>`
- **Body**:
  ```json
  {
    "user_id": number,
    "subject_id": number,
    "status": string,
    "date": string
  }
  ```
- **Odpowiedź**:
  - `201`: `{ "message": "Attendance added successfully" }`
  - `400`: `{ "message": "Invalid input" }` lub `{ "message": "User ID, subject ID, status, and date are required" }`
  - `500`: `{ "message": "Error saving attendance" }`

#### POST /api/teacher/exam (TokenAuthMiddleware, TeacherAuthMiddleware)
- **Opis**: Dodaje nowy egzamin.
- **Nagłówek**: `Authorization: Bearer <token>`
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
- **Odpowiedź**:
  - `201`: `{ "message": "Exam created successfully" }`
  - `400`: `{ "message": "Invalid input" }` lub `{ "message": "Class name, teacher ID, subject ID, date, and type are required" }`
  - `500`: `{ "message": "Error saving exam" }`

#### POST /api/teacher/class (TokenAuthMiddleware, TeacherAuthMiddleware)
- **Opis**: Pobiera listę członków klasy.
- **Nagłówek**: `Authorization: Bearer <token>`
- **Body**: `{ "name": string }`
- **Odpowiedź**:
  - `200`: `[{ "id": number, "user_id": number, "class_name": string }, ...]`
  - `400`: `{ "message": "Invalid input" }`
  - `500`: `{ "message": "Error retrieving class members" }` lub `{ "message": "Error scanning class member" }`

#### POST /api/teacher/student-grades (TokenAuthMiddleware, TeacherAuthMiddleware)
- **Opis**: Pobiera oceny konkretnego ucznia.
- **Nagłówek**: `Authorization: Bearer <token>`
- **Body**: `{ "uid": number }`
- **Odpowiedź**:
  - `200`: `[{ "id": number, "user_id": number, "subject_id": number, "grade": string, "grade_type": string, "date": string }, ...]`
  - `400`: `{ "message": "Invalid input" }`
  - `500`: `{ "message": "Error retrieving grades" }` lub `{ "message": "Error scanning grade" }`

#### POST /api/teacher/student-attendance (TokenAuthMiddleware, TeacherAuthMiddleware)
- **Opis**: Pobiera obecności konkretnego ucznia.
- **Nagłówek**: `Authorization: Bearer <token>`
- **Body**: `{ "uid": number }`
- **Odpowiedź**:
  - `200`: `[{ "id": number, "user_id": number, "subject_id": number, "status": string, "date": string }, ...]`
  - `400`: `{ "message": "Invalid input" }`
  - `500`: `{ "message": "Error retrieving attendance" }` lub `{ "message": "Error scanning attendance" }`

#### POST /api/teacher/student-info (TokenAuthMiddleware, TeacherAuthMiddleware)
- **Opis**: Pobiera dane osobowe konkretnego ucznia.
- **Nagłówek**: `Authorization: Bearer <token>`
- **Body**: `{ "uid": number }`
- **Odpowiedź**:
  - `200`: `{ "first_name": string, "last_name": string, "birth_date": string, "address": string, "phone": string }`
  - `400`: `{ "message": "Invalid input" }`
  - `404`: `{ "message": "User details not found" }`

### Endpointy studenckie (wymagają roli student)
#### GET /api/student/grades (TokenAuthMiddleware, StudentAuthMiddleware)
- **Opis**: Pobiera oceny zalogowanego ucznia.
- **Nagłówek**: `Authorization: Bearer <token>`
- **Odpowiedź**:
  - `200`: `[{ "id": number, "user_id": number, "subject_id": number, "grade": string, "grade_type": string, "date": string }, ...]`
  - `404`: `{ "message": "User not found" }`
  - `500`: `{ "message": "Error retrieving grades" }` lub `{ "message": "Error scanning grade" }`

#### GET /api/student/subjects (TokenAuthMiddleware, StudentAuthMiddleware)
- **Opis**: Pobiera przedmioty dla klasy zalogowanego ucznia.
- **Nagłówek**: `Authorization: Bearer <token>`
- **Odpowiedź**:
  - `200`: `[{ "id": number, "name": string, "class_name": string, "teacher_id": number }, ...]`
  - `404`: `{ "message": "User not found" }`
  - `500`: `{ "message": "Error retrieving subjects" }` lub `{ "message": "Error scanning subject" }`

#### GET /api/student/attendance (TokenAuthMiddleware, StudentAuthMiddleware)
- **Opis**: Pobiera obecności zalogowanego ucznia.
- **Nagłówek**: `Authorization: Bearer <token>`
- **Odpowiedź**:
  - `200`: `[{ "id": number, "user_id": number, "subject_id": number, "status": string, "date": string }, ...]`
  - `404`: `{ "message": "User not found" }`
  - `500`: `{ "message": "Error retrieving attendance" }` lub `{ "message": "Error scanning attendance" }`

## 6. Middleware
Aplikacja używa czterech middleware do uwierzytelniania i autoryzacji:
- **TokenAuthMiddleware**:
  - Weryfikuje token JWT w nagłówku `Authorization` (format: `Bearer <token>`).
  - Ustawia email i rolę w kontekście żądania.
  - Używany dla wszystkich chronionych tras.
- **AdminAuthMiddleware**:
  - Sprawdza, czy użytkownik ma rolę `admin` (na podstawie JWT i bazy danych).
  - Używany dla tras w grupie `/api/admin`.
- **TeacherAuthMiddleware**:
  - Sprawdza, czy użytkownik ma rolę `teacher` (na podstawie JWT i bazy danych).
  - Używany dla tras w grupie `/api/teacher`.
- **StudentAuthMiddleware**:
  - Sprawdza, czy użytkownik ma rolę `student` (na podstawie JWT i bazy danych).
  - Używany dla tras w grupie `/api/student`.

**Dodatkowo**:
- **LoggerMiddleware**: Loguje szczegóły żądań HTTP (metoda, ścieżka, status, czas odpowiedzi).
- **CORS**: Pozwala na żądania z dowolnego źródła z nagłówkami `Authorization`, `Content-Type`.

## 7. Konfiguracja
Aplikacja wymaga ustawienia zmiennych środowiskowych:
- `JWT_KEY`: Klucz do podpisywania tokenów JWT.
- `ADMIN_EMAIL`: E-mail administratora (tworzony przy inicjalizacji).
- `ADMIN_PASSWORD`: Hasło administratora.
- `DB_PATH` (opcjonalne): Ścieżka do pliku bazy danych (domyślnie `./database.db`).
- `PORT` (opcjonalne): Port serwera (domyślnie `:10800`).
- `CERT_PATH` (opcjonalne): Ścieżka do certyfikatu SSL (domyślnie `cert.pem`).
- `KEY_PATH` (opcjonalne): Ścieżka do klucza SSL (domyślnie `key.pem`).

**Przykładowy plik `.env`**:
```
JWT_KEY=your-secret-key
ADMIN_EMAIL=admin@example.com
ADMIN_PASSWORD=secret123
DB_PATH=./database.db
PORT=:10800
CERT_PATH=cert.pem
KEY_PATH=key.pem
```

## 8. Uruchomienie aplikacji
1. **Zainstaluj zależności**:
   ```bash
   go mod tidy
   ```
   Wymagane pakiety:
   - `github.com/gin-gonic/gin`
   - `github.com/gin-contrib/cors`
   - `github.com/golang-jwt/jwt/v4`
   - `golang.org/x/crypto/bcrypt`
   - `modernc.org/sqlite`

2. **Utwórz plik `schema.sql`**:
   - Skopiuj schemat SQL (z tabelami `users`, `persons`, `classes`, `subjects`, `grades`, `timetable`, `attendance`, `exams`, `class_members`) do pliku `schema.sql` w katalogu projektu.

3. **Ustaw zmienne środowiskowe**:
   - Użyj pliku `.env` z pakietem `godotenv` lub ustaw zmienne w systemie:
     ```bash
     export JWT_KEY=your-secret-key
     export ADMIN_EMAIL=admin@example.com
     export ADMIN_PASSWORD=secret123
     ```

4. **Uruchom aplikację**:
   ```bash
   go run .
   ```
   Serwer uruchomi się na `http://localhost:10800` (lub HTTPS, jeśli podano certyfikaty).

5. **Testowanie**:
   - Użyj narzędzia jak Postman lub curl, np.:
     ```bash
     curl -X POST http://localhost:10800/api/login -H "Content-Type: application/json" -d '{"email":"admin@example.com","password":"secret123"}'
     ```

## 9. Bezpieczeństwo
- **Hasła**: Hasła są hashowane za pomocą `bcrypt` przed zapisem do bazy.
- **JWT**: Tokeny JWT są podpisywane kluczem `JWT_KEY` i mają 7-dniowy okres ważności.
- **Role**: Middleware `AdminAuthMiddleware`, `TeacherAuthMiddleware`, i `StudentAuthMiddleware` ograniczają dostęp do odpowiednich ról.
- **CORS**: Ustawienia pozwalają na żądania z dowolnego źródła, co może wymagać zaostrzenia w produkcji.
- **HTTPS**: Opcjonalne wsparcie dla HTTPS (wymaga certyfikatów).

## 10. Uwagi do implementacji
- Endpoint `/api/grades/:user_id` wymaga parametru `user_id` w ścieżce URL, co jest obsługiwane w kliencie CLI poprzez dynamiczne budowanie adresu.
- Funkcje takie jak `/api/admin/class`, `/api/teacher/class`, `/api/admin/student-grades`, `/api/teacher/student-grades` itp. wymagają wysyłania danych w formacie JSON z polem `name` lub `uid`.
- Kod CLI zakłada, że użytkownik zna swoją rolę (`admin`, `teacher`, `student`), ponieważ nie dekoduje tokena JWT lokalnie. W razie potrzeby można dodać dekodowanie tokena, aby automatycznie wybierać odpowiednie ścieżki (`/api/admin/*` lub `/api/teacher/*`).

