package services

import (
	"fmt"
	"strings"
	"tracker/app/models"
)

/**
tutor_id INT REFERENCES stu_tracker.Tutors(id) ON DELETE CASCADE,
    session_date TIMESTAMP NOT NULL,
    location_id INT REFERENCES stu_tracker.Locations(id) ON DELETE SET NULL,
    substitute BOOLEAN DEFAULT FALSE,
    substitute_id INT REFERENCES stu_tracker.Tutors(id) ON DELETE SET NULL,
    start_time VARCHAR(10),
    subject VARCHAR(100),
    notes TEXT,
    edited_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
**/

func (s *AuthService) SessionSearch(ss models.SearchQuery) ([]models.ServiceSession, error) {
	query, args := buildSearchQuery(ss)
	fmt.Println(query)
	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("error querying locations: %w", err)
	}
	defer rows.Close()
	var sessions []models.ServiceSession
	for rows.Next() {
		var session models.ServiceSession
		err := rows.Scan(
			&session.FirstName,
			&session.LastName,
			&session.ID,
			&session.TutorId,
			&session.Location,
			&session.Substitute,
			&session.SubstituteId,
			&session.StartTime,
			&session.Subject,
			&session.Notes,
			&session.EditedAt,
			&session.CreatedAt,
			&session.ProgramName,
			&session.SubjectName,
			&session.StudentCount,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}
		sessions = append(sessions, session)
	}
	// Check for any errors encountered during iteration
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}
	return sessions, nil
}

func (s *AuthService) StudentSessionSearch(ss models.SearchQuery) ([]models.StudentSessions, error) {
	query, args := buildStudentSearchQuery(ss)
	fmt.Println(query)
	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("error querying locations: %w", err)
	}
	defer rows.Close()
	var sessions []models.StudentSessions
	for rows.Next() {
		var session models.StudentSessions
		err := rows.Scan(
			&session.ID,
			&session.FirstName,
			&session.LastName,
			&session.SessionCount,
			&session.AssessmentCount,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}
		sessions = append(sessions, session)
	}
	// Check for any errors encountered during iteration
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}
	return sessions, nil
}

func (s *AuthService) SessionInfo(session_id int64) ([]models.SessionInfoStudent, error) {
	query := `
	SELECT ss.duration, st.id as student_id, st.first_name, st.last_name, COALESCE(st.middle_name, '') as middle_name, st.email, st.grade_level AS grade, st.period
	FROM 
		stu_tracker.Session_students ss 
	JOIN 
		stu_tracker.Students st 
	ON 
		st.id = ss.student_id 
	WHERE 
		ss.session_id = $1`
	fmt.Println(query)

	rows, err := s.db.Query(query, session_id)
	if err != nil {
		return nil, fmt.Errorf("error querying Session_students: %w", err)
	}
	defer rows.Close()
	var sessionsInfo []models.SessionInfoStudent
	for rows.Next() {
		var session models.SessionInfoStudent
		err := rows.Scan(
			&session.Duration,
			&session.ID,
			&session.FirstName,
			&session.LastName,
			&session.MiddleName,
			&session.Email,
			&session.Grade,
			&session.Period,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}
		sessionsInfo = append(sessionsInfo, session)
	}
	// Check for any errors encountered during iteration
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}
	return sessionsInfo, nil
}

func (s *AuthService) StudentInfo(student_id int64) ([]models.SessionInfoStudent, error) {
	query := `
	SELECT ss.duration, st.id as student_id, st.first_name, st.last_name, COALESCE(st.middle_name, '') as middle_name, st.email, st.grade_level AS grade, st.period
	FROM 
		stu_tracker.Session_students ss 
	JOIN 
		stu_tracker.Students st 
	ON 
		st.id = ss.student_id 
	WHERE 
		ss.session_id = $1`
	fmt.Println(query)

	rows, err := s.db.Query(query, student_id)
	if err != nil {
		return nil, fmt.Errorf("error querying Session_students: %w", err)
	}
	defer rows.Close()
	var sessionsInfo []models.SessionInfoStudent
	for rows.Next() {
		var session models.SessionInfoStudent
		err := rows.Scan(
			&session.Duration,
			&session.ID,
			&session.FirstName,
			&session.LastName,
			&session.MiddleName,
			&session.Email,
			&session.Grade,
			&session.Period,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}
		sessionsInfo = append(sessionsInfo, session)
	}
	// Check for any errors encountered during iteration
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}
	return sessionsInfo, nil
}

