CREATE SCHEMA stu_tracker;

CREATE TABLE stu_tracker.Permissions (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    primary_role TEXT,
    secondary_role TEXT,
    description TEXT
);

CREATE TABLE stu_tracker.Organization(
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) UNIQUE,
    address VARCHAR(255),
    zip_code VARCHAR(10),
    state VARCHAR(100),
    city VARCHAR(255)
);

CREATE TABLE stu_tracker.Admin_root (
    id SERIAL PRIMARY KEY,
    password_hash VARCHAR(255) NOT NULL,
    email VARCHAR(200) NOT NULL UNIQUE CHECK (email ~* '^[A-Za-z0-9._%-]+@[A-Za-z0-9.-]+[.][A-Za-z]+$'),
    fullname VARCHAR (100) DEFAULT NULL,
    last_name VARCHAR(100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    stripe_id VARCHAR(100) DEFAULT NULL,
    organization_id INT REFERENCES stu_tracker.Organization(id) ON DELETE CASCADE,
    organization_name VARCHAR(255) DEFAULT NULL
);

CREATE TABLE stu_tracker.Admin_staff (
    id SERIAL PRIMARY KEY,
    fullname VARCHAR(255),
    organization_id INT REFERENCES stu_tracker.Organization(id) ON DELETE CASCADE,
    email VARCHAR(100) NOT NULL UNIQUE,
    job_title VARCHAR(255),
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    region VARCHAR(100) DEFAULT NULL,
    state VARCHAR(100) NOT NULL
);

CREATE TABLE stu_tracker.Admin_Permissions (
    id SERIAL PRIMARY KEY,
    admin_id INT NOT NULL,
    permission_id INT NOT NULL,
    FOREIGN KEY (admin_id) REFERENCES stu_tracker.Admin_staff(id) ON DELETE CASCADE,
    FOREIGN KEY (permission_id) REFERENCES stu_tracker.Permissions(id) ON DELETE CASCADE,
    UNIQUE (admin_id, permission_id)
);

CREATE TABLE stu_tracker.District (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    city VARCHAR(255) NOT NULL,
    state VARCHAR(100) NOT NULL,
    region VARCHAR(100) NOT NULL,
    admin_id INT REFERENCES stu_tracker.Admin_root(id) ON DELETE CASCADE,
    organization_id INT REFERENCES stu_tracker.Organization(id) ON DELETE CASCADE
);

CREATE TABLE stu_tracker.Locations (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    district_id INT REFERENCES stu_tracker.District(id) ON DELETE SET NULL,
    organization_id INT REFERENCES stu_tracker.Organization(id) ON DELETE CASCADE,
    address VARCHAR(255) NOT NULL,
    city VARCHAR(255) NOT NULL,
    state VARCHAR(100),
    zip_code VARCHAR(20),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE stu_tracker.Subjects (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    organization_id INT REFERENCES stu_tracker.Organization(id) ON DELETE SET NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE stu_tracker.Location_subjects (
    subject_id INT REFERENCES stu_tracker.Subjects(id) ON DELETE CASCADE,
    location_id INT REFERENCES stu_tracker.Locations(id) ON DELETE CASCADE,
    PRIMARY KEY (subject_id, location_id)
);


CREATE TABLE stu_tracker.Programs (
    id SERIAL PRIMARY KEY,
    program_name VARCHAR(150) NOT NULL,
    organization_id INT REFERENCES stu_tracker.Organization(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE stu_tracker.Tutors(
    id SERIAL PRIMARY KEY,
    organization_id INT REFERENCES stu_tracker.Organization(id) ON DELETE CASCADE,
    location_id INT REFERENCES stu_tracker.Locations(id) ON DELETE SET NULL,
    email VARCHAR(150) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE stu_tracker.Location_contacts (
    location_id INT REFERENCES stu_tracker.Locations(id) ON DELETE CASCADE,
    program_id INT REFERENCES stu_tracker.Programs(id) ON DELETE CASCADE,
    description VARCHAR(255),
    first_name VARCHAR(255),
    last_name VARCHAR(255),
    email VARCHAR(255),
    phone TEXT CHECK(phone ~ '^(\+\d{1,2}\s)?\(?\d{3}\)?[\s.-]?\d{3}[\s.-]?\d{4}$')
);

CREATE TABLE stu_tracker.Location_programs (
    location_id INT REFERENCES stu_tracker.Locations(id) ON DELETE CASCADE,
    program_id INT REFERENCES stu_tracker.Programs(id) ON DELETE CASCADE,
    organization_id INT REFERENCES stu_tracker.Organization(id) ON DELETE CASCADE,
    PRIMARY KEY (location_id, program_id)
);

CREATE TABLE stu_tracker.Tutor_locations (
    tutor_id INT REFERENCES stu_tracker.Tutors(id) ON DELETE CASCADE,
    location_id INT REFERENCES stu_tracker.Locations(id) ON DELETE CASCADE,
    attendance_link VARCHAR(255) DEFAULT NULL,
    organization_id INT REFERENCES stu_tracker.Organization(id) ON DELETE CASCADE,
    PRIMARY KEY (tutor_id, location_id, organization_id)
);
-- Index location_id

CREATE TABLE stu_tracker.Tutor_Permissions (
    id SERIAL PRIMARY KEY,
    tutor_id INT NOT NULL,
    permission_id INT NOT NULL,
    FOREIGN KEY (tutor_id) REFERENCES stu_tracker.Tutors(id) ON DELETE CASCADE,
    FOREIGN KEY (permission_id) REFERENCES stu_tracker.Permissions(id) ON DELETE CASCADE,
    UNIQUE (tutor_id, permission_id)
);

CREATE TABLE stu_tracker.Materials (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    external_link TEXT,
    description VARCHAR(255),
    location_id INT REFERENCES stu_tracker.Locations(id) ON DELETE SET NULL,
    organization_id INT REFERENCES stu_tracker.Organization(id) ON DELETE CASCADE,
    program_id INT REFERENCES stu_tracker.Programs(id) ON DELETE SET NULL,
    version DECIMAL,
    pre BOOLEAN DEFAULT FALSE,
    mid BOOLEAN DEFAULT FALSE,
    post BOOLEAN DEFAULT FALSE,
    visible BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Added easy_score boolean for Question support
CREATE TABLE stu_tracker.Assessments (
    id SERIAL PRIMARY KEY,
    title VARCHAR(100) NOT NULL,
    description TEXT,
    letter VARCHAR(10) NOT NULL,
    cycle VARCHAR(100) NOT NULL,
    alpha_identifier VARCHAR(10),
    external_link TEXT,
    max_score INT,
    subject_id INT REFERENCES stu_tracker.Subjects(id) ON DELETE SET NULL,
    program_id INT REFERENCES stu_tracker.Programs(id) ON DELETE SET NULL,
    material_id INT REFERENCES stu_tracker.Materials(id) ON DELETE SET NULL,
    pre BOOLEAN DEFAULT FALSE,
    mid BOOLEAN DEFAULT FALSE,
    post BOOLEAN DEFAULT FALSE,
    visible BOOLEAN DEFAULT FALSE,
    version DECIMAL,
    easy_score BOOLEAN DEFAULT FALSE,
    organization_id INT REFERENCES stu_tracker.Organization(id) ON DELETE CASCADE,
    edited_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE stu_tracker.Semester (
    id SERIAL PRIMARY KEY,
    organization_id INT REFERENCES stu_tracker.Organization(id) ON DELETE CASCADE,
    year INT,
    title VARCHAR(100),
    date_start TIMESTAMP NOT NULL,
    date_end TIMESTAMP NOT NULL,
    active BOOLEAN DEFAULT FALSE
);

CREATE TABLE stu_tracker.Semester_Location (
    id SERIAL PRIMARY KEY,
    semester_id INT,
    location_id INT,
    organization_id INT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(semester_id, location_id, organization_id),
    FOREIGN KEY (semester_id) REFERENCES stu_tracker.Semester(id) ON DELETE SET NULL,
    FOREIGN KEY (location_id) REFERENCES stu_tracker.Locations(id) ON DELETE SET NULL,
    FOREIGN KEY (organization_id) REFERENCES stu_tracker.Organization(id) ON DELETE SET NULL
);


/** PERIOD CHANGED TO VAR CHAR */
/** DIRECT-PARTNERSHIP ADDED APR 19 2025 */ 
CREATE TABLE stu_tracker.Students (
    id SERIAL PRIMARY KEY,
    location_id INT REFERENCES stu_tracker.Locations(id) ON DELETE SET NULL,
    tutor_id INT REFERENCES stu_tracker.Tutors(id) ON DELETE SET NULL,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    middle_name VARCHAR(255) DEFAULT NULL,
    semester_id INT REFERENCES stu_tracker.Semester(id) ON DELETE SET NULL,
    period VARCHAR(200), 
    email VARCHAR(200) CHECK (
        email IS NULL OR 
        email = '' OR 
        email ~* '^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$'
    ),
    grade_level INT CHECK (grade_level BETWEEN 0 AND 12),
    active BOOLEAN DEFAULT TRUE,
    direct_partnership BOOLEAN DEFAULT FALSE,
    created_by TEXT,
    created_by_id INT,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE stu_tracker.Sessions (
    id SERIAL PRIMARY KEY,
    tutor_id INT REFERENCES stu_tracker.Tutors(id) ON DELETE SET NULL,
    session_date TIMESTAMP NOT NULL,
    location_id INT REFERENCES stu_tracker.Locations(id) ON DELETE CASCADE,
    organization_id INT REFERENCES stu_tracker.Organization(id) ON DELETE SET NULL,
    program_id INT REFERENCES stu_tracker.Programs(id) ON DELETE CASCADE,
    substitute BOOLEAN DEFAULT FALSE,
    substitute_id INT REFERENCES stu_tracker.Tutors(id) ON DELETE SET NULL,
    semester_id INT REFERENCES stu_tracker.Semester(id) ON DELETE SET NULL,
    student_count INT,
    in_school BOOLEAN DEFAULT FALSE, 
    start_time VARCHAR(10),
    duration INT,
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
    subject_id INT DEFAULT NULL,
    absent BOOLEAN DEFAULT FALSE,
    duration INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (session_id, student_id),
    FOREIGN KEY (session_id) REFERENCES stu_tracker.Sessions(id) ON DELETE CASCADE,
    FOREIGN KEY (student_id) REFERENCES stu_tracker.Students(id) ON DELETE CASCADE
);


-- Create a start and end time?
CREATE TABLE stu_tracker.Assessments_students (
    id SERIAL PRIMARY KEY,
    session_id INT NOT NULL,
    student_id INT NOT NULL,
    score FLOAT NOT NULL CHECK (score >= 0),
    assessment_id INT NOT NULL,
    subject_id INT REFERENCES stu_tracker.Subjects(id) ON DELETE SET NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (student_id, assessment_id, session_id),
    FOREIGN KEY (student_id) REFERENCES stu_tracker.Students(id) ON DELETE CASCADE,
    FOREIGN KEY (assessment_id) REFERENCES stu_tracker.Assessments(id) ON DELETE CASCADE,
    FOREIGN KEY (session_id) REFERENCES stu_tracker.Sessions(id) ON DELETE CASCADE
);

CREATE TABLE stu_tracker.Notifications (
    id SERIAL PRIMARY KEY,
    organization_id INT REFERENCES stu_tracker.Organization(id) ON DELETE CASCADE,
    location_id INT REFERENCES stu_tracker.Locations(id) DEFAULT NULL,
    district_id INT REFERENCES stu_tracker.District(id) ON DELETE SET NULL,
    title TEXT,
    body TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE stu_tracker.Announcements (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    body TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    location_id INT REFERENCES stu_tracker.Locations(id) ON DELETE SET NULL,
    severity VARCHAR(20),
    organization_id INT NOT NULL REFERENCES stu_tracker.Organization(id) ON DELETE CASCADE,
    program_id INT REFERENCES stu_tracker.Programs(id) ON DELETE SET NULL,
    admin_id INT REFERENCES stu_tracker.Admin_root(id) ON DELETE SET NULL,
    staff_id INT REFERENCES stu_tracker.Admin_staff(id) ON DELETE CASCADE
);

CREATE TABLE stu_tracker.User_Acknowledgments (
    id SERIAL PRIMARY KEY,
    tutor_id INT NOT NULL,
    announcement_id INT NOT NULL,
    acknowledged BOOLEAN DEFAULT FALSE,
	organization_id INT NOT NULL,
    acknowledged_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (organization_id) REFERENCES stu_tracker.Organization(id) ON DELETE SET NULL,
    FOREIGN KEY (tutor_id) REFERENCES stu_tracker.Tutors(id) ON DELETE CASCADE,
    FOREIGN KEY (announcement_id) REFERENCES stu_tracker.Announcements(id) ON DELETE CASCADE,
    UNIQUE (tutor_id, announcement_id)
);

CREATE TABLE stu_tracker.Tutor_schedules (
    id SERIAL PRIMARY KEY,
    tutor_id INTEGER NOT NULL REFERENCES stu_tracker.Tutors(id) ON DELETE CASCADE,
    program_id INTEGER NOT NULL REFERENCES stu_tracker.Programs(id),
    schedule_type VARCHAR(20) NOT NULL CHECK (schedule_type IN ('inclusion', 'exclusion')),
    start_date DATE NOT NULL,
    end_date DATE,
    recurring BOOLEAN DEFAULT FALSE,
    recurrence_pattern JSONB,
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE stu_tracker.Questions (
    id SERIAL PRIMARY KEY,
    assessment_id INT REFERENCES stu_tracker.Assessments(id) ON DELETE CASCADE,
    image_url TEXT,
    question_text TEXT NOT NULL,
    question_type VARCHAR(50),
    points INT DEFAULT 1,
    order_number INT,
    is_required BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE stu_tracker.Choices (
    id SERIAL PRIMARY KEY,
    question_id INT REFERENCES stu_tracker.Questions(id) ON DELETE CASCADE,
    choice_text TEXT NOT NULL,
    is_correct BOOLEAN DEFAULT FALSE,
    order_number INT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);


 --- INDEXES ---
CREATE INDEX idx_tutor_schedules_tutor_id ON stu_tracker.Tutor_schedules(tutor_id);
CREATE INDEX idx_tutor_schedules_program_id ON stu_tracker.Tutor_schedules(program_id);
CREATE INDEX idx_tutor_schedules_dates ON stu_tracker.Tutor_schedules(start_date, end_date);
CREATE INDEX idx_tutor_location ON stu_tracker.Tutors(location_id);

CREATE INDEX idx_sessions_tutor ON stu_tracker.Sessions(tutor_id);
CREATE INDEX idx_student_semester ON stu_tracker.Sessions(semester_id);
CREATE INDEX idx_session_date ON stu_tracker.Sessions(session_date);

CREATE INDEX idx_session_assessments ON stu_tracker.Assessments_students(session_id);

CREATE INDEX idx_session_students_session ON stu_tracker.Session_students(session_id);
CREATE INDEX idx_session_students_student ON stu_tracker.Session_students(student_id);

CREATE INDEX idx_student_location ON stu_tracker.Students(location_id);
-- NEW SCHEMA --


-- April 29 --  -- NEW -- NEW -- NEW -- NEW
CREATE INDEX idx_student_semester_id ON stu_tracker.Students(semester_id);

CREATE TABLE stu_tracker.Assessment_sessions (
    id SERIAL PRIMARY KEY,
    tutor_id INT REFERENCES stu_tracker.Tutors(id) ON DELETE SET NULL,
    student_id INT REFERENCES stu_tracker.Students(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    first_name TEXT,
    last_name TEXT,
    assessment_id INT REFERENCES stu_tracker.Assessments(id) ON DELETE CASCADE,
    semester_id INT REFERENCES stu_tracker.Semester(id) ON DELETE CASCADE,
    session_token UUID NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    completed BOOLEAN DEFAULT FALSE,
    grade_assessment BOOLEAN DEFAULT FALSE
);


CREATE TABLE stu_tracker.Assessment_answers (
    id SERIAL PRIMARY KEY,
    assessment_student_id INT NOT NULL REFERENCES stu_tracker.Assessments_students(id) ON DELETE CASCADE,
    question_id INT NOT NULL REFERENCES stu_tracker.Questions(id) ON DELETE CASCADE,
    choice_id INT REFERENCES stu_tracker.Choices(id) ON DELETE SET NULL,
    answer_text TEXT,
    is_correct BOOLEAN,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE stu_tracker.Session_answers (
    id SERIAL PRIMARY KEY,
    assessment_id INT REFERENCES stu_tracker.Assessments(id) ON DELETE CASCADE,
    student_id INT REFERENCES stu_tracker.Students(id) ON DELETE CASCADE,
    session_token UUID,
    question_id INT REFERENCES stu_tracker.Questions(id) ON DELETE CASCADE,
    choice_id INT REFERENCES stu_tracker.Choices(id) ON DELETE SET NULL,
    answer_text TEXT,
    answered_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE stu_tracker.Locations_teachers (
    id SERIAL PRIMARY KEY,
    name TEXT,
    room TEXT,
    grade_level INT,
    substitute BOOLEAN DEFAULT FALSE,
    location_id INT REFERENCES stu_tracker.Locations(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);


-- NEW -- NEW -- NEW -- NEW
ALTER TABLE stu_tracker.Students ADD COLUMN teacher_id INT REFERENCES stu_tracker.Locations_teachers(id) ON DELETE SET NULL;
ALTER TABLE stu_tracker.Questions ADD COLUMN standard_text TEXT;

-- May 26 -- Adjust of timeframes ETO ...
ALTER TABLE stu_tracker.Programs ADD timeframe_required BOOLEAN DEFAULT FALSE;
ALTER TABLE stu_tracker.Students ADD timeframe BOOLEAN DEFAULT FALSE;
ALTER TABLE stu_tracker.Students ADD timeframe_start TEXT DEFAULT NULL;
ALTER TABLE stu_tracker.Students ADD timeframe_end TEXT DEFAULT NULL;


ALTER TABLE stu_tracker.Session_students ADD timeframe BOOLEAN DEFAULT FALSE;
ALTER TABLE stu_tracker.Session_students ADD timeframe_start TEXT DEFAULT NULL;
ALTER TABLE stu_tracker.Session_students ADD timeframe_end TEXT DEFAULT NULL;

-- May 30 --- 

ALTER TABLE stu_tracker.Assessments_students ALTER COLUMN score TYPE FLOAT;
-- ALTER TABLE stu_tracker.Assessment_answers DROP CONSTRAINT assessment_answers_assessment_student_id_question_id_key;

ALTER TABLE stu_tracker.Materials ADD s3_reference TEXT DEFAULT NULL;

-- JULY 8 --
ALTER TABLE stu_tracker.Students ADD duration_required BOOLEAN DEFAULT FALSE;
ALTER TABLE stu_tracker.Students ADD tardy BOOLEAN DEFAULT FALSE;
ALTER TABLE stu_tracker.Assessments ADD grade_level INT;
ALTER TABLE stu_tracker.Assessments ADD questionnaire BOOLEAN DEFAULT TRUE;


-- Does not exist ALTER TABLE stu_tracker.Assessments drop constraint assessments_students_assessment_id_fkey;


CREATE TABLE stu_tracker.Generate_questions_task (
    input_key UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    status TEXT,
    retry_count INT DEFAULT 0,
    s3_output_key TEXT,
    organization_id INT REFERENCES stu_tracker.Organization(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE stu_tracker.Pre_assessment_questionnaire (
    id SERIAL PRIMARY KEY,
    student_id INT REFERENCES stu_tracker.Students(id) ON DELETE CASCADE,
    assessment_id INT REFERENCES stu_tracker.Assessments(id) ON DELETE CASCADE,
    session_token UUID NOT NULL,
    sleep_hours FLOAT DEFAULT 0,
    study_hours FLOAT DEFAULT 0,
    effort_score smallint NOT NULL CHECK (effort_score BETWEEN 0 AND 10),
    tutor_sessions INT DEFAULT 0,
    parental_help smallint NOT NULL CHECK (parental_help BETWEEN 0 AND 3),
    sports_hours INT DEFAULT 0,
    peer_influence smallint NOT NULL CHECK (peer_influence BETWEEN 0 AND 3),  
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);




ALTER TABLE stu_tracker.Students
ADD COLUMN gender VARCHAR(50) CHECK (
    gender IS NULL OR gender IN ('Male', 'Female', 'NA')
);

ALTER TABLE stu_tracker.Students
ADD race VARCHAR(100) CHECK (
    race IS NULL OR race IN ('American Indian or Alaska Native', 'Asian', 
    'Black or African American','Hispanic or Latino', 'Native Hawaiian or Other Pacific Islander',
    'White', 'Two or More Races', 'Other', 'Prefer not to say')
);


CREATE TABLE stu_tracker.Student_report (
    id SERIAL PRIMARY KEY,
    input_key UUID DEFAULT gen_random_uuid(),
    student_id INT REFERENCES stu_tracker.Students(id) ON DELETE CASCADE,
    semester_id INT REFERENCES stu_tracker.Semester(id) ON DELETE SET NULL,
    s3_output_key TEXT,
    retry_count INT DEFAULT 0,
    status VARCHAR(100) CHECK (status IS NULL OR status IN ('DONE', 'PENDING', 'ERROR', 'STARTED', 'RETRY')),
    organization_id INT REFERENCES stu_tracker.Organization(id) ON DELETE CASCADE,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE stu_tracker.Organization_report (
    id SERIAL PRIMARY KEY,
    input_key UUID DEFAULT gen_random_uuid(),
    organization_id INT REFERENCES stu_tracker.Organization(id) ON DELETE CASCADE,
    sorted_by TEXT,
	entity varchar(50),
    s3_output_key TEXT,
    retry_count INT DEFAULT 0,
    status VARCHAR(100) CHECK (status IS NULL OR status IN ('DONE', 'PENDING', 'ERROR', 'STARTED', 'RETRY')),
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE stu_tracker.Generate_materials_task (
    input_key UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    status VARCHAR(20) CHECK (status IS NULL OR status IN ('DONE', 'PENDING', 'ERROR', 'STARTED', 'RETRY')),
    retry_count INT DEFAULT 0,
    s3_output_key TEXT,
    assessment_id INT  REFERENCES stu_tracker.Assessments(id) ON DELETE CASCADE,
    organization_id INT REFERENCES stu_tracker.Organization(id) ON DELETE CASCADE,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

ALTER TABLE stu_tracker.Assessments_students ADD questionnaire_id INT REFERENCES stu_tracker.Pre_assessment_questionnaire(id) ON DELETE SET NULL;


ALTER TABLE stu_tracker.Generate_materials_task ADD input_tokens INT DEFAULT 0;
ALTER TABLE stu_tracker.Generate_materials_task ADD output_tokens INT DEFAULT 0; 


ALTER TABLE stu_tracker.Generate_questions_task ADD input_tokens INT DEFAULT 0;
ALTER TABLE stu_tracker.Generate_questions_task ADD output_tokens INT DEFAULT 0; 


ALTER TABLE stu_tracker.Admin_staff ADD district_id INT REFERENCES stu_tracker.District(id) ON DELETE SET NULL;
ALTER TABLE stu_tracker.Tutor_schedules ADD COLUMN workweek TEXT[];
ALTER TABLE stu_tracker.Admin_staff ADD active BOOLEAN DEFAULT FALSE;

ALTER TABLE stu_tracker.Semester ADD archive BOOLEAN DEFAULT FALSE;
ALTER TABLE stu_tracker.Locations ADD archive BOOLEAN DEFAULT FALSE;


CREATE TABLE stu_tracker.Surveys (
    id SERIAL PRIMARY KEY,
    organization_id INT REFERENCES stu_tracker.Organization(id) ON DELETE SET NULL,
    title TEXT,
    description TEXT,
    order_by INT DEFAULT 1,
    is_active BOOLEAN NOT NULL DEFAULT FALSE,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE stu_tracker.survey_questions (
  id               SERIAL PRIMARY KEY,
  survey_id        INT NOT NULL REFERENCES stu_tracker.surveys(id) ON DELETE CASCADE,
  order_index      INT NOT NULL,                            -- for dynamic ordering
  question_text    TEXT NOT NULL,
  question_type    TEXT NOT NULL CHECK (question_type IN ('short_text','long_text','single_choice','multi_choice','number','yes_no')),
  created_at       TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at       TIMESTAMPTZ NOT NULL DEFAULT now(),
  UNIQUE (survey_id, order_index)                           
);

CREATE TABLE stu_tracker.Survey_choice (
    id SERIAL PRIMARY KEY,
    question_survey_id INT NOT NULL REFERENCES stu_tracker.survey_questions(id) ON DELETE CASCADE,
    choice_text TEXT,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE stu_tracker.Program_survey (
    id SERIAL PRIMARY KEY,
    survey_id INT NOT NULL REFERENCES stu_tracker.surveys(id) ON DELETE CASCADE,
    program_id INT REFERENCES stu_tracker.Programs(id) ON DELETE CASCADE,
    organization_id INT REFERENCES stu_tracker.Organization(id) ON DELETE CASCADE
);

ALTER TABLE stu_tracker.Programs ADD survey_required BOOLEAN DEFAULT FALSE;

CREATE TABLE stu_tracker.Survey_response (
    id SERIAL PRIMARY KEY,
    organization_id INT REFERENCES stu_tracker.Organization(id) ON DELETE SET NULL,
    session_id INT REFERENCES stu_tracker.Sessions(id) ON DELETE CASCADE,
    question_survey_id INT NOT NULL REFERENCES stu_tracker.survey_questions(id) ON DELETE SET NULL,
    response_text TEXT,
    response_choice INT REFERENCES stu_tracker.Survey_choice(id) DEFAULT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

ALTER TABLE stu_tracker.Survey_response add question_text TEXT;
ALTER TABLE stu_tracker.Survey_response add sentiment_score REAL DEFAULT NULL;
ALTER TABLE stu_tracker.Survey_response add sentiment_label SMALLINT DEFAULT NULL;


-- Maybe
ALTER TABLE stu_tracker.Programs ADD start_time varchar(8) DEFAULT NULL;
ALTER TABLE stu_tracker.Programs ADD end_time varchar(8) DEFAULT NULL;
ALTER TABLE stu_tracker.Location_contacts ADD notes VARCHAR(255);
ALTER TABLE stu_tracker.Location_contacts ADD organization_id INT REFERENCES stu_tracker.Organization(id) ON DELETE SET NULL;
ALTER TABLE stu_tracker.Location_contacts ADD id SERIAL PRIMARY KEY;
ALTER table stu_tracker.Surveys ADD internal BOOLEAN DEFAULT TRUE;

-- What you sell
CREATE TABLE stu_tracker.subscription_plan (
  id SERIAL PRIMARY KEY,
  code TEXT UNIQUE NOT NULL,            -- 'free', 'pro', 'business'
  name TEXT NOT NULL,                   -- "Pro"
  stripe_price_id TEXT UNIQUE,          -- map to Stripe price (optional if you sell per-seat or multiple prices)
  is_active BOOLEAN DEFAULT TRUE
);

-- The limits/features for each plan
CREATE TABLE stu_tracker.plan_entitlement (
  id SERIAL PRIMARY KEY,
  plan_id INT NOT NULL REFERENCES stu_tracker.subscription_plan(id) ON DELETE CASCADE,
  key TEXT NOT NULL,                    -- 'max_locations', 'max_users', 'feature_export_csv', 'requests_per_day'
  limit_value BIGINT,                   -- null if it's a flag feature, otherwise a numeric limit
  enabled BOOLEAN DEFAULT TRUE,
  enterprise BOOLEAN DEFAULT FALSE,
  UNIQUE(plan_id, key)
);


ALTER TABLE stu_tracker.subscription_plan ADD cost_monthly FLOAT;
ALTER TABLE stu_tracker.subscription_plan ADD cost_yearly FLOAT;
ALTER TABLE stu_tracker.Admin_root ADD stripe_session_id TEXT;
ALTER TABLE stu_tracker.Admin_root ADD stripe_customer_id TEXT;


CREATE TABLE stu_tracker.admin_location_access (
    admin_id INT NOT NULL REFERENCES stu_tracker.admin_staff(id) ON DELETE CASCADE,
    location_id INT NOT NULL REFERENCES stu_tracker.locations(id) ON DELETE CASCADE,
    organization_id INT NOT NULL REFERENCES stu_tracker.organization(id) ON DELETE CASCADE,
    PRIMARY KEY (admin_id, location_id)
);

ALTER TABLE stu_tracker.Semester ADD workweek TEXT[];

CREATE TABLE stu_tracker.Semester_schedule (
    id SERIAL PRIMARY KEY,
    semester_id INT NOT NULL REFERENCES stu_tracker.Semester(id),
    schedule_type VARCHAR(20) NOT NULL CHECK (schedule_type IN ('inclusion', 'exclusion')),
    start_date DATE NOT NULL,
    end_date DATE,
    recurring BOOLEAN DEFAULT FALSE,
    recurrence_pattern JSONB,
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE stu_tracker.Student_groups (
    id SERIAL PRIMARY KEY,
    tutor_id INT NOT NULL REFERENCES stu_tracker.Tutors(id) ON DELETE CASCADE,
    location_id INT NOT NULL REFERENCES stu_tracker.Locations(id) ON DELETE CASCADE,
    title TEXT NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE stu_tracker.Student_group_attendees(
    student_id INT REFERENCES stu_tracker.Students(id) ON DELETE CASCADE,
    student_group_id INT REFERENCES stu_tracker.Student_groups(id) ON DELETE CASCADE,
    create_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

ALTER TABLE stu_tracker.Student_groups ADD semester_id INT REFERENCES stu_tracker.Semester(id) ON DELETE CASCADE;
ALTER TABLE stu_tracker.Student_group_attendees ADD CONSTRAINT pk_student_group UNIQUE (student_id, student_group_id);

-- NEEDE TO ADD THIS
CREATE TABLE stu_tracker.organization_subscription (
    id SERIAL PRIMARY KEY,
    organization_id INT REFERENCES stu_tracker.Organization(id) ON DELETE SET NULL,
    plan_id INT REFERENCES stu_tracker.subscription_plan(id) ON DELETE CASCADE,
    status TEXT,
    current_period_start TIMESTAMP DEFAULT NULL,
    current_period_end TIMESTAMP DEFAULT NULL,
    stripe_customer_id TEXT,
    stripe_subscription_id TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    stripe_session_id TEXT,
    event_id TEXT,
    lastest_invoice_url TEXT,
    subscription_status TEXT,
    canceled_at TIMESTAMP DEFAULT NULL
);

ALTER TABLE stu_tracker.Assessment_answers ADD feedback TEXT;


CREATE TABLE stu_tracker.Assessment_grader_task(
    id BIGSERIAL PRIMARY KEY,
    session_token TEXT UNIQUE,
    model_id TEXT,
    status TEXT NOT NULL DEFAULT 'PENDING',
    attempts INT NOT NULL DEFAULT 0,
    max_attempts INT NOT NULL DEFAULT 5,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE stu_tracker.Grader_task_item (
    id BIGSERIAL PRIMARY KEY,
    item_key INT NOT NULL, -- session_answers_id
    task_id BIGINT REFERENCES stu_tracker.assessment_grader_task(ID) ON DELETE CASCADE,
    status TEXT NOT NULL DEFAULT 'PENDING',
    attempts INT NOT NULL DEFAULT 0,
    last_error TEXT,
    idempotency_key TEXT, -- MODEL_ID + ITEM_KEYS
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(task_id, item_key)
);

ALTER TABLE stu_tracker.Assessment_answers ADD points FLOAT;
ALTER TABLE stu_tracker.Sessions ADD session_token UUID DEFAULT NULL;

CREATE UNIQUE INDEX IF NOT EXISTS ux_answers_student_question ON stu_tracker.Assessment_answers(assessment_student_id, question_id, choice_id);

CREATE UNIQUE INDEX IF NOT EXISTS ux_grader_task_session_model ON stu_tracker.Assessment_grader_task(session_token, model_id);


CREATE TABLE stu_tracker.LLM_usage (
    id BIGSERIAL PRIMARY KEY,
    organization_id INT REFERENCES stu_tracker.Organization(id) ON DELETE CASCADE,
    input_tokens INTEGER NOT NULL DEFAULT 0,
    output_tokens INTEGER NOT NULL DEFAULT 0,
    model TEXT NOT NULL,
    provider TEXT NOT NULL,
    status TEXT NOT NULL,
    error_type TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

ALTER TABLE stu_tracker.Tutor_schedules ADD COLUMN start_time TIME, ADD COLUMN end_time TIME;
ALTER TABLE stu_tracker.Tutor_schedules ADD COLUMN location_id INT REFERENCES stu_tracker.Locations(id) ON DELETE CASCADE;


ALTER TABLE stu_tracker.Assessment_sessions ADD COLUMN join_code VARCHAR(8) NOT NULL;

CREATE UNIQUE INDEX ON stu_tracker.Assessment_sessions (join_code);