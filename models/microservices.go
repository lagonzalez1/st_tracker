package models

type RequestEventGeneration struct {
	GenerationType              *string                      `json:"generate_type"`
	RequestQuestions            *RequestQuestions            `json:"generate_questions"`
	RequestQuestionsDoMaterials *RequestQuestionsDoMaterials `json:"generate_questions_do_materials"`
	RequestMaterials            *RequestMaterials            `json:"generate_materials"`
	OrganizationID              *int64                       `json:"organization_id"`
}

type RequestQuestions struct {
	S3OutputKey        *string `json:"s3_output_key"`
	DistrictID         *int64  `json:"district_id"`
	SubjectID          *int64  `json:"subject_id"`
	Description        *string `json:"description"`
	Difficulty         *string `json:"difficulty"`
	GradeLevel         *int64  `json:"grade_level"`
	MaxPoints          *int64  `json:"max_points"`
	QuestionsCount     *int64  `json:"question_count"`
	CustomInstructions *string `json:"custom_instructions"`
}
type RequestQuestionsDoMaterials struct {
	S3OutputKey        *string `json:"s3_output_key"`
	DistrictID         *int64  `json:"district_id"`
	SubjectID          *int64  `json:"subject_id"`
	MaterailId         *int64  `json:"material_id"`
	Description        *string `json:"description"`
	Difficulty         *string `json:"difficulty"`
	GradeLevel         *int64  `json:"grade_level"`
	MaxPoints          *int64  `json:"max_points"`
	QuestionsCount     *int64  `json:"question_count"`
	CustomInstructions *string `json:"custom_instructions"`
}

type RequestMaterials struct {
	S3OutputKey        *string `json:"s3_output_key"`
	AssessmentId       *int64  `json:"assessment_id"`
	BiasType           *string `json:"bias_type"`
	CustomInstructions *string `json:"custom_instructions"`
}

type RemoveGeneratedQuestion struct {
	InputKey       *string `json:"input_key"`
	OrganizationID int64   `json:"organization_id"`
}

type RequestSurveyPayload struct {
	SessionID      *int64 `json:"session_id"`
	OrganizationID *int64 `json:"organization_id"`
}

type RequestStudentReport struct {
	StudentID   *int64  `json:"student_id"`
	SemesterID  *int64  `json:"semester_id"`
	S3OutputKey *string `json:"s3_output_key"`
}

type StudentReport struct {
	GeneratedAt                   int `json:"generated_at"`
	AllScores                     any `json:"all_scores"`
	SubjectScores                 any `json:"subject_scores"`
	ScoresLinearRegression        any `json:"scores_linear_regression"`
	SubjectScoresLinearRegression any `json:"subject_scores_linear_regression"`
}
