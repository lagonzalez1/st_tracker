package services

import (
	"context"
	"fmt"
	"tracker/app/models"
)

/**
CREATE TABLE stu_tracker.Pre_assessment_questionnaire (
    student_id INT REFERENCES stu_tracker.Students(id) ON DELETE CASCADE,
    assessment_id INT REFERENCES stu_tracker.Assessments(id) ON DELETE CASCADE,
    session_token UUID NOT NULL,
    sleep_hours FLOAT DEFAULT 0,
    effort_score smallint NOT NULL CHECK (effort_score BETWEEN 0 AND 10),
    tutor_sessions INT DEFAULT 0,
    parental_help smallint NOT NULL CHECK (parental_help BETWEEN 0 AND 3),
    sports_hours INT DEFAULT 0,
    peer_influence smallint NOT NULL CHECK (peer_influence BETWEEN 0 AND 3),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

**/

func (s *AuthService) CreatePreAssessment(ctx context.Context, req models.RegisterPreAssessment) (*models.ResponsePreAssessment, error) {
	// Input validation
	if req.StudentId == nil || *req.SessionToken == "" || req.AssessmentId == nil {
		return nil, fmt.Errorf("missing required fields")
	}
	var newID int64
	query := `INSERT INTO stu_tracker.Pre_assessment_questionnaire (
				student_id, assessment_id, session_token, sleep_hours,
				effort_score, tutor_sessions, parental_help, sports_hours, peer_influence) 
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id;`
	err := s.db.QueryRow(query, *req.StudentId, *req.AssessmentId, req.SessionToken, req.SleepHours, req.EffortScore,
		req.TutorSessions, req.ParentalHelp, req.SportHours, req.PeerInfluence).Scan(&newID)
	if err != nil {
		return nil, fmt.Errorf("failed to insert student: %w", err)
	}
	return &models.ResponsePreAssessment{
		Status: "OK",
		ID:     newID,
	}, nil
}
