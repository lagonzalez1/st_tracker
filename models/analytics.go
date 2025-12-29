package models

import (
	"database/sql"
	"time"
)

type RequestSessionAnalytics struct {
	ID             *int      `json:"id"`
	Email          string    `json:"email"`
	OrganizationID *int      `json:"organization_id"`
	SemesterID     *int      `json:"semester_id"`
	LocationID     *int      `json:"location_id"`
	ProgramID      *int      `json:"program_id"`
	StartDate      time.Time `json:"start_date"`
	EndDate        time.Time `json:"end_date"`
}

type SessionAnalytics struct {
	SessionCount    *int     `json:"session_count"`
	AssessmentCount *int     `json:"assessment_count"`
	StudentCount    *int     `json:"student_count"`
	SessionDuration *float32 `json:"session_duration"`
}

type SessionsAnalyticsLocal struct {
	SessionCount    *int     `json:"session_count"`
	AssessmentCount *int     `json:"assessment_count"`
	SessionDuration *float32 `json:"session_duration"`
}

type RequestSessionBChart struct {
	ID             *int64    `json:"id"`
	Email          string    `json:"email"`
	OrganizationID *int64    `json:"organization_id"`
	SemesterID     *int64    `json:"semester_id"`
	LocationID     *int64    `json:"location_id"`
	ProgramID      *int64    `json:"program_id"`
	StartDate      time.Time `json:"start_date"`
	EndDate        time.Time `json:"end_date"`
}

type RequestCycleGrowth struct {
	ID             *int64    `json:"id"`
	Email          string    `json:"email"`
	OrganizationID *int64    `json:"organization_id"`
	SemesterID     *int64    `json:"semester_id"`
	LocationID     *int64    `json:"location_id"`
	ProgramID      *int64    `json:"program_id"`
	StartDate      time.Time `json:"start_date"`
	EndDate        time.Time `json:"end_date"`
}

type ResponseSessionBChart struct {
	SessionDate           time.Time `json:"session_date"`
	TotalSessions         int       `json:"total_sessions"`
	StudentAverage        float32   `json:"student_average"`
	TotalStudents         int       `json:"total_students"`
	MinStudents           int       `json:"min_students"`
	MaxStudents           int       `json:"max_students"`
	StudentAverageRounded float32   `json:"student_average_rounded"`
	Month                 string    `json:"month"`
}

type ResponseAbsentPresent struct {
	Present *int64 `json:"present_count"`
	Absent  *int64 `json:"absent_count"`
}

type ResponseAssessmentCompletion struct {
	EnrolledCount   *int64 `json:"enrolled_count"`
	AssessedCount   *int64 `json:"assessed_count"`
	PostAssessed    *int64 `json:"post_assessed_count"`
	IntremAssessend *int64 `json:"intrem_assessed_count"`
}

type ResponseAssessmentVScore struct {
	StudentID    *int64   `json:"student_id"`
	FirstName    *string  `json:"first_name"`
	LastName     *string  `json:"last_name"`
	SessionCount *int64   `json:"session_count"`
	ScoreAverage *float32 `json:"score_average"`
}

type ResponseStudentVAssessments struct {
	StudentID       *int64   `json:"student_id"`
	SessionID       *int64   `json:"session_id"`
	FirstName       *string  `json:"first_name"`
	LastName        *string  `json:"last_name"`
	Status          *string  `json:"status"`
	AssessmentTitle *string  `json:"assessment_title"`
	AssessmentID    *int64   `json:"assessment_id"`
	Score           *float32 `json:"score"`
	MaxScore        *int64   `json:"max_score"`
}

type ResponseSentiment struct {
	ProgramName    *string   `json:"program_name"`
	SessionDate    time.Time `json:"session_date"`
	TutorID        *int64    `json:"tutor_id"`
	FirstName      *string   `json:"first_name"`
	QuestionText   *string   `json:"question_text"`
	ResponseText   *string   `json:"response_text"`
	SentimentScore *float32  `json:"sentiment_score"`
	SessionID      *int64    `json:"session_id"`
}

