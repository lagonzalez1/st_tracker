package models

import "database/sql"

type User struct {
	ID             int64    `json:"id"`
	Email          string   `json:"email"`
	Password       string   `json:"password"`
	Type           string   `json:"type"`
	Permissions    []string `json:"permissions"`
	OrganizationId *int64   `json:"organization_id"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Type     string `json:"type"`
}

type LoginResponse struct {
	User         User                     `json:"user"`
	Token        string                   `json:"token"`
	RefreshToken string                   `json:"refresh_token"`
	Permissions  LoginResponsePermissions `json:"permissions"`
}

type LoginResponsePermissions struct {
	DisableUpdate bool `json:"disable_update"`
	DisableCreate bool `json:"disable_create"`
	DisableDelete bool `json:"disable_delete"`
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
	ID         *int64 `json:"id"`
	FirstName  string `json:"firstname"`
	MiddleName string `json:"middle_name"`
	LastName   string `json:"last_name"`
	Email      string `json:"email"`
	GradeLevel int    `json:"grade_level"`
	Active     bool   `json:"active"`
	CreatedAt  string `json:"created_at"`
	LocationId int64  `json:"location_id"`
}

type ResponseRequestStudentList struct {
	ID         int64  `json:"id"`
	FirstName  string `json:"firstname"`
	MiddleName string `json:"middle_name"`
	LastName   string `json:"last_name"`
	Email      string `json:"email"`
	GradeLevel int    `json:"grade_level"`
	Active     bool   `json:"active"`
	CreatedAt  string `json:"created_at"`
	Period     int64  `json:"period"`
	SemesterId *int64 `json:"semester_id"`
	LocationId *int64 `json:"location_id"`
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
	ZipCode    int64  `json:"zip_code"`
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
	ZipCode        int64  `json:"zip_code"`
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
	OrganizationId *int64 `json:"organization_id"`
}
type ResponseRequestAdminList struct {
	ID       int64  `json:"id"`
	Fullname string `json:"fullname"`
	Email    string `json:"email"`
	Region   string `json:"region"`
	State    string `json:"state"`
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
	ID          int64  `json:"id"`
	ProgramName string `json:"program_name"`
	CreatedAt   string `json:"created_at"`
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
	ID             *int64 `json:"id"`
	ProgramName    string `json:"program_name"`
	AdminID        int64  `json:"admin_id"`
	OrganizationId *int64 `json:"organization_id"`
}

type ResponseRequestProgram struct {
	Status    string `json:"status"`
	ProgramId int64  `json:"id"`
}

type ResponseRequestMaterialsList struct {
	ID           int64  `json:"id"`
	Title        string `json:"title"`
	ExternalLink string `json:"external_link"`
	Description  string `json:"description"`
	Version      string `json:"version"`
	Pre          bool   `json:"pre"`
	Mid          bool   `json:"mid"`
	Post         bool   `json:"post"`
	Visable      bool   `json:"visable"`
	CreatedAt    string `json:"created_at"`
}

type RegisterRequestMaterials struct {
	ID             *int64 `json:"id"`
	Title          string `json:"title"`
	ExternalLink   string `json:"external_link"`
	Description    string `json:"description"`
	Version        string `json:"version"`
	Pre            bool   `json:"pre"`
	Mid            bool   `json:"mid"`
	Post           bool   `json:"post"`
	Visable        bool   `json:"visable"`
	CreatedAt      string `json:"created_at"`
	LocationId     *int64 `json:"location_id"`
	OrganizationId *int64 `json:"organization_id"`
}

type ResponseRequestMaterials struct {
	Status     string `json:"status"`
	MaterialId int64  `json:"id"`
}

type ResponseRequestTutorsList struct {
	ID         int64  `json:"id"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	Email      string `json:"email"`
	CreatedAt  string `json:"created_at"`
	LocationId *int64 `json:"location_id"`
}

type RegisterRequestTutor struct {
	ID             int64  `json:"id"`
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	Password       string `json:"password"`
	Email          string `json:"email"`
	CreatedAt      string `json:"created_at"`
	LocationId     *int64 `json:"location_id"`
	OrganizationId *int64 `json:"organization_id"`
}

type ResponseRequestTutor struct {
	Status  string `json:"status"`
	TutorId int64  `json:"id"`
}

type ResponseRequestSemesterList struct {
	ID    int64  `json:"id"`
	Title string `json:"title"`
	Year  int64  `json:"year"`
}

type RegisterRequestSemester struct {
	ID             *int64 `json:"id"`
	Title          string `json:"title"`
	Year           *int64 `json:"year"`
	OrganizationId *int64 `json:"organization_id"`
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
	Status string `json:"status"`
}

type RemoveAdmin struct {
	ID *int64 `json:"id"`
}

