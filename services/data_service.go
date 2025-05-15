package services

import (
	"context"
	"fmt"
	"strings"
	"tracker/app/models"

	"github.com/lib/pq"
)

func (s *AuthService) StudentAssessmentSearch(c context.Context, student_assessment_id *int64, easyScore bool) ([]models.StudentAssessmentSearch, error) {
	query := `SELECT
			ans.question_id, 
			q.question_text,
			q.points,
			q.question_type,
			ans.choice_id,
			ans.answer_text,
			ans.is_correct,
			c.choice_text
			FROM 
				stu_tracker.Assessment_answers ans
			LEFT JOIN
				stu_tracker.Choices c
			ON
				ans.choice_id = c.id
			LEFT JOIN
				stu_tracker.Questions q
			ON
				q.id = ans.question_id
			WHERE
				ans.assessment_student_id = $1;`

	rows, err := s.db.QueryContext(c, query, student_assessment_id)
	if err != nil {
		return nil, fmt.Errorf("error querying Assessment_answers: %w", err)
	}
	defer rows.Close()
	var studentAssessments []models.StudentAssessmentSearch
	for rows.Next() {
		var assessment models.StudentAssessmentSearch
		err := rows.Scan(
			&assessment.QuestionID,
			&assessment.Question,
			&assessment.Points,
			&assessment.QuestionType,
			&assessment.ChoiceID,
			&assessment.AnswerText,
			&assessment.IsCorrect,
			&assessment.ChoiceText,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)

		}
		studentAssessments = append(studentAssessments, assessment)
	}
	return studentAssessments, nil
}