type ResponseAssessments struct {
	StudentID        *int64    `json:"student_id"`
	StudentFirstName *string   `json:"student_first_name"`
	ProgramName      *string   `json:"program_name"`
	Assessment       *string   `json:"assessment"`
	SessionDate      time.Time `json:"session_date"`
	TutorID          *int64    `json:"tutor_id"`
	SessionID        *int64    `json:"session_id"`
	MaxScore         *float32  `json:"max_score"`
	Score            *float32  `json:"score"`
	ProgramID        *int64    `json:"program_id"`
}

type ResponseCycleGrowthDelim struct {
	LocationID        *int64   `json:"location_id"`
	ProgramID         *int64   `json:"program_id"`
	ProgramName       *string  `json:"program_name"`
	LocationName      *string  `json:"location_name"`
	AverageScore      *float32 `json:"average_score"`
	MaxScore          *float32 `json:"max_score"`
	MinScore          *float32 `json:"min_score"`
	Pre               bool     `json:"pre"`
	Mid               bool     `json:"mid"`
	Post              bool     `json:"post"`
	StandardDeviation *float32 `json:"standard_deviation"`
	TotalAssessments  *int64   `json:"total_assessments"`
	Cycle             *int64   `json:"cycle"`
}

type ResponseCycleGrowth struct {
	LocationID        *int64   `json:"location_id"`
	ProgramID         *int64   `json:"program_id"`
	ProgramName       *string  `json:"program_name"`
	LocationName      *string  `json:"location_name"`
	AverageScore      *float32 `json:"average_score"`
	MaxScore          *float32 `json:"max_score"`
	MinScore          *float32 `json:"min_score"`
	StandardDeviation *float32 `json:"standard_deviation"`
	TotalAssessments  *int64   `json:"total_assessments"`
	Cycle             *int64   `json:"cycle"`
}

type ResponseAssessmentsBChart struct {
	LocationName      string  `json:"location_name"`
	AssessemtsTotal   int     `json:"assessments_total"`
	AssessemtsAverage float32 `json:"assessments_average"`
	MinScore          float32 `json:"min_score"`
	MaxScore          float32 `json:"max_score"`
	AssessmentName    string  `json:"assessment_name"`
	AssessmentCycle   string  `json:"cycle"`
	AssessmentLetter  string  `json:"letter"`
}

type ResponseProgramsBChart struct {
	ProgramName string `json:"program_name"`
	Count       int    `json:"count"`
}

type ResponseTutorsBChart struct {
	ID              int64   `json:"id"`
	StudentCount    int     `json:"total_student_count"`
	AverageStudents float32 `json:"average_student_count"`
	TotalSessions   int     `json:"total_sessions"`
	AssessmentCount int     `json:"assessments_count"`
	FirstName       string  `json:"first_name"`
	LastName        string  `json:"last_name"`
}

type ResponseAssessmentTrendline struct {
	ID              int `json:"id"`
	AssessmentCount int `json:"assessment_count"`
	Year            int `json:"year"`
	Month           int `json:"month"`
}
type ResponseSessionTrendline struct {
	ID           int `json:"id"`
	SessionCount int `json:"session_total"`
	Year         int `json:"year"`
	Month        int `json:"month"`
}

type RequestSemestersVAssessmentChart struct {
	OrganizationID *int64 `json:"organization_id"`
	Semester1ID    *int64 `json:"semester1_id"`
	Semester2ID    *int64 `json:"semester2_id"`
	ProgramID      *int64 `json:"program_id"`
	Assessment1ID  *int64 `json:"assessment1_id"`
	LocationID     *int64 `json:"location_id"`
}

