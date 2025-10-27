package models

import (
	"database/sql"
	"time"
)

type User struct {
	ID             int64    `json:"id"`
	Email          string   `json:"email"`
	Password       string   `json:"password"`
	Type           string   `json:"type"`
	Permissions    []string `json:"permissions"`
	OrganizationId int64    `json:"organization_id"`
	FirstName      string   `json:"first_name"`
	LastName       string   `json:"last_name"`
	Active         bool     `json:"active"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Type     string `json:"type"`
}

type LoginResponse struct {
	User           User                         `json:"user"`
	Token          *string                      `json:"token"`
	RefreshToken   *string                      `json:"refresh_token"`
	Permissions    LoginResponsePermissions     `json:"permissions"`
	TutorLocations []TutorLocationList          `json:"tutor_locations"`
	TutorPrograms  []ResponseRequestProgramList `json:"tutor_programs"`
}

type LoginResponsePermissions struct {
	DisableUpdate bool `json:"disable_update"`
	DisableCreate bool `json:"disable_create"`
	DisableDelete bool `json:"disable_delete"`
}

type RegisterOrganizationResponse struct {
	Status string `json:"status,omitempty"`
}

type ResponseGenerationOrganization struct {
	Id      *int64  `json:"id"`
	Title   *string `json:"title"`
	Address *string `json:"address"`
	ZipCode *string `json:"zip_code"`
	City    *string `json:"city"`
	State   *string `json:"state"`
}

type ResponseGenerationMaterials struct {
	InputSum  *int64  `json:"input_sum"`
	OutputSum *int64  `json:"output_sum"`
	YearMonth *string `json:"year_month"`
}

type ResponseGenerationQuestions struct {
	InputSum  *int64  `json:"input_sum"`
	OutputSum *int64  `json:"output_sum"`
	YearMonth *string `json:"year_month"`
}

type RegisterResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

type RegisterRequestAdminRoot struct {
	Password         string `json:"password"`
	Email            string `json:"email"`
	OrganizationName string `json:"organization_name"`
	Address          string `json:"address"`
	State            string `json:"state"`
	City             string `json:"city"`
	ZipCode          string `json:"zip_code"`
	OrganizationId   *int64 `json:"organization_id"`
}

type RegisterResponseAdminRoot struct {
	ID             int64  `json:"id"`
	Email          string `json:"email"`
	OrganizationId *int64 `json:"organization_id"`
}
type RegisterRequestStudents struct {
	ID                *int64  `json:"id"`
	FirstName         string  `json:"firstname"`
	MiddleName        string  `json:"middle_name"`
	LastName          string  `json:"last_name"`
	Period            *int64  `json:"period,omitempty"`
	SemesterID        *int64  `json:"semester_id"`
	Email             *string `json:"email"`
	GradeLevel        int64   `json:"grade_level"`
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
	Race              *string `json:"race"`
	Gender            *string `json:"gender"`
}

type ResponseRequestStudentList struct {
	ID                int64   `json:"id"`
	FirstName         string  `json:"first_name"`
	MiddleName        string  `json:"middle_name"`
	LastName          string  `json:"last_name"`
	Email             *string `json:"email"`
	GradeLevel        int64   `json:"grade_level"`
	Active            bool    `json:"active"`
	CreatedAt         string  `json:"created_at"`
	CreatedBy         string  `json:"created_by"`
	Period            *int64  `json:"period,omitempty"`
	SemesterId        *int64  `json:"semester_id"`
	DirectPartnership bool    `json:"direct_partnership"`
	TeacherID         *int64  `json:"teacher_id"`
	TeacherName       *string `json:"teacher_name"`
	LocationId        *int64  `json:"location_id"`
	DurationRequired  *bool   `json:"duration_required"`
	Timeframe         *bool   `json:"timeframe"`
	TimeframeStart    *string `json:"timeframe_start"`
	TimeframeEnd      *string `json:"timeframe_end"`
	Race              *string `json:"race"`
	Gender            *string `json:"gender"`
}

type ResponseRequestStudents struct {
	Status    string `json:"status"`
	StudentID int64  `json:"id"`
}
type ResponseRequestLocations struct {
	ID         int64  `json:"id"`
	Name       string `json:"name"`
	Address    string `json:"address"`
	City       string `json:"city"`
	State      string `json:"state"`
	Archive    bool   `json:"archive"`
	ZipCode    string `json:"zip_code"`
	CreatedAt  string `json:"created_at"`
	DistrictId *int64 `json:"district_id"`
}

type RegisterRequestLocation struct {
	ID             int64  `json:"id"`
	Name           string `json:"name"`
	Address        string `json:"address"`
	DistrictId     int64  `json:"district_id"`
	City           string `json:"city"`
	State          string `json:"state"`
	Archive        bool   `json:"archive"`
	ZipCode        string `json:"zip_code"`
	CreatedAt      string `json:"created_at"`
	OrganizationId *int64 `json:"organization_id"`
}

type ResponseRequestLocation struct {
	Status     string `json:"status"`
	LocationId int64  `json:"id"`
}

type RegisterRequestAdmin struct {
	ID             *int64 `json:"id"`
	Fullname       string `json:"fullname"`
	Email          string `json:"email"`
	Region         string `json:"region"`
	State          string `json:"state"`
	Password       string `json:"password_hash"`
	RootId         int64  `json:"root_id"`
	Active         bool   `json:"active"`
	DistrictId     *int64 `json:"district_id"`
	OrganizationId *int64 `json:"organization_id"`
}
type ResponseRequestAdminList struct {
	ID         int64  `json:"id"`
	Fullname   string `json:"fullname"`
	Email      string `json:"email"`
	Active     bool   `json:"active"`
	Region     string `json:"region"`
	DistrictId *int64 `json:"district_id"`
	State      string `json:"state"`
}

type ResponseRequestAdmin struct {
	Status   string `json:"status"`
	Admin_id int64  `json:"id"`
}

type ResponseRequestDistrictList struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	City      string `json:"city"`
	State     string `json:"state"`
	Region    string `json:"region"`
	CreatedAt string `json:"created_at"`
}

type ResponseRequestProgramList struct {
	ID                int64  `json:"id"`
	LocationID        *int64 `json:"location_id"`
	ProgramName       string `json:"program_name"`
	TimeFrameRequired bool   `json:"timeframe_required"`
	SurveyRequired    bool   `json:"survey_required"`
	CreatedAt         string `json:"created_at"`
}

type RegisterRequestDistrict struct {
	ID             *int64 `json:"id"`
	Name           string `json:"name"`
	City           string `json:"city"`
	Region         string `json:"region"`
	State          string `json:"state"`
	AdminId        int64  `json:"admin_id"`
	OrganizationId *int64 `json:"organization_id"`
}

type ResponseRequestDistrict struct {
	Status     string `json:"status"`
	DistrictId int64  `json:"id"`
}

type RegisterRequestProgram struct {
	ID                *int64 `json:"id"`
	ProgramName       string `json:"program_name"`
	AdminID           int64  `json:"admin_id"`
	OrganizationId    *int64 `json:"organization_id"`
	SurveyRequired    bool   `json:"survey_required"`
	TimeFrameRequired bool   `json:"timeframe_required"`
}

type RegisterRequestSurvey struct {
	ID              *int64            `json:"id"`
	Title           *string           `json:"title"`
	Description     *string           `json:"description"`
	IsActive        bool              `json:"is_active"`
	OrderBy         *int64            `json:"order_by"`
	Questions       []SurveyQuestions `json:"questions"`
	RemoveQuestions []int64           `json:"remove_questions"`
}

type SurveyQuestions struct {
	ID           *int64                 `json:"id"`
	SurveyID     *int64                 `json:"survey_id"`
	OrderIndex   *int64                 `json:"order_index"`
	QuestionText *string                `json:"question_text"`
	QuestionType *string                `json:"question_type"`
	Choices      []SurveyQuestionChoice `json:"choices"`
}

type SurveyQuestionChoice struct {
	ID               *int64  `json:"id"`
	SurveyQuestionId *int64  `json:"survey_question_id"`
	ChoiceText       *string `json:"choice_text"`
}

type PermissionsMap struct {
	ID             *int64 `json:"id"`
	PermissionName string `json:"permission_name"`
}

type RegisterPermissionRequest struct {
	ID                *int64   `json:"id"`
	Role              string   `json:"role"`
	User              string   `json:"user"`
	Permissions       []string `json:"permissions"`
	UpdatePermissions []string `json:"updatePermissions"`
	OrganizationId    *int64   `json:"organization_id"`
}
type RegisterPermissionResponse struct {
	Status string `json:"status"`
}

type ResponseRequestSurvey struct {
	Status    string `json:"status"`
	ProgramId int64  `json:"id"`
}
type ResponseRequestProgram struct {
	Status    string `json:"status"`
	ProgramId int64  `json:"id"`
}

type ResponseRequestMaterialsList struct {
	ID           int64   `json:"id"`
	Title        string  `json:"title"`
	ExternalLink string  `json:"external_link"`
	Description  string  `json:"description"`
	ProgramName  *string `json:"program_name"`
	Version      float64 `json:"version"`
	Pre          bool    `json:"pre"`
	Mid          bool    `json:"mid"`
	Post         bool    `json:"post"`
	Visible      bool    `json:"visible"`
	ProgramId    *int64  `json:"program_id"`
	SReference   *string `json:"s3_reference"`
	CreatedAt    string  `json:"created_at"`
}

type RegisterRequestMaterials struct {
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
	File             bool    `json:"file"`
	FileType         *string `json:"file_type"`
}

type ResponseRequestMaterials struct {
	Status     string  `json:"status"`
	MaterialId int64   `json:"id"`
	UploadUrl  *string `json:"upload_url"`
}

type ResponseRequestTutorsList struct {
	ID         int64  `json:"id"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	Email      string `json:"email"`
	Active     bool   `json:"active"`
	CreatedAt  string `json:"created_at"`
	LocationId *int64 `json:"location_id"`
}

