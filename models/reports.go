package models

type StudentReportResponse struct {
	Status     *string            `json:"status"`
	OutputKey  *string            `json:"s3_output_key"`
	JsonReport *StudentReportData `json:"json_report"`
}

// StudentReportData represents the full JSONB content from the database
type StudentReportData struct {
	AllScores            AllScores             `json:"all_scores"`
	GeneratedAt          float64               `json:"generated_at"`
	SubjectBias          *[]SubjectBias        `json:"subject_bias"`
	LearningDisability   *[]LearningDisability `json:"learning_disability"` // Pointer for nulls
	AssessmentComparison *[]AssessmentCompare  `json:"assessment_comparison"`
	LDLinearRegression   LDLinearRegression    `json:"learning_disability_linear_regression"`
}

type AllScores struct {
	Data   []float32      `json:"data"`
	Labels []string       `json:"labels"`
	Scores MovingAverages `json:"scores"`
}

type LearningDisability struct {
	Attendance       *float32 `json:"attendance"`
	PreviousScores   *float32 `json:"previous_Scores"`
	ExamScore        *float32 `json:"exam_score"`
	TutoringSessions *float32 `json:"tutoring_sessions"`
}

type MovingAverages struct {
	CMA []float64 `json:"CMA"`
	EMA []float64 `json:"EMA"`
	SMA []float64 `json:"SMA"`
}

type SubjectBias struct {
	Mean          float64 `json:"mean"`
	Subject       string  `json:"subject"`
	PercentChange float64 `json:"percent_change"`
}

type AssessmentCompare struct {
	Mid             float64 `json:"mid"`
	Pre             float64 `json:"pre"`
	Post            float64 `json:"post"`
	AlphaIdentifier string  `json:"alpha_identifier"`
}

type LDLinearRegression struct {
	ScoresLR *[]LDLinearRegressionResults `json:"scores_linear_regression"`
}

type LDLinearRegressionResults struct {
	Prediction *float32 `json:"prediction"`
	Actual     *float32 `json:"actual"`
	Title      *string  `json:"title"`
}
