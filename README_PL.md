# Dokumentacja aplikacji Mercury Backend

## 1. Cel aplikacji
Mercury Backend to aplikacja serwerowa napisana w Go, która dostarcza REST API do zarządzania systemem szkolnym. Umożliwia rejestrację użytkowników, logowanie, zarządzanie planem lekcji, ocenami oraz danymi osobowymi uczniów i nauczycieli. Aplikacja używa bazy danych SQLite oraz uwierzytelniania opartego na JWT z rolami użytkowników (student, teacher, admin).

## 2. Struktura projektu
Aplikacja składa się z jednego głównego pliku `main.go`, który zawiera:
- Inicjalizację bazy danych i serwera.
- Definicję tras API i middleware.
- Funkcje pomocnicze, modele danych oraz handlery.

### Pliki i ich role:
- `main.go`: Główny plik konfigurujący serwer Gin, middleware, trasy API oraz inicjalizację bazy danych.
- `schema.sql`: Plik SQL definiujący schemat bazy danych (tabele users, persons, classes, subjects, itp.).
- Modele Go: Struktury danych (np. User, Person, Grade, TimetableEntry) mapujące tabele SQL.
- Handlery: Funkcje obsługujące żądania HTTP (np. RegisterUser, Login, AddGrade).
- Middleware: Funkcje uwierzytelniania i autoryzacji (TokenAuthMiddleware, AdminAuthMiddleware, TeacherAuthMiddleware).

## 3. Schemat bazy danych
Baza danych SQLite zawiera następujące tabele:
- `users`: Przechowuje dane użytkowników (uid, email, password, role).
- `persons`: Dane osobowe użytkowników (user_id, first_name, last_name, birth_date, address, phone).
- `classes`: Klasy szkolne (id, name).
- `subjects`: Przedmioty szkolne (id, name, class_name, teacher_id).
- `students_subjects`: Powiązania uczniów z przedmiotami (user_id, subject_id).
- `teachers_subjects`: Powiązania nauczycieli z przedmiotami (user_id, subject_id).
- `grades`: Oceny, uwagi i wartości niestandardowe (user_id, subject_id, grade, grade_type, date).
- `class_members`: Powiązania użytkowników z klasami (user_id, class_name).
- `timetable`: Plan lekcji (day, subject_id, time_start, time_end, room, teacher_id, class_name).

Szczegółowy schemat znajduje się w pliku `schema.sql`.

## 4. Modele danych
Modele Go mapują tabele SQL i są używane w handlerach oraz żądaniach HTTP:
- `User`: `{ UID, Email, Password, Role }` – dane użytkownika.
- `Person`: `{ ID, UserID, FirstName, LastName, BirthDate, Address, Phone }` – dane osobowe.
- `Class`: `{ ID, Name }` – klasa szkolna.
- `Subject`: `{ ID, Name, ClassName, TeacherID }` – przedmiot.
- `StudentSubject`: `{ ID, UserID, SubjectID }` – powiązanie ucznia z przedmiotem.
- `TeacherSubject`: `{ ID, UserID, SubjectID }` – powiązanie nauczyciela z przedmiotem.
- `Grade`: `{ ID, UserID, SubjectID, Grade, GradeType, Date }` – ocena/uwaga.
- `ClassMember`: `{ ID, UserID, ClassName }` – powiązanie z klasą.
- `TimetableEntry`: `{ ID, Day, SubjectID, StartTime, EndTime, Room, TeacherID, ClassName }` – wpis w planie lekcji.
- `AccessRequest`: `{ Email, Password, Argument }` – dane logowania/rejestracji.
- `Claims`: `{ Email, Role, StandardClaims }` – dane JWT.
- `Input`: `{ OldPassword, NewPassword }` – zmiana hasła.

## 5. Endpointy API
API jest dostępne pod adresem `http://localhost:10800/api` (lub HTTPS, jeśli skonfigurowano certyfikaty). Poniżej opis endpointów:

