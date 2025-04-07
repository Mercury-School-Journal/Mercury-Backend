-- schema.sql

CREATE TABLE IF NOT EXISTS users (
    uid INTEGER PRIMARY KEY AUTOINCREMENT,
    email TEXT UNIQUE,
    password TEXT,
    role TEXT CHECK(role IN ('student', 'teacher', 'admin')) NOT NULL
);

CREATE TABLE IF NOT EXISTS students (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER,
    first_name TEXT,
    last_name TEXT,
    birth_date TEXT,
    address TEXT,
    phone TEXT,
    email TEXT,
    FOREIGN KEY(user_id) REFERENCES users(uid)
);

CREATE TABLE IF NOT EXISTS teachers (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER,
    first_name TEXT,
    last_name TEXT,
    birth_date TEXT,
    address TEXT,
    phone TEXT,
    email TEXT,
    FOREIGN KEY(user_id) REFERENCES users(uid)
);

CREATE TABLE IF NOT EXISTS students_subjects (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER,
    subject TEXT,
    FOREIGN KEY(user_id) REFERENCES users(uid),
    FOREIGN KEY(subject) REFERENCES subjects(name)
);

CREATE TABLE IF NOT EXISTS teachers_subjects (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER,
    subject TEXT,
    FOREIGN KEY(user_id) REFERENCES users(uid),
    FOREIGN KEY(subject) REFERENCES subjects(name)
);

CREATE TABLE IF NOT EXISTS grades (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER,
    subject TEXT,
    grade INTEGER,
    FOREIGN KEY(user_id) REFERENCES users(uid),
    FOREIGN KEY(subject) REFERENCES subjects(name)
);

CREATE TABLE IF NOT EXISTS subjects (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT,
    class TEXT,
    teacher_id INTEGER,
    FOREIGN KEY(teacher_id) REFERENCES users(uid),
    FOREIGN KEY(class) REFERENCES classes(name)
);

CREATE TABLE IF NOT EXISTS classes (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT
);

CREATE TABLE IF NOT EXISTS members_of_class (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER,
    role TEXT CHECK(role IN ('student', 'teacher')) NOT NULL,
    class TEXT,
    FOREIGN KEY(user_id) REFERENCES users(uid),
    FOREIGN KEY(class) REFERENCES classes(name)
);

CREATE TABLE IF NOT EXISTS timetable (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    day TEXT,
    subject INTEGER,
    time_start TEXT,
    time_end TEXT,
    room TEXT,
    teacher INTEGER,
    class TEXT,
    FOREIGN KEY(class) REFERENCES classes(name),
    FOREIGN KEY(teacher) REFERENCES users(uid),
    FOREIGN KEY(subject) REFERENCES subjects(id)
);