type RegisterRequestTutor struct {
	ID             int64  `json:"id"`
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	Password       string `json:"password"`
	Email          string `json:"email"`
	Active         bool   `json:"active"`
	EmailChange    string `json:"email_change"`
	CreatedAt      string `json:"created_at"`
	LocationId     *int64 `json:"location_id"`
	OrganizationId *int64 `json:"organization_id"`
}

type ResponseRequestTutor struct {
	Status  string `json:"status"`
	TutorId int64  `json:"id"`
}

type ResponseRequestSemesterList struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	Year      int64     `json:"year"`
	Archive   bool      `json:"archive"`
	DateStart time.Time `json:"date_start"`
	DateEnd   time.Time `json:"date_end"`
	Active    bool      `json:"active"`
	WorkWeek  *[]string `json:"work_week"`
}

type ResponseRequestSemesterLocationList struct {
	ID             int64  `json:"id"`
	LocationId     *int64 `json:"location_id"`
	OrganizationId *int64 `json:"organization_id"`
	SemesterID     int64  `json:"semester_id"`
	SemesterName   int64  `json:"semester_name"`
	CreatedAt      string `json:"created_at"`
	Title          string `json:"title"`
	Year           int64  `json:"year"`
	DateStart      string `json:"date_start"`
	DateEnd        string `json:"date_end"`
}

