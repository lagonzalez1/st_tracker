package models

type AssessmentGrader struct {
	QuestionID *int64 `json:"question_id"`
	Points     int    `json:"points"`
	ChoiceID   *int64 `json:"choice_id"`
	IsCorrect  bool   `json:"is_correct"`
}

type AssessmentScore struct {
	Points          int              `json:"points"`
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

type DeleteAssessmentSession struct {
	SessionID string `json:"session_id"`
}

type DeleteAssessmentSessionResponse struct {
	Status string `json:"status"`
}

type StudentAssessmentSession struct {
	ID              *int64 `json:"id"`
	FirstName       string `json:"first_name"`
	LastName        string `json:"last_name"`
	AssessmentID    *int64 `json:"assessment_id"`
	IsActive        bool   `json:"is_active"`
	Completed       bool   `json:"completed"`
	GradeAssessment bool   `json:"grade_assessment"`
	AssessmentTitle string `json:"assessment_title"`
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
	Answers      map[string]interface{} `json:"answers"`
	AssessmentID *int64                 `json:"assessment_id"`
	TutorID      *int64                 `json:"tutor_id"`
	SessionID    string                 `json:"session_id"`
	StudentID    *int64                 `json:"student_id"`
}

type StudentAssessmentChoices struct {
	QuestionID *int64  `json:"question_id"`
	ChoiceID   *int64  `json:"choice_id"`
	AnswerText *string `json:"answer_text"`
}
