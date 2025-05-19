package services

import (
	"context"
	"fmt"
	"strings"
	"tracker/app/models"
)

func buildQuerySessionBChart(baseQuery string, req models.RequestSessionBChart) (string, []interface{}) {
	var clauses []string
	var args []interface{}
	argCount := 1

	if req.LocationID != nil {
		clauses = append(clauses, fmt.Sprintf("ss.location_id = $%d", argCount))
		args = append(args, req.LocationID)
		argCount += 1
	}
	if req.OrganizationID != nil {
		clauses = append(clauses, fmt.Sprintf("ss.organization_id = $%d", argCount))
		args = append(args, req.OrganizationID)
		argCount += 1
	}
	if req.ProgramID != nil {
		clauses = append(clauses, fmt.Sprintf("ss.program_id = $%d", argCount))
		args = append(args, req.ProgramID)
		argCount += 1
	}
	if req.SemesterID != nil {
		clauses = append(clauses, fmt.Sprintf("ss.semester_id = $%d", argCount))
		args = append(args, req.SemesterID)
		argCount += 1
	}
	if !req.StartDate.IsZero() && !req.EndDate.IsZero() {
		clauses = append(clauses, fmt.Sprintf("ss.session_date BETWEEN $%d AND $%d", argCount, argCount+1))
		args = append(args, req.StartDate, req.EndDate)
		argCount += 2
	} else if !req.StartDate.IsZero() {
		clauses = append(clauses, fmt.Sprintf("ss.session_date >= $%d", argCount))
		args = append(args, req.StartDate)
		argCount++
	}
	if len(clauses) > 0 {
		baseQuery += " WHERE " + strings.Join(clauses, " AND ")
	}
	baseQuery += `
	GROUP BY
		DATE(session_date), month_name
	ORDER BY
		DATE(session_date)
	`
	return baseQuery, args
}

func buildQueryAssessmentBChart(baseQuery string, req models.RequestSessionBChart) (string, []interface{}) {
	var clauses []string
	var args []interface{}
	argCount := 1

	if req.LocationID != nil {
		clauses = append(clauses, fmt.Sprintf("ss.location_id = $%d", argCount))
		args = append(args, req.LocationID)
		argCount += 1
	}
	if req.OrganizationID != nil {
		clauses = append(clauses, fmt.Sprintf("ss.organization_id = $%d", argCount))
		args = append(args, req.OrganizationID)
		argCount += 1
	}
	if req.ProgramID != nil {
		clauses = append(clauses, fmt.Sprintf("ss.program_id = $%d", argCount))
		args = append(args, req.ProgramID)
		argCount += 1
	}
	if req.SemesterID != nil {
		clauses = append(clauses, fmt.Sprintf("ss.semester_id = $%d", argCount))
		args = append(args, req.SemesterID)
		argCount += 1
	}
	if !req.StartDate.IsZero() && !req.EndDate.IsZero() {
		clauses = append(clauses, fmt.Sprintf("ss.session_date BETWEEN $%d AND $%d", argCount, argCount+1))
		args = append(args, req.StartDate, req.EndDate)
		argCount += 2
	} else if !req.StartDate.IsZero() {
		clauses = append(clauses, fmt.Sprintf("ss.session_date >= $%d", argCount))
		args = append(args, req.StartDate)
		argCount++
	}
	if len(clauses) > 0 {
		baseQuery += " WHERE " + strings.Join(clauses, " AND ")
	}
	baseQuery += `
	GROUP BY
   		loc.name, asts.title, asts.cycle, asts.letter
	ORDER BY
		loc.name;
	`
	return baseQuery, args
}