type RequestAssessmentGrowth struct {
	OrganizationID *int64 `json:"organization_id"`
	Semester1ID    *int64 `json:"semester1_id"`
	ProgramID      *int64 `json:"program_id"`
	Assessment2ID  *int64 `json:"assessment2_id"`
	Assessment1ID  *int64 `json:"assessment1_id"`
	LocationID     *int64 `json:"location_id"`
}

// How to return {SingleDate, sem1, sem2, scores...}
type ResponseSemestersVAssessmentChart struct {
	DataSet1 []AssessmentComparison `json:"data1_set"`
}

type SemestersVAssessmentChart struct {
	Date         time.Time `json:"date"`
	Score        int       `json:"score"`
	AverageScore float32   `json:"average_score"`
	MinScore     int       `json:"min_score"`
	MaxScore     int       `json:"max_score"`
}

type AssessmentComparison struct {
	AssessmentS1        int     `json:"assessment_s1"`
	AssessmentS2        int     `json:"assessment_s2"`
	CountS1             int     `json:"count_s1"`
	CountS2             int     `json:"count_s2"`
	MinScoreS1          int     `json:"min_score_s1"`
	MaxScoreS1          int     `json:"max_score_s1"`
	MinScoreS2          int     `json:"min_score_s2"`
	MaxScoreS2          int     `json:"max_score_s2"`
	Semester1Avg        float64 `json:"semester_1_avg"`
	Semester2Avg        float64 `json:"semester_2_avg"`
	ScoreDifference     float64 `json:"score_difference"`
	RateOfChangePercent float64 `json:"rate_of_change_percent"`
}

type ResponseAssessmentGrowth struct {
	DataSet1 AssessmentGrowth `json:"assessment_1"`
	DataSet2 AssessmentGrowth `json:"assessment_2"`
}

type AssessmentGrowth struct {
	Count        *int    `json:"count"`
	AssessmentId *int    `json:"assessment_id"`
	MinScore     *int    `json:"min_score"`
	MaxScore     *int    `json:"max_score"`
	Average      float32 `json:"average"`
}

type RequestTutorsSessions struct {
	ID             *int64 `json:"id"`
	Email          string `json:"email"`
	OrganizationID *int64 `json:"organization_id"`
	SemesterID     *int64 `json:"semester_id"`
	LocationID     *int64 `json:"location_id"`
	ProgramID      *int64 `json:"program_id"`
	SurveyRequired bool   `json:"survey_required"`
}

type StudentsSession struct {
	ID         *int64 `json:"id"`
	Email      string `json:"email"`
	FirstName  string `json:"first_name"`
	MiddleName string `json:"middle_name,omitempty"`
	LastName   string `json:"last_name"`
}

type ResponseTutorSessions struct {
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
	Subject      sql.NullInt32  `json:"subject"`
	Substitute   sql.NullBool   `json:"substitute"`
	TutorId      sql.NullInt64  `json:"tutor_id"`
	CreatedAt    time.Time      `json:"created_at"`
	EditedAt     string         `json:"edited_at"`
	SubjectName  string         `json:"subject_name"`
	StudentCount int64          `json:"student_count"`
}

type ResponseTutorLowPerformance struct {
	ID                  *int64   `json:"id"`
	TutorName           *string  `json:"fullname"`
	SessionCount        *int64   `json:"session_count"`
	UniqueStudentCount  *int64   `json:"unique_student_count"`
	AverageStudentScore *float64 `json:"avg_student_score"`
	SessionPercentile   *float64 `json:"session_percentile"`
	StudentPercentile   *float64 `json:"student_percentile"`
	PerformanceStatus   *string  `json:"performance_status"`
}

type RequestSentiment struct {
	ID             *int64    `json:"id"`
	Email          string    `json:"email"`
	OrganizationID *int64    `json:"organization_id"`
	SemesterID     *int64    `json:"semester_id"`
	LocationID     *int64    `json:"location_id"`
	ProgramID      *int64    `json:"program_id"`
	StartDate      time.Time `json:"start_date"`
	EndDate        time.Time `json:"end_date"`
}
