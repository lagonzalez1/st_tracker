PresentifyClone API

A Go-based backend service for session storage and assessment creation using Gorilla Mux.

⸻

Table of Contents
	1.	Features
	2.	Tech Stack
	3.	Prerequisites
	4.	Installation & Setup
	•	Clone Repository
	•	Environment Variables
	•	Running Locally
	•	Running with Docker
	5.	API Reference
	•	Base URL
	•	Authentication
	•	Endpoints
	•	Register
	•	Login
	6.	Analytics & Reporting
	7.	License

⸻

Features
	•	Multi-tenant support: Organizations, administrators, tutors, and students.
	•	Session tracking: Record start/end times, durations, and scores.
	•	Assessment management: Create and deliver assessments and live sessions.
	•	Role-based dashboards: Customized views for root users, admins, and tutors.
	•	Analytics & Export: Light analytics and downloadable XLSX reports per tutor/student.

Tech Stack
	•	Language: Go (1.23+)
	•	Router: Gorilla Mux
	•	Database: PostgreSQL
	•	Containerization: Docker & Docker Compose

Prerequisites
	•	Go 1.23+ installed on your machine
	•	Docker & Docker Compose (for containerized run)
	•	PostgreSQL instance (local or managed)

Installation & Setup

Clone Repository

git clone https://github.com/your-org/presentifyclone.git
cd presentifyclone

Environment Variables

Copy the .env.example to .env and adjust the values:

# .env
JWT_SECRET=your_jwt_secret
B_PORT=3333
POSTGRES_USER=postgres
POSTGRES_PASSWORD=pgpassword
POSTGRES_URL=localhost
POSTGRES_PORT=5433
POSTGRES_HOST=db
DB_NAME=postgres
RABBIT_HOST=rabbitmq
RABBIT_PORT=5672
RABBIT_USERNAME=guest
RABBIT_PASSWORD=guest
ORG_ADD_KEY="2025-CREATE-0"
ORG_ADD_KEY_1="2025-CREATE-1"
ORG_ADD_KEY_2="2025-CREATE-2"
ORG_ADD_KEY_3="2025-CREATE-3"
SQL_URL="HTTP://"
AWS_ACCESS_KEY_ID=AWS_SECRET_ACCESS_KEY
AWS_SECRET_ACCESS_KEY=AWS_SECRET_ACCESS_KEY
AWS_REGION=AWS_REGION

Running Locally

go mod tidy
go run app/main.go

The service will start on http://localhost:3333 by default.

Running with Docker
	1.	Ensure Docker daemon is running.
	2.	Build and start containers:



docker-compose up –build

3. The API will be available at `http://localhost:3333`.


## API Reference

### Base URL
```text
{NODE_ENV == production ? AWS_LINK : http://localhost:3333}/api

Authentication
	•	Uses JWT stored in an HTTP-only cookie (_auth).
	•	Response headers include x-access-token for token refresh.
	•	Interceptors handle 401/403 and token renewal.

Endpoints

Register

Create a new root user and organization.
	•	URL: /register
	•	Method: POST
	•	Auth required: No

Request Body:

{
  "email": "admin@example.com",
  "password": "securePassword123!",
  "organization_name": "My Company Inc",
  "address": "1234 Main St",
  "city": "Los Angeles",
  "state": "CA",
  "zip_code": "90001"
}

Responses:

Status	Description	Body
201	Created	{ "message": "User created" }
400	Validation error	{ "error": "Invalid request" }
500	Server error	{ "error": "Server error" }


⸻

Login

Authenticate a user and obtain JWT tokens.
	•	URL: /login
	•	Method: POST
	•	Auth required: No

Request Body:

{
  "email": "admin@example.com",
  "password": "SomeSupperSecretPassword22$$"
}

Responses:

Status	Description	Body
200	OK	{ "refresh_token": "...", "token": "...", "user": { ... } }
401	Unauthorized	{ "error": "Invalid credentials" }
500	Server error	{ "error": "Server error" }

Analytics & Reporting
	•	Access light analytics on root/admin dashboards.
	•	Download session and assessment data as .xlsx for further analysis.


Create student

Create a new root student.
	•	URL: /create_student
	•	Method: POST
	•	Auth required: Yes


Request Body:

{
    ID                *int64  `json:"id"`
    FirstName         string  `json:"firstname"`
    MiddleName        string  `json:"middle_name"`
    LastName          string  `json:"last_name"`
    Period            *int64  `json:"period,omitempty"`
    SemesterID        *int64  `json:"semester_id"`
    Email             *string `json:"email"`
    GradeLevel        int     `json:"grade_level"`
    Active            bool    `json:"active"`
    CreatedAt         string  `json:"created_at"`
    LocationId        *int64  `json:"location_id"`
    DirectPartnership bool    `json:"direct_partnership"`
    TutorID           *int64  `json:"tutor_id"`
    TeacherID         *int64  `json:"teacher_id"`
    CreatedBy         string  `json:"created_by"`
    Timeframe         *bool   `json:"timeframe"`
    DurationRequired  *bool   `json:"duration_required"`
    TimeframeStart    *string `json:"timeframe_start"`
    TimeframeEnd      *string `json:"timeframe_end"`
}

Responses:

Status	Description	Body
201	Created	{ "message": "User created" }
400	Validation error { "error": "Invalid request" }
500	Server error { "error": "Server error" }





License

This project is licensed under the MIT License.

⸻

Built with ❤️ 