func buildQueryTutorsBChart(baseQuery string, req models.RequestSessionBChart) (string, []interface{}) {
	var clauses []string
	var args []interface{}
	argCount := 1

	if req.LocationID != nil {
		clauses = append(clauses, fmt.Sprintf("ss.location_id = $%d", argCount))
		args = append(args, req.LocationID)
		argCount += 1
	}
	if req.OrganizationID != nil {
		clauses = append(clauses, fmt.Sprintf("ss.organization_id = $%d", argCount))
		args = append(args, req.OrganizationID)
		argCount += 1
	}
	if req.ProgramID != nil {
		clauses = append(clauses, fmt.Sprintf("ss.program_id = $%d", argCount))
		args = append(args, req.ProgramID)
		argCount += 1
	}
	if req.SemesterID != nil {
		clauses = append(clauses, fmt.Sprintf("ss.semester_id = $%d", argCount))
		args = append(args, req.SemesterID)
		argCount += 1
	}
	if !req.StartDate.IsZero() && !req.EndDate.IsZero() {
		clauses = append(clauses, fmt.Sprintf("ss.session_date BETWEEN $%d AND $%d", argCount, argCount+1))
		args = append(args, req.StartDate, req.EndDate)
		argCount += 2
	} else if !req.StartDate.IsZero() {
		clauses = append(clauses, fmt.Sprintf("ss.session_date >= $%d", argCount))
		args = append(args, req.StartDate)
		argCount++
	}
	if len(clauses) > 0 {
		baseQuery += " WHERE " + strings.Join(clauses, " AND ")
	}
	baseQuery += `
			GROUP BY
				first_name, tr.id
			ORDER BY
				first_name
	`
	return baseQuery, args
}