func (s *AuthService) AssessmentInfo(session_id int64) ([]models.AssessmentInfoStudent, error) {
	query := `
	SELECT ats.title, ats.letter,
	ats.cycle, ats.max_score, aas.score, aas.created_at, st.id
	FROM 
		stu_tracker.Assessments_students aas
	JOIN 
		stu_tracker.Students st 
	ON 
		st.id = aas.student_id
	JOIN 
		stu_tracker.Assessments ats 
	ON 
		ats.id = aas.assessment_id
	WHERE 
		aas.session_id = $1;`
	rows, err := s.db.Query(query, session_id)
	if err != nil {
		return nil, fmt.Errorf("error querying AssessmentInfo: %w", err)
	}
	defer rows.Close()
	var assessmentInfo []models.AssessmentInfoStudent
	for rows.Next() {
		var assessment models.AssessmentInfoStudent
		err := rows.Scan(
			&assessment.Title,
			&assessment.Letter,
			&assessment.Cycle,
			&assessment.MaxScore,
			&assessment.Score,
			&assessment.CreatedAt,
			&assessment.StudentID,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}
		assessmentInfo = append(assessmentInfo, assessment)
	}
	// Check for any errors encountered during iteration
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}
	return assessmentInfo, nil
}

func (s *AuthService) StudentSessionInfo(student_id int64, organization_id int64) ([]models.StudentSessionInfo, error) {
	query := `
	SELECT
		ast.created_at, ast.absent, ast.subject_id, ast.duration
	FROM 
		stu_tracker.Sessions ss
	LEFT JOIN 
		stu_tracker.Session_students ast
	ON	
		ast.session_id = ss.id
	WHERE 
		ast.student_id = 44 AND ss.organization_id = $2`
	rows, err := s.db.Query(query, student_id, organization_id)
	if err != nil {
		return nil, fmt.Errorf("error querying AssessmentInfo: %w", err)
	}
	defer rows.Close()
	var assessmentInfo []models.StudentSessionInfo
	for rows.Next() {
		var assessment models.StudentSessionInfo
		err := rows.Scan(
			&assessment.CreatedAt,
			&assessment.Absent,
			&assessment.SubjectID,
			&assessment.Duration,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}
		assessmentInfo = append(assessmentInfo, assessment)
	}
	// Check for any errors encountered during iteration
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}
	return assessmentInfo, nil
}

func (s *AuthService) StudentAssessmentInfo(student_id int64, organization_id int64) ([]models.StudentAssessmentInfo, error) {
	query := `
		SELECT
			ast.created_at, ast.absent, ast.subject_id, ast.duration
		FROM 
			stu_tracker.Sessions ss
		LEFT JOIN 
			stu_tracker.Session_students ast
		ON	
			ast.session_id = ss.id
		WHERE 
			ast.student_id = 44 AND ss.organization_id = $2`
	rows, err := s.db.Query(query, student_id, organization_id)
	if err != nil {
		return nil, fmt.Errorf("error querying AssessmentInfo: %w", err)
	}
	defer rows.Close()
	var assessmentInfo []models.StudentAssessmentInfo
	for rows.Next() {
		var assessment models.StudentAssessmentInfo
		err := rows.Scan(
			&assessment.CreatedAt,
			&assessment.Absent,
			&assessment.AssessmentID,
			&assessment.MaxScore,
			&assessment.Score,
			&assessment.Cycle,
			&assessment.Letter,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}
		assessmentInfo = append(assessmentInfo, assessment)
	}
	// Check for any errors encountered during iteration
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}
	return assessmentInfo, nil
}