type RegisterSemesterDates struct {
	SemesterID   *int64    `json:"semester_id"`
	Notes        *string   `json:"notes"`
	ScheduleType *string   `json:"schedule_type"`
	StartDate    time.Time `json:"start_date"`
	EndDate      time.Time `json:"end_date"`
}

type RegisterStudentGroup struct {
	ID          *int64  `json:"id"`
	Title       *string `json:"title"`
	Description *string `json:"description"`
	SemesterID  *int64  `json:"semester_id"`
	LocationID  *int64  `json:"location_id"`
	TutorID     *int64  `json:"tutor_id"`
}

type ResponseSemesterDates struct {
	ID           *int64    `json:"id"`
	SemesterID   *int64    `json:"semester_id"`
	Notes        *string   `json:"notes"`
	ScheduleType *string   `json:"schedule_type"`
	StartDate    time.Time `json:"start_date"`
	EndDate      time.Time `json:"end_date"`
	CreatedAt    time.Time `json:"created_at"`
}

type RegisterRequestSemester struct {
	ID             *int64    `json:"id"`
	Title          string    `json:"title"`
	Year           *int64    `json:"year"`
	OrganizationId *int64    `json:"organization_id"`
	DateStart      string    `json:"date_start"`
	DateEnd        string    `json:"date_end"`
	Active         bool      `json:"active"`
	Archive        bool      `json:"archive"`
	WorkWeek       *[]string `json:"work_week"`
}

type RegisterRequestSemesterLocation struct {
	ID             *int64 `json:"id"`
	LocationID     *int64 `json:"location_id"`
	OrganizationId *int64 `json:"organization_id"`
	SemesterID     *int64 `json:"semester_id"`
}

type ResponseRequestSemester struct {
	Status string `json:"status"`
	ID     *int64 `json:"id"`
}

type ResponseUpdateAdmin struct {
	ID     *int64 `json:"id"`
	Status string `json:"status"`
}

type ResponseUpdate struct {
	Status    string  `json:"status"`
	UploadUrl *string `json:"upload_url"`
}

type RemoveAdmin struct {
	ID *int64 `json:"id"`
}

type RemoveResponse struct {
	Status string `json:"status"`
}

type RemoveRequest struct {
	ID             *int64 `json:"id"`
	OrganizationId int64  `json:"organization_id"`
}

type RemoveAdminLocation struct {
	AdminID        *int64 `json:"admin_id"`
	LocationID     *int64 `json:"location_id"`
	OrganizationId int64  `json:"organization_id"`
}

type RegisterSurveyProgram struct {
	ID        *int64 `json:"id"`
	ProgramID *int64 `json:"program_id"`
	SurveyID  *int64 `json:"survey_id"`
}
type ResponseSurveyProgram struct {
	ID          *int64  `json:"id"`
	ProgramID   *int64  `json:"program_id"`
	SurveyTitle *string `json:"survey_title"`
	SurveyID    *int64  `json:"survey_id"`
}

// ss.id, ss.title, ss.description, sq.id, sq.order_index, sq.question_text, sq.question_type, sc.id, sc.choice_text
type ResponseQuestionsSurvey struct {
	SurveyTitle        *string `json:"survey_title"`
	SurveyID           *int64  `json:"id"`
	SurveyDescription  *string `json:"survey_description"`
	QuestionID         *int64  `json:"question_id"`
	QuestionIndex      *int64  `json:"question_index"`
	QuestionText       *string `json:"question_text"`
	QuestionType       *string `json:"question_type"`
	QuestionChoiceID   *int64  `json:"question_choice_id"`
	QuestionChoiceText *string `json:"question_choice_text"`
}

