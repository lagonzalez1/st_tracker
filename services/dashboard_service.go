package services

import (
	"context"
	"database/sql"
	"fmt"
	"time"
	"tracker/app/models"

	"github.com/lib/pq"
)

func (s *AuthService) GetLocationsByID(ctx context.Context, id int64, role string) ([]models.ResponseRequestLocations, error) {
	var query string
	query += `
		SELECT loc.id, loc.name, loc.address, loc.city, loc.state, loc.zip_code, loc.created_at, loc.district_id
		FROM stu_tracker.Locations loc WHERE loc.organization_id = $1;`
	rows, err := s.db.QueryContext(ctx, query, id)
	if err != nil {
		return nil, fmt.Errorf("error querying locations: %w", err)
	}
	defer rows.Close()

	var locations []models.ResponseRequestLocations
	for rows.Next() {
		var location models.ResponseRequestLocations
		err := rows.Scan(
			&location.ID,
			&location.Name,
			&location.Address,
			&location.City,
			&location.State,
			&location.ZipCode,
			&location.CreatedAt,
			&location.DistrictId,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}
		// Check for any errors encountered during iteration
		if err = rows.Err(); err != nil {
			return nil, fmt.Errorf("error iterating over rows: %w", err)
		}
		locations = append(locations, location)
	}

	return locations, nil
}

func (s *AuthService) GetSubjectById(ctx context.Context, id int64, role string) ([]models.SubjectList, error) {
	var query string
	query += `
		SELECT sb.id, sb.title, sb.description, sb.organization_id, sb.created_at
		FROM stu_tracker.Subjects sb WHERE sb.organization_id = $1;`
	rows, err := s.db.QueryContext(ctx, query, id)
	if err != nil {
		return nil, fmt.Errorf("error querying locations: %w", err)
	}
	defer rows.Close()

	var subjects []models.SubjectList
	for rows.Next() {
		var subject models.SubjectList
		err := rows.Scan(
			&subject.ID,
			&subject.Title,
			&subject.Description,
			&subject.OrganizationId,
			&subject.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}
		subjects = append(subjects, subject)
	}

	// Check for any errors encountered during iteration
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}

	return subjects, nil
}

func (s *AuthService) GetSubjectByLocation(ctx context.Context, org_id int64, loc_id int64) ([]models.SubjectList, error) {
	var query string
	query += `
		SELECT s.id, s.title, s.description, s.organization_id, s.created_at
		FROM stu_tracker.Location_subjects ls
		LEFT JOIN stu_tracker.Subjects s 
		ON s.id = ls.subject_id
		WHERE ls.location_id = $1 AND s.organization_id = $2;`
	rows, err := s.db.QueryContext(ctx, query, loc_id, org_id)
	if err != nil {
		return nil, fmt.Errorf("error querying locations: %w", err)
	}
	defer rows.Close()

	var subjects []models.SubjectList
	for rows.Next() {
		var subject models.SubjectList
		err := rows.Scan(
			&subject.ID,
			&subject.Title,
			&subject.Description,
			&subject.OrganizationId,
			&subject.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}
		subjects = append(subjects, subject)
	}

	// Check for any errors encountered during iteration
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}

	return subjects, nil
}

