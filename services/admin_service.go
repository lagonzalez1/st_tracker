package services

import (
	"context"
	"fmt"
	"strings"
	"time"
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
				first_name, tr.id
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

func (s *AuthService) GetStudentGroups(ctx context.Context, lid *int64, tid *int64, sid *int64) ([]models.RegisterStudentGroup, error) {
	base := `SELECT sg.id, sg.title, sg.description, sg.location_id, sg.semester_id, sg.tutor_id
	FROM stu_tracker.Student_groups sg
	JOIN stu_tracker.Semester s 
	ON s.id = sg.semester_id
	WHERE s.active = TRUE AND sg.location_id = $1 AND sg.tutor_id = $2 AND sg.semester_id = $3;`
	rows, err := s.db.QueryContext(ctx, base, lid, tid, sid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var g []models.RegisterStudentGroup
	for rows.Next() {
		var sm models.RegisterStudentGroup
		if err := rows.Scan(
			&sm.ID,
			&sm.Title,
			&sm.Description,
			&sm.LocationID,
			&sm.SemesterID,
			&sm.TutorID,
		); err != nil {
			return nil, err
		}
		g = append(g, sm)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return g, nil
}

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
    	SELECT DISTINCT
			tr.id,
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
			&s.ID,
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

func (s *AuthService) GetCycleGrowth(req models.RequestCycleGrowth) ([]models.ResponseCycleGrowth, error) {
	base := `
        SELECT 
			ss.location_id,
			ss.program_id,
			l.name AS location_name,
			pg.program_name,
			k.cycle,
			ROUND( (SUM(ast.score)::numeric / NULLIF(SUM(k.max_score),0)) * 100, 2 ) AS avg_score,
			ROUND(MIN((ast.score / NULLIF(k.max_score,0)) * 100.0)::numeric, 2)    AS min_score,
			ROUND(MAX((ast.score / NULLIF(k.max_score,0)) * 100.0)::numeric, 2)    AS max_score,
			ROUND(STDDEV(ast.score)::NUMERIC, 2) AS stddev_score,
			COUNT(ast.id) AS total_assessments
		FROM stu_tracker.Assessments_students ast
		JOIN stu_tracker.Sessions ss 
			ON ast.session_id = ss.id
		JOIN stu_tracker.Assessments k 
			ON k.id = ast.assessment_id
		JOIN stu_tracker.Programs pg
			ON pg.id = ss.program_id
		JOIN stu_tracker.Locations l
			ON l.id = ss.location_id
	`
	query, args := buildQuery(base, req)
	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var sessions []models.ResponseCycleGrowth
	for rows.Next() {
		var s models.ResponseCycleGrowth
		if err := rows.Scan(
			&s.LocationID,
			&s.ProgramID,
			&s.LocationName,
			&s.ProgramName,
			&s.Cycle,
			&s.AverageScore,
			&s.MinScore,
			&s.MaxScore,
			&s.StandardDeviation,
			&s.TotalAssessments,
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

func buildQuery(baseQuery string, req models.RequestCycleGrowth) (string, []interface{}) {
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
		GROUP BY ss.location_id, ss.program_id, k.cycle, l.name, pg.program_name
		ORDER BY ss.location_id, ss.program_id;
	`
	return baseQuery, args
}

func (s *AuthService) GetCycleGrowthDelim(req models.RequestCycleGrowth) ([]models.ResponseCycleGrowthDelim, error) {
	base := `
        SELECT 
			s.location_id,
			s.program_id,
			k.cycle,
			k.pre,
			k.mid,
			k.post,
			l.name AS location_name,
			pg.program_name,
			ROUND( (SUM(ast.score)::numeric / NULLIF(SUM(k.max_score),0)) * 100, 2 ) AS avg_score,
			ROUND(MIN((ast.score / NULLIF(k.max_score,0)) * 100.0)::numeric, 2)    AS min_score,
			ROUND(MAX((ast.score / NULLIF(k.max_score,0)) * 100.0)::numeric, 2)    AS max_score,
			ROUND(STDDEV(ast.score)::NUMERIC, 2) AS stddev_score,
			COUNT(ast.id) AS total_assessments
		FROM stu_tracker.Assessments_students ast
		JOIN stu_tracker.Sessions s 
			ON ast.session_id = s.id
		JOIN stu_tracker.Assessments k 
			ON k.id = ast.assessment_id
		JOIN stu_tracker.Programs pg
			ON pg.id = s.program_id
		JOIN stu_tracker.Locations l
			ON l.id = s.location_id
	`
	query, args := buildQueryDelim(base, req)
	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var sessions []models.ResponseCycleGrowthDelim
	for rows.Next() {
		var s models.ResponseCycleGrowthDelim
		if err := rows.Scan(
			&s.LocationID,
			&s.ProgramID,
			&s.Cycle,
			&s.Pre,
			&s.Mid,
			&s.Post,
			&s.LocationName,
			&s.ProgramName,
			&s.AverageScore,
			&s.MinScore,
			&s.MaxScore,
			&s.StandardDeviation,
			&s.TotalAssessments,
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

func buildQueryDelim(baseQuery string, req models.RequestCycleGrowth) (string, []interface{}) {
	var clauses []string
	var args []interface{}
	argCount := 1

	if req.LocationID != nil {
		clauses = append(clauses, fmt.Sprintf("s.location_id = $%d", argCount))
		args = append(args, req.LocationID)
		argCount += 1
	}
	if req.OrganizationID != nil {
		clauses = append(clauses, fmt.Sprintf("s.organization_id = $%d", argCount))
		args = append(args, req.OrganizationID)
		argCount += 1
	}
	if req.ProgramID != nil {
		clauses = append(clauses, fmt.Sprintf("s.program_id = $%d", argCount))
		args = append(args, req.ProgramID)
		argCount += 1
	}
	if req.SemesterID != nil {
		clauses = append(clauses, fmt.Sprintf("s.semester_id = $%d", argCount))
		args = append(args, req.SemesterID)
		argCount += 1
	}
	if !req.StartDate.IsZero() && !req.EndDate.IsZero() {
		clauses = append(clauses, fmt.Sprintf("s.session_date BETWEEN $%d AND $%d", argCount, argCount+1))
		args = append(args, req.StartDate, req.EndDate)
		argCount += 2
	} else if !req.StartDate.IsZero() {
		clauses = append(clauses, fmt.Sprintf("s.session_date >= $%d", argCount))
		args = append(args, req.StartDate)
		argCount++
	}
	if len(clauses) > 0 {
		baseQuery += " WHERE " + strings.Join(clauses, " AND ")
	}
	baseQuery += `
		GROUP BY s.location_id, s.program_id, k.cycle, k.pre, k.post, k.mid, l.name, pg.program_name
		ORDER BY s.location_id, s.program_id;
	`
	return baseQuery, args
}

func (s *AuthService) GetAbsentPresent(c context.Context, req models.RequestCycleGrowth) (*models.ResponseAbsentPresent, error) {
	base := `
        SELECT
		SUM(CASE WHEN ss.absent = FALSE THEN 1 ELSE 0 END) AS present_count,
		SUM(CASE WHEN ss.absent = TRUE  THEN 1 ELSE 0 END) AS absent_count
		FROM stu_tracker.session_students ss
		JOIN stu_tracker.sessions s ON s.id = ss.session_id
	`
	var model models.ResponseAbsentPresent
	query, args := buildQueryAbsentPresent(base, req)
	err := s.db.QueryRowContext(c, query, args...).Scan(&model.Present, &model.Absent)
	if err != nil {
		return nil, err
	}
	return &model, nil
}

func buildQueryAbsentPresent(baseQuery string, req models.RequestCycleGrowth) (string, []interface{}) {
	var clauses []string
	var args []interface{}
	argCount := 1

	if req.LocationID != nil {
		clauses = append(clauses, fmt.Sprintf("s.location_id = $%d", argCount))
		args = append(args, req.LocationID)
		argCount += 1
	}
	if req.OrganizationID != nil {
		clauses = append(clauses, fmt.Sprintf("s.organization_id = $%d", argCount))
		args = append(args, req.OrganizationID)
		argCount += 1
	}
	if req.ProgramID != nil {
		clauses = append(clauses, fmt.Sprintf("s.program_id = $%d", argCount))
		args = append(args, req.ProgramID)
		argCount += 1
	}
	if req.SemesterID != nil {
		clauses = append(clauses, fmt.Sprintf("s.semester_id = $%d", argCount))
		args = append(args, req.SemesterID)
		argCount += 1
	}
	if !req.StartDate.IsZero() && !req.EndDate.IsZero() {
		clauses = append(clauses, fmt.Sprintf("s.session_date BETWEEN $%d AND $%d", argCount, argCount+1))
		args = append(args, req.StartDate, req.EndDate)
		argCount += 2
	} else if !req.StartDate.IsZero() {
		clauses = append(clauses, fmt.Sprintf("s.session_date >= $%d", argCount))
		args = append(args, req.StartDate)
		argCount++
	}
	if len(clauses) > 0 {
		baseQuery += " WHERE " + strings.Join(clauses, " AND ")
	}
	return baseQuery, args
}

func (s *AuthService) GetAssessmentCompletion(c context.Context, req models.RequestCycleGrowth) (*models.ResponseAssessmentCompletion, error) {
	base := `
		WITH params AS (
			SELECT
				$1::int          AS semester_id,         -- required
				$2::int          AS location_id,         -- optional (nullable)
				$3::int          AS organization_id,     -- optional (nullable)
				$4::date  AS start_date,          -- optional (nullable)
				$5::date  AS end_date            -- optional (nullable)
			),

			enrolled AS (
			SELECT st.id AS student_id
			FROM stu_tracker.students st
			CROSS JOIN params p
			WHERE st.semester_id = p.semester_id
				AND (p.location_id IS NULL OR st.location_id = p.location_id)
				AND (
				p.organization_id IS NULL OR
				EXISTS (
					SELECT 1
					FROM stu_tracker.locations l
					WHERE l.id = st.location_id
					AND l.organization_id = p.organization_id
				)
				)
			),

			assessed AS (
			SELECT DISTINCT ast.student_id
			FROM stu_tracker.assessments_students ast
			JOIN stu_tracker.sessions s ON s.id = ast.session_id
			CROSS JOIN params p
			WHERE s.semester_id = p.semester_id
				AND (p.location_id   IS NULL OR s.location_id   = p.location_id)
				AND (p.organization_id IS NULL OR s.organization_id = p.organization_id)
				AND (p.start_date    IS NULL OR s.session_date >= p.start_date)
				AND (p.end_date      IS NULL OR s.session_date <  p.end_date)
			),

			interim_assessed AS (
			SELECT DISTINCT ast.student_id
			FROM stu_tracker.assessments_students ast
			JOIN stu_tracker.sessions s ON s.id = ast.session_id
			JOIN stu_tracker.assessments a ON a.id = ast.assessment_id
			CROSS JOIN params p
			WHERE s.semester_id = p.semester_id
				AND (p.location_id   IS NULL OR s.location_id   = p.location_id)
				AND (p.organization_id IS NULL OR s.organization_id = p.organization_id)
				AND (p.start_date    IS NULL OR s.session_date >= p.start_date)
				AND (p.end_date      IS NULL OR s.session_date <  p.end_date)
				AND (a.mid IS TRUE)  -- or a.cycle = 'mid'
			),

			post_assessed AS (
			SELECT DISTINCT ast.student_id
			FROM stu_tracker.assessments_students ast
			JOIN stu_tracker.sessions s ON s.id = ast.session_id
			JOIN stu_tracker.assessments a ON a.id = ast.assessment_id
			CROSS JOIN params p
			WHERE s.semester_id = p.semester_id
				AND (p.location_id   IS NULL OR s.location_id   = p.location_id)
				AND (p.organization_id IS NULL OR s.organization_id = p.organization_id)
				AND (p.start_date    IS NULL OR s.session_date >= p.start_date)
				AND (p.end_date      IS NULL OR s.session_date <  p.end_date)
				AND (a.post IS TRUE) -- or a.cycle = 'post'
			)

			SELECT
			COALESCE((SELECT COUNT(*) FROM enrolled), 0)         AS enrolled_count,
			COALESCE((SELECT COUNT(*) FROM assessed), 0)         AS assessed_count,
			COALESCE((SELECT COUNT(*) FROM interim_assessed), 0) AS interim_assessed_count,
			COALESCE((SELECT COUNT(*) FROM post_assessed), 0)    AS post_assessed_count;
	`
	var model models.ResponseAssessmentCompletion
	var startTime *time.Time = nil
	var endTime *time.Time = nil
	if !req.StartDate.IsZero() {
		startTime = &req.StartDate
	}
	if !req.EndDate.IsZero() {
		endTime = &req.EndDate
	}
	err := s.db.QueryRowContext(c, base, req.SemesterID, req.LocationID, *req.OrganizationID, startTime, endTime).Scan(&model.EnrolledCount, &model.AssessedCount, &model.IntremAssessend, &model.PostAssessed)
	if err != nil {
		return nil, err
	}
	return &model, nil
}

func (s *AuthService) GetSentimentAnalysis(c context.Context, req models.RequestSentiment) ([]models.ResponseSentiment, error) {
	query := `
		SELECT
			p.program_name AS program_name,
			s.session_date,
			t.id AS tutor_id,
			t.first_name,
			sr.session_id,
			sr.question_text,
			sr.response_text,
			sr.sentiment_score
		FROM stu_tracker.Sessions s
		LEFT JOIN stu_tracker.Tutors t ON s.tutor_id = t.id
		LEFT JOIN stu_tracker.Programs p ON s.program_id = p.id
		RIGHT JOIN stu_tracker.Survey_response sr ON sr.session_id = s.id
	`
	q, args := buildQuerySentimentAnalysis(query, &req)
	fmt.Printf("query string:  %s", q)
	rows, err := s.db.QueryContext(c, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var res []models.ResponseSentiment
	for rows.Next() {
		var model models.ResponseSentiment
		if err := rows.Scan(
			&model.ProgramName,
			&model.SessionDate,
			&model.TutorID,
			&model.FirstName,
			&model.SessionID,
			&model.QuestionText,
			&model.ResponseText,
			&model.SentimentScore,
		); err != nil {
			return nil, err
		}
		res = append(res, model)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("unable to scan from db response")
	}

	return res, nil
}

func buildQuerySentimentAnalysis(baseQuery string, req *models.RequestSentiment) (string, []interface{}) {
	var clauses []string
	var args []interface{}
	argCount := 1

	if req.LocationID != nil {
		clauses = append(clauses, fmt.Sprintf("s.location_id = $%d", argCount))
		args = append(args, req.LocationID)
		argCount += 1
	}
	if req.OrganizationID != nil {
		clauses = append(clauses, fmt.Sprintf("s.organization_id = $%d", argCount))
		args = append(args, req.OrganizationID)
		argCount += 1
	}
	if req.ProgramID != nil {
		clauses = append(clauses, fmt.Sprintf("s.program_id = $%d", argCount))
		args = append(args, req.ProgramID)
		argCount += 1
	}
	if req.SemesterID != nil {
		clauses = append(clauses, fmt.Sprintf("s.semester_id = $%d", argCount))
		args = append(args, req.SemesterID)
		argCount += 1
	}
	if !req.StartDate.IsZero() && !req.EndDate.IsZero() {
		clauses = append(clauses, fmt.Sprintf("s.session_date BETWEEN $%d AND $%d", argCount, argCount+1))
		args = append(args, req.StartDate, req.EndDate)
		argCount += 2
	} else if !req.StartDate.IsZero() {
		clauses = append(clauses, fmt.Sprintf("s.session_date >= $%d", argCount))
		args = append(args, req.StartDate)
		argCount++
	}
	if len(clauses) > 0 {
		baseQuery += " WHERE " + strings.Join(clauses, " AND ")
	}
	baseQuery += `
		GROUP BY
			program_name, t.first_name, s.session_date, t.id, sr.id, sr.session_id
		ORDER BY
			sr.id, sr.session_id, s.session_date ASC;`
	return baseQuery, args

}

func (s *AuthService) GetSentimentAnalysisByTutor(c context.Context, tid *int64, orgid *int64) ([]models.ResponseSentiment, error) {
	query := `
		SELECT
			p.program_name AS program_name,
			s.session_date,
			t.id AS tutor_id,
			t.first_name,
			sr.session_id,
			sr.question_text,
			sr.response_text,
			sr.sentiment_score
		FROM stu_tracker.Sessions s
		JOIN stu_tracker.Tutors t ON s.tutor_id = t.id
		JOIN stu_tracker.Programs p ON s.program_id = p.id
		RIGHT JOIN stu_tracker.Survey_response sr
			ON sr.session_id = s.id
		WHERE s.organization_id = $1 AND s.tutor_id = $2
		ORDER BY s.session_date ASC;
	`
	rows, err := s.db.QueryContext(c, query, tid, orgid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var res []models.ResponseSentiment
	for rows.Next() {
		var model models.ResponseSentiment
		if err := rows.Scan(
			&model.ProgramName,
			&model.SessionDate,
			&model.TutorID,
			&model.FirstName,
			&model.SessionID,
			&model.QuestionText,
			&model.ResponseText,
			&model.SentimentScore,
		); err != nil {
			return nil, err
		}
		res = append(res, model)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("unable to scan from db response")
	}

	return res, nil
}

func (s *AuthService) GetAssessmentsByTutor(c context.Context, tid *int64, orgid *int64) ([]models.ResponseAssessments, error) {
	query := `
	SELECT st.session_id, ss.session_date, st.student_id, st.score, ast.title, p.program_name, ast.max_score, p.id, s.first_name
	from stu_tracker.Assessments_students st
	JOIN stu_tracker.Sessions ss ON ss.id = st.session_id 
	JOIN stu_tracker.Assessments ast ON ast.id = st.assessment_id 
	JOIN stu_tracker.Tutors t ON ss.tutor_id = t.id
	JOIN stu_tracker.Programs p ON ss.program_id = p.id
	JOIN stu_tracker.Students s ON st.student_id = s.id
	WHERE t.id = $1 AND ss.organization_id = $2
	GROUP BY st.student_id, ss.session_date, st.session_id, st.score, ast.title, p.program_name, ast.max_score,  p.id, s.first_name
	ORDER BY ss.session_date desc;
`
	rows, err := s.db.QueryContext(c, query, tid, orgid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var res []models.ResponseAssessments
	for rows.Next() {
		var model models.ResponseAssessments
		if err := rows.Scan(
			&model.SessionID,
			&model.SessionDate,
			&model.StudentID,
			&model.Score,
			&model.Assessment,
			&model.ProgramName,
			&model.MaxScore,
			&model.ProgramID,
			&model.StudentFirstName,
		); err != nil {
			return nil, err
		}
		res = append(res, model)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("unable to scan from db response")
	}

	return res, nil
}

func (s *AuthService) GetAssessmentPivotTable(c context.Context, sid *int64, orgid *int64) (*models.ResponseAssessmentPivotTable, error) {
	query := `SELECT t.first_name, t.last_name, p.program_name AS program,
	 p.id, s.duration, t.id, s.session_date, s.start_time
	 FROM stu_tracker.Sessions s
	 LEFT JOIN stu_tracker.Programs p ON p.id = s.program_id
	 JOIN stu_tracker.Tutors t ON t.id = s.tutor_id
	 WHERE s.id = $1 AND s.organization_id = $2;`
	var m models.ResponseAssessmentPivotTable
	err := s.db.QueryRowContext(c, query, sid, orgid).Scan(&m.FirstName,
		&m.LastName,
		&m.Program,
		&m.ProgramID,
		&m.Duration,
		&m.TutorID,
		&m.SessionDate,
		&m.StartTime)
	if err != nil {
		return nil, err
	}

	return &m, nil
}

func (s *AuthService) GetSubscriptions(c context.Context, orgid *int64) ([]models.AvilableSubscriptions, error) {
	query := `SELECT id, code, name, stripe_price_id, is_active, cost_yearly, cost_monthly FROM stu_tracker.subscription_plan;`
	rows, err := s.db.QueryContext(c, query)
	if err != nil {
		return nil, err
	}
	var res []models.AvilableSubscriptions
	for rows.Next() {
		var s models.AvilableSubscriptions
		if err := rows.Scan(
			&s.ID,
			&s.Code,
			&s.Name,
			&s.PriceID,
			&s.IsActive,
			&s.CostYearly,
			&s.CostMonthly,
		); err != nil {
			return nil, err
		}
		res = append(res, s)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("unable to scan from db response")
	}
	return res, nil
}

func (s *AuthService) GetSubscriptionsEntitlements(c context.Context, orgid *int64) ([]models.SubscriptionsEntitlements, error) {
	query := `SELECT id, plan_id, key, limit_value, enabled, enterprise FROM stu_tracker.plan_entitlement;`
	rows, err := s.db.QueryContext(c, query)
	if err != nil {
		return nil, err
	}
	var res []models.SubscriptionsEntitlements
	for rows.Next() {
		var s models.SubscriptionsEntitlements
		if err := rows.Scan(
			&s.ID,
			&s.PlanID,
			&s.Key,
			&s.LimitValue,
			&s.Enabled,
			&s.Enterprise,
		); err != nil {
			return nil, err
		}
		res = append(res, s)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("unable to scan from db response")
	}
	return res, nil
}

func (s *AuthService) GetSubscriptionsByOrganization(c context.Context, orgid *int64) ([]models.OrganizationSubscription, error) {
	query := `SELECT id, plan_id, status, current_period_start, current_period_end, canceled_at, stripe_customer_id, stripe_subscription_id
			FROM stu_tracker.organization_subscription
			WHERE organization_id = $1;`
	rows, err := s.db.QueryContext(c, query, orgid)
	if err != nil {
		return nil, err
	}
	var res []models.OrganizationSubscription
	for rows.Next() {
		var s models.OrganizationSubscription
		if err := rows.Scan(
			&s.ID,
			&s.PlanID,
			&s.Status,
			&s.CurrentPeriodStart,
			&s.CurrentPeriodEnd,
			&s.CanceledAt,
			&s.StripeCustomerID,
			&s.StripeSubscriptionID,
		); err != nil {
			return nil, err
		}
		res = append(res, s)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("unable to scan from db response")
	}
	return res, nil
}

func (s *AuthService) GetSessionVScore(c context.Context, req models.RequestSessionBChart) ([]models.ResponseAssessmentVScore, error) {
	query := `
		SELECT
		sstu.student_id,
		st.first_name,
		st.last_name,
		COUNT(sstu.session_id)::int AS x_sessions,
		ROUND(AVG( (asj.score::numeric / NULLIF(a.max_score, 0)) * 100 ), 1) AS y_avg_score_pct
		FROM stu_tracker.Session_students sstu
		JOIN stu_tracker.Sessions ss ON ss.id = sstu.session_id
		LEFT JOIN stu_tracker.Assessments_students asj ON asj.session_id = ss.id
		LEFT JOIN stu_tracker.Assessments a ON a.id = asj.assessment_id
		JOIN stu_tracker.Students st ON st.id = sstu.student_id
	`
	q, args := buildSessionVScore(query, &req)
	rows, err := s.db.QueryContext(c, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var res []models.ResponseAssessmentVScore
	for rows.Next() {
		var m models.ResponseAssessmentVScore
		if err := rows.Scan(
			&m.StudentID,
			&m.FirstName,
			&m.LastName,
			&m.SessionCount,
			&m.ScoreAverage,
		); err != nil {
			return nil, err
		}
		res = append(res, m)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("unable to scan from db response")
	}
	return res, nil
}

func buildSessionVScore(baseQuery string, req *models.RequestSessionBChart) (string, []interface{}) {
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
		GROUP BY sstu.student_id, st.first_name, st.last_name
		HAVING COUNT(DISTINCT sstu.session_id) > 0;`
	return baseQuery, args

}

func (s *AuthService) GetStudentVAssessments(c context.Context, req models.RequestSessionBChart) ([]models.ResponseStudentVAssessments, error) {
	query := `
		SELECT
			st.id AS student_id,
			st.first_name,
			st.last_name,
			am.id AS assessment_id,
			am.title AS assessment_title,
			COALESCE(ast.score, 0) AS assessment_score,
			am.max_score AS assessment_max_score,
			CASE 
				WHEN ast.score IS NULL THEN 'Absent'
				ELSE 'Present'
			END AS status,
			ss.id
		FROM stu_tracker.Sessions ss
		INNER JOIN stu_tracker.Assessments_students ast ON ast.session_id = ss.id
		LEFT JOIN stu_tracker.Students st ON st.id = ast.student_id
		JOIN stu_tracker.Assessments am ON ast.assessment_id = am.id
	`
	q, args := buildGetStudentVAssessments(query, &req)
	rows, err := s.db.QueryContext(c, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var res []models.ResponseStudentVAssessments
	for rows.Next() {
		var m models.ResponseStudentVAssessments
		if err := rows.Scan(
			&m.StudentID,
			&m.FirstName,
			&m.LastName,
			&m.AssessmentID,
			&m.AssessmentTitle,
			&m.Score,
			&m.MaxScore,
			&m.Status,
			&m.SessionID,
		); err != nil {
			return nil, err
		}
		res = append(res, m)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("unable to scan from db response")
	}
	return res, nil
}

func buildGetStudentVAssessments(baseQuery string, req *models.RequestSessionBChart) (string, []interface{}) {
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
		GROUP BY st.first_name, st.last_name, am.id, am.max_score, ast.score, st.id, ss.id
		ORDER BY st.first_name, st.last_name, am.id`
	return baseQuery, args

}
