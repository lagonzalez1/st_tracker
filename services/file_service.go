package services

import (
	"bytes"
	"context"
	"fmt"
	"sort"
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
			&session.ProgramID,
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
	f, err := buildSessionTutorFile(sessions, ss.SortKey)
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
	fmt.Println("queryAssessmentSessions length", len(assessmentSessions))
	if err != nil {
		return nil, fmt.Errorf("unable to query student sessions %v", err)
	}

	// If file has some len build and return
	f, err := buildStudentSessionFile(&studentSessions, &assessmentSessions, ss.SortKey)
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

// Get assessments per session given session ids
func (s *AuthService) queryAssessmentSessions(c context.Context, sessions []int64) ([]models.AssessmentsData, error) {
	fmt.Println("queryAssessmentSessions the sessions", sessions)
	query := `
	SELECT a.title, a.max_score, ast.score, 
	sn.session_date, a.letter, a.cycle, 
	a.pre, a.mid, a.post, a.version, ss.first_name, ss.last_name, ss.id, ast.session_id
	FROM 
		stu_tracker.Assessments_students ast
	JOIN
		stu_tracker.Assessments a
	ON
		a.id = ast.assessment_id
	LEFT JOIN
		stu_tracker.Students ss
	ON
		ss.id = ast.student_id
	LEFT JOIN
		stu_tracker.Sessions sn
	ON
		sn.id = ast.session_id
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
		sk.session_date,
		CASE
			WHEN ss.subject_id IS NULL THEN 'NA'
			ELSE sj.title
		END AS subject,
		st.grade_level,
		st.timeframe,
		st.timeframe_start,
		st.timeframe_end
		FROM 
			stu_tracker.Session_students ss 
		JOIN
			stu_tracker.Sessions sk
		ON
			sk.id = ss.session_id
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
			&student.Timeframe,
			&student.TimeframeStart,
			&student.TimeframeEnd,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}
		students = append(students, student)
	}

	return students, nil
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
			pg.program_name,
			pg.id
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

func buildStudentSessionFile(studentSessions *[]models.StudentSession, studentAssessments *[]models.AssessmentsData, sortKey string) (*excelize.File, error) {
	switch sortKey {
	case "group_students":
		f, err := buildSessionStudentFileByStudents(studentSessions, studentAssessments)
		if err != nil {
			return nil, err
		}
		return f, nil
	case "all":
		f := excelize.NewFile()
		sheet := "Student report dataframe"
		if err := f.SetSheetName("Sheet1", sheet); err != nil {
			return nil, fmt.Errorf("failed to rename sheet: %v", err)
		}
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
				if err := f.SetColWidth(sheet, "G", "G", 20); err != nil {
					return nil, fmt.Errorf("failed to set column width: %v", err)
				}
				if err := f.SetColWidth(sheet, "E", "E", 20); err != nil {
					return nil, fmt.Errorf("failed to set column width: %v", err)
				}
				if err := f.SetColWidth(sheet, "L", "L", 20); err != nil {
					return nil, fmt.Errorf("failed to set column width: %v", err)
				}
				if err := f.SetColWidth(sheet, "K", "K", 20); err != nil {
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
					fmt.Sprintf("K%d", rowNum): studentAssessmentsData.Title,
					fmt.Sprintf("L%d", rowNum): studentAssessmentsData.Letter,
					fmt.Sprintf("M%d", rowNum): studentAssessmentsData.Cycle,
					fmt.Sprintf("N%d", rowNum): studentAssessmentsData.Pre,
					fmt.Sprintf("O%d", rowNum): studentAssessmentsData.Mid,
					fmt.Sprintf("P%d", rowNum): studentAssessmentsData.Post,
					fmt.Sprintf("Q%d", rowNum): studentAssessmentsData.Version,
					fmt.Sprintf("R%d", rowNum): studentAssessmentsData.Score,
					fmt.Sprintf("S%d", rowNum): studentAssessmentsData.MaxScore,
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
	default:
		return nil, fmt.Errorf("no sort key provided")
	}

}

func getLastFirstSessionDate(sessions []models.TutorSessionData) ([]string, error) {
	if len(sessions) == 0 {
		return nil, fmt.Errorf("no sessions found")
	}
	// Sort the session list by sessionDate
	sort.Slice(sessions, func(i, j int) bool {
		return sessions[i].SessionDate.Before(sessions[j].SessionDate)
	})
	start := sessions[0].SessionDate
	end := sessions[len(sessions)-1].SessionDate
	// Normalize time to midnight
	start = time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, start.Location())
	end = time.Date(end.Year(), end.Month(), end.Day(), 0, 0, 0, 0, end.Location())

	var dates []string
	for d := start; !d.After(end); d = d.AddDate(0, 0, 1) {
		dates = append(dates, d.Format("2006-01-02"))
	}
	return dates, nil
}
func getLastFirstSessionDateStudents(sessions []models.StudentSession) ([]string, error) {
	if len(sessions) == 0 {
		return nil, fmt.Errorf("no sessions found")
	}
	// Sort the session list by sessionDate
	sort.Slice(sessions, func(i, j int) bool {
		return sessions[i].SessionDate.Before(sessions[j].SessionDate)
	})
	start := sessions[0].SessionDate
	end := sessions[len(sessions)-1].SessionDate
	// Normalize time to midnight
	start = time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, start.Location())
	end = time.Date(end.Year(), end.Month(), end.Day(), 0, 0, 0, 0, end.Location())

	var dates []string
	for d := start; !d.After(end); d = d.AddDate(0, 0, 1) {
		dates = append(dates, d.Format("2006-01-02"))
	}
	return dates, nil
}

// Build file from sQL query for tutor data
func buildSessionTutorFile(sessions []models.TutorSessionData, sortKey string) (*excelize.File, error) {

	switch sortKey {
	case "group_tutors":
		f, err := buildSessionTutorFileGroupByTutor(sessions)
		if err != nil {
			return nil, err
		}
		return f, nil
	case "all":
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
	default:
		return nil, fmt.Errorf("no sort key provided")
	}
}

func buildSessionTutorFileGroupByTutor(sessions []models.TutorSessionData) (*excelize.File, error) {
	f := excelize.NewFile()
	sheetName := "Tutor grouped sessions dataframe"
	if err := f.SetSheetName("Sheet1", sheetName); err != nil {
		return nil, fmt.Errorf("failed to rename sheet: %v", err)
	}
	startHeaders := []string{"Tutor name", "Program", "Substitute", "Student count"}
	// Generate the first and last and inbetween dates for headers
	datesHeaders, err := getLastFirstSessionDate(sessions)
	if err != nil {
		return nil, err
	}
	headers := append(startHeaders, datesHeaders...)
	for i, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		if err := f.SetCellValue(sheetName, cell, header); err != nil {
			return nil, fmt.Errorf("failed to set header %s: %v", header, err)
		}
	}
	dateColMap := make(map[string]int)
	for i, date := range datesHeaders {
		colIndex := len(startHeaders) + i + 1 // 1-based index for Excel
		dateColMap[date] = colIndex
	}
	for i := range headers {
		colLetter, _ := excelize.ColumnNumberToName(i + 1)
		_ = f.SetColWidth(sheetName, colLetter, colLetter, 12) // Set width to 15
	}
	// Group by
	grouped := make(map[int64]map[int64][]models.TutorSessionData)
	for _, tutor := range sessions {
		tutorID := int64(*tutor.TutorID)
		programID := int64(*tutor.ProgramID)
		if _, ok := grouped[tutorID]; !ok {
			grouped[tutorID] = make(map[int64][]models.TutorSessionData)
		}
		grouped[tutorID][programID] = append(grouped[tutorID][programID], tutor)
	}
	row := 2
	for _, programs := range grouped {
		for _, sessions := range programs {
			studentRunningCount := 0
			tutorName := sessions[0].TutorName
			programName := sessions[0].Program
			substitute := sessions[0].Substitute
			programStudentCount := 0
			for _, s := range sessions {
				programStudentCount += s.StudentCount
			}
			studentRunningCount += programStudentCount
			baseCol := []interface{}{tutorName, programName, substitute, studentRunningCount}
			for i, val := range baseCol {
				cell, _ := excelize.CoordinatesToCellName(i+1, row)
				f.SetCellValue(sheetName, cell, val)
			}
			for _, session := range sessions {
				sessionDate := session.SessionDate.Format("2006-01-02")
				startTime := session.StartTime
				if colIdx, ok := dateColMap[sessionDate]; ok {
					cell, _ := excelize.CoordinatesToCellName(colIdx, row)
					existing, _ := f.GetCellValue(sheetName, cell)
					if existing != "" {
						startTime = existing + ", " + startTime
					}
					f.SetCellValue(sheetName, cell, startTime)
				}
			}
			row++
		}
	}
	return f, nil
}

// This can be optimized please fix later.
func buildSessionStudentFileByStudents(studentSessions *[]models.StudentSession, studentAssessments *[]models.AssessmentsData) (*excelize.File, error) {
	f := excelize.NewFile()
	sheetName := "Student sessions dataframe"
	if err := f.SetSheetName("Sheet1", sheetName); err != nil {
		return nil, fmt.Errorf("failed to rename sheet: %v", err)
	}
	startHeaders := []string{"First name", "Last name", "Subjects", "Duration", "Timeframe"}
	// Generate the first and last and inbetween dates for headers
	datesHeaders, err := getLastFirstSessionDateStudents(*studentSessions)
	if err != nil {
		return nil, err
	}
	headers := append(startHeaders, datesHeaders...)
	for i, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		if err := f.SetCellValue(sheetName, cell, header); err != nil {
			return nil, fmt.Errorf("failed to set header %s: %v", header, err)
		}
	}
	// build the date map
	dateColMap := make(map[string]int)
	for i, date := range datesHeaders {
		colIndex := len(startHeaders) + i + 1 // 1-based index for Excel
		dateColMap[date] = colIndex
	}
	// Give width to the cols
	for i := range headers {
		colLetter, _ := excelize.ColumnNumberToName(i + 1)
		_ = f.SetColWidth(sheetName, colLetter, colLetter, 12) // Set width to 15
	}
	grouped := make(map[int64][]models.StudentSession)
	groupedAssessments := make(map[int64][]models.AssessmentsData)
	for _, session := range *studentSessions {
		studentID := int64(*session.StudentID)
		grouped[studentID] = append(grouped[studentID], session)
	}
	for _, assessments := range *studentAssessments {
		studentID := int64(*assessments.StudentID)
		groupedAssessments[studentID] = append(groupedAssessments[studentID], assessments)
	}
	row := 2
	for _, sessions := range grouped {
		// Assume all sessions belong to the same student
		if len(sessions) == 0 {
			continue
		}
		firstSession := sessions[0]
		fmt.Println(firstSession)
		firstName := firstSession.FirstName
		lastName := firstSession.LastName
		var subjectList []string
		var durationTotal int
		var timeframe string

		if firstSession.Timeframe {
			timeframe = fmt.Sprintf("%s-%s", *firstSession.TimeframeStart, *firstSession.TimeframeEnd)
		}

		for _, session := range sessions {
			subjectList = append(subjectList, session.Subject)
			durationTotal += session.Duration

			sessionDate := session.SessionDate.Format("2006-01-02")
			if colIdx, ok := dateColMap[sessionDate]; ok {
				cell, _ := excelize.CoordinatesToCellName(colIdx, row)
				existing, _ := f.GetCellValue(sheetName, cell)
				val := "T"
				if existing != "" {
					val = existing + ", T"
				}
				f.SetCellValue(sheetName, cell, val)
			}
		}
		// Remove duplicates and join subjects
		subjects := strings.Join(removeDuplicates(subjectList), ", ")
		baseCol := []interface{}{firstName, lastName, subjects, durationTotal, timeframe}
		for i, val := range baseCol {
			cell, _ := excelize.CoordinatesToCellName(i+1, row)
			f.SetCellValue(sheetName, cell, val)
		}

		row++ // Now we move to the next row
	}

	assessmentHeaders := []string{"First name", "Last name", "Assessment title", "Version", "Cycle", "Score", "Max score", "pre", "mid", "post", "Created at"}

	// Set assessment headers
	for i, header := range assessmentHeaders {
		cell, _ := excelize.CoordinatesToCellName(i+1, row)
		if err := f.SetCellValue(sheetName, cell, header); err != nil {
			return nil, fmt.Errorf("unable to set assessment headers")
		}
	}
	row += 1
	for _, assessments := range groupedAssessments {
		if len(assessments) == 0 {
			continue
		}
		firstAssessment := assessments[0]
		firstName := firstAssessment.StudentName
		lastName := firstAssessment.StudentLastName
		for _, val := range assessments {
			assessmentTitle := val.Title
			assessmentVersion := val.Version
			assessmentCycle := val.Cycle
			score := val.Score
			maxScore := val.MaxScore
			pre := val.Pre
			mid := val.Mid
			post := val.Post
			createdAt := val.CreatedAt
			baseCol := []interface{}{firstName, lastName, assessmentTitle, assessmentVersion, assessmentCycle, score, maxScore, pre, mid, post, createdAt}
			for i, item := range baseCol {
				cell, _ := excelize.CoordinatesToCellName(i+1, row)
				f.SetCellValue(sheetName, cell, item)
			}
			row++
		}

	}
	return f, nil
}
func removeDuplicates(input []string) []string {
	seen := make(map[string]bool)
	var result []string
	for _, val := range input {
		if !seen[val] {
			seen[val] = true
			result = append(result, val)
		}
	}
	return result
}