func (s *AuthService) GetStudentsByID(c context.Context, id int64, role string, locationId int64, tutorId int64) ([]models.ResponseRequestStudentList, error) {
	var query string
	var rows *sql.Rows
	var err error
	// If Admin or root return all students given organization and location id
	if role == "ADMIN" || role == "ROOT" {
		query += `
		SELECT stu.id, stu.first_name, stu.last_name, stu.middle_name, 
		stu.location_id, stu.email, stu.grade_level, stu.active, stu.period, 
		stu.created_at, stu.semester_id, direct_partnership
		FROM stu_tracker.Students stu
		JOIN stu_tracker.Locations loc
		ON stu.location_id = loc.id
		WHERE loc.organization_id = $1 AND loc.id = $2;`
		rows, err = s.db.QueryContext(c, query, id, locationId)
	}
	// If tutor return all students and if tutor_id is marked return such students
	if role == "TUTOR" {
		query += `
		SELECT stu.id, stu.first_name, stu.last_name, stu.middle_name, 
		stu.location_id, stu.email, stu.grade_level, stu.active, stu.period, 
		stu.created_at, stu.semester_id, direct_partnership
		FROM stu_tracker.Students stu
		JOIN stu_tracker.Locations loc
		ON stu.location_id = loc.id
		WHERE loc.id = $1
		AND ( 
			direct_partnership = FALSE 
			OR (direct_partnership = TRUE AND tutor_id = $2)
		);`
		rows, err = s.db.QueryContext(c, query, locationId, tutorId)
	}

	if err != nil {
		return nil, fmt.Errorf("error querying Students: %w", err)
	}
	defer rows.Close()

	var students []models.ResponseRequestStudentList
	for rows.Next() {
		var student models.ResponseRequestStudentList
		var middleName sql.NullString
		var gradeLevel sql.NullInt64
		err = rows.Scan(
			&student.ID,
			&student.FirstName,
			&student.LastName,
			&middleName,
			&student.LocationId,
			&student.Email,
			&gradeLevel,
			&student.Active,
			&student.Period,
			&student.CreatedAt,
			&student.SemesterId,
			&student.DirectPartnership,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %v", err)
		}
		// Convert nullable fields
		if middleName.Valid {
			student.MiddleName = middleName.String
		}
		if gradeLevel.Valid {
			grade := int(gradeLevel.Int64)
			student.GradeLevel = grade
		}

		students = append(students, student)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}
	return students, nil
}

func (s *AuthService) GetAdminStaffById(c context.Context, id int64, role string) ([]models.ResponseRequestAdminList, error) {
	var query string
	query += `SELECT id,fullname, email, region, state 
			  FROM stu_tracker.Admin_staff
			  WHERE organization_id = $1;`

	rows, err := s.db.QueryContext(c, query, id)
	if err != nil {
		return nil, fmt.Errorf("error querying Admin: %w", err)
	}
	defer rows.Close()

	var admins []models.ResponseRequestAdminList
	for rows.Next() {
		var admin models.ResponseRequestAdminList
		err := rows.Scan(
			&admin.ID,
			&admin.Fullname,
			&admin.Email,
			&admin.Region,
			&admin.State,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}
		admins = append(admins, admin)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}
	return admins, nil
}

func (s *AuthService) GetDistrictsById(c context.Context, id int64, role string) ([]models.ResponseRequestDistrictList, error) {
	var query string
	query += `SELECT id, name, city ,state, region, created_at
			  FROM stu_tracker.district ds WHERE ds.organization_id = $1;`
	rows, err := s.db.QueryContext(c, query, id)
	if err != nil {
		return nil, fmt.Errorf("error quering get districts by id %w", err)
	}
	defer rows.Close()
	var locations []models.ResponseRequestDistrictList
	for rows.Next() {
		var location models.ResponseRequestDistrictList
		err := rows.Scan(
			&location.ID,
			&location.Name,
			&location.City,
			&location.State,
			&location.Region,
			&location.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}
		locations = append(locations, location)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error scanning for districts")
	}
	return locations, nil
}

func (s *AuthService) GetProgramsId(c context.Context, id int64, role string) ([]models.ResponseRequestProgramList, error) {
	var query string
	query += `SELECT id, program_name
			  FROM stu_tracker.programs pg
			  WHERE pg.organization_id = $1;`
	rows, err := s.db.QueryContext(c, query, id)
	if err != nil {
		return nil, fmt.Errorf("error quering get districts by id %w", err)
	}
	defer rows.Close()
	var programs []models.ResponseRequestProgramList
	for rows.Next() {
		var program models.ResponseRequestProgramList
		err := rows.Scan(
			&program.ID,
			&program.ProgramName,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}
		programs = append(programs, program)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error scanning for districts")
	}
	return programs, nil
}

