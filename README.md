PresentifyClone API


A Go-based backend service for session storage and assessment creation using Gorilla Mux.

⸻

Table of Contents
1.	Features
2.	Tech Stack
3.	Prerequisites
4.	Installation & Setup
*	Clone Repository
*	Environment Variables
*	Running Locally
*	Running with Docker
5.	API Reference
*	Base URL
*	Authentication
*	Endpoints
*	Register
*	Login
6.	Upcoming Features *
7.	License

⸻

## 1
Features
*	Multi-tenant support: Organizations, administrators, tutors, and students.
*	Session tracking: Record start/end times, durations, and scores.
*	Assessment management: Create and deliver assessments and live sessions.
*	Role-based dashboards: Customized views for root users, admins, and tutors.
*	Analytics & Export: Light analytics and downloadable XLSX reports per tutor/student.


## 2
Tech Stack
*	Language: Go (1.23+)
*	Router: Gorilla Mux
*	Database: PostgreSQL
*	Containerization: Docker & Docker Compose

## 3
Prerequisites
*	Go 1.23+ installed on your machine
*	Docker & Docker Compose (for containerized run)
*	PostgreSQL instance (local or managed)


- - - -
Frontend application is avilable on
- https://presentifyclone.click
- - - -

## 4
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






## 5
## API Reference

