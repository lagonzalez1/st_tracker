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
	GradeLevel        int     `json:"grade_level"`
	Active            bool    `json:"active"`
	CreatedAt         string  `json:"created_at"`
	LocationId        *int64  `json:"location_id"`
	DirectPartnership bool    `json:"direct_partnership"`
	TutorID           *int64  `json:"tutor_id"`
	CreatedBy         string  `json:"created_by"`
}

type ResponseRequestStudentList struct {
	ID                int64   `json:"id"`
	FirstName         string  `json:"first_name"`
	MiddleName        string  `json:"middle_name"`
	LastName          string  `json:"last_name"`
	Email             *string `json:"email"`
	GradeLevel        int     `json:"grade_level"`
	Active            bool    `json:"active"`
	CreatedAt         string  `json:"created_at"`
	Period            *int64  `json:"period,omitempty"`
	SemesterId        *int64  `json:"semester_id"`
	DirectPartnership bool    `json:"direct_partnership"`
	LocationId        *int64  `json:"location_id"`
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

type PermissionsMap struct {
	ID             *int64 `json:"id"`
	PermissionName string `json:"permission_name"`
}

type RegisterPermissionRequest struct {
	ID             *int64   `json:"id"`
	Role           string   `json:"role"`
	User           string   `json:"user"`
	Permissions    []string `json:"permissions"`
	OrganizationId *int64   `json:"organization_id"`
}
type RegisterPermissionResponse struct {
	Status string `json:"status"`
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
	Version      float64 `json:"version"`
	Pre          bool    `json:"pre"`
	Mid          bool    `json:"mid"`
	Post         bool    `json:"post"`
	Visible      bool    `json:"visible"`
	ProgramId    *int64  `json:"program_id"`
	CreatedAt    string  `json:"created_at"`
}

type RegisterRequestMaterials struct {
	ID             *int64  `json:"id"`
	Title          string  `json:"title"`
	ExternalLink   string  `json:"external_link"`
	Description    string  `json:"description"`
	Version        float64 `json:"version"`
	Pre            bool    `json:"pre"`
	Mid            bool    `json:"mid"`
	Post           bool    `json:"post"`
	Visible        bool    `json:"visible"`
	CreatedAt      string  `json:"created_at"`
	LocationId     *int64  `json:"location_id"`
	ProgramId      *int64  `json:"program_id"`
	OrganizationId *int64  `json:"organization_id"`
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
	DateStart time.Time `json:"date_start"`
	DateEnd   time.Time `json:"date_end"`
	Active    bool      `json:"active"`
}

type ResponseRequestSemesterLocationList struct {
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

type RegisterRequestSemester struct {
	ID             *int64 `json:"id"`
	Title          string `json:"title"`
	Year           *int64 `json:"year"`
	OrganizationId *int64 `json:"organization_id"`
	DateStart      string `json:"date_start"`
	DateEnd        string `json:"date_end"`
	Active         bool   `json:"active"`
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
	Session        RegisterTutorSession     `json:"session"`
	SessionList    []RegisterStudentSession `json:"student_sessions"`
	OrganizationID *int64                   `json:"organization_id"`
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
	Duration     *int   `json:"duration"`
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
	SessionDate     string        `json:"session_date"`
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
	SessionCount    *int64 `json:"session_count"`
	AssessmentCount *int64 `json:"assessment_count"`
}

type RegisterStudentSession struct {
	ID              *int64 `json:"id"`
	Absent          bool   `json:"absent"`
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
	SubjectId       *int64 `json:"subject_id"`
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
}

type RegisterAssessment struct {
	ID              *int64  `json:"id"`
	Title           string  `json:"title"`
	Description     string  `json:"description,omitempty"`
	Letter          string  `json:"letter"`
	Cycle           *int64  `json:"cycle"`
	Visible         bool    `json:"visible"`
	AlphaIdentifier string  `json:"alpha_identifier,omitempty"`
	ExternalLink    string  `json:"external_link,omitempty"`
	MaxScore        *int64  `json:"max_score,omitempty"`
	SubjectId       *int64  `json:"subject_id,omitempty"`
	OrganizationID  *int64  `json:"organization_id"`
	MaterialID      *int    `json:"material_id,omitempty"`
	ProgramId       *int64  `json:"program_id"`
	CreatedAt       string  `json:"created_at"`
	Version         float64 `json:"version"`
	Pre             bool    `json:"pre"`
	Mid             bool    `json:"mid"`
	Post            bool    `json:"post"`
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
	ID         *int64  `json:"student_id"`
	FirstName  string  `json:"first_name"`
	LastName   string  `json:"last_name"`
	Email      *string `json"email"`
	MiddleName string  `json:"middle_name"`
	Duration   *int64  `json:"duration"`
	Period     *int64  `json:"period,omitempty"`
	Grade      *int64  `json:"grade"`
}

type AssessmentInfoStudent struct {
	Title     string    `json:"title"`
	Letter    string    `json:"letter"`
	Cycle     string    `json:"cycle"`
	MaxScore  *int64    `json:"max_score"`
	Subject   *int64    `json:"subject"`
	Score     *int64    `json:"score"`
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
	CreatedAt time.Time `json:"created_at"`
	Score     *int64    `json:"score"`
	MaxScore  *int64    `json:"max_score"`
	Letter    string    `json:"letter"`
	Cycle     string    `json:"cycle"`
	SessionID *int64    `json:"session_id"`
	Title     string    `json:"title"`
	Pre       bool      `json:"pre"`
	Mid       bool      `json:"mid"`
	Post      bool      `json:"post"`
	Version   float64   `json:"version"`
}

type SessionTrail struct {
	ID              int64     `json:"id"`
	FirstName       string    `json:"first_name"`
	LastName        string    `json:"last_name"`
	ProgramName     string    `json:"program_name"`
	SessionDuration int       `json:"session_duration"`
	StartTime       string    `json:"start_time"`
	Notes           string    `json:"notes"`
	CreatedAt       time.Time `json:"created_at"`
	StudentCount    int       `json:"student_count"`
	Substitute      bool      `json:"substitute"`
	Absent          bool      `json:"absent"`
	StudentDuration int       `json:"student_duration"`
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
	ID             int    `json:"id"`
	Title          string `json:"title"`
	Body           string `json:"body"`
	CreatedAt      string `json:"created_at"`
	LocationID     []*int `json:"location_id"`     // Nullable
	Severity       string `json:"severity"`        // Default: "info"
	OrganizationID int    `json:"organization_id"` // Required
	ProgramID      *int   `json:"program_id"`      // Nullable
	AdminID        *int   `json:"admin_id"`        // Required
	StaffID        *int   `json:"staff_id"`
}

type RegisterUpdateAnnouncements struct {
	ID             *int   `json:"id"`
	Title          string `json:"title"`
	Body           string `json:"body"`
	CreatedAt      string `json:"created_at"`
	LocationID     *int   `json:"location_id"`     // Nullable
	Severity       string `json:"severity"`        // Default: "info"
	OrganizationID int    `json:"organization_id"` // Required
	ProgramID      *int   `json:"program_id"`      // Nullable
	AdminID        *int   `json:"admin_id"`        // Required
	StaffID        *int   `json"staff_id"`
}

type ResponseAnnouncement struct {
	ID     int64  `json:"id"`
	Status string `json:"status"`
}

type AnnouncementsList struct {
	ID             int    `json:"id"`
	Title          string `json:"title"`
	Body           string `json:"body"`
	CreatedAt      string `json:"created_at"`
	LocationID     *int   `json:"location_id"`     // Nullable
	Severity       string `json:"severity"`        // Default: "info"
	OrganizationID int    `json:"organization_id"` // Required
	ProgramID      *int   `json:"program_id"`      // Nullable
	AdminID        *int   `json:"admin_id"`        // Required
	StaffID        *int   `json:"staff_id"`
	StaffName      string `json:"staff_name"`
	LocationName   string `json:"location_name"`
	ProgramName    string `json:"program_name"`
}

type RegisterSubjectLocation struct {
	ID             int  `json:"id"`
	SubjectID      *int `json:"subject_id"`
	LocationID     *int `json:"location_id"`
	OrganizationID *int `json:"organization_id"` // Required
}

type RemoveSubjectLocation struct {
	SubjectID  *int `json:"subject_id"`
	LocationID *int `json:"location_id"`
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
	SessionID  *int64 `json:"session_id"`
	StudentID  *int64 `json:"student_id"`
	FirstName  string `json:"first_name"`
	MiddleName string `json:"middle_name"`
	LastName   string `json:"last_name"`
	Grade      int    `json:"grade_level"`
}

type TutorSessionsList struct {
	SessionID    *int64     `json:"session_id"`
	ProgramName  string     `json:"program_name"`
	LocationID   *int64     `json:"location_id"`
	SemesterID   *int64     `json:"semester_id"`
	ProgramID    *int64     `json:"program_id"`
	SessionDate  time.Time  `json:"session_date"`
	StartTime    string     `json:"start_time"`
	Students     []Students `json:"students"`
	StudentCount int        `json:"student_count"`
	InSchool     bool       `json:"in_school"`
	Substitute   bool       `json:"substitute"`
	Semester     string     `json:"semester"`
	LocationName string     `json:"location_name"`
}