func (s *AuthService) GetSemestersVAssessmentChart(req models.RequestSemestersVAssessmentChart) ([]*models.AssessmentComparison, error) {
	if req.Semester1ID == req.Semester2ID {
		return nil, fmt.Errorf("semester id cannot be the same")
	}
	semester_1_dataset, err := s.GetSemestersVAssessmentChartData(&req)
	if err != nil {
		return nil, err
	}
	return semester_1_dataset, nil

}

func (s *AuthService) GetSemestersVAssessmentChartData(req *models.RequestSemestersVAssessmentChart) ([]*models.AssessmentComparison, error) {
	query := `WITH SemesterAverages AS (
		SELECT
			ast.assessment_id,
			ss.semester_id,
			MIN(ast.score) AS min_score,
			MAX(ast.score) AS max_score,
			AVG(ast.score) AS average_score,
			ss.organization_id,
			COUNT(ast.id) AS instances
		FROM
			stu_tracker.Assessments_students ast
		INNER JOIN
			stu_tracker.Sessions ss		
		ON
			ast.session_id = ss.id
		GROUP BY
			ast.assessment_id, ss.semester_id, ss.organization_id
	)
	SELECT
		s1.assessment_id AS assessment_s1,
		s2.assessment_id AS assessment_s2,
		s1.instances AS count_s1,
		s2.instances AS count_s2,
		s1.min_score,
		s1.max_score,
		s2.max_score,
		s2.min_score,
		COALESCE(s1.average_score, 0) AS semester_1_avg,
		COALESCE(s2.average_score, 0) AS semester_2_avg,
		(s2.average_score - s1.average_score) AS score_difference,
		ROUND((s2.average_score - s1.average_score) / s1.average_score * 100, 2) AS rate_of_change_percent
	FROM
		SemesterAverages s1
	JOIN
		SemesterAverages s2
	ON
		s1.semester_id = $1 AND s2.semester_id = $2
	WHERE
		s1.average_score IS NOT NULL AND s2.average_score IS NOT NULL
	AND
		s1.assessment_id = $3
	AND 
		s2.assessment_id = $4
	AND
		s1.organization_id = $5
	AND 
		s2.organization_id = $6
		
	GROUP BY
		semester_1_avg, semester_2_avg, score_difference, rate_of_change_percent,
		s1.assessment_id, 
		s1.min_score,
		s1.max_score,
		s2.max_score,
		s2.min_score,
		s2.assessment_id,count_s1,count_s2`

	rows, err := s.db.Query(query, req.Semester1ID, req.Semester2ID, req.Assessment1ID, req.Assessment1ID, req.OrganizationID, req.OrganizationID)
	if err != nil {
		return nil, fmt.Errorf("error querying AssessmentInfo: %w", err)
	}
	defer rows.Close()

	// Initialize a slice to store all rows
	var dataSet []*models.AssessmentComparison

	// Iterate over each row
	for rows.Next() {
		// Create a new instance of AssessmentComparison for each row
		var row models.AssessmentComparison

		// Scan the row into the struct fields
		err := rows.Scan(
			&row.AssessmentS1,
			&row.AssessmentS2,
			&row.CountS1,
			&row.CountS2,
			&row.MinScoreS1,
			&row.MaxScoreS1,
			&row.MinScoreS2,
			&row.MaxScoreS2,
			&row.Semester1Avg,
			&row.Semester2Avg,
			&row.ScoreDifference,
			&row.RateOfChangePercent,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}

		// Append the row to the dataSet slice
		dataSet = append(dataSet, &row)
	}

	// Check for errors after iteration
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	// Return the populated dataSet
	return dataSet, nil

}

func (s *AuthService) GetAssessmentGrowth(req models.RequestAssessmentGrowth) (*models.ResponseAssessmentGrowth, error) {
	if req.Assessment1ID == req.Assessment2ID {
		return nil, fmt.Errorf("assessment 1 cannot be assessment 2")
	}
	assessment1, err := s.QueryAssessmentGrowth(req.LocationID, req.OrganizationID, req.Assessment1ID)
	if err != nil {
		return nil, fmt.Errorf("unable to get assessment 1 data %w", err)
	}
	assessment2, err := s.QueryAssessmentGrowth(req.LocationID, req.OrganizationID, req.Assessment2ID)
	if err != nil {
		return nil, fmt.Errorf("unable to get assessment 2 data %w", err)
	}
	return &models.ResponseAssessmentGrowth{
		DataSet1: *assessment1,
		DataSet2: *assessment2,
	}, nil
}