type ResponseRecentTutorSessions struct {
	SessionID    *int64    `json:"session_id"`
	StudentCount *int64    `json:"student_count"`
	LocationId   *int64    `json:"location_id"`
	SubstituteId *int64    `json:"substitute_id"`
	ProgramId    *int64    `json:"program_id"`
	SemesterId   *int64    `json:"semester_id"`
	SubjectId    *int64    `json:"subject_id"`
	SessionDate  time.Time `json:"session_date"`
	Program      *string   `json:"program"`
	Location     *string   `json:"location"`
	FirstName    *string   `json:"first_name"`
	LastName     *string   `json:"last_name"`
	StartTime    string    `json:"start_time"`
	Substitute   bool      `json:"substitute"`
	TutorId      *int64    `json:"tutor_id"`
	Duration     *int64    `json:"duration"`
}

// SELECT t.first_name, t.last_name, s.tutor_id, p.program_name, SUM(s.duration), COUNT(s.id)
type ResponseLocalSessionAverage struct {
	Program      *string  `json:"program"`
	FirstName    *string  `json:"first_name"`
	LastName     *string  `json:"last_name"`
	TutorID      *int64   `json:"tutor_id"`
	ProgramID    *int64   `json:"program_id"`
	DurationSum  *float32 `json:"duration_sum"`
	SessionCount *float32 `json:"session_count"`
}

type ResponseAssessmentPivotTable struct {
	Program     *string    `json:"program"`
	FirstName   *string    `json:"first_name"`
	LastName    *string    `json:"last_name"`
	TutorID     *int64     `json:"tutor_id"`
	ProgramID   *int64     `json:"program_id"`
	StartTime   *string    `json:"start_time"`
	Duration    *int64     `json:"duration"`
	SessionDate *time.Time `json:"session_date"`
}

type ResponseAdminLocations struct {
	AdminID    *int64  `json:"admin_id"`
	LocationID *int64  `json:"location_id"`
	Location   *string `json:"location"`
}

type RegisterAdminLocation struct {
	AdminID    *int64 `json:"admin_id"`
	LocationID *int64 `json:"location_id"`
}

type RegisterOrganization struct {
	Fullname         *string `json:"fullname"`
	Email            *string `json:"email"`
	Password         *string `json:"password"`
	OrganizationName *string `json:"organization_name"`
	Address          *string `json:"address"`
	ZipCode          *string `json:"zip_code"`
	State            *string `json:"state"`
	City             *string `json:"city"`
	Code             *string `json:"code"`
}

type ResponseRegisterOrganization struct {
	Status *string `json:"status"`
}

type RegisterStudentSessionList struct {
	Session        RegisterTutorSession          `json:"session"`
	SessionList    []RegisterStudentSession      `json:"student_sessions"`
	Assessments    map[string]*AssessmentPayload `json:"assessments"`
	SessionToken   *string                       `json:"session_token"`
	OrganizationID *int64                        `json:"organization_id"`
}

type AssessmentPayload struct {
	AssessmentID *int64                 `json:"assessment_id"`
	Choices      map[string]interface{} `json:"choices,omitempty"`
	Grader       map[string]bool        `json:"grader,omitempty"`
}

type RegisterTutorSession struct {
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
	Duration     *int64 `json:"duration"`
	InSchool     bool   `json:"in_school"`
}

type SubjectList struct {
	ID             *int64 `json:"id"`
	Title          string `json:"title"`
	Description    string `json:"description"`
	OrganizationId *int64 `json:"organization_id"`
	CreatedAt      string `json:"created_at"`
}

type RegisterSubject struct {
	ID             *int64 `json:"id"`
	Title          string `json:"title"`
	Description    string `json:"description"`
	OrganizationId *int64 `json:"organization_id"`
}

type RegisterStudentGroupAttendies struct {
	GroupID  *int64      `json:"group_id"`
	Students []*Students `json:"students"`
}

type RegisterLocationContact struct {
	ID             *int64  `json:"id"`
	ProgramID      *int64  `json:"program_id"`
	LocationID     *int64  `json:"location_id"`
	Description    *string `json:"description"`
	FirstName      *string `json:"first_name"`
	LastName       *string `json:"last_name"`
	Email          *string `json:"email"`
	Phone          *string `json:"phone"`
	Notes          *string `json:"notes"`
	OrganizationId *int64  `json:"organization_id"`
}

type ResponseLocationContact struct {
	ID             *int64  `json:"id"`
	Description    *string `json:"description"`
	LocationID     *int64  `json:"location_id"`
	FirstName      *string `json:"first_name"`
	LastName       *string `json:"last_name"`
	Email          *string `json:"email"`
	Phone          *string `json:"phone"`
	Notes          *string `json:"notes"`
	OrganizationId *int64  `json:"organization_id"`
}

type SimpleReponse struct {
	ID     *int64 `json:"id"`
	Status string `json:"status"`
}

type ResponseRegisterSubject struct {
	Status    string `json:"status"`
	SubjectID int64  `json:"id"`
}

