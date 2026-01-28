package models

import "time"

type AssessmentGrader struct {
	QuestionID   *int64  `json:"question_id"`
	Points       int     `json:"points"`
	ChoiceID     *int64  `json:"choice_id"`
	IsCorrect    bool    `json:"is_correct"`
	QuestionType *string `json:"question_type"`
}

type AssessmentScore struct {
	Points          float32          `json:"points"`
	Correct         int              `json:"correct"`
	MaxScore        int              `json:"total"`
	QuestionEntries []AnswerFeedback `json:"correct_map"`
}

type AnswerFeedback struct {
	QuestionID int64   `json:"question_id"`
	ChoiceID   *int64  `json:"choice_id,omitempty"`
	IsCorrect  bool    `json:"is_correct"`
	AnswerText *string `json:"answer_text,omitempty"`
}

type RegisterStudentAssessmentSession struct {
	StudentAssessmentSession []StudentAssesmentSession `json:"students"`
}

type StudentAssesmentSession struct {
	ID           *int64  `json:"id"`
	FirstName    string  `json:"first_name"`
	LastName     string  `json:"last_name"`
	ProgramId    *int64  `json:"program_id"`
	LocationId   *int64  `json:"location_id"`
	TutorId      *int64  `json:"tutor_id"`
	SubjectId    *int64  `json:"subject_id"`
	AssessmentId *int64  `json:"assessment_id"`
	SemesterId   *int64  `json:"semester_id"`
	EasyScoreID  bool    `json:"easy_score"`
	JoinCode     *string `json:"join_code"`
}

type DeleteAssessmentSession struct {
	SessionID string `json:"session_id"`
}

type DeleteStudentSession struct {
	SessionID string `json:"session_id"`
	StudentID *int64 `json:"student_id"`
}

type DeleteStudentSessionResponse struct {
	Status string `json:"status"`
}

type DeleteAssessmentSessionResponse struct {
	Status string `json:"status"`
}

type StudentAssessmentSession struct {
	ID              *int64  `json:"id"`
	FirstName       string  `json:"first_name"`
	LastName        string  `json:"last_name"`
	AssessmentID    *int64  `json:"assessment_id"`
	IsActive        bool    `json:"is_active"`
	Completed       bool    `json:"completed"`
	JoinCode        *string `json:"join_code"`
	GradeAssessment bool    `json:"grade_assessment"`
	AssessmentTitle string  `json:"assessment_title"`
}

type ResponseStudentAssessmentSession struct {
	Status         string `json:"status"`
	SessionsActive int    `json:"sessions_active"`
	SessionID      string `json:"session_id"`
}

type StudentSubmitResponse struct {
	Status  string `json:"status"`
	Answers int    `json:"answered"`
}

type RegisterStudentAssessment struct {
	Answers         map[string]interface{} `json:"answers"`
	AssessmentID    *int64                 `json:"assessment_id"`
	TutorID         *int64                 `json:"tutor_id"`
	SessionID       string                 `json:"session_id"`
	StudentID       *int64                 `json:"student_id"`
	QuestionnaireID *int64                 `json:"questionnaire_id"`
}

type StudentAssessmentChoices struct {
	QuestionID   *int64  `json:"question_id"`
	ChoiceID     *int64  `json:"choice_id"`
	AnswerText   *string `json:"answer_text"`
	AssessmentID *int64  `json:"assessment_id"`
}

type SQSAssessmentPayload struct {
	ID           *int64                 `json:"id"`
	FirstName    string                 `json:"first_name"`
	LastName     string                 `json:"last_name"`
	AssessmentID *int64                 `json:"assessment_id"`
	Answers      map[string]interface{} `json:"answers"`
}

type Choice struct {
	ChoiceID    *int   `json:"choice_id,omitempty" db:"choice_id"`
	ChoiceText  string `json:"choice_text" db:"choice_text"`
	IsCorrect   bool   `json:"is_correct" db:"is_correct"`
	OrderNumber int    `json:"order_number" db:"order_number"`
}

// Question represents a single assessment question
type Question struct {
	QuestionID   *int     `json:"question_id,omitempty" db:"question_id"`
	StandardText string   `json:"standard_text" db:"standard_text"`
	ImageURL     *string  `json:"image_url,omitempty" db:"image_url"`
	QuestionText string   `json:"question_text" db:"question_text"`
	QuestionType string   `json:"question_type" db:"question_type"`
	Points       float64  `json:"points" db:"points"`
	OrderNumber  int      `json:"order_number" db:"order_number"`
	IsRequired   bool     `json:"is_required" db:"is_required"`
	Choices      []Choice `json:"choices" db:"choices"`

	// Optional metadata fields
	CreatedAt *time.Time `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt *time.Time `json:"updated_at,omitempty" db:"updated_at"`
}

// Assessment represents the full assessment with questions
type Assessment struct {
	Questions []Question `json:"questions"`
}