### Publiczne endpointy
- **POST /api/login**
  - Opis: Loguje użytkownika i zwraca token JWT.
  - Body: `{ "email": string, "password": string }`
  - Odpowiedź:
    - `200`: `{ "token": string }`
    - `400`: `{ "message": "Invalid input" }`
    - `401`: `{ "message": "Invalid credentials" }`
  - Przykład: 
    ```json
    POST /api/login 
    { "email": "admin@example.com", "password": "secret" }
    ```

- **GET /api/ping**
  - Opis: Zwraca status 204, potwierdzając działanie serwera.
  - Odpowiedź: `204 (No Content)`

- **GET /api/timetable**
  - Opis: Pobiera cały plan lekcji.
  - Odpowiedź:
    - `200`: `[{ "id": number, "day": string, "subject_id": number, "start_time": string, "end_time": string, "room": string, "teacher_id": number, "class_name": string }, ...]`
    - `500`: `{ "message": "Error retrieving timetable" }`

### Endpointy chronione (wymagają JWT)
- **PUT /api/change-password** (TokenAuthMiddleware)
  - Opis: Zmienia hasło użytkownika.
  - Nagłówek: `Authorization: Bearer <token>`
  - Body: `{ "old_password": string, "new_password": string }`
  - Odpowiedź:
    - `200`: `{ "message": "Password changed successfully" }`
    - `400`: `{ "message": "Invalid input" }`
    - `401`: `{ "message": "Incorrect old password" }`
    - `404`: `{ "message": "User not found" }`

- **DELETE /api/delete-account** (TokenAuthMiddleware)
  - Opis: Usuwa konto użytkownika.
  - Nagłówek: `Authorization: Bearer <token>`
  - Odpowiedź:
    - `200`: `{ "message": "User deleted successfully" }`
    - `500`: `{ "message": "Error deleting user" }`

### Endpointy administracyjne (wymagają roli admin)
- **POST /api/register** (TokenAuthMiddleware, AdminAuthMiddleware)
  - Opis: Rejestruje nowego użytkownika i jego dane osobowe.
  - Nagłówek: `Authorization: Bearer <token>`
  - Body: 
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
  - Odpowiedź:
    - `201`: `{ "message": "User created successfully" }`
    - `400`: `{ "message": "Invalid input" }`
    - `409`: `{ "message": "Email already taken" }`
    - `500`: `{ "message": "Error saving user" }`

- **POST /api/timetable** (TokenAuthMiddleware, AdminAuthMiddleware)
  - Opis: Dodaje nowy wpis do planu lekcji.
  - Nagłówek: `Authorization: Bearer <token>`
  - Body: 
    ```json
    { 
      "day": string, 
      "subject_id": number, 
      "start_time": string, 
      "end_time": string, 
      "room": string, 
      "teacher_id": number, 
      "class_name": string 
    }
    ```
  - Odpowiedź:
    - `201`: `{ "message": "Timetable entry created successfully" }`
    - `400`: `{ "message": "Invalid input" }`
    - `500`: `{ "message": "Error saving timetable entry" }`

### Endpointy nauczycielskie (wymagają roli teacher)
- **POST /api/grades/:user_id** (TokenAuthMiddleware, TeacherAuthMiddleware)
  - Opis: Dodaje ocenę, uwagę lub wartość niestandardową dla ucznia.
  - Nagłówek: `Authorization: Bearer <token>`
  - Parametr: `user_id` (ID ucznia)
  - Body: 
    ```json
    { 
      "user_id": number, 
      "subject_id": number, 
      "grade": string, 
      "grade_type": string, 
      "date": string 
    }
    ```
  - Odpowiedź:
    - `201`: `{ "message": "Grade created successfully" }`
    - `400`: `{ "message": "Invalid input" }`
    - `500`: `{ "message": "Error saving grade" }`