func (s *AuthService) GetMaterialsById(c context.Context, id int64, role string) ([]models.ResponseRequestMaterialsList, error) {
	var query string
	query += `SELECT id, title, external_link, description, version, pre, mid, post, visible, created_at, program_id
			  FROM stu_tracker.materials mt
			  WHERE mt.organization_id = $1;`
	rows, err := s.db.QueryContext(c, query, id)
	if err != nil {
		return nil, fmt.Errorf("error quering get districts by id %w", err)
	}
	defer rows.Close()
	var materials []models.ResponseRequestMaterialsList
	for rows.Next() {
		var material models.ResponseRequestMaterialsList
		err := rows.Scan(
			&material.ID,
			&material.Title,
			&material.ExternalLink,
			&material.Description,
			&material.Version,
			&material.Pre,
			&material.Mid,
			&material.Post,
			&material.Visible,
			&material.CreatedAt,
			&material.ProgramId,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}
		materials = append(materials, material)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error scanning for districts")
	}
	return materials, nil
}

func (s *AuthService) GetTutorsById(c context.Context, id int64, role string, locid int64) ([]models.ResponseRequestTutorsList, error) {
	var query string
	var args []interface{}
	if locid <= 0 {
		query += `SELECT id, first_name, last_name, email, created_at, location_id
			  FROM stu_tracker.Tutors tr
			  WHERE tr.organization_id = $1;`
		args = append(args, id)
	} else {
		query += `SELECT tr.id, tr.first_name, tr.last_name, tr.email, tr.created_at, tr.location_id
					FROM stu_tracker.Tutors tr
					WHERE tr.location_id = $1
					OR EXISTS (
						SELECT 1
						FROM stu_tracker.Tutor_locations o
						WHERE o.tutor_id = tr.id AND o.location_id = $2);`

		args = append(args, locid, locid)
	}
	rows, err := s.db.QueryContext(c, query, args...)
	if err != nil {
		return nil, fmt.Errorf("error quering get districts by id %w", err)
	}
	defer rows.Close()
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
			return nil, fmt.Errorf("error scanning row: %w", err)
		}
		tutors = append(tutors, tutor)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error scanning for districts")
	}
	return tutors, nil
}

func (s *AuthService) GetSemestersById(ctx context.Context, id int64, role string) ([]models.ResponseRequestSemesterList, error) {
	var query string
	query += `SELECT id, year, title, date_start, date_end, active
			  FROM stu_tracker.Semester sm
			  WHERE sm.organization_id = $1;`
	rows, err := s.db.QueryContext(ctx, query, id)
	if err != nil {
		return nil, fmt.Errorf("error quering get districts by id %w", err)
	}
	defer rows.Close()
	var semesters []models.ResponseRequestSemesterList
	for rows.Next() {
		var semester models.ResponseRequestSemesterList
		err := rows.Scan(
			&semester.ID,
			&semester.Year,
			&semester.Title,
			&semester.DateStart,
			&semester.DateEnd,
			&semester.Active,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}
		semesters = append(semesters, semester)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error scanning for districts")
	}
	return semesters, nil
}

func (s *AuthService) GetSemesterLocationById(c context.Context, role string, location_id int64, idd int64) ([]models.ResponseRequestSemesterLocationList, error) {
	var query string
	query += `SELECT 
				sl.location_id, sl.semester_id,
				sl.created_at, ss.title, ss.year, 
				ss.date_start, ss.date_end
			  FROM 
			  	stu_tracker.Semester_Location sl
			  JOIN
			  	stu_tracker.Semester ss
			  ON
			  	ss.id = sl.semester_id
			  WHERE 
			  	sl.organization_id = $1 AND sl.location_id = $2;`
	rows, err := s.db.QueryContext(c, query, idd, location_id)
	if err != nil {
		return nil, fmt.Errorf("error quering get districts by id %w", err)
	}
	defer rows.Close()
	var semesters []models.ResponseRequestSemesterLocationList
	for rows.Next() {
		var semester models.ResponseRequestSemesterLocationList
		err := rows.Scan(
			&semester.LocationId,
			&semester.SemesterID,
			&semester.CreatedAt,
			&semester.Title,
			&semester.Year,
			&semester.DateStart,
			&semester.DateEnd,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}
		semesters = append(semesters, semester)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error scanning for districts")
	}
	return semesters, nil
}

