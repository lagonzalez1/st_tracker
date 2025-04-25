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
	Total           int              `json:"total"`
	QuestionEntries []AnswerFeedback `json:"correct_map"`
}

type AnswerFeedback struct {
	QuestionID int64 `json:"question_id"`
	ChoiceID   int64 `json:"choice_id"`
	IsCorrect  bool  `json:"is_correct"`
}
