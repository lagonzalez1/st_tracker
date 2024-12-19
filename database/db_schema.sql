CREATE SCHEMA stu_tracker;

CREATE TABLE stu_tracker.Admin_root (
    id SERIAL PRIMARY KEY,
    password_hash VARCHAR(255) NOT NULL,
    email VARCHAR(100) NOT NULL UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    stripe_id VARCHAR(100) DEFAULT NULL,
    organization_name VARCHAR(255) DEFAULT NULL
);

CREATE TABLE stu_tracker.Admin_staff (
    id SERIAL PRIMARY KEY,
    fullname VARCHAR(255),
    root_id INT REFERENCES stu_tracker.Admin_root(id) ON DELETE CASCADE,
    email VARCHAR(100) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    region VARCHAR(100) DEFAULT NULL,
    state VARCHAR(100) NOT NULL
);


CREATE TABLE stu_tracker.District (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    created_at VARCHAR(255) NOT NULL,
    city VARCHAR(255) NOT NULL,
    state VARCHAR(100) NOT NULL,
    region VARCHAR(100) NOT NULL
);

CREATE TABLE stu_tracker.Location_root (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    district_id INT REFERENCES stu_tracker.District(id) ON DELETE CASCADE,
    address VARCHAR(255) NOT NULL,
    city VARCHAR(100),
    state VARCHAR(100),
    zip_code VARCHAR(100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE stu_tracker.Program_root (
    id SERIAL PRIMARY KEY,
    program_name VARCHAR(150) NOT NULL
);

CREATE TABLE stu_tracker.Tutors (
    id SERIAL PRIMARY KEY,
    location_id INT REFERENCES stu_tracker.Location_root(id) ON DELETE SET NULL,
    program_id INT REFERENCES  stu_tracker.Program_root(id) ON DELETE SET NULL,
    admin_id int REFERENCES stu_tracker.Program_root(id) ON DELETE SET NULL,
    email VARCHAR(100) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    first_name VARCHAR(50) NOT NULL,
    last_name VARCHAR(50) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE stu.stu_tracker.Material (
    id SERIAL PRIMARY KEY
    title VARCHAR(255) NOT NULL,
    external_link VARCHAR(255),
    description VARCHAR(255),
    version VARCHAR(255),
    pre BOOLEAN DEFAULT FALSE,
    post BOOLEAN DEFAULT FALSE,
    visable BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
)

CREATE TABLE stu_tracker.Assessments (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    letter VARCHAR(255) NOT NULL,
    cycle VARCHAR(100) NOT NULL,
    step VARCHAR(100) NOT NULL,
    grade_level VARCHAR(10),
    external_link VARCHAR(255),
    questions_max INT,
    subject VARCHAR(150),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE stu_tracker.Assessments_log (
    id SERIAL PRIMARY KEY,
    assesment_id INT REFERENCES stu_tracker.Assessments(id) ON DELETE SET NULL,
    incorrect INT,
    correct INT,
    score INT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE stu_tracker.Students (
    id SERIAL PRIMARY KEY,
    location_id INT REFERENCES stu_tracker.Location_root(id) ON DELETE SET NULL,
    assesment_log INT REFERENCES stu_tracker.Assessments_log(id) ON DELETE SET NULL,
    tutor_id INT REFERENCES stu_tracker.Tutors(id) ON DELETE SET NULL,
    school_id VARCHAR(255),
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    middle_name VARCHAR(255) DEFAULT NULL,
    email VARCHAR(100),
    grade_level VARCHAR(20),
    active BOOLEAN,
    tardy BOOLEAN,
    subject VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE stu_tracker.Sessions (
    id SERIAL PRIMARY KEY,
    tutor_id INT REFERENCES stu_tracker.Tutors(id) ON DELETE CASCADE,
    student_id INT REFERENCES stu_tracker.Students(id) ON DELETE CASCADE,
    session_date TIMESTAMP NOT NULL,
    assesment_id INT REFERENCES stu_tracker.Assessments_log(id) ON DELETE SET NULL,
    start_time INT NOT NULL,
    duration_minutes INT NOT NULL,
    subject VARCHAR(100),
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