func (s *AuthService) GetAssessmentsById(c context.Context, id int64, role string) ([]models.ResponseAssessmentList, error) {
	var query string
	query += `SELECT 
			aas.id, aas.title, aas.description, aas.letter,
			aas.cycle, aas.alpha_identifier, aas.external_link, 
			aas.max_score, aas.subject_id, aas.material_id, 
			aas.created_at, 
			COALESCE(sb.title, 'NA') AS subject_name, 
			COALESCE(pg.program_name, 'NA') AS program_name, 
			pg.id as program_id,
			aas.version, aas.pre, aas.mid, aas.post, aas.visible, aas.easy_score
			FROM stu_tracker.Assessments aas 
			LEFT JOIN stu_tracker.Subjects sb
			ON sb.id = aas.subject_id
			LEFT JOIN stu_tracker.Programs pg
			ON pg.id = aas.program_id
			WHERE aas.organization_id = $1;`
	rows, err := s.db.QueryContext(c, query, id)
	if err != nil {
		return nil, fmt.Errorf("error quering get districts by id %w", err)
	}
	defer rows.Close()
	var assessments []models.ResponseAssessmentList
	for rows.Next() {
		var assessment models.ResponseAssessmentList
		err := rows.Scan(
			&assessment.ID,
			&assessment.Title,
			&assessment.Description,
			&assessment.Letter,
			&assessment.Cycle,
			&assessment.AlphaIdentifier,
			&assessment.ExternalLink,
			&assessment.MaxScore,
			&assessment.SubjectId,
			&assessment.MaterialID,
			&assessment.CreatedAt,
			&assessment.SubjectName,
			&assessment.ProgramName,
			&assessment.ProgramId,
			&assessment.Version,
			&assessment.Pre,
			&assessment.Mid,
			&assessment.Post,
			&assessment.Visible,
			&assessment.EasyScore,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}
		assessments = append(assessments, assessment)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error scanning for districts")
	}
	return assessments, nil
}

func (s *AuthService) GetAssessmentsQuestionsChoice(c context.Context, assessment_id int64) ([]models.ResponseAssessmentQuestionsChoice, error) {
	query := `
		SELECT 
			q.id as question_id,
			q.assessment_id,
			q.image_url,
			q.question_text,
			q.question_type,
			q.points,
			q.order_number,
			c.id as choice_id,
			c.choice_text,
			CAST(c.is_correct AS TEXT) as is_correct,
			COALESCE(c.order_number, 0) as choice_order
		FROM 
			stu_tracker.Questions q
		LEFT JOIN 
			stu_tracker.Choices c ON c.question_id = q.id
		WHERE q.assessment_id = $1
		ORDER BY q.order_number ASC, c.order_number ASC;`
	rows, err := s.db.QueryContext(c, query, assessment_id)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	var results []models.ResponseAssessmentQuestionsChoice

	for rows.Next() {
		var r models.ResponseAssessmentQuestionsChoice
		err := rows.Scan(
			&r.QuestionID,
			&r.AssessmentID,
			&r.ImageURL,
			&r.QuestionText,
			&r.QuestionType,
			&r.Points,
			&r.OrderNumber,
			&r.ChoiceID,
			&r.ChoiceText,
			&r.IsCorrect,
			&r.ChoiceOrderNumber,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		results = append(results, r)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return results, nil
}

func (s *AuthService) GetProgramsByLocation(c context.Context, locId int64, org_id int64) ([]models.ResponseRequestProgramList, error) {
	var query string
	query += `
		SELECT p.id, p.program_name, p.created_at, lp.location_id
		FROM stu_tracker.Location_programs lp
		JOIN stu_tracker.Programs p ON lp.program_id = p.id
		WHERE lp.location_id = $1 AND p.organization_id = $2`

	rows, err := s.db.QueryContext(c, query, locId, org_id)
	if err != nil {
		return nil, fmt.Errorf("error querying program: %w", err)
	}
	defer rows.Close()

	var programs []models.ResponseRequestProgramList
	for rows.Next() {
		var program models.ResponseRequestProgramList
		err := rows.Scan(
			&program.ID,
			&program.ProgramName,
			&program.CreatedAt,
			&program.LocationID,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}
		programs = append(programs, program)
	}

	// Check for any errors encountered during iteration
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}

	return programs, nil
}

func (s *AuthService) GetProgramsByIds(c context.Context, locId []int64, org_id int64) ([]models.ResponseRequestProgramList, error) {
	var query string
	query += `
		SELECT p.id, p.program_name, p.created_at
		FROM stu_tracker.Location_programs lp
		JOIN stu_tracker.Programs p 
		ON lp.program_id = p.id
		WHERE lp.location_id = ANY($1) AND p.organization_id = $2`

	rows, err := s.db.QueryContext(c, query, pq.Array(locId), org_id)
	if err != nil {
		return nil, fmt.Errorf("error querying program: %w", err)
	}
	defer rows.Close()

	var programs []models.ResponseRequestProgramList
	for rows.Next() {
		var program models.ResponseRequestProgramList
		err := rows.Scan(
			&program.ID,
			&program.ProgramName,
			&program.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}
		programs = append(programs, program)
	}

	// Check for any errors encountered during iteration
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}

	return programs, nil
}