type ServiceSession struct {
	FirstName       string        `json:"first_name"`
	LastName        string        `json:"last_name"`
	ID              *int64        `json:"id"`
	Location        string        `json:"location_name"`
	SubstituteId    sql.NullInt64 `json:"substitute_id"`
	ProgramId       sql.NullInt64 `json:"program_id"`
	ProgramName     string        `json:"program_name"`
	Notes           string        `json:"notes"`
	SessionDate     time.Time     `json:"session_date"`
	StartTime       string        `json:"start_time"`
	Subject         sql.NullInt32 `json:"subject"`
	Substitute      sql.NullBool  `json:"substitute"`
	TutorId         sql.NullInt64 `json:"tutor_id"`
	CreatedAt       time.Time     `json:"created_at"`
	EditedAt        string        `json:"edited_at"`
	SubjectName     string        `json:"subject_name"`
	StudentCount    int64         `json:"student_count"`
	AssessmentCount *int64        `json:"assessment_count"`
	SubFirstName    *string       `json:"substitute_first_name"`
	SubLastName     *string       `json:"substitute_last_name"`
}

type StudentSessions struct {
	ID              *int64 `json:"id"`
	FirstName       string `json:"first_name"`
	LastName        string `json:"last_name"`
	MiddleName      string `json:"middle_name"`
	SessionCount    *int64 `json:"session_count"`
	AssessmentCount *int64 `json:"assessment_count"`
}