type RemoveResponse struct {
	Status string `json:"status"`
}

type RemoveRequest struct {
	ID *int64 `json:"id"`
}

type RegisterStudentSessionList struct {
	Session     RegisterTutorSession     `json:"session"`
	SessionList []RegisterStudentSession `json:"student_sessions"`
}

type RegisterTutorSession struct {
	ID           *int64         `json:"id"`
	StudentCount *int64         `json:"student_count"`
	LocationId   *int64         `json:"location_id"`
	SubstituteId sql.NullInt64  `json:"substitute_id"`
	ProgramId    sql.NullInt64  `json:"program_id"`
	Notes        string         `json:"notes"`
	SessionDate  string         `json:"session_date"`
	StartTime    sql.NullString `json:"start_time"`
	Subject      string         `json:"subject"`
	Substitute   sql.NullBool   `json:"substitute"`
	TutorId      sql.NullInt64  `json:"tutor_id"`
	CreatedAt    string         `json:"created_at"`
	EditedAt     string         `json:"edited_at"`
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

type ResponseRegisterSubject struct {
	Status    string `json:"status"`
	SubjectID int64  `json:"id"`
}

type ServiceSession struct {
	FirstName    string         `json:"first_name"`
	LastName     string         `json:"last_name"`
	ID           *int64         `json:"id"`
	Location     string         `json:"location_name"`
	SubstituteId sql.NullInt64  `json:"substitute_id"`
	ProgramId    sql.NullInt64  `json:"program_id"`
	ProgramName  string         `json:"program_name"`
	Notes        string         `json:"notes"`
	SessionDate  string         `json:"session_date"`
	StartTime    sql.NullString `json:"start_time"`
	Subject      string         `json:"subject"`
	Substitute   sql.NullBool   `json:"substitute"`
	TutorId      sql.NullInt64  `json:"tutor_id"`
	CreatedAt    string         `json:"created_at"`
	EditedAt     string         `json:"edited_at"`
	SubjectName  string         `json:"subject_name"`
}

type RegisterStudentSession struct {
	ID              *int64 `json:"id"`
	FirstName       string `json:"first_name"`
	LastName        string `json:"last_name"`
	SessionDate     string `json:"session_date"`
	Duration        *int64 `json:"duration"`
	StartTime       string `json:"start_time"`
	Notes           string `json:"notes"`
	OrganizationId  *int64 `json:"organization_id"`
	ProgramId       *int64 `json:"program_id"`
	LocationId      *int64 `json:"location_id"`
	TutorId         *int64 `json:"tutor_id"`
	Subject         string `json:"subject"`
	AssessmentId    *int64 `json:"assessment_id"`
	AssessmentScore *int64 `json:"score"`
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
	ID              *int64 `json:"id"`
	Title           string `json:"title"`
	Description     string `json:"description,omitempty"`
	Letter          string `json:"letter"`
	Cycle           *int64 `json:"cycle"`
	AlphaIdentifier string `json:"alpha_identifier,omitempty"`
	ExternalLink    string `json:"external_link,omitempty"`
	MaxScore        *int64 `json:"max_score,omitempty"`
	Subject         string `json:"subject,omitempty"`
	MaterialID      *int64 `json:"material_id,omitempty"`
	OrganizationID  *int64 `json:"organization_id"`
	CreatedAt       string `json:"created_at"`
}

type RegisterAssessment struct {
	ID              *int64 `json:"id"`
	Title           string `json:"title"`
	Description     string `json:"description,omitempty"`
	Letter          string `json:"letter"`
	Cycle           *int64 `json:"cycle"`
	Visable         bool   `json:"visable"`
	AlphaIdentifier string `json:"alpha_identifier,omitempty"`
	ExternalLink    string `json:"external_link,omitempty"`
	MaxScore        *int64 `json:"max_score,omitempty"`
	Subject         string `json:"subject,omitempty"`
	OrganizationID  *int64 `json:"organization_id"`
	MaterialID      *int   `json:"material_id,omitempty"`
	ProgramId       *int64 `json:"program_id"`
	CreatedAt       string `json:"created_at"`
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

type RemoveLocationProgram struct {
	ProgramId      *int64 `json:"program_id"`
	LocationId     *int64 `json:"location_id"`
	OrganizationID *int64 `json:"organization_id"`
}

type RegisterResponseLocationProgram struct {
	Status string `json:"status"`
}

type SearchQuery struct {
	SearchTerm string `json:"search_term,omitempty"`
	LocationId *int64 `json:"location_id,omitempty"`
	ProgramId  *int64 `json:"program_id,omitempty"`
	DateStart  string `json:"date,omitempty"`
	DateEnd    string `json:"date_end,omitempty"`
	SubjectId  *int64 `json:"subject_id,omitempty"`
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