func (s *AuthService) GetTutorLocations(c context.Context, tutor_id int64, org_id int64) ([]models.TutorLocationList, error) {
	var query string
	// I can join program id to return programs as well
	query += `
		SELECT 
			ls.name AS location_name,
			ls.id AS id
		FROM stu_tracker.Tutors t
		LEFT JOIN stu_tracker.Locations ls ON t.location_id = ls.id
		WHERE t.id = $1 AND t.location_id IS NOT NULL

		UNION

		SELECT 
			l.name AS location_name,
			tl.location_id AS id
		FROM stu_tracker.Tutor_locations tl
		JOIN stu_tracker.Locations l ON tl.location_id = l.id
		WHERE tl.tutor_id = $1;`

	rows, err := s.db.QueryContext(c, query, tutor_id)
	if err != nil {
		return nil, fmt.Errorf("error querying Tutor_locations: %w", err)
	}
	defer rows.Close()

	var tutors_locations []models.TutorLocationList
	for rows.Next() {
		var tutor models.TutorLocationList
		err := rows.Scan(
			&tutor.LocationName,
			&tutor.ID,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}
		tutors_locations = append(tutors_locations, tutor)
	}

	// Check for any errors encountered during iteration
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}

	return tutors_locations, nil
}

func (s *AuthService) GetTutorSchedules(c context.Context, tutor_id int64) ([]models.RegisterScheduleList, error) {
	var query string
	// I can join program id to return programs as well
	query += `
		SELECT 
		ts.id, ts.tutor_id, p.program_name AS program_name, schedule_type, 
		ts.start_date, ts.end_date, ts.recurring, ts.notes, ts.created_at
		FROM stu_tracker.Tutor_schedules ts
		JOIN stu_tracker.Programs p
		ON p.id = ts.program_id
		WHERE ts.tutor_id = $1;
		`

	rows, err := s.db.QueryContext(c, query, tutor_id)
	if err != nil {
		return nil, fmt.Errorf("error querying Tutor_locations: %w", err)
	}
	defer rows.Close()

	var schedules []models.RegisterScheduleList
	for rows.Next() {
		var schedule models.RegisterScheduleList
		err := rows.Scan(
			&schedule.ID,
			&schedule.TutorID,
			&schedule.ProgramName,
			&schedule.ScheduleType,
			&schedule.StartDate,
			&schedule.EndDate,
			&schedule.Recurring,
			&schedule.Notes,
			&schedule.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}
		schedules = append(schedules, schedule)
	}

	// Check for any errors encountered during iteration
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}

	return schedules, nil
}

