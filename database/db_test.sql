CREATE SCHEMA stu_test_2;

CREATE TABLE stu_test_2.Organization(
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) UNIQUE,
    address VARCHAR(255),
    zip_code VARCHAR(10),
    state VARCHAR(100),
    city VARCHAR(255)
);

CREATE TABLE stu_test_2.Admin_root (
    id SERIAL PRIMARY KEY,
    password_hash VARCHAR(255) NOT NULL,
    email VARCHAR(200) NOT NULL UNIQUE CHECK (email ~* '^[A-Za-z0-9._%-]+@[A-Za-z0-9.-]+[.][A-Za-z]+$'),
    fullname VARCHAR (100) DEFAULT NULL,
    last_name VARCHAR(100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    stripe_id VARCHAR(100) DEFAULT NULL,
    organization_id INT REFERENCES stu_test_2.Organization(id) ON DELETE CASCADE,
    organization_name VARCHAR(255) DEFAULT NULL
);

CREATE TABLE stu_test_2.Permissions (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    primary_role TEXT,
    secondary_role TEXT,
    description TEXT
);

CREATE TABLE stu_test_2.Admin_staff (
    id SERIAL PRIMARY KEY,
    fullname VARCHAR(255),
    organization_id INT REFERENCES stu_test_2.Organization(id) ON DELETE CASCADE,
    email VARCHAR(100) NOT NULL UNIQUE,
    job_title VARCHAR(255),
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    region VARCHAR(100) DEFAULT NULL,
    state VARCHAR(100) NOT NULL
);

CREATE TABLE stu_test_2.Admin_Permissions (
    id SERIAL PRIMARY KEY,
    admin_id INT NOT NULL,
    permission_id INT NOT NULL,
    FOREIGN KEY (admin_id) REFERENCES stu_test_2.Admin_staff(id) ON DELETE CASCADE,
    FOREIGN KEY (permission_id) REFERENCES stu_test_2.Permissions(id) ON DELETE CASCADE,
    UNIQUE (admin_id, permission_id)
);

CREATE TABLE stu_test_2.District (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    city VARCHAR(255) NOT NULL,
    state VARCHAR(100) NOT NULL,
    region VARCHAR(100) NOT NULL,
    admin_id INT REFERENCES stu_test_2.Admin_root(id) ON DELETE CASCADE,
    organization_id INT REFERENCES stu_test_2.Organization(id) ON DELETE CASCADE
);

CREATE TABLE stu_test_2.Locations (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    district_id INT REFERENCES stu_test_2.District(id) ON DELETE SET NULL,
    organization_id INT REFERENCES stu_test_2.Organization(id) ON DELETE CASCADE,
    address VARCHAR(255) NOT NULL,
    city VARCHAR(255) NOT NULL,
    state VARCHAR(100),
    zip_code VARCHAR(20),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE stu_test_2.Subjects(
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    organization_id INT REFERENCES stu_test_2.Organization(id) ON DELETE SET NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE stu_test_2.Location_subjects (
    subject_id INT REFERENCES stu_test_2.Subjects(id) ON DELETE CASCADE,
    location_id INT REFERENCES stu_test_2.Locations(id) ON DELETE CASCADE,
    PRIMARY KEY (subject_id, location_id)
);


CREATE TABLE stu_test_2.Programs (
    id SERIAL PRIMARY KEY,
    program_name VARCHAR(150) NOT NULL,
    organization_id INT REFERENCES stu_test_2.Organization(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE stu_test_2.Tutors(
    id SERIAL PRIMARY KEY,
    organization_id INT REFERENCES stu_test_2.Organization(id) ON DELETE CASCADE,
    location_id INT REFERENCES stu_test_2.Locations(id) ON DELETE SET NULL,
    email VARCHAR(150) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE stu_test_2.Location_contacts (
    location_id INT REFERENCES stu_test_2.Locations(id) ON DELETE CASCADE,
    program_id INT REFERENCES stu_test_2.Programs(id) ON DELETE CASCADE,
    description VARCHAR(255),
    first_name VARCHAR(255),
    last_name VARCHAR(255),
    email VARCHAR(255),
    phone TEXT CHECK(phone ~ '^(\+\d{1,2}\s)?\(?\d{3}\)?[\s.-]?\d{3}[\s.-]?\d{4}$')
);

CREATE TABLE stu_test_2.Location_programs (
    location_id INT REFERENCES stu_test_2.Locations(id) ON DELETE CASCADE,
    program_id INT REFERENCES stu_test_2.Programs(id) ON DELETE CASCADE,
    organization_id INT REFERENCES stu_test_2.Organization(id) ON DELETE CASCADE,
    PRIMARY KEY (location_id, program_id)
);

CREATE TABLE stu_test_2.Tutor_locations (
    tutor_id INT REFERENCES stu_test_2.Tutors(id) ON DELETE CASCADE,
    location_id INT REFERENCES stu_test_2.Locations(id) ON DELETE CASCADE,
    attendance_link VARCHAR(255) DEFAULT NULL,
    organization_id INT REFERENCES stu_test_2.Organization(id) ON DELETE CASCADE,
    PRIMARY KEY (tutor_id, location_id, organization_id)
);
-- Index location_id

CREATE TABLE stu_test_2.Tutor_Permissions (
    id SERIAL PRIMARY KEY,
    tutor_id INT NOT NULL,
    permission_id INT NOT NULL,
    FOREIGN KEY (tutor_id) REFERENCES stu_test_2.Tutors(id) ON DELETE CASCADE,
    FOREIGN KEY (permission_id) REFERENCES stu_test_2.Permissions(id) ON DELETE CASCADE,
    UNIQUE (tutor_id, permission_id)
);

CREATE TABLE stu_test_2.Materials (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    external_link TEXT,
    description VARCHAR(255),
    location_id INT REFERENCES stu_test_2.Locations(id) ON DELETE SET NULL,
    organization_id INT REFERENCES stu_test_2.Organization(id) ON DELETE CASCADE,
    program_id INT REFERENCES stu_test_2.Programs(id) ON DELETE SET NULL,
    version DECIMAL,
    pre BOOLEAN DEFAULT FALSE,
    mid BOOLEAN DEFAULT FALSE,
    post BOOLEAN DEFAULT FALSE,
    visible BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Added easy_score boolean for Question support
CREATE TABLE stu_test_2.Assessments (
    id SERIAL PRIMARY KEY,
    title VARCHAR(100) NOT NULL,
    description TEXT,
    letter VARCHAR(10) NOT NULL,
    cycle VARCHAR(100) NOT NULL,
    alpha_identifier VARCHAR(10) UNIQUE,
    external_link TEXT,
    max_score INT,
    subject_id INT REFERENCES stu_test_2.Subjects(id) ON DELETE SET NULL,
    program_id INT REFERENCES stu_test_2.Programs(id) ON DELETE SET NULL,
    material_id INT REFERENCES stu_test_2.Materials(id) ON DELETE SET NULL,
    pre BOOLEAN DEFAULT FALSE,
    mid BOOLEAN DEFAULT FALSE,
    post BOOLEAN DEFAULT FALSE,
    visible BOOLEAN DEFAULT FALSE,
    version DECIMAL,
    easy_score BOOLEAN DEFAULT FALSE,
    organization_id INT REFERENCES stu_test_2.Organization(id) ON DELETE CASCADE,
    edited_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE stu_test_2.Semester(
    id SERIAL PRIMARY KEY,
    organization_id INT REFERENCES stu_test_2.Organization(id) ON DELETE CASCADE,
    year INT,
    title VARCHAR(100),
    date_start TIMESTAMP NOT NULL,
    date_end TIMESTAMP NOT NULL,
    active BOOLEAN DEFAULT FALSE
);

CREATE TABLE stu_test_2.Semester_Location (
    id SERIAL PRIMARY KEY,
    semester_id INT,
    location_id INT,
    organization_id INT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(semester_id, location_id, organization_id),
    FOREIGN KEY (semester_id) REFERENCES stu_test_2.Semester(id) ON DELETE SET NULL,
    FOREIGN KEY (location_id) REFERENCES stu_test_2.Locations(id) ON DELETE SET NULL,
    FOREIGN KEY (organization_id) REFERENCES stu_test_2.Organization(id) ON DELETE SET NULL
);


/** PERIOD CHANGED TO VAR CHAR */
/** DIRECT-PARTNERSHIP ADDED APR 19 2025 */ 
CREATE TABLE stu_test_2.Students (
    id SERIAL PRIMARY KEY,
    location_id INT REFERENCES stu_test_2.Locations(id) ON DELETE SET NULL,
    tutor_id INT REFERENCES stu_test_2.Tutors(id) ON DELETE SET NULL,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    middle_name VARCHAR(255) DEFAULT NULL,
    semester_id INT REFERENCES stu_test_2.Semester(id) ON DELETE SET NULL,
    period VARCHAR(200), 
    email VARCHAR(200) CHECK (
        email IS NULL OR 
        email = '' OR 
        email ~* '^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$'
    ),,
    grade_level INT CHECK (grade_level BETWEEN 0 AND 12),
    active BOOLEAN DEFAULT TRUE,
    direct_partnership BOOLEAN DEFAULT FALSE,
    created_by TEXT,
    created_by_id INT,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE stu_test_2.Sessions (
    id SERIAL PRIMARY KEY,
    tutor_id INT REFERENCES stu_test_2.Tutors(id) ON DELETE SET NULL,
    session_date TIMESTAMP NOT NULL,
    location_id INT REFERENCES stu_test_2.Locations(id) ON DELETE CASCADE,
    organization_id INT REFERENCES stu_test_2.Organization(id) ON DELETE SET NULL,
    program_id INT REFERENCES stu_test_2.Programs(id) ON DELETE CASCADE,
    substitute BOOLEAN DEFAULT FALSE,
    substitute_id INT REFERENCES stu_test_2.Tutors(id) ON DELETE SET NULL,
    semester_id INT REFERENCES stu_test_2.Semester(id) ON DELETE SET NULL,
    student_count INT,
    in_school BOOLEAN DEFAULT FALSE, 
    start_time VARCHAR(10),
    duration INT,
    subject VARCHAR(100),
    subject_id INT REFERENCES stu_test_2.Subjects(id) ON DELETE SET NULL,
    notes TEXT,
    edited_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE stu_test_2.Session_students (
    id SERIAL PRIMARY KEY,
    session_id INT NOT NULL,
    student_id INT NOT NULL,
    subject_id INT DEFAULT NULL,
    absent BOOLEAN DEFAULT FALSE,
    duration INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (session_id, student_id),
    FOREIGN KEY (session_id) REFERENCES stu_test_2.Sessions(id) ON DELETE CASCADE,
    FOREIGN KEY (student_id) REFERENCES stu_test_2.Students(id) ON DELETE CASCADE
);


-- Create a start and end time?
CREATE TABLE stu_test_2.Assessments_students (
    id SERIAL PRIMARY KEY,
    session_id INT NOT NULL,
    student_id INT NOT NULL,
    score INT NOT NULL CHECK (score >= 0),
    assessment_id INT NOT NULL,
    subject_id INT REFERENCES stu_test_2.Subjects(id) ON DELETE SET NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (student_id, assessment_id, session_id),
    FOREIGN KEY (student_id) REFERENCES stu_test_2.Students(id) ON DELETE CASCADE,
    FOREIGN KEY (assessment_id) REFERENCES stu_test_2.Assessments(id) ON DELETE CASCADE,
    FOREIGN KEY (session_id) REFERENCES stu_test_2.Sessions(id) ON DELETE CASCADE
);

CREATE TABLE stu_test_2.Notifications (
    id SERIAL PRIMARY KEY,
    organization_id INT REFERENCES stu_test_2.Organization(id) ON DELETE CASCADE,
    location_id INT REFERENCES stu_test_2.Locations(id) DEFAULT NULL,
    district_id INT REFERENCES stu_test_2.District(id) ON DELETE SET NULL,
    title TEXT,
    body TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE stu_test_2.Announcements (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    body TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    location_id INT REFERENCES stu_test_2.Locations(id) ON DELETE SET NULL,
    severity VARCHAR(20),
    organization_id INT NOT NULL REFERENCES stu_test_2.Organization(id) ON DELETE CASCADE,
    program_id INT REFERENCES stu_test_2.Programs(id) ON DELETE SET NULL,
    admin_id INT REFERENCES stu_test_2.Admin_root(id) ON DELETE SET NULL,
    staff_id INT REFERENCES stu_test_2.Admin_staff(id) ON DELETE CASCADE
);

CREATE TABLE stu_test_2.User_Acknowledgments (
    id SERIAL PRIMARY KEY,
    tutor_id INT NOT NULL,
    announcement_id INT NOT NULL,
    acknowledged BOOLEAN DEFAULT FALSE,
	organization_id INT NOT NULL,
    acknowledged_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (organization_id) REFERENCES stu_test_2.Organization(id) ON DELETE SET NULL,
    FOREIGN KEY (tutor_id) REFERENCES stu_test_2.Tutors(id) ON DELETE CASCADE,
    FOREIGN KEY (announcement_id) REFERENCES stu_test_2.Announcements(id) ON DELETE CASCADE,
    UNIQUE (tutor_id, announcement_id)
);


 --- INDEXES ---
CREATE INDEX idx_tutor_schedules_tutor_id ON stu_test_2.Tutor_schedules(tutor_id);
CREATE INDEX idx_tutor_schedules_program_id ON stu_test_2.Tutor_schedules(program_id);
CREATE INDEX idx_tutor_schedules_dates ON stu_test_2.Tutor_schedules(start_date, end_date);
CREATE INDEX idx_tutor_location ON stu_test_2.Tutors(location_id);
CREATE INDEX idx_sessions_tutor ON stu_test_2.Sessions(tutor_id);
CREATE INDEX idx_student_semester ON stu_test_2.Sessions(semester_id);
CREATE INDEX idx_session_date ON stu_test_2.Sessions(session_date);
CREATE INDEX idx_session_assessments ON stu_test_2.Assessments_students(session_id);
CREATE INDEX idx_session_students_session ON stu_test_2.Session_students(session_id);
CREATE INDEX idx_session_students_student ON stu_test_2.Session_students(student_id);
CREATE INDEX idx_student_location ON stu_test_2.Students(location_id);
-- NEW SCHEMA --


-- April 29 --  -- NEW -- NEW -- NEW -- NEW
CREATE INDEX idx_student_semester_id ON stu_test_2.Students(semester_id);

CREATE TABLE stu_test_2.Assessment_sessions (
    id SERIAL PRIMARY KEY,
    tutor_id INT REFERENCES stu_test_2.Tutors(id) ON DELETE CASCADE,
    student_id INT REFERENCES stu_test_2.Students(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    first_name TEXT,
    last_name TEXT,
    assessment_id INT REFERENCES stu_test_2.Assessments(id) ON DELETE CASCADE,
    semester_id INT REFERENCES stu_test_2.Semester(id) ON DELETE CASCADE,
    session_token UUID NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    completed BOOLEAN DEFAULT FALSE,
    grade_assessment BOOLEAN DEFAULT FALSE
);

CREATE TABLE stu_test_2.Session_answers (
    id SERIAL PRIMARY KEY,
    assessment_id INT REFERENCES stu_test_2.Assessments(id) ON DELETE CASCADE,
    student_id INT REFERENCES stu_test_2.Students(id) ON DELETE CASCADE,
    session_token UUID,
    question_id INT REFERENCES stu_test_2.Questions(id) ON DELETE CASCADE,
    choice_id INT REFERENCES stu_test_2.Choices(id) ON DELETE NULL,
    answer_text TEXT,
    answered_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

ALTER TABLE stu_test_2.Assessment_sessions ADD COLUMN completed BOOLEAN DEFAULT FALSE;
ALTER TABLE stu_test_2.Assessment_sessions ADD COLUMN grade_assessment BOOLEAN DEFAULT FALSE;
ALTER TABLE stu_test_2.Assessment_sessions ADD COLUMN answer_text TEXT;
ALTER table stu_test_2.Session_answers ALTER COLUMN choice_id DROP NOT NULL;
ALTER TABLE stu_test_2.Session_answers ADD COLUMN answer_text TEXT;

CREATE TABLE stu_test_2.Locations_teachers (
    id SERIAL PRIMARY KEY,
    name TEXT,
    room TEXT,
    grade_level INT,
    substitute BOOLEAN DEFAULT FALSE,
    location_id INT REFERENCES stu_test_2.Locations(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

ALTER TABLE stu_test_2.Students ADD COLUMN teacher_id INT REFERENCES stu_test_2.Locations_teachers(id) ON DELETE SET NULL;
INSERT INTO stu_test_2.Permissions( name, description, primary_role, secondary_role) VALUES ( 'write:teacher', 'write teacher data', 'tutor', 'admin');
INSERT INTO stu_test_2.Permissions( name, description, primary_role, secondary_role) VALUES ( 'view:teacher', 'view teacher data', 'tutor', 'admin');
INSERT INTO stu_test_2.Permissions( name, description, primary_role, secondary_role) VALUES ( 'delete:teacher', 'delete big teacher', 'tutor', 'admin');

-- NEW -- NEW -- NEW -- NEW

CREATE TABLE stu_test_2.Tutor_schedules (
    id SERIAL PRIMARY KEY,
    tutor_id INTEGER NOT NULL REFERENCES stu_test_2.Tutors(id) ON DELETE CASCADE,
    program_id INTEGER NOT NULL REFERENCES stu_test_2.Programs(id),
    schedule_type VARCHAR(20) NOT NULL CHECK (schedule_type IN ('inclusion', 'exclusion')),
    start_date DATE NOT NULL,
    end_date DATE,
    recurring BOOLEAN DEFAULT FALSE,
    recurrence_pattern JSONB, -- For complex recurrence rules
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE stu_test_2.Questions (
    id SERIAL PRIMARY KEY,
    assessment_id INT REFERENCES stu_test_2.Assessments(id) ON DELETE CASCADE,
    image_url TEXT,
    question_text TEXT NOT NULL,
    question_type VARCHAR(50) NOT NULL, -- e.g., 'multiple_choice', 'true_false', 'short_answer'
    points INT DEFAULT 1,
    order_number INT, -- optional, if you want to sort questions
    is_required BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE stu_test_2.Choices (
    id SERIAL PRIMARY KEY,
    question_id INT REFERENCES stu_test_2.Questions(id) ON DELETE CASCADE,
    choice_text TEXT NOT NULL,
    is_correct BOOLEAN DEFAULT FALSE,
    order_number INT, -- useful for displaying options in order
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE stu_test_2.Assessment_answers (
    id SERIAL PRIMARY KEY,
    assessment_student_id INT NOT NULL REFERENCES stu_test_2.Assessments_students(id) ON DELETE CASCADE,
    question_id INT NOT NULL REFERENCES stu_test_2.Questions(id) ON DELETE CASCADE,
    choice_id INT REFERENCES stu_test_2.Choices(id) ON DELETE SET NULL,
    answer_text TEXT, -- for short_answer type
    is_correct BOOLEAN,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (assessment_student_id, question_id)
);

