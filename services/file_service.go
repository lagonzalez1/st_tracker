package services

import (
	"bytes"
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"
	"tracker/app/models"

	"github.com/lib/pq"
	"github.com/xuri/excelize/v2"
)

// For each tutor session_id : 1
// Return each student under such session
// AND
// Return each assessment under such session

func (s *AuthService) GetTutorFileData(c context.Context, ss models.RequestDownloadData) (*excelize.File, error) {
	query, args := buildSessionQueryTutorData(ss)
	rows, err := s.db.QueryContext(c, query, args...)
	if err != nil {
		return nil, fmt.Errorf("error querying locations: %w", err)
	}
	defer rows.Close()
	var sessions []models.TutorSessionData
	for rows.Next() {
		var session models.TutorSessionData
		err := rows.Scan(
			&session.SessionID,
			&session.TutorID,
			&session.SessionDate,
			&session.Substitute,
			&session.StudentCount,
			&session.StartTime,
			&session.Duration,
			&session.Notes,
			&session.TutorName,
			&session.TutorLastName,
			&session.Program,
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
	if len(sessions) <= 0 {
		return nil, nil
	}
	// If file has some len build and return
	f, err := buildSessionTutorFile(sessions)
	if err != nil {
		return nil, fmt.Errorf("unable to build session tutor file %v", err)
	}
	// Verify file can be written to a buffer
	var buf bytes.Buffer
	if err := f.Write(&buf); err != nil {
		return nil, fmt.Errorf("failed to verify Excel file creation: %v", err)
	}
	return f, nil
}

func (s *AuthService) GetStudentFileData(c context.Context, ss models.RequestDownloadData) (*excelize.File, error) {
	// Get all SQL session ids by params
	query, args := buildSessionQuery(ss)
	rows, err := s.db.QueryContext(c, query, args...)
	if err != nil {
		return nil, fmt.Errorf("error querying locations: %w", err)
	}
	defer rows.Close()
	var sessionIDs []int64
	for rows.Next() {
		var id int64
		err := rows.Scan(
			&id,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}
		sessionIDs = append(sessionIDs, id)
		// Check for any errors encountered during iteration
		if err = rows.Err(); err != nil {
			return nil, fmt.Errorf("error iterating over rows: %w", err)
		}
	}
	// Check if any session exist
	if len(sessionIDs) <= 0 {
		return nil, nil
	}
	studentSessions, err := s.queryStudentsSessions(c, sessionIDs)
	if err != nil {
		return nil, fmt.Errorf("unable to query student sessions %v", err)
	}
	assessmentSessions, err := s.queryAssessmentSessions(c, sessionIDs)
	if err != nil {
		return nil, fmt.Errorf("unable to query student sessions %v", err)
	}

	// If file has some len build and return
	f, err := buildStudentSessionFile(&studentSessions, &assessmentSessions)
	if err != nil {
		return nil, fmt.Errorf("unable to build session tutor file %v", err)
	}
	// Verify file can be written to a buffer
	var buf bytes.Buffer
	if err := f.Write(&buf); err != nil {
		return nil, fmt.Errorf("failed to verify Excel file creation: %v", err)
	}
	return f, nil
}

func (s *AuthService) querySessionsTutors(sessions []int64) ([]models.AssessmentsData, error) {
	query := `
	SELECT 
	FROM 
		stu_tracker.Assessments_students ast
	JOIN
		stu_tracker.Assessments a
	ON
		a.id = ast.assessment_id
	JOIN
		stu_tracker.Students ss
	ON
		ss.id = ast.student_id
	WHERE 
		ast.session_id = ANY($1);`

	rows, err := s.db.Query(query, pq.Array(sessions))
	if err != nil {
		return nil, fmt.Errorf("error querying locations: %w", err)
	}
	defer rows.Close()
	var assessments []models.AssessmentsData
	for rows.Next() {
		var assessment models.AssessmentsData
		err := rows.Scan(
			&assessment.Title,
			&assessment.MaxScore,
			&assessment.Score,
			&assessment.CreatedAt,
			&assessment.Letter,
			&assessment.Cycle,
			&assessment.Pre,
			&assessment.Mid,
			&assessment.Post,
			&assessment.Version,
			&assessment.StudentName,
			&assessment.StudentLastName,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}
		assessments = append(assessments, assessment)
	}

	return assessments, nil
}

// Get assessments per session given session ids
func (s *AuthService) queryAssessmentSessions(c context.Context, sessions []int64) ([]models.AssessmentsData, error) {
	query := `
	SELECT a.title, a.max_score, ast.score, 
	ast.created_at, a.letter, a.cycle, 
	a.pre, a.mid, a.post, a.version, ss.first_name, ss.last_name, ss.id, ast.session_id
	FROM 
		stu_tracker.Assessments_students ast
	LEFT JOIN
		stu_tracker.Assessments a
	ON
		a.id = ast.assessment_id
	LEFT JOIN
		stu_tracker.Students ss
	ON
		ss.id = ast.student_id
	WHERE 
		ast.session_id = ANY($1);`

	rows, err := s.db.QueryContext(c, query, pq.Array(sessions))
	if err != nil {
		return nil, fmt.Errorf("error querying locations: %w", err)
	}
	defer rows.Close()
	var assessments []models.AssessmentsData
	for rows.Next() {
		var assessment models.AssessmentsData
		err := rows.Scan(
			&assessment.Title,
			&assessment.MaxScore,
			&assessment.Score,
			&assessment.CreatedAt,
			&assessment.Letter,
			&assessment.Cycle,
			&assessment.Pre,
			&assessment.Mid,
			&assessment.Post,
			&assessment.Version,
			&assessment.StudentName,
			&assessment.StudentLastName,
			&assessment.StudentID,
			&assessment.SessionID,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}
		assessments = append(assessments, assessment)
	}

	return assessments, nil
}

// Get all session from a student given sessionid
func (s *AuthService) queryStudentsSessions(c context.Context, sessions []int64) ([]models.StudentSession, error) {
	query := `
		SELECT ss.session_id, 
		st.id,
		st.first_name, 
		st.last_name, 
		ss.absent, 
		ss.duration, 
		ss.created_at,
		CASE
			WHEN ss.subject_id IS NULL THEN 'NA'
			ELSE sj.title
		END AS subject,
		st.grade_level
		FROM 
			stu_tracker.Session_students ss 
		LEFT JOIN 
			stu_tracker.Students st
		ON 
			st.id = ss.student_id
		LEFT JOIN
			stu_tracker.Subjects sj
		ON
			sj.id = ss.subject_id
		WHERE 
			ss.session_id = ANY($1);`
	rows, err := s.db.QueryContext(c, query, pq.Array(sessions))
	if err != nil {
		return nil, fmt.Errorf("error querying locations: %w", err)
	}
	defer rows.Close()
	var students []models.StudentSession
	for rows.Next() {
		var student models.StudentSession
		err := rows.Scan(
			&student.SessionID,
			&student.StudentID,
			&student.FirstName,
			&student.LastName,
			&student.Absent,
			&student.Duration,
			&student.SessionDate,
			&student.Subject,
			&student.Grade,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}
		students = append(students, student)
	}

	return students, nil
}

// Get all sessions given search params
func buildQuery(ss models.RequestDownloadData) (string, []interface{}) {
	argIndex := 1
	var args []interface{}
	var conditions []string
	query := ` `

	if ss.LocationID != nil {
		conditions = append(conditions, fmt.Sprintf("xs.location_id = $%d", argIndex))
		args = append(args, ss.LocationID)
		argIndex++
	}
	if ss.ProgramID != nil {
		conditions = append(conditions, fmt.Sprintf("xs.program_id = $%d", argIndex))
		args = append(args, ss.ProgramID)
		argIndex++
	}
	if ss.SemesterID != nil {
		conditions = append(conditions, fmt.Sprintf("xs.semester_id = $%d", argIndex))
		args = append(args, ss.SemesterID)
		argIndex++
	}
	if !ss.DateStart.IsZero() && !ss.DateEnd.IsZero() {
		conditions = append(conditions, fmt.Sprintf("DATE(xs.session_date) BETWEEN $%d AND $%d", argIndex, argIndex+1))
		args = append(args, ss.DateStart, ss.DateEnd)
		argIndex += 2
	} else if !ss.DateStart.IsZero() {
		conditions = append(conditions, fmt.Sprintf("DATE(xs.session_date) >= $%d", argIndex))
		args = append(args, ss.DateStart)
		argIndex++
	}

	if ss.SubjectID != nil {
		conditions = append(conditions, fmt.Sprintf("xs.subject_id = $%d", argIndex))
		args = append(args, ss.SubjectID)
		argIndex++
	}

	if len(conditions) > 0 {
		query += "WHERE " + strings.Join(conditions, " AND ")
	}
	return query, args
}

// Get all sessions given search params
func buildSessionQuery(ss models.RequestDownloadData) (string, []interface{}) {
	argIndex := 1
	var args []interface{}
	var conditions []string
	query := `SELECT ss.id FROM stu_tracker.Sessions ss `

	if ss.LocationID != nil {
		conditions = append(conditions, fmt.Sprintf("ss.location_id = $%d", argIndex))
		args = append(args, ss.LocationID)
		argIndex++
	}
	if ss.ProgramID != nil {
		conditions = append(conditions, fmt.Sprintf("ss.program_id = $%d", argIndex))
		args = append(args, ss.ProgramID)
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

	if ss.SubjectID != nil {
		conditions = append(conditions, fmt.Sprintf("ss.subject_id = $%d", argIndex))
		args = append(args, ss.SubjectID)
		argIndex++
	}

	if len(conditions) > 0 {
		query += "WHERE " + strings.Join(conditions, " AND ")
	}
	return query, args
}

// Build SQL query for tutor data
func buildSessionQueryTutorData(ss models.RequestDownloadData) (string, []interface{}) {
	argIndex := 1
	var args []interface{}
	var conditions []string
	query := `SELECT 
			ss.id, 
			ss.tutor_id, 
			ss.session_date, 
			ss.substitute, 
			ss.student_count, 
			ss.start_time, 
			ss.duration, 
			ss.notes, 
			t.first_name, 
			t.last_name,
			pg.program_name
		FROM stu_tracker.Sessions ss
		JOIN stu_tracker.Tutors t
		ON t.id = ss.tutor_id 
		JOIN stu_tracker.Programs pg
		ON pg.id = ss.program_id `

	if ss.LocationID != nil {
		conditions = append(conditions, fmt.Sprintf("ss.location_id = $%d", argIndex))
		args = append(args, ss.LocationID)
		argIndex++
	}
	if ss.ProgramID != nil {
		conditions = append(conditions, fmt.Sprintf("ss.program_id = $%d", argIndex))
		args = append(args, ss.ProgramID)
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
	if ss.SubjectID != nil {
		conditions = append(conditions, fmt.Sprintf("ss.subject_id = $%d", argIndex))
		args = append(args, ss.SubjectID)
		argIndex++
	}

	if len(conditions) > 0 {
		query += "WHERE " + strings.Join(conditions, " AND ")
	}
	return query, args
}

// Helper function for safe date formatting
func formatSessionDate(date time.Time) string {
	if date.IsZero() {
		return "N/A" // Or whatever default you prefer
	}
	return date.Format("2006-01-02") // YYYY-MM-DD format
}

func buildStudentSessionFile(studentSessions *[]models.StudentSession, studentAssessments *[]models.AssessmentsData) (*excelize.File, error) {
	f := excelize.NewFile()
	sheet := "Student report dataframe"
	if err := f.SetSheetName("Sheet1", sheet); err != nil {
		return nil, fmt.Errorf("failed to rename sheet: %v", err)
	}
	fmt.Println("StudentSession Size: ", len(*studentSessions))
	// Col values A , B, C, D ...
	headers := []string{"SID", "First name", "Last name", "Session id", "Subject",
		"Duration", "Session Date", "Absent", "Notes", "Grade", "Assessment title", "Letter", "Cycle", "Pre", "Mid", "Post", "Version", "Score", "Max score"}
	for i, header := range headers {
		cell := fmt.Sprintf("%c1", 'A'+i)
		if err := f.SetCellValue(sheet, cell, header); err != nil {
			return nil, fmt.Errorf("failed to set header %s: %v", header, err)
		}
	}
	for i, sessions := range *studentSessions {
		rowNum := i + 2
		if i == 0 {
			if err := f.SetColWidth(sheet, "E", "E", 20); err != nil {
				return nil, fmt.Errorf("failed to set column width: %v", err)
			}
			if err := f.SetColWidth(sheet, "H", "H", 20); err != nil {
				return nil, fmt.Errorf("failed to set column width: %v", err)
			}
			if err := f.SetColWidth(sheet, "L", "L", 20); err != nil {
				return nil, fmt.Errorf("failed to set column width: %v", err)
			}
		}
		cells := map[string]interface{}{
			fmt.Sprintf("A%d", rowNum): *sessions.StudentID,
			fmt.Sprintf("B%d", rowNum): sessions.FirstName,
			fmt.Sprintf("C%d", rowNum): sessions.LastName,
			fmt.Sprintf("D%d", rowNum): *sessions.SessionID,
			fmt.Sprintf("E%d", rowNum): sessions.Subject,
			fmt.Sprintf("F%d", rowNum): sessions.Duration,
			fmt.Sprintf("G%d", rowNum): formatSessionDate(sessions.SessionDate),
			fmt.Sprintf("H%d", rowNum): sessions.Absent,
			fmt.Sprintf("I%d", rowNum): sessions.Notes,
			fmt.Sprintf("J%d", rowNum): sessions.Grade,
		}
		// check each cell for errors and insert
		for cell, value := range cells {
			if err := f.SetCellValue(sheet, cell, value); err != nil {
				return nil, fmt.Errorf("unable to insert cell into sheet %v", err)
			}
		}
	}
	// Create a map of sessionID
	studentAssessmentsMap := make(map[string]models.AssessmentsData)
	for _, data := range *studentAssessments {
		compositeKey := fmt.Sprintf("%d:%d", *data.SessionID, *data.StudentID)
		studentAssessmentsMap[compositeKey] = data
	}
	// Get all rows from the sheet
	rows, err := f.GetRows(sheet)
	if err != nil {
		return nil, fmt.Errorf("unable to get rows %v", err)
	}
	for i, row := range rows {
		if i == 0 {
			continue
		}
		if len(row) == 0 {
			continue
		}
		sessionID, err := strconv.ParseInt(row[3], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("unable to parse sessionID")
		}
		studentID, err := strconv.ParseInt(row[0], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("unable to parse sessionID")
		}
		compositeKey := fmt.Sprintf("%d:%d", sessionID, studentID)
		// Check if sessionID found in assessmentMAP
		if studentAssessmentsData, exist := studentAssessmentsMap[compositeKey]; exist {
			rowNum := i + 1
			cells := map[string]interface{}{
				fmt.Sprintf("L%d", rowNum): studentAssessmentsData.Title,
				fmt.Sprintf("M%d", rowNum): studentAssessmentsData.Letter,
				fmt.Sprintf("N%d", rowNum): studentAssessmentsData.Cycle,
				fmt.Sprintf("O%d", rowNum): studentAssessmentsData.Pre,
				fmt.Sprintf("P%d", rowNum): studentAssessmentsData.Mid,
				fmt.Sprintf("Q%d", rowNum): studentAssessmentsData.Post,
				fmt.Sprintf("R%d", rowNum): studentAssessmentsData.Version,
				fmt.Sprintf("S%d", rowNum): studentAssessmentsData.Score,
				fmt.Sprintf("T%d", rowNum): studentAssessmentsData.MaxScore,
			}
			// check each cell for errors and insert
			for cell, value := range cells {
				if err := f.SetCellValue(sheet, cell, value); err != nil {
					return nil, fmt.Errorf("unable to insert cell into sheet %v", err)
				}
			}
		}
	}

	return f, nil
}

// Build file from sQL query for tutor data
func buildSessionTutorFile(sessions []models.TutorSessionData) (*excelize.File, error) {
	f := excelize.NewFile()
	sheetName := "Tutor sessions dataframe"
	if err := f.SetSheetName("Sheet1", sheetName); err != nil {
		return nil, fmt.Errorf("failed to rename sheet: %v", err)
	}
	// Set headers with error checking
	headers := []string{"First name", "Last name", "SessionID", "TutorID", "StudentCount", "SessionDate",
		"Duration", "Substitute", "Program", "Start time"}
	for i, header := range headers {
		cell := fmt.Sprintf("%c1", 'A'+i)
		if err := f.SetCellValue(sheetName, cell, header); err != nil {
			return nil, fmt.Errorf("failed to set header %s: %v", header, err)
		}
	}
	// Insert data with error checking
	for i, tutor := range sessions {
		rowNum := i + 2
		// Set column width (with error checking)
		if i == 0 {
			if err := f.SetColWidth(sheetName, "A", "A", 20); err != nil {
				return nil, fmt.Errorf("failed to set column width: %v", err)
			}
			if err := f.SetColWidth(sheetName, "B", "B", 20); err != nil {
				return nil, fmt.Errorf("failed to set column width: %v", err)
			}
			if err := f.SetColWidth(sheetName, "F", "F", 15); err != nil {
				return nil, fmt.Errorf("failed to set column width: %v", err)
			}
		}
		// Write data with error checking for each cell
		cells := map[string]interface{}{
			fmt.Sprintf("A%d", rowNum): tutor.TutorName,
			fmt.Sprintf("B%d", rowNum): tutor.TutorLastName,
			fmt.Sprintf("C%d", rowNum): *tutor.SessionID,
			fmt.Sprintf("D%d", rowNum): *tutor.TutorID,
			fmt.Sprintf("E%d", rowNum): tutor.StudentCount,
			fmt.Sprintf("F%d", rowNum): formatSessionDate(tutor.SessionDate),
			fmt.Sprintf("G%d", rowNum): tutor.Duration,
			fmt.Sprintf("H%d", rowNum): tutor.Substitute,
			fmt.Sprintf("I%d", rowNum): tutor.Program,
			fmt.Sprintf("J%d", rowNum): tutor.StartTime,
		}
		for cell, value := range cells {
			if err := f.SetCellValue(sheetName, cell, value); err != nil {
				return nil, fmt.Errorf("failed to set cell %s: %v", cell, err)
			}
		}
	}
	return f, nil
}