func (s *AuthService) QueryAssessmentGrowth(LocationID *int64, OrganizationID *int64, Assessment1ID *int64) (*models.AssessmentGrowth, error) {

	query, args := buildSearchQueryAssessmentGrowth(LocationID, OrganizationID, Assessment1ID)
	var AssessmentGrowth models.AssessmentGrowth
	err := s.db.QueryRow(query, args...).Scan(&AssessmentGrowth.AssessmentId, &AssessmentGrowth.MinScore, &AssessmentGrowth.MaxScore, &AssessmentGrowth.Average, &AssessmentGrowth.Count)
	if err != nil {
		return nil, err
	}
	return &AssessmentGrowth, nil
}

func buildSearchQueryAssessmentGrowth(LocationID *int64, OrganizationID *int64, Assessment1ID *int64) (string, []interface{}) {
	argIndex := 1
	var args []interface{}
	var conditions []string
	query := `
		SELECT
			ast.assessment_id,
			MIN(ast.score) AS min_score,
			MAX(ast.score) AS max_score,
			AVG(ast.score) AS average_score,
			COUNT(ast.id)
		FROM
			stu_tracker.Assessments_students ast
		INNER JOIN
			stu_tracker.Sessions ss
		ON
			ss.id = ast.session_id`
	if LocationID != nil {
		conditions = append(conditions, fmt.Sprintf("ss.location_id = $%d", argIndex))
		args = append(args, LocationID)
		argIndex++
	}
	if OrganizationID != nil {
		conditions = append(conditions, fmt.Sprintf("ss.organization_id = $%d", argIndex))
		args = append(args, OrganizationID)
		argIndex++
	}
	if Assessment1ID != nil {
		conditions = append(conditions, fmt.Sprintf("ast.assessment_id = $%d", argIndex))
		args = append(args, Assessment1ID)
		argIndex++
	}
	if len(conditions) > 0 {
		query += ` WHERE ` + strings.Join(conditions, " AND ")
	}
	query += `
		GROUP BY
			ast.assessment_id`

	return query, args
}

func buildSearchQuery(ss models.SearchQuery) (string, []interface{}) {
	argIndex := 1
	var args []interface{}
	var conditions []string
	query := `SELECT  
			t.first_name AS tutor_first_name, 
			t.last_name AS tutor_last_name, 
			ss.id,
			ss.tutor_id, 
			ll.name AS location_name,
			ss.substitute, 
			ss.substitute_id, 
			ss.start_time, 
			COALESCE(ss.subject_id, null) AS subject_id,
			ss.notes, 
			ss.edited_at, 
			ss.created_at,
			COALESCE(pg.program_name, 'No program') AS program_name,
    		COALESCE(sb.title, 'No Subject') AS subject_name,
			ss.student_count
		FROM 
			stu_tracker.Sessions ss
		JOIN 
    		stu_tracker.Tutors t ON ss.tutor_id = t.id 
		JOIN 
			stu_tracker.Locations ll ON ll.id = ss.location_id
		LEFT JOIN 
			stu_tracker.Programs pg ON pg.id = ss.program_id 
		LEFT JOIN 
    		stu_tracker.Subjects sb ON ss.subject_id = sb.id 
		`

	if ss.SearchTerm != "" {
		conditions = append(conditions, fmt.Sprintf("ss.notes ILIKE $%d", argIndex))
		args = append(args, "%"+ss.SearchTerm+"%")
		argIndex++
	}
	if ss.LocationId != nil {
		conditions = append(conditions, fmt.Sprintf("ss.location_id = $%d", argIndex))
		args = append(args, ss.LocationId)
		argIndex++
	}
	if ss.OrganizationID != nil {
		conditions = append(conditions, fmt.Sprintf("ss.organization_id = $%d", argIndex))
		args = append(args, ss.OrganizationID)
		argIndex++
	}
	if ss.ProgramId != nil {
		conditions = append(conditions, fmt.Sprintf("ss.program_id = $%d", argIndex))
		args = append(args, ss.ProgramId)
		argIndex++
	}
	if ss.SemesterID != nil {
		conditions = append(conditions, fmt.Sprintf("ss.semester_id = $%d", argIndex))
		args = append(args, ss.SemesterID)
		argIndex++
	}
	if !ss.DateStart.IsZero() && !ss.DateEnd.IsZero() {
		conditions = append(conditions, fmt.Sprintf("DATE(ss.session_date) BETWEEN $%d AND $%d", argIndex, argIndex+1))
		args = append(args, ss.DateStart, ss.DateEnd)
		argIndex += 2
	} else if !ss.DateStart.IsZero() {
		conditions = append(conditions, fmt.Sprintf("DATE(ss.session_date) >= $%d", argIndex))
		args = append(args, ss.DateStart)
		argIndex++
	}

	if ss.SubjectId != nil {
		conditions = append(conditions, fmt.Sprintf("ss.subject_id = $%d", argIndex))
		args = append(args, ss.SubjectId)
		argIndex++
	}

	if len(conditions) > 0 {
		query += "WHERE " + strings.Join(conditions, " AND ")
	}
	return query, args
}