### Base URL
{NODE_ENV == production ? AWS_LINK : http://localhost:3333}/api

Authentication
	•	Uses JWT stored in an HTTP-only cookie (_auth).
	•	Response headers include x-access-token for token refresh.
	•	Interceptors handle 401/403 and token renewal.

## Endpoints



### Register

Create a new root user and organization.
	•	URL: /register
	•	Method: POST
	•	Auth required: No

Request Body:
```
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
```

⸻

### Login

Authenticate a user and obtain JWT tokens.
	•	URL: /login
	•	Method: POST
	•	Auth required: No

Request Body:
```
{
  "email": "admin@example.com",
  "password": "SomeSupperSecretPassword22$$"
}

Responses:

JWT encapsulates email, organizationid .. to authorize all request.

Status	Description	Body
200	OK	{ "refresh_token": "...", "token": "...", "user": { ... } }
401	Unauthorized	{ "error": "Invalid credentials" }
500	Server error	{ "error": "Server error" }
```

### Create student

Create a new root student.
	•	URL: /create_student
	•	Method: POST
	•	Auth required: Yes

Request Body:
```
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
```
### Create location

Create a new location in reference to a school location.
	•	URL: /create_location
	•	Method: POST
	•	Auth required: Yes
```
Request Body:
{
    ID             int64  `json:"id"`
    Name           string `json:"name"`
    Address        string `json:"address"`
    DistrictId     int64  `json:"district_id"`
    City           string `json:"city"`
    State          string `json:"state"`
    ZipCode        string `json:"zip_code"`
    CreatedAt      string `json:"created_at"`
    OrganizationId *int64 `json:"organization_id"`
}

Responses:

Status	Description	Body
201	Created	{ "message": "Location created" }
400	Validation error { "error": "Invalid request" }
500	Server error { "error": "Server error" }
```

### Create admin

Create a new administrator 
	•	URL: /create_admin_staff
	•	Method: POST
	•	Auth required: Yes

Request Body:
```
{
    ID             *int64 `json:"id"`
    Fullname       string `json:"fullname"`
    Email          string `json:"email"`
    Region         string `json:"region"`
    State          string `json:"state"`
    Password       string `json:"password_hash"`
    RootId         int64  `json:"root_id"`
    OrganizationId *int64 `json:"organization_id"`
}

Responses:

Status	Description	Body
201	Created	{ "message": "Admin created" }
400	Validation error { "error": "Invalid request" }
500	Server error { "error": "Server error" }
```

### Create district

Create a new region/ district, as school locations are typically sorted by region, district. 
	•	URL: /create_district
	•	Method: POST
	•	Auth required: Yes

Request Body:
```
{
    ID             *int64 `json:"id"`
    Name           string `json:"name"`
    City           string `json:"city"`
    Region         string `json:"region"`
    State          string `json:"state"`
    AdminId        int64  `json:"admin_id"`
    OrganizationId *int64 `json:"organization_id"`
}

Responses:

Status	Description	Body
201	Created	{ "message": "District created" }
400	Validation error { "error": "Invalid request" }
500	Server error { "error": "Server error" }
```

### Create program

Create a new program, typically used to containerize requirements.
	•	URL: /create_program
	•	Method: POST
	•	Auth required: Yes

Request Body:
```
{
    ID                *int64 `json:"id"`
    ProgramName       string `json:"program_name"`
    OrganizationId    *int64 `json:"organization_id"`
    TimeFrameRequired bool   `json:"timeframe_required"`
}

Responses:

Status	Description	Body
201	Created	{ "message": "Program created" }
400	Validation error { "error": "Invalid request" }
500	Server error { "error": "Server error" }
```

### Create material

Create materials, typically used together with assessments, or standalone reference materials.

	•	URL: /create_material
	•	Method: POST
	•	Headers: "Content-Type": "multipart/form-data" 
	•	Auth required: Yes
	•	MultipartForm: Yes


Request Body:

```
{
	file: Binary (Accepts file type .jpg, .png, jpeg Images)
	data: {
		ID               *int64  `json:"id"`
		Title            string  `json:"title"`
		ExternalLink     string  `json:"external_link"`
		Description      string  `json:"description"`
		Version          float64 `json:"version"`
		Pre              bool    `json:"pre"`
		Mid              bool    `json:"mid"`
		Post             bool    `json:"post"`
		Visible          bool    `json:"visible"`
		CreatedAt        string  `json:"created_at"`
		LocationId       *int64  `json:"location_id"`
		ProgramId        *int64  `json:"program_id"`
		SReference       *string `json:"s3_reference"`
		OrganizationId   *int64  `json:"organization_id"`
		SReferenceDelete bool    `json:"s3_reference_remove"`
	}	
}

Responses:

Status	Description	Body
201	Created	{ "message": "Material created" }
400	Validation error { "error": "Invalid request" }
500	Server error { "error": "Server error" }
```
### Create tutor

Create a new tutor/support specialist. 

	•	URL: /create_tutor
	•	Method: POST
	•	Auth required: Yes

Request Body:
```
{
    ID             int64  `json:"id"`
    FirstName      string `json:"first_name"`
    LastName       string `json:"last_name"`
    Password       string `json:"password"`
    Email          string `json:"email"`
    EmailChange    string `json:"email_change"`
    CreatedAt      string `json:"created_at"`
    LocationId     *int64 `json:"location_id"`
    OrganizationId *int64 `json:"organization_id"`
}

Responses:

Status	Description	Body
201	Created	{ "message": "Tutor created" }
400	Validation error { "error": "Invalid request" }
500	Server error { "error": "Server error" }
```

### Create semester

Create a new semester, used to track and sort by date. 

	•	URL: /create_semester
	•	Method: POST
	•	Auth required: Yes

Request Body:
```
{
    ID             *int64 `json:"id"`
    Title          string `json:"title"`
    Year           *int64 `json:"year"`
    OrganizationId *int64 `json:"organization_id"`
    DateStart      string `json:"date_start"`
    DateEnd        string `json:"date_end"`
    Active         bool   `json:"active"`
}

Responses:

Status	Description	Body
201	Created	{ "message": "Semester created" }
400	Validation error { "error": "Invalid request" }
500	Server error { "error": "Server error" }
```


### Create semester location

Create a semester location joins a location to a semester time range

	•	URL: /create_semester_location
	•	Method: POST
	•	Auth required: Yes

Request Body:
```
{
    ID             *int64 `json:"id"`
    LocationID     *int64 `json:"location_id"`
    OrganizationId *int64 `json:"organization_id"`
    SemesterID     *int64 `json:"semester_id"`
}

Responses:

Status	Description	Body
201	Created	{ "message": "semester location created" }
400	Validation error { "error": "Invalid request" }
500	Server error { "error": "Server error" }
```


### Create semester location

Create a semester location joins a location to a semester time range

	•	URL: /create_semester_location
	•	Method: POST
	•	Auth required: Yes

Request Body:
```
{
    ID             *int64 `json:"id"`
    LocationID     *int64 `json:"location_id"`
    OrganizationId *int64 `json:"organization_id"`
    SemesterID     *int64 `json:"semester_id"`
}

Responses:

Status	Description	Body
201	Created	{ "message": "semester location created" }
400	Validation error { "error": "Invalid request" }
500	Server error { "error": "Server error" }
```


### Create program location

Bind a program to a location

	•	URL: /create_program_location
	•	Method: POST
	•	Auth required: Yes

Request Body:
```
{
    ProgramId      *int64 `json:"program_id"`
    LocationId     *int64 `json:"location_id"`
    OrganizationID *int64 `json:"organization_id"`
}

Responses:

Status	Description	Body
201	Created	{ "message": "Program location created" }
400	Validation error { "error": "Invalid request" }
500	Server error { "error": "Server error" }

```
### Create s3 object

* In development/ experimental
Used to create images with no association returns a ID to used to link association to tables including, materails, assessments, questions.. etc.

	•	URL: /create_s3_object
	•	Method: POST
	•	Headers: "Content-Type": "multipart/form-data" 
	•	Auth required: Yes
	•	MultiformData: Yes 

Request Body:
```
{
    file: Binary (Accepts file type .jpg, .png, jpeg Images)
}

Responses:

Status	Description	Body
201	Created	{ "message": "UUID" }
400	Validation error { "error": "Invalid request" }
500	Server error { "error": "Server error" }
```

### Create teacher

Create a teacher and link to location by the required location_id paramater.
	•	URL: /create_teacher
	•	Method: POST
	•	Auth required: Yes

Request Body:
```
{
    ID         *int64 `json:"id"`
    Name       string `json:"name"`
    Room       string `json:"room"`
    GradeLevel int64  `json:"grade_level"`
    Substitute bool   `json:"substitute"`
    LocationID *int64 `json:"location_id"`
}

Responses:

Status	Description	Body
201	Created	{ "message": "Teacher created" }
400	Validation error { "error": "Invalid request" }
500	Server error { "error": "Server error" }
```


### Create session
	Create a session with student and tutors. Each session can also include some assessments and or assessment sessions.
	•	URL: /create_session
	•	Method: POST
	•	Auth required: Yes

Request Body:
```
{
    session: {
		ID           *int64 `json:"id"`
		StudentCount *int64 `json:"student_count"`
		LocationId   *int64 `json:"location_id"`
		SubstituteId *int64 `json:"substitute_id"`
		ProgramId    *int64 `json:"program_id"`
		SemesterId   *int64 `json:"semester_id"`
		SubjectId    *int64 `json:"subject_id"`
		Notes        string `json:"notes"`
		SessionDate  string `json:"session_date"`
		StartTime    string `json:"start_time"`
		Substitute   bool   `json:"substitute"`
		TutorId      *int64 `json:"tutor_id"`
		CreatedAt    string `json:"created_at"`
		EditedAt     string `json:"edited_at"`
		Duration     *int   `json:"duration"`
		InSchool     bool   `json:"in_school"`
	}
    student_sessions: {
		ID              *int64  `json:"id"`
		Absent          bool    `json:"absent"`
		FirstName       string  `json:"first_name"`
		LastName        string  `json:"last_name"`
		SessionDate     string  `json:"session_date"`
		Duration        *int64  `json:"duration"`
		StartTime       string  `json:"start_time"`
		Notes           string  `json:"notes"`
		OrganizationId  *int64  `json:"organization_id"`
		ProgramId       *int64  `json:"program_id"`
		LocationId      *int64  `json:"location_id"`
		TutorId         *int64  `json:"tutor_id"`
		SubjectId       *int64  `json:"subject_id"`
		AssessmentId    *int64  `json:"assessment_id"`
		AssessmentScore *int64  `json:"score"`
		EasyScoreID     bool    `json:"easy_score"`
		Timeframe       *bool   `json:"timeframe"`
		TimeframeStart  *string `json:"timeframe_start"`
		TimeframeEnd    *string `json:"timeframe_end"`
	}
    assessments: {
		AssessmentID *int64                 `json:"assessment_id"`
		Choices      map[string]interface{} `json:"choices,omitempty"`
		Grader       map[string]bool        `json:"grader,omitempty"`
	}
    SessionToken   *string                       `json:"session_token"`
    OrganizationID *int64                        `json:"organization_id"`
}

Responses:

Status	Description	Body
201	Created	{ "message": "Session created" }
400	Validation error { "error": "Invalid request" }
500	Server error { "error": "Server error" }
```

### Create assessment

Create an assessment
In development/experimental
	•	URL: /create_assessment
	•	Method: POST
	•	Auth required: Yes

Request Body:
```
{
    ID              *int64                `json:"id"`
    Title           string                `json:"title"`
    Description     string                `json:"description,omitempty"`
    Letter          string                `json:"letter"`
    Cycle           *int64                `json:"cycle"`
    Visible         bool                  `json:"visible"`
    AlphaIdentifier string                `json:"alpha_identifier,omitempty"`
    ExternalLink    string                `json:"external_link,omitempty"`
    MaxScore        *int64                `json:"max_score,omitempty"`
    SubjectId       *int64                `json:"subject_id,omitempty"`
    OrganizationID  *int64                `json:"organization_id"`
    MaterialID      *int                  `json:"material_id,omitempty"`
    ProgramId       *int64                `json:"program_id"`
    GradeLevel      *int                  `json:"grade_level"`
    CreatedAt       string                `json:"created_at"`
    Version         float64               `json:"version"`
    Pre             bool                  `json:"pre"`
    Mid             bool                  `json:"mid"`
    Post            bool                  `json:"post"`
    EasyScore       bool                  `json:"easy_score"`
    questions: {
		QuestionID   *int64              `json:"question_id"`
		ImageURL     string              `json:"image_url"`
		Required     bool                `json:"is_required"`
		OrderNumber  int                 `json:"order_number"`
		Standard     *string             `json:"standard_text"`
		Points       int                 `json:"points"`
		QuestionText string              `json:"question_text"`
		QuestionType string              `json:"question_type"`
		choices: {
			ChoiceID    *int64 `json:"choice_id"`
			ChoiceText  string `json:"choice_text"`
			IsCorrect   bool   `json:"is_correct"`
			OrderNumber int    `json:"order_number"`
		} 
	}
    RemoveQuestions []int64               `json:"remove_questions"`
}

Responses:

Status	Description	Body
201	Created	{ "message": "Session created" }
400	Validation error { "error": "Invalid request" }
500	Server error { "error": "Server error" }
```


### Create tutor location

Bind tutor to location
	•	URL: /create_tutor_location
	•	Method: POST
	•	Auth required: Yes

Request Body:
```
{
    LocationId     *int64 `json:"location_id"`
    TutorId        *int64 `json:"tutor_id"`
    OrganizationID *int64 `json:"organization_id"`
}

Responses:

Status	Description	Body
201	Created	{ "message": "Tutor location created" }
400	Validation error { "error": "Invalid request" }
500	Server error { "error": "Server error" }

```
### Create subject

Create a subject.

	•	URL: /create_subject
	•	Method: POST
	•	Auth required: Yes

Request Body:
```
{
    ID             *int64 `json:"id"`
    Title          string `json:"title"`
    Description    string `json:"description"`
    OrganizationId *int64 `json:"organization_id"`
}

Responses:

Status	Description	Body
201	Created	{ "message": "Subject created" }
400	Validation error { "error": "Invalid request" }
500	Server error { "error": "Server error" }
```

### Create permission

Grant user persmissions

	•	URL: /create_permission
	•	Method: POST
	•	Auth required: Yes

Request Body:
```
{
    ID                *int64   `json:"id"`
    Role              string   `json:"role"`
    User              string   `json:"user"`
    Permissions       []string `json:"permissions"`
    UpdatePermissions []string `json:"updatePermissions"`
    OrganizationId    *int64   `json:"organization_id"`
}

Responses:

Status	Description	Body
201	Created	{ "message": "Permission created" }
400	Validation error { "error": "Invalid request" }
500	Server error { "error": "Server error" }
```

### Create Subject location

Bind a subject to a location

	•	URL: /create_subject_location
	•	Method: POST
	•	Auth required: Yes

Request Body:
```
{
    ID             int  `json:"id"`
    SubjectID      *int `json:"subject_id"`
    LocationID     *int `json:"location_id"`
    OrganizationID *int `json:"organization_id"` // Required
}

Responses:

Status	Description	Body
201	Created	{ "message": "subject location created" }
400	Validation error { "error": "Invalid request" }
500	Server error { "error": "Server error" }
```

### Create announcments

Create an announcment

	•	URL: /create_announcement
	•	Method: POST
	•	Auth required: Yes

Request Body:
```
{
    ID             int    `json:"id"`
    Title          string `json:"title"`
    Body           string `json:"body"`
    CreatedAt      string `json:"created_at"`
    LocationID     []*int `json:"location_id"`     
    Severity       string `json:"severity"`        
    OrganizationID int    `json:"organization_id"` 
    ProgramID      *int   `json:"program_id"`      
    AdminID        *int   `json:"admin_id"`      
    StaffID        *int   `json:"staff_id"`
}

Responses:

Status	Description	Body
201	Created	{ "message": "announcment created" }
400	Validation error { "error": "Invalid request" }
500	Server error { "error": "Server error" }
```

### Create assessment session

Create online assessments per student. 
Returns a session_id, number of session active, status.

	•	URL: /create_assessment_sessions
	•	Method: POST
	•	Auth required: Yes

Request Body:
```
{
    students: {
		ID           *int64 `json:"id"`
		FirstName    string `json:"first_name"`
		LastName     string `json:"last_name"`
		ProgramId    *int64 `json:"program_id"`
		LocationId   *int64 `json:"location_id"`
		TutorId      *int64 `json:"tutor_id"`
		SubjectId    *int64 `json:"subject_id"`
		AssessmentId *int64 `json:"assessment_id"`
		SemesterId   *int64 `json:"semester_id"`
		EasyScoreID  bool   `json:"easy_score"`
	}
}

Responses:

Status	Description	Body
201	Created	{ "Status", "session_active", "session_id" }
400	Validation error { "error": "Invalid request" }
500	Server error { "error": "Server error" }
```

### Additional information
> [!NOTE]
> If endpoint starts as /create_program 
> ID parameter is not required in the body.
> It's update endpoint is /update_program
> ID parameter is required.
> It's delete endpoint is /delete_program
> ID parameter is required.


## 6
### Some exciting features being implemented

1. Attach files to materials [x]
	- A useful feature when pre-defined documents are available. These files are easily viewable by relevant staff members.
2. Add Images to Questions [x]
	- Improvement from previous implementation, that is before users must have a static https link to image
	- Improvement is users can dynamically attach images or files to each question, no static link required.
3. Artifical intelligence [Testing]
	- Leverages LLM models to generate assessments based on various user-defined factors. 
	- Utilizes database context such as location, district, and academic standards to create contextually appropriate questions.
	- Context aware/ context injection using google gemini API.
4. Artifical intelligence [Testing]
	- Grade assessments, eleminate potential user error.



### How features above are implemented

1. Attach files to materials
	- With the use of S3, the user can upload files, images, etc. Each upload is marked with a UUID and stored in the database for retrieval.
	- When a user needs to fetch such a document, an endpoint is triggered with the UUID as a parameter and returns a signed URL.

2. Add images to questions
	- Since the user can create N number of questions, uploading N number of images at once can become an issue if N is large.
	- To circumvent this, the user will upload each image as they create a new question.
	- Once the user uploads an image, a UUID is stored on the frontend. If the screen is exited at any point, each stored UUID is deleted.
	- Otherwise, on submit, link the UUID to each question.
	- This solution works well since it allows for the implementation of AI-generated assessments. Each question is queued with the ability to upload an image, as well as review and update generated assessment questions.

3. Artifical intelligence
	- To prevent from API calls to Bedrock to delay its best to ofload this task to a microservice.
	- Using RabbitMQ to send request to a python function/listener would be a viable solution. 
	- We also ensure each request is re-queued if failure occurs. 
	- Use an input/output key UUID to reference the S3 bucket object and the RDS input key row.
	- Update the RDS row to indicate "complete"; the frontend polls and fetches from S3 using the output/input key.



## 7
License

This project is licensed under the MIT License.

⸻

Built with ❤️ 