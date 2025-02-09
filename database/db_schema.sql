CREATE SCHEMA stu_tracker;

CREATE TABLE stu_tracker.Admin_root (
    id SERIAL PRIMARY KEY,
    password_hash VARCHAR(255) NOT NULL,
    email VARCHAR(100) NOT NULL UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    stripe_id VARCHAR(100) DEFAULT NULL,
    organization_id INT REFERENCES stu_tracker.Organization(id) ON DELETE CASCADE,
    organization_name VARCHAR(255) DEFAULT NULL
);

CREATE TABLE stu_tracker.Organization (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255),
    address VARCHAR(255),
    zip_code VARCHAR(10),
    state VARCHAR(10),
    city VARCHAR(255),
);

CREATE TABLE stu_tracker.permissions (
    id SERIAL PRIMARY KEY,
    organization_id INT REFERENCES stu_tracker.Organization(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL UNIQUE,
    description TEXT
);

CREATE TABLE stu_tracker.admin_role_permissions (
    user_id INT REFERENCES stu_tracker.Admin_staff(id) NOT NULL,
    permission_id INT REFERENCES stu_tracker.permissions(id) ON DELETE CASCADE,
    PRIMARY KEY (user_id, permission_id)
);

CREATE TABLE stu_tracker.tutor_role_permissions (
    tutor_id INT REFERENCES stu_tracker.Tutors(id) NOT NULL,
    permission_id INT REFERENCES stu_tracker.permissions(id) ON DELETE CASCADE,
    PRIMARY KEY (tutor_id, permission_id)
);

CREATE TABLE stu_tracker.Admin_staff (
    id SERIAL PRIMARY KEY,
    fullname VARCHAR(255),
    organization_id INT REFERENCES stu_tracker.Organization(id) ON DELETE CASCADE
    email VARCHAR(100) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    region VARCHAR(100) DEFAULT NULL,
    state VARCHAR(100) NOT NULL
);

CREATE TABLE stu_tracker.District (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    city VARCHAR(255) NOT NULL,
    state VARCHAR(100) NOT NULL,
    region VARCHAR(100) NOT NULL,
    admin_id INT REFERENCES stu_tracker.Admin_root(id) ON DELETE CASCADE
    organization_id INT REFERENCES stu_tracker.Organization(id) ON DELETE CASCADE
);

CREATE TABLE stu_tracker.Locations (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    district_id INT REFERENCES stu_tracker.District(id) ON DELETE SET NULL,
    admin_id INT REFERENCES stu_tracker.Admin_root(id) ON DELETE SET NULL,
    organization_id INT REFERENCES stu_tracker.Organization(id) ON DELETE SET NULL,
    address VARCHAR(255) NOT NULL,
    city VARCHAR(100),
    state VARCHAR(100),
    zip_code VARCHAR(100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE stu_tracker.Subjects(
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    organization_id INT REFERENCES stu_tracker.Organization(id) ON DELETE SET NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
)

CREATE TABLE stu_tracker.Location_subjects (
    subject_id INT REFERENCES stu_tracker.Subjects(id) ON DELETE CASCADE,
    location_id INT REFERENCES stu_tracker.Locations(id) ON DELETE CASCADE,
    PRIMARY KEY (subject_id, location_id)
);

CREATE TABLE stu_tracker.Location_programs (
    location_id INT REFERENCES stu_tracker.Locations(id) ON DELETE CASCADE,
    program_id INT REFERENCES stu_tracker.Programs(id) ON DELETE CASCADE,
    organization_id INT REFERENCES stu_tracker.Organization(id) ON DELETE CASCADE,
    PRIMARY KEY (location_id, program_id)
);

CREATE TABLE stu_tracker.Location_contacts {
    location_id INT REFERENCES stu_tracker.Locations(id) ON DELETE CASCADE,
    program_id INT REFERENCES stu_tracker.Programs(id) ON DELETE CASCADE,
    description VARCHAR(255),
    first_name VARCHAR(255),
    last_name VARCHAR(255),
    email VARCHAR(255),
    phone AS TEXT CHECK(VALUE ~ '^(\+\d{1,2}\s)?\(?\d{3}\)?[\s.-]?\d{3}[\s.-]?\d{4}$'),
}

CREATE TABLE stu_tracker.Programs (
    id SERIAL PRIMARY KEY,
    program_name VARCHAR(150) NOT NULL,
    organization_id INT REFERENCES stu_tracker.Organization(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE stu_tracker.Tutor_locations (
    tutor_id INT REFERENCES stu_tracker.Tutors(id) ON DELETE CASCADE,
    location_id INT REFERENCES stu_tracker.Locations(id) ON DELETE CASCADE,
    attendance_link VARCHAR(255) DEFAULT NULL,
    organization_id INT REFERENCES stu_tracker.Organization(id) ON DELETE CASCADE,
    PRIMARY KEY (tutor_id, location_id, organization_id)
);


CREATE TABLE stu_tracker.Tutors (
    id SERIAL PRIMARY KEY,
    organization_id INT REFERENCES stu_tracker.Organization(id) ON DELETE CASCADE,
    location_id INT REFERENCES stu_tracker.Locations(id) ON DELETE SET NULL,
    email VARCHAR(100) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE stu_tracker.Materials (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    external_link TEXT,
    description VARCHAR(255),
    admin_id int REFERENCES stu_tracker.Admin_root(id) ON DELETE SET NULL,
    location_id INT REFERENCES stu_tracker.Locations(id) ON DELETE SET NULL,
    organization_id INT REFERENCES stu_tracker.Organization(id) ON DELETE CASCADE,
    version VARCHAR(255),
    pre BOOLEAN DEFAULT FALSE,
    mid BOOLEAN DEFAULT FALSE,
    post BOOLEAN DEFAULT FALSE,
    visable BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
)

CREATE TABLE stu_tracker.Assessments (
    id SERIAL PRIMARY KEY,
    title VARCHAR(100) NOT NULL,
    description TEXT,
    letter VARCHAR(10) NOT NULL,
    cycle VARCHAR(100) NOT NULL,
    alpha_identifier VARCHAR(10) UNIQUE,
    external_link TEXT,
    max_score INT,
    subject VARCHAR(150),
    program_id INT REFERENCES stu_tracker.Programs(id) ON DELETE SET NULL,
    material_id INT REFERENCES stu_tracker.Materials(id) ON DELETE SET NULL,
    organization_id INT REFERENCES stu_tracker.Organization(id) ON DELETE CASCADE,
    edited_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE stu_tracker.Semester (
    id SERIAL PRIMARY KEY,
    organization_id INT REFERENCES stu_tracker.Organization(id) ON DELETE CASCADE,
    year INT,
    title VARCHAR(100)
);


/** PERIOD CHANGED TO VAR CHAR*/
CREATE TABLE stu_tracker.Students (
    id SERIAL PRIMARY KEY,
    location_id INT REFERENCES stu_tracker.Locations(id) ON DELETE SET NULL,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    middle_name VARCHAR(255) DEFAULT NULL,
    semester_id INT REFERENCES stu_tracker.Semester(id) ON DELETE SET NULL,
    period VARCHAR(200), 
    email VARCHAR(100),
    grade_level INT,
    active BOOLEAN,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE stu_tracker.Sessions (
    id SERIAL PRIMARY KEY,
    tutor_id INT REFERENCES stu_tracker.Tutors(id) ON DELETE CASCADE,
    session_date TIMESTAMP NOT NULL,
    location_id INT REFERENCES stu_tracker.Locations(id) ON DELETE SET NULL,
    program_id INT REFERENCES stu_tracker.Programs(id) ON DELETE SET NULL,
    substitute BOOLEAN DEFAULT FALSE,
    substitute_id INT REFERENCES stu_tracker.Tutors(id) ON DELETE SET NULL,
    student_count INT,
    start_time VARCHAR(10),
    subject VARCHAR(100),
    subject_id INT REFERENCES stu_tracker.Subjects(id) ON DELETE SET NULL,
    notes TEXT,
    edited_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE stu_tracker.Session_students (
    id SERIAL PRIMARY KEY,
    session_id INT NOT NULL,
    student_id INT NOT NULL,
    duration INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (session_id, student_id),
    FOREIGN KEY (session_id) REFERENCES stu_tracker.Sessions(id) ON DELETE CASCADE,
    FOREIGN KEY (student_id) REFERENCES stu_tracker.Students(id) ON DELETE CASCADE
);

CREATE TABLE stu_tracker.Assessments_students (
    id SERIAL PRIMARY KEY,
    session_id INT NOT NULL,
    student_id INT NOT NULL,
    score INT NOT NULL,
    assessment_id INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (student_id, assessment_id),
    FOREIGN KEY (student_id) REFERENCES stu_tracker.Students(id) ON DELETE CASCADE,
    FOREIGN KEY (assessment_id) REFERENCES stu_tracker.Assessments(id) ON DELETE CASCADE,
    FOREIGN KEY (session_id) REFERENCES stu_tracker.Sessions(id) ON DELETE CASCADE
);

CREATE TABLE stu_tracker.Notifications (
    id SERIAL PRIMARY KEY,
    organization_id INT REFERENCES stu_tracker.Students(id),
    location_id INT REFERENCES stu_tracker.Locations(id) DEFAULT NULL,
    district_id INT REFERENCES stu_tracker.District(id) ON DELETE SET NULL,
    title TEXT,
    body TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
)

CREATE TABLE stu_tracker.Google_sheets(
    location_id INT REFERENCES stu_tracker.Locations(id) DEFAULT NULL,
    tutor_id INT REFERENCES stu_tracker.Tutors(id) ON DELETE CASCADE, 
    google_url TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
)