## 6. Middleware
Aplikacja używa trzech middleware do uwierzytelniania i autoryzacji:
- **TokenAuthMiddleware**:
  - Weryfikuje token JWT w nagłówku Authorization (format: Bearer `<token>`).
  - Ustawia email i rolę w kontekście żądania.
  - Używany dla wszystkich chronionych tras.

- **AdminAuthMiddleware**:
  - Sprawdza, czy użytkownik ma rolę admin (na podstawie JWT i bazy danych).
  - Używany dla tras `/api/register` i `/api/timetable` (POST).

- **TeacherAuthMiddleware**:
  - Sprawdza, czy użytkownik ma rolę teacher (na podstawie JWT i bazy danych).
  - Używany dla trasy `/api/grades/:user_id` (POST).

Dodatkowo:
- **LoggerMiddleware**: Loguje szczegóły żądań HTTP (metoda, ścieżka, status, czas odpowiedzi).
- **CORS**: Pozwala na żądania z dowolnego źródła z nagłówkami Authorization, Content-Type.

## 7. Konfiguracja
Aplikacja wymaga ustawienia zmiennych środowiskowych:
- `JWT_KEY`: Klucz do podpisywania tokenów JWT.
- `ADMIN_EMAIL`: E-mail administratora (tworzony przy inicjalizacji).
- `ADMIN_PASSWORD`: Hasło administratora.
- `DB_PATH` (opcjonalne): Ścieżka do pliku bazy danych (domyślnie `./database.db`).
- `PORT` (opcjonalne): Port serwera (domyślnie `:10800`).
- `CERT_PATH` (opcjonalne): Ścieżka do certyfikatu SSL (domyślnie `cert.pem`).
- `KEY_PATH` (opcjonalne): Ścieżka do klucza SSL (domyślnie `key.pem`).

### Przykładowy plik `.env`:
```plaintext
JWT_KEY=your-secret-key
ADMIN_EMAIL=admin@example.com
ADMIN_PASSWORD=secret123
DB_PATH=./database.db
PORT=:10800
CERT_PATH=cert.pem
KEY_PATH=key.pem
```

## 8. Uruchomienie aplikacji
1. Zainstaluj zależności:
   ```bash
   go mod tidy
   ```
   Wymagane pakiety: 
   - `github.com/gin-gonic/gin`
   - `github.com/gin-contrib/cors`
   - `github.com/golang-jwt/jwt/v4`
   - `golang.org/x/crypto/bcrypt`
   - `modernc.org/sqlite`.

2. Utwórz plik `schema.sql`:
   Skopiuj poprawiony schemat SQL (z tabelą `persons`, `subject_id`, `grade_type`) do pliku `schema.sql` w katalogu projektu.

3. Ustaw zmienne środowiskowe:
   Użyj pliku `.env` z pakietem `godotenv` lub ustaw zmienne w systemie:
   ```bash
   export JWT_KEY=your-secret-key
   export ADMIN_EMAIL=admin@example.com
   export ADMIN_PASSWORD=secret123
   ```

4. Uruchom aplikację:
   ```bash
   go run .
   ```
   Serwer uruchomi się na `http://localhost:10800` (lub HTTPS, jeśli podano certyfikaty).

### Testowanie:
Użyj narzędzia jak Postman lub curl, np.:
```bash
curl -X POST http://localhost:10800/api/login -H "Content-Type: application/json" -d '{"email":"admin@example.com","password":"secret123"}'
```

## 9. Bezpieczeństwo
- **Hasła**: Hasła są hashowane za pomocą bcrypt przed zapisem do bazy.
- **JWT**: Tokeny JWT są podpisywane kluczem `JWT_KEY` i mają 7-dniowy okres ważności.
- **Role**: Middleware `AdminAuthMiddleware` i `TeacherAuthMiddleware` ograniczają dostęp do odpowiednich ról.
- **CORS**: Ustawienia pozwalają na żądania z dowolnego źródła, co może wymagać zaostrzenia w produkcji.
- **HTTPS**: Opcjonalne wsparcie dla HTTPS (wymaga certyfikatów).