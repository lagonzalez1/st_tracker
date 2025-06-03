package models

import "time"

type RequestStudentData struct {
	LocationID *int64    `json:"location_id"`
	ProgramID  *int64    `json:"program_id"`
	SubjectID  *int64    `json:"subject_id"`
	SemesterID *int64    `json:"semester_id"`
	DateStart  time.Time `json:"date"`
	DateEnd    time.Time `json:"date_end"`
}

type RequestDownloadData struct {
	LocationID *int64    `json:"location_id"`
	ProgramID  *int64    `json:"program_id"`
	SubjectID  *int64    `json:"subject_id"`
	SemesterID *int64    `json:"semester_id"`
	DateStart  time.Time `json:"date"`
	DateEnd    time.Time `json:"date_end"`
	SortKey    string    `json:"sort_key"`
}
type AssessmentsData struct {
	Title           string    `json:"title"`
	MaxScore        int       `json:"max_score"`
	CreatedAt       time.Time `json:"created_at"`
	Letter          string    `json:"letter"`
	Cycle           string    `json:"cycle"`
	Pre             bool      `json:"pre"`
	Mid             bool      `json:"mid"`
	Post            bool      `json:"port"`
	Version         float32   `json:"version"`
	Score           float64   `json:"score"`
	StudentID       *int64    `json:"student_id"`
	StudentName     string    `json:"student_name"`
	StudentLastName string    `json:"student_last_name"`
	SessionID       *int64    `json:"session_id"`
}
type StudentData struct {
	FirstName string    `json:"student_first_name"`
	LastName  string    `json:"student_last_name"`
	Absent    bool      `json:"absent"`
	Duration  int       `json:"duration"`
	CreatedAt time.Time `json:"created_at"`
}

type StudentRow struct {
	StudentID      *int64            `json:"student_id"`
	Student        StudentData       `json:"student"`
	SessionData    []StudentSession  `json:"student_sessions"`
	AssessmentData []AssessmentsData `json:"assessment_data"`
}

type StudentSession struct {
	FirstName      string    `json:"student_first_name"`
	LastName       string    `json:"student_last_name"`
	SessionID      *int64    `json:"session_id"`
	StudentID      *int64    `json:"student_id"`
	Subject        string    `json:"subject"`
	Duration       int       `json:"duration"`
	CreatedAt      time.Time `json:"created_at"`
	StartTime      string    `json:"start_time"`
	SessionDate    time.Time `json:"session_date"`
	Notes          string    `json:"notes"`
	Program        string    `json:"program"`
	ProgramID      *int64    `json:"program_id"`
	Absent         bool      `json:"absent"`
	Grade          int       `json:"grade_level"`
	Timeframe      bool      `json:"timeframe"`
	TimeframeStart *string   `json:"timeframe_start"`
	TimeframeEnd   *string   `json:"timeframe_end"`
}

type TutorSessionData struct {
	SessionID     *int64    `json:"session_id"`
	TutorID       *int64    `json:"tutor_id"`
	TutorName     string    `json:"tutor_name"`
	TutorLastName string    `json:"tutor_last_name"`
	Substitute    bool      `json:"substitute"`
	StudentCount  int       `json:"student_count"`
	Duration      int       `json:"duration"`
	CreatedAt     time.Time `json:"created_at"`
	StartTime     string    `json:"start_time"`
	SessionDate   time.Time `json:"session_date"`
	Notes         string    `json:"notes"`
	Program       string    `json:"program"`
	ProgramID     *int64    `json:"program_id"`
}