type RegisterStudentSession struct {
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

type StudentList struct {
	ID         *int64 `json:"id"`
	MaterialId *int64 `json:"material_id"`
	Score      *int64 `json:"score"`
	Subject    *int64 `json:"subject"`
	ProgramId  *int64 `json:"program_id"`
}

type ResponseStudentSession struct {
	Status          string `json:"status"`
	StudentCount    int64  `json:"student_count"`
	AssessmentCount int64  `json:"assessment_count"`
}

type ResponseAssessmentList struct {
	ID              *int64  `json:"id"`
	Title           string  `json:"title"`
	Description     string  `json:"description,omitempty"`
	Letter          string  `json:"letter"`
	Cycle           *int64  `json:"cycle"`
	AlphaIdentifier string  `json:"alpha_identifier,omitempty"`
	ExternalLink    string  `json:"external_link,omitempty"`
	MaxScore        *int64  `json:"max_score,omitempty"`
	SubjectId       *int64  `json:"subject_id,omitempty"`
	ProgramId       *int64  `json:"program_id"`
	MaterialID      *int64  `json:"material_id,omitempty"`
	GradeLevel      *int64  `json:"grade_level"`
	SubjectName     string  `json:"subject_name,omitempty"`
	ProgramName     string  `json:"program_name,omitempty"`
	MaterialName    string  `json:"material_name,omitempty"`
	OrganizationID  *int64  `json:"organization_id"`
	CreatedAt       string  `json:"created_at"`
	Visible         bool    `json:"visible"`
	Version         float64 `json:"version"`
	Pre             bool    `json:"pre"`
	Mid             bool    `json:"mid"`
	Post            bool    `json:"post"`
	EasyScore       bool    `json:"easy_score"`
	Questionnaire   bool    `json:"questionnaire"`
}

type ResponseAssessmentQuestionsChoice struct {
	QuestionID        *int64  `json:"question_id"`
	AssessmentID      *int64  `json:"assessment_id"`
	ImageURL          *string `json:"image_url"`
	QuestionText      string  `json:"question_text"`
	QuestionType      string  `json:"question_type"`
	Points            int64   `json:"points"`
	OrderNumber       int64   `json:"order_number"`
	ChoiceID          *int64  `json:"choice_id"`
	ChoiceText        *string `json:"choice_text"`
	ChoiceOrderNumber int64   `json:"choice_order"`
	IsCorrect         *bool   `json:"is_correct"`
}

type ResponseAssessmentQuestionsChoiceExternal struct {
	QuestionID        *int64  `json:"question_id"`
	AssessmentID      *int64  `json:"assessment_id"`
	ImageURL          *string `json:"image_url"`
	QuestionText      string  `json:"question_text"`
	QuestionType      string  `json:"question_type"`
	Points            int64   `json:"points"`
	OrderNumber       int64   `json:"order_number"`
	ChoiceID          *int64  `json:"choice_id"`
	ChoiceText        *string `json:"choice_text"`
	ChoiceOrderNumber int64   `json:"choice_order"`
}

/**
	student_id int64 REFERENCES stu_tracker.Students(id) ON DELETE CASCADE,
    assessment_id int64 REFERENCES stu_tracker.Assessments(id) ON DELETE CASCADE,
    session_token UUID NOT NULL,
    semester_id int64 REFERENCES stu_tracker.Semester(id) ON DELETE SET NULL,
	sleep_hours FLOAT DEFAULT 0,
    effort_score smallint NOT NULL CHECK (effort_score BETWEEN 0 AND 10),
    tutor_sessions int64 DEFAULT 0,
    parental_help smallint NOT NULL CHECK (parental_help BETWEEN 0 AND 3),
    sports_hours int64 DEFAULT 0,
    organization_id int64 REFERENCES stu_tracker.Organization(id) ON DELETE CASCADE,
    peer_influence smallint NOT NULL CHECK (peer_influence BETWEEN 0 AND 3),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
**/

type RegisterPreAssessment struct {
	StudentId     *int64  `json:"student_id"`
	AssessmentId  *int64  `json:"assessment_id"`
	SessionToken  *string `json:"session_token"`
	SleepHours    float32 `json:"sleep_hours"`
	StudyHours    float32 `json:"study_hours"`
	ParentalHelp  int64   `json:"parental_help"`
	EffortScore   int64   `json:"effort_score"`
	TutorSessions int64   `json:"tutor_sessions"`
	SportHours    int64   `json:"sport_hours"`
	PeerInfluence int64   `json:"peer_influence"`
}
type RegisterAnnouncementsAck struct {
	AnnouncmentID *int64 `json:"announcement_id"`
}

type RegisterAssessment struct {
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
	MaterialID      *int64                `json:"material_id,omitempty"`
	ProgramId       *int64                `json:"program_id"`
	GradeLevel      *int64                `json:"grade_level"`
	CreatedAt       string                `json:"created_at"`
	Version         float64               `json:"version"`
	Pre             bool                  `json:"pre"`
	Mid             bool                  `json:"mid"`
	Post            bool                  `json:"post"`
	EasyScore       bool                  `json:"easy_score"`
	Questions       []AssessmentQuestions `json:"questions"`
	RemoveQuestions []int64               `json:"remove_questions"`
	Images          []ImagePayload        `json:"images"`
	Questionnaire   bool                  `json:"questionnaire"`
}

type ImagePayload struct {
	Index int64  `json:"index"`
	ID    string `json:"id"`
}

type AssessmentQuestions struct {
	QuestionID   *int64              `json:"question_id"`
	ImageURL     string              `json:"image_url"`
	Required     bool                `json:"is_required"`
	OrderNumber  int64               `json:"order_number"`
	Standard     *string             `json:"standard_text"`
	Points       int64               `json:"points"`
	QuestionText string              `json:"question_text"`
	QuestionType string              `json:"question_type"`
	Choice       []AssessmentChoices `json:"choices"`
}

type AssessmentChoices struct {
	ChoiceID    *int64 `json:"choice_id"`
	ChoiceText  string `json:"choice_text"`
	IsCorrect   bool   `json:"is_correct"`
	OrderNumber int64  `json:"order_number"`
}

type RegisterResponseAssessment struct {
	Status       string `json:"status"`
	AssessmentId int64  `json:"assessment_id"`
}

type RegisterLocationProgram struct {
	ProgramId      *int64 `json:"program_id"`
	LocationId     *int64 `json:"location_id"`
	OrganizationID *int64 `json:"organization_id"`
}

type RegisterProgramSurvey struct {
	ProgramId      *int64 `json:"program_id"`
	SurveyId       *int64 `json:"survey_id"`
	OrganizationID *int64 `json:"organization_id"`
}

type RemoveLocationProgram struct {
	ProgramId      *int64 `json:"program_id"`
	LocationId     *int64 `json:"location_id"`
	OrganizationID *int64 `json:"organization_id"`
}

type RegisterResponseProgramSurvey struct {
	Status string `json:"status"`
}

type RegisterResponseLocationProgram struct {
	Status string `json:"status"`
}

type SearchQuery struct {
	OrganizationID *int64    `json:"organization_id"`
	SearchTerm     string    `json:"search_term,omitempty"`
	LocationId     *int64    `json:"location_id,omitempty"`
	ProgramId      *int64    `json:"program_id,omitempty"`
	SemesterID     *int64    `json:"semester_id,omitempty"`
	DateStart      time.Time `json:"date"`
	DateEnd        time.Time `json:"date_end"`
	SubjectId      *int64    `json:"subject_id,omitempty"`
}

type SearchQueryTutor struct {
	OrganizationID *int64 `json:"organization_id"`
	SearchTerm     string `json:"search_term"`
}

type SessionSearchResponse struct {
	Status string    `json:"status"`
	Data   []Session `json:"data"`
}

type Session struct {
	ID              *int64 `json:"id"`
	FirstName       string `json:"first_name"`
	LastName        string `json:"last_name"`
	SessionDate     string `json:"session_date"`
	Duration        *int64 `json:"duration"`
	StartTime       string `json:"start_time"`
	Notes           string `json:"notes"`
	OrganizationId  *int64 `json:"organization_id"`
	LocationId      *int64 `json:"location_id"`
	TutorId         *int64 `json:"tutor_id"`
	Subject         string `json:"subject"`
	AssessmentId    *int64 `json:"assessment_id"`
	AssessmentScore *int64 `json:"score"`
}

type SessionInfoStudent struct {
	ID             *int64  `json:"student_id"`
	FirstName      string  `json:"first_name"`
	LastName       string  `json:"last_name"`
	Email          *string `json:"email"`
	MiddleName     string  `json:"middle_name"`
	Duration       *int64  `json:"duration"`
	Period         *int64  `json:"period,omitempty"`
	Grade          *int64  `json:"grade"`
	Timeframe      *bool   `json:"timeframe"`
	TimeframeStart *string `json:"timeframe_start"`
	TimeframeEnd   *string `json:"timeframe_end"`
}

type AssessmentInfoStudent struct {
	Title     string    `json:"title"`
	Letter    string    `json:"letter"`
	Cycle     string    `json:"cycle"`
	MaxScore  *int64    `json:"max_score"`
	Subject   *int64    `json:"subject"`
	Score     *float64  `json:"score"`
	CreatedAt time.Time `json:"created_at"`
	StudentID *int64    `json:"student_id"`
}
type StudentSessionInfo struct {
	CreatedAt time.Time `json:"created_at"`
	SubjectID *int64    `json:"subject_id"`
	Duration  *int64    `json:"duration"`
	Absent    bool      `json:"absent"`
}

type StudentAssessmentInfo struct {
	ID          *int64    `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	SessionDate time.Time `json:"session_date"`
	Score       *float64  `json:"score"`
	MaxScore    *int64    `json:"max_score"`
	Letter      string    `json:"letter"`
	Cycle       string    `json:"cycle"`
	SessionID   *int64    `json:"session_id"`
	Title       string    `json:"title"`
	Pre         bool      `json:"pre"`
	Mid         bool      `json:"mid"`
	Post        bool      `json:"post"`
	Version     float64   `json:"version"`
	EasyScore   bool      `json:"easy_score"`
}

type SessionTrail struct {
	ID              int64     `json:"id"`
	FirstName       string    `json:"first_name"`
	LastName        string    `json:"last_name"`
	ProgramName     string    `json:"program_name"`
	SessionDuration int64     `json:"session_duration"`
	StartTime       string    `json:"start_time"`
	Notes           string    `json:"notes"`
	CreatedAt       time.Time `json:"created_at"`
	StudentCount    int64     `json:"student_count"`
	Substitute      bool      `json:"substitute"`
	Absent          bool      `json:"absent"`
	StudentDuration int64     `json:"student_duration"`
	SubstituteName  string    `json:"substitute_name,omitempty"`
	SessionDate     time.Time `json:"session_date"`
	Timeframe       *bool     `json:"timeframe"`
	TimeframeStart  *string   `json:"timeframe_start"`
	TimeframeEnd    *string   `json:"timeframe_end"`
}

type TutorLocationList struct {
	LocationName string `json:"location_name"`
	ID           int64  `json:"id"`
}
type TutorProgramList struct {
	ProgramName string `json:"program_name"`
	ID          int64  `json:"id"`
}

type RegisterTutorLocation struct {
	LocationId     *int64 `json:"location_id"`
	TutorId        *int64 `json:"tutor_id"`
	OrganizationID *int64 `json:"organization_id"`
}

type RegisterTutorResponse struct {
	Status string `json:"status"`
}

type RemoveTutorLocation struct {
	LocationId     *int64 `json:"location_id"`
	TutorId        *int64 `json:"tutor_id"`
	OrganizationID *int64 `json:"organization_id"`
}

type PermissionsList struct {
	ID          *int64 `json:"id"`
	Description string `json:"description"`
	Name        string `json:"name"`
}

type AnnouncementRequest struct {
	ID             int64   `json:"id"`
	Email          string  `json:"email"`
	Role           string  `json:"role"`
	OrganizationID int64   `json:"organization_id"`
	ProgramID      []int64 `json:"program_ids"`
	LocationIDs    []int64 `json:"location_ids"`
}

type RegisterAnnouncements struct {
	ID             int64  `json:"id"`
	Title          string `json:"title"`
	Body           string `json:"body"`
	CreatedAt      string `json:"created_at"`
	LocationID     *int64 `json:"location_id"`     // Nullable
	Severity       string `json:"severity"`        // Default: "info"
	OrganizationID *int64 `json:"organization_id"` // Required
	ProgramID      *int64 `json:"program_id"`      // Nullable
	AdminID        *int64 `json:"admin_id"`        // Required
	StaffID        *int64 `json:"staff_id"`
}

type RegisterUpdateAnnouncements struct {
	ID             *int64 `json:"id"`
	Title          string `json:"title"`
	Body           string `json:"body"`
	CreatedAt      string `json:"created_at"`
	LocationID     *int64 `json:"location_id"`     // Nullable
	Severity       string `json:"severity"`        // Default: "info"
	OrganizationID *int64 `json:"organization_id"` // Required
	ProgramID      *int64 `json:"program_id"`      // Nullable
	AdminID        *int64 `json:"admin_id"`        // Required
	StaffID        *int64 `json:"staff_id"`
}

type ResponseAnnouncement struct {
	ID     int64  `json:"id"`
	Status string `json:"status"`
}

type ResponseAnnouncementAck struct {
	ID     *int64 `json:"id"`
	Status string `json:"status"`
}

type AnnouncementsList struct {
	ID             int64  `json:"id"`
	Title          string `json:"title"`
	Body           string `json:"body"`
	CreatedAt      string `json:"created_at"`
	LocationID     *int64 `json:"location_id"`     // Nullable
	Severity       string `json:"severity"`        // Default: "info"
	OrganizationID int64  `json:"organization_id"` // Required
	ProgramID      *int64 `json:"program_id"`      // Nullable
	AdminID        *int64 `json:"admin_id"`        // Required
	StaffID        *int64 `json:"staff_id"`
	StaffName      string `json:"staff_name"`
	LocationName   string `json:"location_name"`
	ProgramName    string `json:"program_name"`
	Ack            bool   `json:"ack"`
}

type AnnouncementsListAck struct {
	ID             *int64    `json:"id"`
	TutotID        *int64    `json:"tutor_id"`
	AnnouncementID *int64    `json:"announcement_id"`
	OrganizationID *int64    `json:"organization_id"`
	FirstName      *string   `json:"first_name"`
	LastName       *string   `json:"last_name"`
	AcknowledgedAt time.Time `json:"acknowledged_at"`
	Acknowledged   bool      `json:"acknowledged"`
}

type RegisterSubjectLocation struct {
	ID             int64  `json:"id"`
	SubjectID      *int64 `json:"subject_id"`
	LocationID     *int64 `json:"location_id"`
	OrganizationID *int64 `json:"organization_id"` // Required
}

type RemoveSubjectLocation struct {
	SubjectID  *int64 `json:"subject_id"`
	LocationID *int64 `json:"location_id"`
}

type ResponseSubjectLocation struct {
	ID     *int64 `json:"id"`
	Status string `json:"status"`
}

type RequestTutorSessions struct {
	ID             *int64 `json:"id"`
	Email          string `json:"email"`
	OrganizationID *int64 `json:"organization_id"`
	SemesterID     *int64 `json:"semester_id"`
	LocationID     *int64 `json:"location_id"`
}

type Students struct {
	ID               *int64  `json:"id"`
	SessionID        *int64  `json:"session_id"`
	StudentID        *int64  `json:"student_id"`
	FirstName        string  `json:"first_name"`
	MiddleName       string  `json:"middle_name"`
	LastName         string  `json:"last_name"`
	Grade            int64   `json:"grade_level"`
	Timeframe        *bool   `json:"timeframe"`
	DurationRequired *bool   `json:"duration_required"`
	TimeframeStart   *string `json:"timeframe_start"`
	TimeframeEnd     *string `json:"timeframe_end"`
}

type TutorSessionsList struct {
	SessionID         *int64                 `json:"session_id"`
	ProgramName       string                 `json:"program_name"`
	LocationID        *int64                 `json:"location_id"`
	SemesterID        *int64                 `json:"semester_id"`
	ProgramID         *int64                 `json:"program_id"`
	SessionDate       time.Time              `json:"session_date"`
	StartTime         string                 `json:"start_time"`
	Students          []Students             `json:"students"`
	SurveyRequirement SessionCompletedSurvey `json:"survey_requirement"`
	StudentCount      int64                  `json:"student_count"`
	InSchool          bool                   `json:"in_school"`
	SurveyRequired    bool                   `json:"survey_required"`
	Substitute        bool                   `json:"substitute"`
	Semester          string                 `json:"semester"`
	LocationName      string                 `json:"location_name"`
}

type StudentAssessmentSearch struct {
	QuestionID   *int64  `json:"question_id"`
	Question     *string `json:"question"`
	QuestionType *string `json:"question_type"`
	Points       int64   `json:"points"`
	MaxPoints    int64   `json:"max_points"`
	IsCorrect    bool    `json:"is_correct"`
	ChoiceID     *int64  `json:"choice_id"`
	AnswerText   *string `json:"answer_text"`
	ChoiceText   *string `json:"choice_text"`
}

type RegisterTeacher struct {
	ID         *int64 `json:"id"`
	Name       string `json:"name"`
	Room       string `json:"room"`
	GradeLevel int64  `json:"grade_level"`
	Substitute bool   `json:"substitute"`
	LocationID *int64 `json:"location_id"`
}

type RegisterTeacherResponse struct {
	Status string `json:"status"`
	ID     *int64 `json:"id"`
}

type ResponseTeachers struct {
	ID         *int64 `json:"id"`
	Name       string `json:"name"`
	Room       string `json:"room"`
	GradeLevel int64  `json:"grade_level"`
	Substitute bool   `json:"substitute"`
	LocationID *int64 `json:"location_id"`
}

type StudentDetails struct {
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	MiddleName string `json:"middle_name"`
	GradeLevel int64  `json:"grade_level"`
}

type SessionCompletedSurvey struct {
	Completed bool `json:"completed"`
}
type DeleteImageRequest struct {
	ID string `json:"id"`
}

type ResponsePreAssessment struct {
	Status string `json:"status"`
	ID     int64  `json:"id"`
}

type Survey struct {
	ID             *int64    `json:"id" db:"id"`
	OrganizationID *int64    `json:"organization_id,omitempty" db:"organization_id"` // nullable
	Title          string    `json:"title" db:"title"`
	Description    string    `json:"description" db:"description"`
	OrderBy        int       `json:"order_by" db:"order_by"`
	IsActive       bool      `json:"is_active" db:"is_active"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
}

type SurveyChoicesList struct {
	ID               *int64  `json:"id"`
	QuestionServeyId *int64  `json:"question_survey_id"` // nullable
	ChoiceText       *string `json:"choice_text"`
}

type RegisterSurvey struct {
	SessionID *int64                    `json:"session_id"`
	Survey    []RegisterSurveyQuestions `json:"survey"`
}

type RegisterSurveyQuestions struct {
	Questions []Questions `json:"questions"`
}

type Questions struct {
	QuestionID     *int64  `json:"question_id"`
	QuestionText   *string `json:"question_text"`
	ResponseChoice *int64  `json:"response_choice"`
	ResponseText   *string `json:"response_text"`
}
