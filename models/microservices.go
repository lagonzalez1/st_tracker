package models

type RequestQuestions struct {
	S3OutputKey    *string `json:"s3_output_key"`
	DistrictID     *int64  `json:"district_id"`
	SubjectID      *int64  `json:"subject_id"`
	OrganizationID *int64  `json:"organization_id"`
	Description    *string `json:"description"`
	Difficulty     *string `json:"difficulty"`
	GradeLevel     *int64  `json:"grade_level"`
	MaxPoints      *int64  `json:"max_points"`
	QuestionsCount *int64  `json:"questions_count"`
}

type RemoveGeneratedQuestion struct {
	InputKey       *string `json:"input_key"`
	OrganizationID int64   `json:"organization_id"`
}