func (s *AuthService) SessionSearch(c context.Context, ss models.SearchQuery) ([]models.ServiceSession, error) {
	query, args := buildSearchQuery(ss)
	rows, err := s.db.QueryContext(c, query, args...)
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
			&session.SubFirstName,
			&session.SubLastName,
			&session.StartTime,
			&session.Subject,
			&session.Notes,
			&session.EditedAt,
			&session.CreatedAt,
			&session.ProgramName,
			&session.SubjectName,
			&session.StudentCount,
			&session.SessionDate,
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

func (s *AuthService) TutorSearch(c context.Context, query models.SearchQueryTutor) ([]models.ResponseRequestTutorsList, error) {
	if query.SearchTerm == "" || query.OrganizationID == nil {
		return nil, fmt.Errorf("Organization or search term cannot be empty")
	}
	q := `SELECT t.id, t.first_name, t.last_name, t.email, t.created_at, t.location_id
			  FROM 
			  	stu_tracker.Tutors t 
			  WHERE 
			  	t.organization_id = $1
			  AND 
			  	to_tsvector('english', first_name || ' ' || last_name) @@ to_tsquery('english', $2 || ':*'); `

	rows, err := s.db.QueryContext(c, q, *query.OrganizationID, query.SearchTerm)
	if err != nil {
		return nil, fmt.Errorf("error querying tutors: %w", err)
	}
	var tutors []models.ResponseRequestTutorsList
	for rows.Next() {
		var tutor models.ResponseRequestTutorsList
		err := rows.Scan(
			&tutor.ID,
			&tutor.FirstName,
			&tutor.LastName,
			&tutor.Email,
			&tutor.CreatedAt,
			&tutor.LocationId,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning tutors: %w", err)
		}
		tutors = append(tutors, tutor)
	}
	return tutors, nil
}

func (s *AuthService) StudentSessionSearch(c context.Context, ss models.SearchQuery) ([]models.StudentSessions, error) {
	query, args := buildStudentSearchQuery(ss)
	fmt.Println(query)
	rows, err := s.db.QueryContext(c, query, args...)
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
			&session.MiddleName,
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

func (s *AuthService) GetSessionsTutors(c context.Context, ss models.RequestTutorsSessions) ([]models.TutorSessionsList, error) {
	// Must check if semester_id is valid. Simply check if within range of current ?
	var query string
	query += `
		SELECT 
		ss.id,
		ss.session_date::TIMESTAMP AS session_date,
		ss.start_time,
		ss.location_id,
		pg.program_name,
		pg.id,
		ss.semester_id,
		ss.student_count,
		ss.in_school,
		ss.substitute,
		sm.title,
		lc.name
		FROM
			stu_tracker.Sessions ss
		JOIN
			stu_tracker.Programs pg
		ON 
			ss.program_id = pg.id
		JOIN
			stu_tracker.Semester sm
		ON 
			ss.semester_id = sm.id
		JOIN
			stu_tracker.Locations lc
		ON 
			ss.location_id = lc.id
		WHERE
			ss.tutor_id = $1 AND ss.semester_id = $2 AND ss.location_id = $3
		ORDER BY 
			ss.session_date desc;
		`
	rows, err := s.db.QueryContext(c, query, ss.ID, ss.SemesterID, ss.LocationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	// Can also return the student count, assessment count ...
	var sessions []models.TutorSessionsList
	for rows.Next() {
		var session models.TutorSessionsList
		err := rows.Scan(
			&session.SessionID,
			&session.SessionDate,
			&session.StartTime,
			&session.LocationID,
			&session.ProgramName,
			&session.ProgramID,
			&session.SemesterID,
			&session.StudentCount,
			&session.InSchool,
			&session.Substitute,
			&session.Semester,
			&session.LocationName,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}
		sessions = append(sessions, session)
	}
	var sessionIDs []int64
	if len(sessions) <= 0 {
		return []models.TutorSessionsList{}, nil
	}
	for i := 0; i < len(sessions); i++ {
		sessionIDs = append(sessionIDs, *sessions[i].SessionID)
	}
	// For each session_id return the students associated
	query2 := `
		SELECT
			ss.session_id,
			st.first_name,
			st.last_name,
			st.middle_name,
			st.grade_level,
			st.id
		FROM 
			stu_tracker.Session_students ss
		INNER JOIN
			stu_tracker.Students st
		ON
			ss.student_id = st.id
		WHERE 
			ss.session_id = ANY($1)
		ORDER BY
			ss.session_id`

	rows, err = s.db.Query(query2, pq.Array(sessionIDs))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	// Can also return the student count, assessment count ...
	studentMap := make(map[int64][]models.Students)
	for rows.Next() {
		var student models.Students
		err := rows.Scan(
			&student.SessionID,
			&student.FirstName,
			&student.LastName,
			&student.MiddleName,
			&student.Grade,
			&student.StudentID,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}
		studentMap[int64(*student.SessionID)] = append(studentMap[int64(*student.SessionID)], student)
	}
	for i := 0; i < len(sessions); i++ {
		sessions[i].Students = studentMap[int64(*sessions[i].SessionID)]
	}
	return sessions, nil
}

func (s *AuthService) GetStudentDetails(c context.Context, session_id *string, student_id *int64) (*models.StudentDetails, error) {

	query := `SELECT 
		s.first_name,
		s.last_name,
		s.middle_name,
		s.grade_level
		FROM stu_tracker.Students s
		INNER JOIN stu_tracker.Assessment_sessions a
		ON a.student_id = s.id
		WHERE a.student_id = $1 AND a.session_token = $2
		LIMIT 1;
	`
	var studentDetails models.StudentDetails
	err := s.db.QueryRowContext(c, query, student_id, session_id).Scan(&studentDetails.FirstName, &studentDetails.LastName, &studentDetails.MiddleName, &studentDetails.GradeLevel)
	if err != nil {
		return nil, err
	}
	return &studentDetails, nil
}

func (s *AuthService) AssessmentInfo(c context.Context, session_id int64) ([]models.AssessmentInfoStudent, error) {
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
	rows, err := s.db.QueryContext(c, query, session_id)
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

func (s *AuthService) StudentAssessmentInfo(c context.Context, student_id int64, organization_id int64) ([]models.StudentAssessmentInfo, error) {
	query := `
		SELECT
		ast.id,
		ast.created_at,
		ast.score,
		ast.session_id,
		a.title,
		a.max_score,
		a.pre,
		a.post,
		a.mid,
		a.cycle,
		a.letter,
		a.version,
		a.easy_score
		FROM 
			stu_tracker.Assessments_students ast
		LEFT JOIN 
			stu_tracker.Students ss
		ON	
			ast.session_id = ss.id
		JOIN
			stu_tracker.Assessments a
		ON
			a.id = ast.assessment_id
		WHERE 
			ast.student_id = $1`
	rows, err := s.db.QueryContext(c, query, student_id)
	if err != nil {
		return nil, fmt.Errorf("error querying AssessmentInfo: %w", err)
	}
	defer rows.Close()
	var assessmentInfo []models.StudentAssessmentInfo
	for rows.Next() {
		var assessment models.StudentAssessmentInfo
		err := rows.Scan(
			&assessment.ID,
			&assessment.CreatedAt,
			&assessment.Score,
			&assessment.SessionID,
			&assessment.Title,
			&assessment.MaxScore,
			&assessment.Pre,
			&assessment.Post,
			&assessment.Mid,
			&assessment.Cycle,
			&assessment.Letter,
			&assessment.Version,
			&assessment.EasyScore,
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

func (s *AuthService) TrailSessions(c context.Context, student_id int64) ([]models.SessionTrail, error) {
	query := `
		SELECT
			ss.id,
			t.first_name, 
			t.last_name,
			COALESCE(p.program_name, 'NA') as program_name, 
			ss.duration as session_duration,
			ss.start_time,
			ss.notes,
			ss.created_at,
			ss.student_count,
			ss.substitute,
			sst.absent,
			sst.duration as student_duration,
			COALESCE(sub.first_name || ' ' || sub.last_name, '') AS substitute_name
		FROM 
			stu_tracker.Session_students sst 
		LEFT JOIN 
			stu_tracker.Sessions ss
		ON 
			sst.session_id = ss.id
		JOIN
			stu_tracker.Tutors t
		ON
			t.id = ss.tutor_id
		JOIN 
			stu_tracker.Programs p
		ON 
			p.id = ss.program_id
		LEFT JOIN
  			stu_tracker.Tutors sub 
		ON 
			sub.id = ss.substitute_id 
		WHERE 
			sst.student_id = $1`

	rows, err := s.db.QueryContext(c, query, student_id)
	if err != nil {
		return nil, fmt.Errorf("error querying Session_students: %w", err)
	}
	defer rows.Close()
	var sessionsInfo []models.SessionTrail
	for rows.Next() {
		var session models.SessionTrail
		err := rows.Scan(
			&session.ID,
			&session.FirstName,
			&session.LastName,
			&session.ProgramName,
			&session.SessionDuration,
			&session.StartTime,
			&session.Notes,
			&session.CreatedAt,
			&session.StudentCount,
			&session.Substitute,
			&session.Absent,
			&session.StudentDuration,
			&session.SubstituteName,
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

func (s *AuthService) SessionInfo(session_id int64) ([]models.SessionInfoStudent, error) {
	query := `
	SELECT
		ss.duration, st.id as student_id, 
		st.first_name, st.last_name, COALESCE(st.middle_name, '') as middle_name, 
		st.email, st.grade_level AS grade, st.period
	FROM 
		stu_tracker.Session_students ss 
	LEFT JOIN 
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
	SELECT 
		ss.duration, st.id as student_id, 
		st.first_name, st.last_name, 
		COALESCE(st.middle_name, '') as middle_name, 
		st.email, st.grade_level AS grade, 
		st.period
	FROM 
		stu_tracker.Session_students ss 
	LEFT JOIN 
		stu_tracker.Students st 
	ON 
		st.id = ss.student_id 
	WHERE 
		ss.student_id = $1`
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
	if *req.Assessment1ID == *req.Assessment2ID {
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
			CASE 
				WHEN ss.substitute = true THEN sub.first_name 
				ELSE NULL 
			END AS substitute_first_name,
			CASE 
				WHEN ss.substitute = true THEN sub.last_name 
				ELSE NULL 
			END AS substitute_last_name, 
			ss.start_time, 
			COALESCE(ss.subject_id, null) AS subject_id,
			ss.notes, 
			ss.edited_at, 
			ss.created_at,
			COALESCE(pg.program_name, 'No program') AS program_name,
    		COALESCE(sb.title, 'No Subject') AS subject_name,
			ss.student_count,
			ss.session_date
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
		LEFT JOIN 
    		stu_tracker.Tutors sub ON ss.substitute_id = sub.id AND ss.substitute = true 
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
			s.middle_name,
			s.last_name,
			COUNT(DISTINCT ss.id) AS session_count,
			COUNT(DISTINCT a.id) AS assessment_count
		FROM 
			stu_tracker.Students s
		JOIN 
			stu_tracker.Session_students ss ON s.id = ss.student_id
		JOIN 
			stu_tracker.Assessments_students a ON s.id = a.student_id
		LEFT JOIN 
			stu_tracker.Sessions st ON st.id = ss.session_id `

	if ss.SearchTerm != "" {
		conditions = append(conditions, fmt.Sprintf("s.first_name ILIKE $%d OR s.last_name ILIKE $%d", argIndex, argIndex+1))
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