func buildSearchQueryTutorSessions(ss *models.RequestTutorSessions) (string, []interface{}) {
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
			ss.notes, 
			ss.edited_at, 
			ss.created_at,
			ss.student_count,
			COUNT(ast.id) AS assessment_count

		FROM 
			stu_tracker.Sessions ss
		JOIN 
    		stu_tracker.Tutors t ON ss.tutor_id = t.id 
		JOIN 
			stu_tracker.Locations ll ON ll.id = ss.location_id
		LEFT JOIN
			stu_tracker.Assessments_students ast ON ast.session_id = ss.id `

	if ss.LocationID != nil {
		conditions = append(conditions, fmt.Sprintf("ss.location_id = $%d", argIndex))
		args = append(args, ss.LocationID)
		argIndex++
	}
	if ss.OrganizationID != nil {
		conditions = append(conditions, fmt.Sprintf("ss.organization_id = $%d", argIndex))
		args = append(args, ss.OrganizationID)
		argIndex++
	}
	if ss.SemesterID != nil {
		conditions = append(conditions, fmt.Sprintf("ss.semester_id = $%d", argIndex))
		args = append(args, ss.SemesterID)
		argIndex++
	}

	if len(conditions) > 0 {
		query += "WHERE " + strings.Join(conditions, " AND ")
	}
	query += ` GROUP BY tutor_first_name, tutor_last_name, ss.id, location_name ORDER BY ss.created_at DESC`
	return query, args
}

func buildLowPerformanceTutors(ss *models.RequestSessionBChart) (string, []interface{}) {
	argIndex := 1
	var args []interface{}
	var conditions []string
	query := `
		WITH tutor_stats AS (
			SELECT 
				t.id,
				t.first_name || ' ' || t.last_name AS tutor_name,
				COUNT(DISTINCT s.id) AS session_count,
				COUNT(DISTINCT aas.student_id) AS unique_student_count,
				AVG(aas.score) AS avg_student_score,
				PERCENT_RANK() OVER (ORDER BY COUNT(DISTINCT s.id)) AS session_percentile,
				PERCENT_RANK() OVER (ORDER BY COUNT(DISTINCT aas.student_id)) AS student_percentile
			FROM 
				stu_tracker.Tutors t
			LEFT JOIN 
				stu_tracker.Sessions s ON t.id = s.tutor_id
			LEFT JOIN 
				stu_tracker.Assessments_students aas ON s.id = aas.session_id `

	if ss.LocationID != nil {
		conditions = append(conditions, fmt.Sprintf("s.location_id = $%d", argIndex))
		args = append(args, ss.LocationID)
		argIndex++
	}
	if ss.SemesterID != nil {
		conditions = append(conditions, fmt.Sprintf("s.semester_id = $%d", argIndex))
		args = append(args, ss.SemesterID)
		argIndex++
	}
	if ss.ProgramID != nil {
		conditions = append(conditions, fmt.Sprintf("s.program_id = $%d", argIndex))
		args = append(args, ss.ProgramID)
		argIndex += 1
	}
	if !ss.StartDate.IsZero() && !ss.EndDate.IsZero() {
		conditions = append(conditions, fmt.Sprintf("s.session_date BETWEEN $%d AND $%d", argIndex, argIndex+1))
		args = append(args, ss.StartDate, ss.EndDate)
		argIndex += 2
	} else if !ss.StartDate.IsZero() {
		conditions = append(conditions, fmt.Sprintf("s.session_date >= $%d", argIndex))
		args = append(args, ss.StartDate)
		argIndex++
	}

	if len(conditions) > 0 {
		query += "WHERE " + strings.Join(conditions, " AND ")
	}
	query += ` AND t.active = TRUE`
	query += ` GROUP BY 
				t.id, t.first_name, t.last_name
			)
			SELECT 
			id,
			tutor_name,
			session_count,
			unique_student_count,
			avg_student_score,
			session_percentile,
			student_percentile,
			CASE 
				WHEN session_percentile < 0.25 AND student_percentile < 0.25 THEN 'High Concern'
				WHEN session_percentile < 0.25 OR student_percentile < 0.25 THEN 'Moderate Concern'
				ELSE 'Performing Adequately'
			END AS performance_status
			FROM 
			tutor_stats
			WHERE 
			session_percentile < 0.25 OR student_percentile < 0.25  -- Bottom 25th percentile
			ORDER BY 
			session_count ASC, 
			unique_student_count ASC;`
	return query, args
}

func (s *AuthService) GetSessionAnalytics(req models.RequestSessionBChart) (*models.SessionAnalytics, error) {
	if req.OrganizationID == nil {
		return &models.SessionAnalytics{
			SessionCount:    nil,
			AssessmentCount: nil,
			StudentCount:    nil,
		}, nil
	}
	sessionCount, studentCount, err := s.GetSessionCounts(int(*req.OrganizationID))
	if err != nil {
		return nil, fmt.Errorf("error getting session counts: %w", err)
	}

	// Get assessment counts
	assessmentCount, err := s.GetAssessmentCounts(int(*req.OrganizationID))
	if err != nil {
		return nil, fmt.Errorf("error getting assessment counts: %w", err)
	}

	// Return the combined results
	return &models.SessionAnalytics{
		SessionCount:    sessionCount,
		AssessmentCount: assessmentCount,
		StudentCount:    studentCount,
	}, nil
}

func (s *AuthService) GetAssessmentCounts(OrganizationID int) (*int, error) {
	var assessmentCount *int

	query := `
        SELECT
			COUNT(ast.id) AS assessments_count
		FROM
			stu_tracker.Sessions ss
		INNER JOIN
			stu_tracker.Assessments_students ast
		ON
			ss.id = ast.session_id
		WHERE
			ss.organization_id = $1
		AND
			DATE(ast.created_at) = CURRENT_DATE;
    	`
	err := s.db.QueryRow(query, OrganizationID).Scan(&assessmentCount)
	if err != nil {
		return nil, fmt.Errorf("error querying assessment counts: %w", err)
	}

	return assessmentCount, nil
}

func (s *AuthService) GetSessionCounts(OrganizationID int) (*int, *int, error) {
	var sessionCount, studentCount *int

	query := `
        SELECT
			COUNT(ss.id) AS session_count,
			SUM(ss.student_count) AS student_count
		FROM
			stu_tracker.Sessions ss
		WHERE
			ss.organization_id = $1
		AND
			DATE(ss.session_date) = CURRENT_DATE;
		`

	err := s.db.QueryRow(query, OrganizationID).Scan(&sessionCount, &studentCount)
	if err != nil {
		return nil, nil, fmt.Errorf("error querying session counts: %w", err)
	}

	return sessionCount, studentCount, nil
}

/*
	Generate sessions given a possible range, location, semester, program
*/

func (s *AuthService) GetSessionBChart(req models.RequestSessionBChart) ([]models.ResponseSessionBChart, error) {
	base := `
        SELECT
		DATE(session_date) AS session_date,
		COUNT(*) AS total_sessions,
		AVG(student_count) AS student_average,
		SUM(student_count) AS total_students,
		MIN(student_count) AS min_students,
		MAX(student_count) AS max_students,
		ROUND(AVG(student_count), 2) AS student_average_rounded,
		TO_CHAR(session_date, 'Month') AS month_name
	FROM
		stu_tracker.Sessions ss
	`
	query, args := buildQuerySessionBChart(base, req)
	fmt.Println(query)
	fmt.Print(query)
	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var sessions []models.ResponseSessionBChart
	for rows.Next() {
		var s models.ResponseSessionBChart
		if err := rows.Scan(
			&s.SessionDate,
			&s.TotalSessions,
			&s.StudentAverage,
			&s.TotalStudents,
			&s.MinStudents,
			&s.MaxStudents,
			&s.StudentAverageRounded,
			&s.Month,
		); err != nil {
			return nil, err
		}
		sessions = append(sessions, s)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return sessions, nil
}

func (s *AuthService) GetAssessmentBChart(req models.RequestSessionBChart) ([]models.ResponseAssessmentsBChart, error) {
	base := `
    SELECT
		loc.name AS location_name,
		asts.title AS assessment_name,
		COUNT(*) AS assessment_total,
		MIN(ats.score) AS min_score,
		MAX(ats.score) AS max_score,
		AVG(ats.score) AS average_score,
		asts.cycle,
		asts.letter
	FROM 
		stu_tracker.Sessions ss
	INNER JOIN	
		stu_tracker.Assessments_students ats
	ON
		ats.session_id = ss.id
	INNER JOIN 
		stu_tracker.Assessments asts
	ON
		ats.assessment_id = asts.id
	INNER JOIN 
		stu_tracker.Locations loc
	ON
		loc.id = ss.location_id`
	query, args := buildQueryAssessmentBChart(base, req)

	fmt.Print(query)
	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var assessments []models.ResponseAssessmentsBChart
	for rows.Next() {
		var s models.ResponseAssessmentsBChart
		if err := rows.Scan(
			&s.LocationName,
			&s.AssessmentName,
			&s.AssessemtsTotal,
			&s.MinScore,
			&s.MaxScore,
			&s.AssessemtsAverage,
			&s.AssessmentCycle,
			&s.AssessmentLetter,
		); err != nil {
			return nil, err
		}
		assessments = append(assessments, s)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return assessments, nil
}

func (s *AuthService) GetProgramsBChart(req models.RequestSessionBChart) ([]models.ResponseProgramsBChart, error) {
	if req.OrganizationID == nil {
		return nil, fmt.Errorf("unable to serach missing org id")
	}
	query := `
    SELECT
		pr.program_name AS program_name,
		COUNT(*)
	FROM 
		stu_tracker.Programs pr
	INNER JOIN
		stu_tracker.Location_programs lp
	ON
		pr.id = lp.program_id
	WHERE
		pr.organization_id = $1
	GROUP BY
		program_name
	ORDER BY
		program_name;`

	rows, err := s.db.Query(query, req.OrganizationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var programs []models.ResponseProgramsBChart
	for rows.Next() {
		var s models.ResponseProgramsBChart
		if err := rows.Scan(
			&s.ProgramName,
			&s.Count,
		); err != nil {
			return nil, err
		}
		programs = append(programs, s)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return programs, nil
}

func (s *AuthService) GetTutorsBChart(req models.RequestSessionBChart) ([]models.ResponseTutorsBChart, error) {
	if req.OrganizationID == nil {
		return nil, fmt.Errorf("unable to serach missing org id")
	}
	base := `
    	SELECT
			SUM(ss.student_count) as total_student_count,
			AVG(ss.student_count) as average_student_count,
			COUNT(ss.id) as total_sessions,
			COUNT(ast.id) as assessments_count,
			tr.first_name as first_name,
			tr.last_name as last_name
		FROM 
			stu_tracker.Tutors tr
		INNER JOIN
			stu_tracker.Sessions ss
		ON
			tr.id = ss.tutor_id
		LEFT JOIN
			stu_tracker.Assessments_students ast
		ON 
			ast.session_id = ss.id`

	query, args := buildQueryTutorsBChart(base, req)
	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var programs []models.ResponseTutorsBChart
	for rows.Next() {
		var s models.ResponseTutorsBChart
		if err := rows.Scan(
			&s.StudentCount,
			&s.AverageStudents,
			&s.TotalSessions,
			&s.AssessmentCount,
			&s.FirstName,
			&s.LastName,
		); err != nil {
			return nil, err
		}
		programs = append(programs, s)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return programs, nil
}

func (s *AuthService) GetAssessmentTrendLine(ctx context.Context, req models.RequestSessionBChart) ([]models.ResponseAssessmentTrendline, error) {
	if req.OrganizationID == nil {
		return nil, fmt.Errorf("unable to serach missing org id")
	}
	base := `
    	SELECT
			count(*) as assessment_total,
			EXTRACT(YEAR FROM session_date) AS YEAR,
			EXTRACT(MONTH FROM session_date) AS MONTH
		FROM 
			stu_tracker.Sessions ss
		INNER JOIN	
			stu_tracker.Assessments_students ats
		ON
			ats.session_id = ss.id
		INNER JOIN 
			stu_tracker.Assessments asts
		ON
			ats.assessment_id = asts.id
		WHERE
			ss.organization_id = $1
		AND
			session_date > (CURRENT_DATE - INTERVAL '6 months')
		GROUP BY
			YEAR, MONTH
		ORDER BY
			YEAR`

	rows, err := s.db.QueryContext(ctx, base, req.OrganizationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	index := 0
	var assessments []models.ResponseAssessmentTrendline
	for rows.Next() {
		var s models.ResponseAssessmentTrendline
		s.ID = index
		if err := rows.Scan(
			&s.AssessmentCount,
			&s.Year,
			&s.Month,
		); err != nil {
			return nil, err
		}
		assessments = append(assessments, s)
		index++
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return assessments, nil
}

func (s *AuthService) GetSessionTrendLine(req models.RequestSessionBChart) ([]models.ResponseSessionTrendline, error) {
	if req.OrganizationID == nil {
		return nil, fmt.Errorf("unable to serach missing org id")
	}
	base := `
    	SELECT
			count(*),
			EXTRACT(YEAR FROM session_date) AS YEAR,
			EXTRACT(MONTH FROM session_date) AS MONTH
		FROM 
			stu_tracker.Sessions ss
		WHERE
			ss.organization_id = $1
		AND
			session_date > (CURRENT_DATE - INTERVAL '6 months')
		GROUP BY
			YEAR, MONTH
		ORDER BY
			YEAR`

	rows, err := s.db.Query(base, req.OrganizationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var sessions []models.ResponseSessionTrendline
	index := 0
	for rows.Next() {
		var s models.ResponseSessionTrendline
		s.ID = index
		if err := rows.Scan(
			&s.SessionCount,
			&s.Year,
			&s.Month,
		); err != nil {
			return nil, err
		}
		sessions = append(sessions, s)
		index++
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return sessions, nil
}

func (s *AuthService) GetTutorSessions(req models.RequestTutorSessions) ([]models.ServiceSession, error) {
	query, args := buildSearchQueryTutorSessions(&req)
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
			&session.Notes,
			&session.EditedAt,
			&session.CreatedAt,
			&session.StudentCount,
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

func (s *AuthService) GetTutorLowPerformance(c context.Context, req models.RequestSessionBChart) ([]models.ResponseTutorLowPerformance, error) {
	query, args := buildLowPerformanceTutors(&req)
	fmt.Println(query)
	rows, err := s.db.QueryContext(c, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var performance []models.ResponseTutorLowPerformance
	for rows.Next() {
		var tutor models.ResponseTutorLowPerformance
		err := rows.Scan(
			&tutor.ID,
			&tutor.TutorName,
			&tutor.SessionCount,
			&tutor.UniqueStudentCount,
			&tutor.AverageStudentScore,
			&tutor.SessionPercentile,
			&tutor.StudentPercentile,
			&tutor.PerformanceStatus,
		)
		if err != nil {
			return nil, err
		}
		performance = append(performance, tutor)
	}

	return performance, nil
}
