package models

// Simplified StudyGuide struct without database tags
type Materials struct {
	GuideType           string               `json:"guide_type"`
	Subject             string               `json:"subject"`
	GradeLevel          string               `json:"grade_level"`
	DurationMinutes     int                  `json:"duration_minutes"`
	LearningObjectives  []string             `json:"learning_objectives"`
	KeyConcepts         []KeyConcept         `json:"key_concepts"`
	Activities          []Activity           `json:"activities"`
	AssessmentQuestions []AssessmentQuestion `json:"assessment_questions"`
	Summary             string               `json:"summary"`
	MaterialsNeeded     []string             `json:"materials_needed"`
	Appendix            *string              `json:"appendix,omitempty"`
}

type KeyConcept struct {
	Title       string   `json:"title"`
	Explanation string   `json:"explanation"`
	Examples    []string `json:"examples"`
}

type Activity struct {
	Title           string   `json:"title"`
	Description     string   `json:"description"`
	Steps           []string `json:"steps"`
	ExpectedOutcome string   `json:"expected_outcome"`
}

type AssessmentQuestion struct {
	Question   string `json:"question"`
	Answer     string `json:"answer"`
	Difficulty string `json:"difficulty"`
}