func buildStudentSearchQuery(ss models.SearchQuery) (string, []interface{}) {
	argIndex := 1
	var args []interface{}
	var conditions []string
	query := `SELECT 
			s.id AS student_id,
			s.first_name,
			s.last_name,
			COUNT(DISTINCT ss.session_id) AS session_count,
			COUNT(DISTINCT a.assessment_id) AS assessment_count
		FROM stu_tracker.Students s
		LEFT JOIN stu_tracker.Session_students ss ON s.id = ss.student_id
		LEFT JOIN stu_tracker.Sessions st ON st.id = ss.session_id
		LEFT JOIN stu_tracker.Assessments_students a ON s.id = a.student_id `

	if ss.SearchTerm != "" {
		conditions = append(conditions, fmt.Sprintf("st.first_name ILIKE $%d OR st.last_name ILIKE $%d", argIndex, argIndex+1))
		args = append(args, "%"+ss.SearchTerm+"%", "%"+ss.SearchTerm+"%")
		argIndex += 2
	}
	if ss.LocationId != nil {
		conditions = append(conditions, fmt.Sprintf("st.location_id = $%d", argIndex))
		args = append(args, ss.LocationId)
		argIndex++
	}
	if ss.OrganizationID != nil {
		conditions = append(conditions, fmt.Sprintf("st.organization_id = $%d", argIndex))
		args = append(args, ss.OrganizationID)
		argIndex++
	}
	if ss.ProgramId != nil {
		conditions = append(conditions, fmt.Sprintf("st.program_id = $%d", argIndex))
		args = append(args, ss.ProgramId)
		argIndex++
	}
	if ss.SemesterID != nil {
		conditions = append(conditions, fmt.Sprintf("st.semester_id = $%d", argIndex))
		args = append(args, ss.SemesterID)
		argIndex++
	}
	if !ss.DateStart.IsZero() && !ss.DateEnd.IsZero() {
		conditions = append(conditions, fmt.Sprintf("DATE(st.session_date) BETWEEN $%d AND $%d", argIndex, argIndex+1))
		args = append(args, ss.DateStart, ss.DateEnd)
		argIndex += 2
	} else if !ss.DateStart.IsZero() {
		conditions = append(conditions, fmt.Sprintf("DATE(st.session_date) >= $%d", argIndex))
		args = append(args, ss.DateStart)
		argIndex++
	}

	if ss.SubjectId != nil {
		conditions = append(conditions, fmt.Sprintf("ss.subject_id = $%d", argIndex))
		args = append(args, ss.SubjectId)
		argIndex++
	}

	if len(conditions) > 0 {
		query += "WHERE " + strings.Join(conditions, " AND ")
	}
	query += ` GROUP BY s.id, s.first_name, s.last_name
				ORDER BY session_count DESC, assessment_count DESC;`
	return query, args
}
