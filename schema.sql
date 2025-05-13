-- schema.sql
-- Table storing system users (students, teachers, admins)
CREATE TABLE IF NOT EXISTS users (
    uid INTEGER PRIMARY KEY AUTOINCREMENT,
    email TEXT UNIQUE NOT NULL, -- Unique email address
    password TEXT NOT NULL, -- User password
    role TEXT NOT NULL CHECK(role IN ('student', 'teacher', 'admin')), -- User role
    UNIQUE(email)
);

-- Table storing personal information for users
CREATE TABLE IF NOT EXISTS persons (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL UNIQUE, -- Reference to user
    first_name TEXT NOT NULL, -- First name
    last_name TEXT NOT NULL, -- Last name
    birth_date TEXT CHECK(birth_date GLOB '[0-9][0-9][0-9][0-9]-[0-1][0-9]-[0-3][0-9]'), -- Birth date in YYYY-MM-DD format
    address TEXT, -- Address
    phone TEXT, -- Phone number
    FOREIGN KEY(user_id) REFERENCES users(uid)
);

-- Table storing classes (student groups)
CREATE TABLE IF NOT EXISTS classes (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE -- Unique class name (e.g., "1A", "2B")
);

-- Table storing school subjects
CREATE TABLE IF NOT EXISTS subjects (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE, -- Unique subject name (e.g., "Mathematics")
    class_name TEXT, -- Name of the class assigned to the subject
    teacher_id INTEGER, -- ID of the teacher assigned to the subject
    FOREIGN KEY(teacher_id) REFERENCES users(uid),
    FOREIGN KEY(class_name) REFERENCES classes(name)
);

-- Table storing student-subject assignments
CREATE TABLE IF NOT EXISTS students_subjects (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL, -- Student ID
    subject_id INTEGER NOT NULL, -- Subject ID
    UNIQUE(user_id, subject_id), -- Prevents duplicates
    FOREIGN KEY(user_id) REFERENCES users(uid),
    FOREIGN KEY(subject_id) REFERENCES subjects(id)
);

-- Table storing teacher-subject assignments
CREATE TABLE IF NOT EXISTS teachers_subjects (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL, -- Teacher ID
    subject_id INTEGER NOT NULL, -- Subject ID
    UNIQUE(user_id, subject_id), -- Prevents duplicates
    FOREIGN KEY(user_id) REFERENCES users(uid),
    FOREIGN KEY(subject_id) REFERENCES subjects(id)
);

-- Table storing grades, comments, and custom values
CREATE TABLE IF NOT EXISTS grades (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL, -- Student ID
    subject_id INTEGER NOT NULL, -- Subject ID
    grade TEXT NOT NULL CHECK(length(grade) <= 255), -- Numeric grade (e.g., "5", "4.5"), comment (e.g., "Missing homework"), or custom value (e.g., "Pass")
    grade_type TEXT NOT NULL CHECK(grade_type IN ('numeric', 'comment', 'custom')), -- Type of entry: numeric, comment, or custom
    date TEXT NOT NULL CHECK(date GLOB '[0-9][0-9][0-9][0-9]-[0-1][0-9]-[0-3][0-9]'), -- Date of entry in YYYY-MM-DD format
    FOREIGN KEY(user_id) REFERENCES users(uid),
    FOREIGN KEY(subject_id) REFERENCES subjects(id)
);

-- Table storing class memberships for users (students and teachers)
CREATE TABLE IF NOT EXISTS class_members (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL, -- User ID
    class_name TEXT NOT NULL, -- Class name
    UNIQUE(user_id, class_name), -- Prevents duplicates
    FOREIGN KEY(user_id) REFERENCES users(uid),
    FOREIGN KEY(class_name) REFERENCES classes(name)
);

-- Table storing the timetable
CREATE TABLE IF NOT EXISTS timetable (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    day TEXT NOT NULL CHECK(day IN ('Monday', 'Tuesday', 'Wednesday', 'Thursday', 'Friday', 'Saturday', 'Sunday')), -- Day of the week
    subject_id INTEGER NOT NULL, -- Subject ID
    time_start TEXT NOT NULL CHECK(time_start GLOB '[0-2][0-9]:[0-5][0-9]'), -- Start time in HH:MM format
    time_end TEXT NOT NULL CHECK(time_end GLOB '[0-2][0-9]:[0-5][0-9]'), -- End time in HH:MM format
    room TEXT, -- Room number or name
    teacher_id INTEGER NOT NULL, -- Teacher ID
    class_name TEXT NOT NULL, -- Class name
    FOREIGN KEY(class_name) REFERENCES classes(name),
    FOREIGN KEY(teacher_id) REFERENCES users(uid),
    FOREIGN KEY(subject_id) REFERENCES subjects(id)
);

-- Indexes for foreign keys to improve query performance
CREATE INDEX idx_persons_user_id ON persons(user_id);
CREATE INDEX idx_subjects_teacher_id ON subjects(teacher_id);
CREATE INDEX idx_subjects_class_name ON subjects(class_name);
CREATE INDEX idx_students_subjects_user_id ON students_subjects(user_id);
CREATE INDEX idx_students_subjects_subject_id ON students_subjects(subject_id);
CREATE INDEX idx_teachers_subjects_user_id ON teachers_subjects(user_id);
CREATE INDEX idx_teachers_subjects_subject_id ON teachers_subjects(subject_id);
CREATE INDEX idx_grades_user_id ON grades(user_id);
CREATE INDEX idx_grades_subject_id ON grades(subject_id);
CREATE INDEX idx_class_members_user_id ON class_members(user_id);
CREATE INDEX idx_class_members_class_name ON class_members(class_name);
CREATE INDEX idx_timetable_teacher_id ON timetable(teacher_id);
CREATE INDEX idx_timetable_class_name ON timetable(class_name);
CREATE INDEX idx_timetable_subject_id ON timetable(subject_id);