func (s *AuthService) GetTutorSessionsAccountability(c context.Context, req models.RequestTutorAccountability) (*models.ResponseTutorAccountability, error) {

	if req.TutorID == nil || req.StartDate.IsZero() || req.EndDate.IsZero() {
		return nil, fmt.Errorf("unable to get sessions missing params")
	}
	var query string
	query += `
		SELECT 
			ss.session_date 
		FROM 
			stu_tracker.Sessions ss
		WHERE 
			ss.tutor_id = $1
		AND 
			DATE(ss.session_date)
		BETWEEN
			$2	
		AND
			$3;`
	rows, err := s.db.QueryContext(c, query, req.TutorID, req.StartDate, req.EndDate)
	if err != nil {
		return nil, fmt.Errorf("error querying Tutor_locations: %w", err)
	}
	defer rows.Close()

	var tutorList []models.TutorAccountability
	for rows.Next() {
		var session models.TutorAccountability
		err := rows.Scan(
			&session.SessionDate,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}
		tutorList = append(tutorList, session)
		// Check for any errors encountered during iteration
		if err = rows.Err(); err != nil {
			return nil, fmt.Errorf("error iterating over rows: %w", err)
		}
	}
	// Find all instances where tutor does not work
	query2 := `
		SELECT start_date, end_date
		FROM stu_tracker.Tutor_schedules
		WHERE schedule_type = 'exclusion'
		AND tutor_id = $1
		AND (
			(start_date <= $3 AND COALESCE(end_date, start_date) >= $2)
		);
	`
	rows, err = s.db.Query(query2, req.TutorID, req.StartDate, req.EndDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []time.Time
	type Schedule struct {
		StartDate time.Time
		EndDate   sql.NullTime
	}
	for rows.Next() {
		var s Schedule
		if err := rows.Scan(&s.StartDate, &s.EndDate); err != nil {
			return nil, err
		}
		rangeStart := s.StartDate
		rangeEnd := s.EndDate.Time
		if !s.EndDate.Valid {
			rangeEnd = s.StartDate
		}

		// Clip the range to the provided window
		if rangeStart.Before(req.StartDate) {
			rangeStart = req.StartDate
		}
		if rangeEnd.After(req.EndDate) {
			rangeEnd = req.EndDate
		}

		// Generate dates in the clipped range
		for d := rangeStart; !d.After(rangeEnd); d = d.AddDate(0, 0, 1) {
			result = append(result, d)
		}
	}
	return &models.ResponseTutorAccountability{
		SessionDates: tutorList,
		NonWorkdays:  result,
	}, nil
}

func (s *AuthService) GetOrganizationPermissions(c context.Context, org_id int64) ([]models.PermissionsList, error) {
	query := `SELECT p.id, p.name, p.description FROM stu_tracker.Permissions p;`
	rows, err := s.db.QueryContext(c, query)
	if err != nil {
		return nil, fmt.Errorf("error querying permissions: %w", err)
	}
	defer rows.Close()
	var permissions []models.PermissionsList
	for rows.Next() {
		var permission models.PermissionsList
		err := rows.Scan(
			&permission.ID,
			&permission.Name,
			&permission.Description,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}
		permissions = append(permissions, permission)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}
	return permissions, nil
}

func (s *AuthService) GetPermissionsById(c context.Context, org_id int64, role string, id int64) ([]models.PermissionsList, error) {
	var query string
	var rows *sql.Rows
	var err error
	switch role {
	case "ROOT":
		query = `
			SELECT p.id, p.name, p.description
			FROM stu_tracker.Permissions p;`
		rows, err = s.db.QueryContext(c, query)
	case "TUTOR":
		query = `
			SELECT tp.permission_id, p.name, p.description
			FROM stu_tracker.Tutor_Permissions tp 
			LEFT JOIN stu_tracker.Permissions p
			ON p.id = tp.permission_id
			WHERE tp.tutor_id = $1;`
		rows, err = s.db.QueryContext(c, query, id)

	case "ADMIN":
		query = `
			SELECT tp.permission_id, p.name, p.description
			FROM stu_tracker.Admin_Permissions tp 
			JOIN stu_tracker.Permissions p
			ON p.id = tp.permission_id
			WHERE tp.admin_id = $1;`
		rows, err = s.db.QueryContext(c, query, id)
	default:
		return nil, fmt.Errorf("invalid role: %s", role)
	}

	if err != nil {
		return nil, fmt.Errorf("error querying permissions: %w", err)
	}
	defer rows.Close()
	var permissions []models.PermissionsList
	for rows.Next() {
		var permission models.PermissionsList
		err := rows.Scan(
			&permission.ID,
			&permission.Name,
			&permission.Description,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}
		permissions = append(permissions, permission)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return permissions, nil
}

func (s *AuthService) GetAnnouncements(c context.Context, req models.AnnouncementRequest) ([]models.AnnouncementsList, error) {
	if req.Role == "ROOT" || req.Role == "ADMIN" {
		announcements, err := s.getAdminAnnouncements(c, int64(req.OrganizationID))
		if err != nil {
			return nil, fmt.Errorf("error on adminAnnouncements %w", err)
		}
		return announcements, nil
	}
	if req.Role == "TUTOR" {
		announcements, err := s.getTutorAnnouncements(c, req.LocationIDs, int64(req.OrganizationID), req.ProgramID)
		if err != nil {
			return nil, err
		}
		return announcements, nil
	}
	return []models.AnnouncementsList{}, nil
}

func (s *AuthService) getAdminAnnouncements(c context.Context, org_id int64) ([]models.AnnouncementsList, error) {
	var query string
	query += `SELECT 
		an.id, an.title, an.body, 
		an.created_at, an.severity, 
		an.program_id, an.admin_id,
		an.staff_id, an.location_id,
		COALESCE(ads.fullname, adr.fullname, 'Unknown') AS staff_name,
		COALESCE(loc.name, 'NA') AS location_name, 
		COALESCE(pr.program_name, 'NA') AS program_name
		FROM stu_tracker.Announcements an 
		LEFT JOIN stu_tracker.Admin_staff ads
		ON ads.id = an.staff_id
		LEFT JOIN stu_tracker.Admin_root adr
		ON adr.id = an.admin_id
		LEFT JOIN stu_tracker.Locations loc
		ON loc.id = an.location_id
		LEFT JOIN stu_tracker.Programs pr
		ON pr.id = an.program_id
		WHERE an.organization_id = $1`
	rows, err := s.db.QueryContext(c, query, org_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var announcements []models.AnnouncementsList
	for rows.Next() {
		var announcement models.AnnouncementsList
		err := rows.Scan(
			&announcement.ID,
			&announcement.Title,
			&announcement.Body,
			&announcement.CreatedAt,
			&announcement.Severity,
			&announcement.ProgramID,
			&announcement.AdminID,
			&announcement.StaffID,
			&announcement.LocationID,
			&announcement.StaffName,
			&announcement.LocationName,
			&announcement.ProgramName,
		)
		if err != nil {
			return nil, err
		}
		announcements = append(announcements, announcement)
	}
	// Check for any errors encountered during iteration
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return announcements, nil

}

func (s *AuthService) getTutorAnnouncements(c context.Context, loc_id []int64, org_id int64, pro_id []int64) ([]models.AnnouncementsList, error) {
	var query string
	query += `
		SELECT an.id, an.title, an.body, 
		an.created_at, an.severity, 
		an.program_id, an.admin_id, 
		an.staff_id, an.location_id, 
		COALESCE(ads.fullname, adr.fullname, 'Unknown') AS staff_name,
		COALESCE(loc.name, 'NA') AS location_name, 
		COALESCE(pr.program_name, 'NA') AS program_name
		FROM stu_tracker.Announcements an 
		LEFT JOIN stu_tracker.Admin_staff ads
		ON ads.id = an.staff_id
		LEFT JOIN stu_tracker.Admin_root adr
		ON adr.id = an.admin_id
		LEFT JOIN stu_tracker.Locations loc
		ON loc.id = an.location_id
		LEFT JOIN stu_tracker.Programs pr
		ON pr.id = an.program_id
		WHERE 
			(an.location_id = ANY($1) OR an.program_id = ANY($2))
		AND an.organization_id = $3
		ORDER BY 
		an.created_at DESC`
	rows, err := s.db.QueryContext(c, query, pq.Array(loc_id), pq.Array(pro_id), org_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var announcements []models.AnnouncementsList
	for rows.Next() {
		var announcement models.AnnouncementsList
		err := rows.Scan(
			&announcement.ID,
			&announcement.Title,
			&announcement.Body,
			&announcement.CreatedAt,
			&announcement.Severity,
			&announcement.ProgramID,
			&announcement.AdminID,
			&announcement.StaffID,
			&announcement.LocationID,
			&announcement.StaffName,
			&announcement.LocationName,
			&announcement.ProgramName,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}
		announcements = append(announcements, announcement)
	}
	// Check for any errors encountered during iteration
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return announcements, nil
}
