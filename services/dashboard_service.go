package services

import (
	"context"
	"database/sql"
	"fmt"
	"time"
	"tracker/app/models"

	"github.com/lib/pq"
)

func (s *AuthService) GetOrganizationById(ctx context.Context, orgid *int64) (*models.ResponseGenerationOrganization, error) {
	var query string
	query += `SELECT id, title, address, zip_code, city, state FROM stu_tracker.Organization WHERE id = $1;`
	var organization models.ResponseGenerationOrganization
	err := s.db.QueryRowContext(ctx, query, *orgid).Scan(&organization.Id, &organization.Title, &organization.Address, &organization.ZipCode, &organization.City, &organization.State)
	if err != nil {
		return nil, fmt.Errorf("unable to query organization %v", err)
	}
	return &organization, nil
}

func (s *AuthService) GetGenerationUsage(ctx context.Context, orgid *int64) ([]models.ResponseGenerationMaterials, []models.ResponseGenerationQuestions, error) {
	var query1 string
	query1 += `SELECT sum(input_tokens),
			sum(output_tokens),
			TO_CHAR(created_at,'YYYY:MM') AS year_month 
			FROM stu_tracker.Generate_materials_task WHERE organization_id = $1 GROUP BY year_month ORDER BY year_month;`
	row1, err := s.db.QueryContext(ctx, query1, *orgid)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to query organization %v", err)
	}
	defer row1.Close()
	var materials []models.ResponseGenerationMaterials
	for row1.Next() {
		var mat models.ResponseGenerationMaterials
		err := row1.Scan(
			&mat.InputSum,
			&mat.OutputSum,
			&mat.YearMonth,
		)
		if err != nil {
			return nil, nil, fmt.Errorf("error scanning row: %w", err)
		}
		// Check for any errors encountered during iteration
		if err = row1.Err(); err != nil {
			return nil, nil, fmt.Errorf("error iterating over rows: %w", err)
		}
		materials = append(materials, mat)
	}
	var query2 string
	query2 += `SELECT sum(input_tokens),
			sum(output_tokens),
			TO_CHAR(created_at,'YYYY:MM') AS year_month 
			FROM stu_tracker.Generate_questions_task WHERE organization_id = $1 GROUP BY year_month ORDER BY year_month;`
	row2, err := s.db.QueryContext(ctx, query2, *orgid)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to query organization %v", err)
	}
	defer row2.Close()
	var questions []models.ResponseGenerationQuestions
	for row2.Next() {
		var que models.ResponseGenerationQuestions
		err := row2.Scan(
			&que.InputSum,
			&que.OutputSum,
			&que.YearMonth,
		)
		if err != nil {
			return nil, nil, fmt.Errorf("error scanning row: %w", err)
		}
		// Check for any errors encountered during iteration
		if err = row2.Err(); err != nil {
			return nil, nil, fmt.Errorf("error iterating over rows: %w", err)
		}
		questions = append(questions, que)
	}

	return materials, questions, nil
}

func (s *AuthService) GetLocationsByID(ctx context.Context, id int64, role string) ([]models.ResponseRequestLocations, error) {
	var query string
	query += `
		SELECT loc.id, loc.name, loc.address, loc.city, loc.state, loc.zip_code, loc.created_at, loc.district_id, archive
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
			&location.Archive,
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

func (s *AuthService) GetSurveysById(ctx context.Context, org_id int64) ([]models.Survey, error) {
	var query string
	query += `
		SELECT id, title, description, is_active, order_by, created_at
		FROM stu_tracker.Surveys WHERE organization_id = $1;`
	rows, err := s.db.QueryContext(ctx, query, org_id)
	if err != nil {
		return nil, fmt.Errorf("error querying locations: %w", err)
	}
	defer rows.Close()

	var surveys []models.Survey
	for rows.Next() {
		var survey models.Survey
		err := rows.Scan(
			&survey.ID,
			&survey.Title,
			&survey.Description,
			&survey.IsActive,
			&survey.OrderBy,
			&survey.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}
		surveys = append(surveys, survey)
	}

	// Check for any errors encountered during iteration
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}

	return surveys, nil
}

func (s *AuthService) GetProgramSurveysById(ctx context.Context, org_id int64, pid int64) ([]models.ResponseQuestionsSurvey, error) {
	var query string
	query += `
		SELECT ss.id AS survey_id, ss.title, ss.description
		FROM stu_tracker.Program_survey ps
		JOIN stu_tracker.Surveys ss ON ps.survey_id = ss.id
		WHERE ps.program_id = $1 AND ps.organization_id = $2
		`
	rows, err := s.db.QueryContext(ctx, query, pid, org_id)
	if err != nil {
		return nil, fmt.Errorf("error querying locations: %w", err)
	}
	defer rows.Close()

	var surveys []models.ResponseQuestionsSurvey
	for rows.Next() {
		var s models.ResponseQuestionsSurvey
		err := rows.Scan(
			&s.SurveyID,
			&s.SurveyTitle,
			&s.SurveyDescription,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}
		surveys = append(surveys, s)
	}

	// Check for any errors encountered during iteration
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}

	return surveys, nil
}

func (s *AuthService) GetSurveyProgramsById(ctx context.Context, org_id int64, pid int64) ([]models.ResponseQuestionsSurvey, error) {
	var query string
	query += `
		SELECT ss.id AS survey_id, ss.title, ss.description, sq.id AS survey_question_id, sq.order_index, sq.question_text, sq.question_type, sc.id AS choice_id, sc.choice_text
		FROM stu_tracker.Program_survey ps
		JOIN stu_tracker.Surveys ss ON ps.survey_id = ss.id
		JOIN stu_tracker.Survey_questions sq ON sq.survey_id = ss.id
		LEFT JOIN stu_tracker.survey_choice sc ON sc.question_survey_id = sq.id
		WHERE ps.program_id = $1 AND ps.organization_id = $2
		`
	rows, err := s.db.QueryContext(ctx, query, pid, org_id)
	if err != nil {
		return nil, fmt.Errorf("error querying locations: %w", err)
	}
	defer rows.Close()

	var surveys []models.ResponseQuestionsSurvey
	for rows.Next() {
		var s models.ResponseQuestionsSurvey
		err := rows.Scan(
			&s.SurveyID,
			&s.SurveyTitle,
			&s.SurveyDescription,
			&s.QuestionID,
			&s.QuestionIndex,
			&s.QuestionText,
			&s.QuestionType,
			&s.QuestionChoiceID,
			&s.QuestionChoiceText,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}
		surveys = append(surveys, s)
	}

	// Check for any errors encountered during iteration
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}

	return surveys, nil
}
func ptrInt64(n sql.NullInt64) *int64 {
	if n.Valid {
		v := n.Int64
		return &v
	}
	return nil
}
func ptrString(s sql.NullString) *string {
	if s.Valid {
		v := s.String
		return &v
	}
	return nil
}
func (s *AuthService) GetSurveyQuestions(ctx context.Context, surveyIDs []int64) ([]models.SurveyQuestions, error) {
	if len(surveyIDs) == 0 {
		return []models.SurveyQuestions{}, nil
	}

	const q = `
		SELECT
		sq.id,
		sq.survey_id,
		sq.order_index,
		sq.question_text,
		sq.question_type,
		sc.id              AS choice_id,
		sc.question_survey_id,
		sc.choice_text
		FROM stu_tracker.survey_questions sq
		LEFT JOIN stu_tracker.survey_choice sc
		ON sc.question_survey_id = sq.id
		WHERE sq.survey_id = ANY($1)
		ORDER BY sq.order_index, sc.id;
`
	rows, err := s.db.QueryContext(ctx, q, pq.Array(surveyIDs))
	if err != nil {
		return nil, fmt.Errorf("query survey questions: %w", err)
	}
	defer rows.Close()

	// Aggregate by question id
	type key = int64
	byQ := make(map[key]*models.SurveyQuestions)
	order := make([]key, 0)

	for rows.Next() {
		var (
			qID, surveyID, orderIdx sql.NullInt64
			qText, qType            sql.NullString
			chID, chQID             sql.NullInt64
			chText                  sql.NullString
		)
		if err := rows.Scan(
			&qID, &surveyID, &orderIdx, &qText, &qType,
			&chID, &chQID, &chText,
		); err != nil {
			return nil, fmt.Errorf("scan row: %w", err)
		}

		id := qID.Int64
		q, ok := byQ[id]
		if !ok {
			q = &models.SurveyQuestions{
				ID:           ptrInt64(qID),
				SurveyID:     ptrInt64(surveyID),
				OrderIndex:   ptrInt64(orderIdx),
				QuestionText: ptrString(qText),
				QuestionType: ptrString(qType),
				Choices:      make([]models.SurveyQuestionChoice, 0, 4), // empty by default
			}
			byQ[id] = q
			order = append(order, id)
		}

		// Append a choice only when the LEFT JOIN row has a real choice
		if chID.Valid {
			q.Choices = append(q.Choices, models.SurveyQuestionChoice{
				ID:               ptrInt64(chID),
				SurveyQuestionId: ptrInt64(chQID),
				ChoiceText:       ptrString(chText),
			})
		}
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows err: %w", err)
	}

	// Preserve question order as returned (by order_index)
	out := make([]models.SurveyQuestions, 0, len(order))
	for _, id := range order {
		out = append(out, *byQ[id])
	}
	return out, nil
}

func (s *AuthService) GetSurveyChoiceById(ctx context.Context, qsid []int64) ([]models.SurveyChoicesList, error) {
	var query string
	query += `
		SELECT id, question_survey_id, choice_text
		FROM stu_tracker.Survey_choice WHERE question_survey_id = ANY($1);`
	rows, err := s.db.QueryContext(ctx, query, pq.Array(qsid))
	if err != nil {
		return nil, fmt.Errorf("error querying locations: %w", err)
	}
	defer rows.Close()

	var surveys []models.SurveyChoicesList
	for rows.Next() {
		var survey models.SurveyChoicesList
		err := rows.Scan(
			&survey.ID,
			&survey.QuestionServeyId,
			&survey.ChoiceText,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}
		surveys = append(surveys, survey)
	}

	// Check for any errors encountered during iteration
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}

	return surveys, nil
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
		stu.created_at, stu.semester_id, stu.direct_partnership,
		COALESCE(stu.created_by, 'NA') AS created_by,
		stu.teacher_id,
		COALESCE(lt.name, '') AS teacher_name,
		stu.timeframe,
		stu.timeframe_start,
		stu.timeframe_end,
		stu.duration_required,
		stu.gender,
		stu.race
		FROM stu_tracker.Students stu
		JOIN stu_tracker.Locations loc
		ON stu.location_id = loc.id
		LEFT JOIN stu_tracker.Locations_teachers lt
		ON lt.id = stu.teacher_id
		WHERE loc.organization_id = $1 AND loc.id = $2;`
		rows, err = s.db.QueryContext(c, query, id, locationId)
	}
	// If tutor return all students and if tutor_id is marked return such students
	if role == "TUTOR" {
		query += `
		SELECT stu.id, stu.first_name, stu.last_name, stu.middle_name, 
		stu.location_id, stu.email, stu.grade_level, stu.active, stu.period, 
		stu.created_at, stu.semester_id, stu.direct_partnership, 
		COALESCE(stu.created_by, 'NA') AS created_by, stu.teacher_id,
		COALESCE(lt.name, '') AS teacher_name,
		stu.timeframe,
		stu.timeframe_start,
		stu.timeframe_end,
		stu.duration_required,
		stu.gender,
		stu.race
		FROM stu_tracker.Students stu
		JOIN stu_tracker.Locations loc
		ON stu.location_id = loc.id
		LEFT JOIN stu_tracker.Locations_teachers lt
		ON lt.id = stu.teacher_id
		WHERE loc.id = $1
		AND ( 
			stu.direct_partnership = FALSE 
			OR (stu.direct_partnership = TRUE AND stu.tutor_id = $2)
		);`
		rows, err = s.db.QueryContext(c, query, locationId, tutorId)
	}
	fmt.Println(query)
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
			&student.CreatedBy,
			&student.TeacherID,
			&student.TeacherName,
			&student.Timeframe,
			&student.TimeframeStart,
			&student.TimeframeEnd,
			&student.DurationRequired,
			&student.Gender,
			&student.Race,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %v", err)
		}

		// Convert nullable fields
		if middleName.Valid {
			student.MiddleName = middleName.String
		}
		if gradeLevel.Valid {
			grade := int64(gradeLevel.Int64)
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
	query += `SELECT id,fullname, email, region, state, district_id, active
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
			&admin.DistrictId,
			&admin.Active,
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
	query += `SELECT id, program_name, timeframe_required, pg.survey_required
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
			&program.TimeFrameRequired,
			&program.SurveyRequired,
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
	query += `SELECT mt.id, mt.title, mt.external_link, mt.description, mt.version,
			 mt.pre, mt.mid, mt.post, mt.visible, mt.created_at, mt.program_id, mt.s3_reference, pg.program_name
			  FROM stu_tracker.materials mt
			  LEFT JOIN stu_tracker.Programs pg
			  ON pg.id = mt.program_id
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
			&material.SReference,
			&material.ProgramName,
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
		query += `SELECT id, first_name, last_name, email, created_at, location_id, active
			  FROM stu_tracker.Tutors tr
			  WHERE tr.organization_id = $1;`
		args = append(args, id)
	} else {
		query += `SELECT tr.id, tr.first_name, tr.last_name, tr.email, tr.created_at, tr.location_id, active
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
			&tutor.Active,
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
	query += `SELECT id, year, title, date_start, date_end, active, archive
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
			&semester.Archive,
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
				ss.date_start, ss.date_end, sl.id
			  FROM 
			  	stu_tracker.Semester_Location sl
			  LEFT JOIN
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
			&semester.ID,
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
			aas.version, aas.pre, aas.mid, aas.post, aas.visible, aas.easy_score, aas.grade_level,
			aas.questionnaire
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
			&assessment.GradeLevel,
			&assessment.Questionnaire,
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
			COALESCE(c.order_number, 0) as choice_order,
			c.is_correct
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
			&r.ChoiceOrderNumber,
			&r.IsCorrect,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		results = append(results, r)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}
	for i := 0; i < len(results); i++ {
		if *results[i].ImageURL != "NA" && *results[i].ImageURL != "" {
			url, err := s.GeneratePresignedUrl(c, "assessment_images/"+*results[i].ImageURL)
			if err != nil {
				return nil, fmt.Errorf("unable to generate presigned url")
			}
			results[i].ImageURL = &url
		}
	}

	return results, nil
}

func (s *AuthService) GetAssessmentSession(ctx context.Context, first_name *string, last_name *string, join_code *string) (*models.ResponseStudentSessionSearch, error) {
	query := `
    SELECT ss.session_token, ss.assessment_id, ss.student_id, ast.title, ast.max_score
    FROM stu_tracker.Assessment_sessions ss
	JOIN stu_tracker.Assessments ast
	ON ast.id = ss.assessment_id
    WHERE ss.join_code = $1 
      AND ss.first_name ILIKE $2 
      AND ss.last_name ILIKE $3
    LIMIT 1;
	`
	var out models.ResponseStudentSessionSearch
	err := s.db.QueryRowContext(
		ctx,
		query,
		join_code,
		"%"+*first_name+"%",
		"%"+*last_name+"%",
	).Scan(&out.SessionToken, &out.AssessmentID, &out.StudentID, &out.Title, &out.MaxScore)

	if err != nil {
		return nil, fmt.Errorf("session not found: %w", err)
	}

	return &out, nil
}

func (s *AuthService) GetAssessmentsQuestionsChoiceExternal(c context.Context, assessment_id int64) ([]models.ResponseAssessmentQuestionsChoiceExternal, error) {
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

	var results []models.ResponseAssessmentQuestionsChoiceExternal

	for rows.Next() {
		var r models.ResponseAssessmentQuestionsChoiceExternal
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

	for i := range results {
		if results[i].ImageURL != nil && *results[i].ImageURL != "NA" && *results[i].ImageURL != "" {
			signedUrl, err := s.GenerateAssessmentsPresignedUrl(c, *results[i].ImageURL)
			fmt.Printf("Signed url %s", signedUrl)
			if err != nil {
				return nil, fmt.Errorf("unable to sign url, possible corrupt image: %v", err)
			}
			results[i].ImageURL = &signedUrl
		}
	}

	return results, nil
}

func (s *AuthService) GetPreAssessment(ctx context.Context, assessmentID int64) (*bool, error) {
	const q = `
        SELECT questionnaire
        FROM stu_tracker.Assessments
        WHERE id = $1
    `
	var questionnaire bool
	if err := s.db.QueryRowContext(ctx, q, assessmentID).Scan(&questionnaire); err != nil {
		if err == sql.ErrNoRows {
			// up to you: return (false, nil) or surface ErrNoRows
			return nil, nil
		}
		return nil, fmt.Errorf("get pre-assessment: %w", err)
	}
	return &questionnaire, nil
}

func (s *AuthService) PreAssessmentCompleted(ctx context.Context, assessmentID, studentID int64, sessionToken *string) (bool, error) {
	// Use EXISTS, and handle NULL session_token with IS NOT DISTINCT FROM
	const q = `
        SELECT EXISTS (
            SELECT 1
            FROM stu_tracker.pre_assessment_questionnaire
            WHERE assessment_id = $1
              AND student_id    = $2
              AND session_token IS NOT DISTINCT FROM $3
        )
    `
	var exists bool
	if err := s.db.QueryRowContext(ctx, q, assessmentID, studentID, sessionToken).Scan(&exists); err != nil {
		return false, fmt.Errorf("check pre-assessment completion: %w", err)
	}
	return exists, nil
}

func (s *AuthService) GetProgramsByLocation(c context.Context, locId int64, org_id int64) ([]models.ResponseRequestProgramList, error) {
	var query string
	query += `
		SELECT p.id, p.program_name, p.created_at, lp.location_id
		FROM stu_tracker.Location_programs lp
		LEFT JOIN stu_tracker.Programs p ON lp.program_id = p.id
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

func (s *AuthService) GetSemesterDates(c context.Context, sid *int64) ([]models.ResponseSemesterDates, error) {
	var query string
	query += `SELECT id, semester_id, schedule_type, start_date, end_date, notes, created_at
	 FROM stu_tracker.Semester_schedule WHERE semester_id = $1;`
	rows, err := s.db.QueryContext(c, query, sid)
	if err != nil {
		return nil, fmt.Errorf("error querying program: %w", err)
	}
	defer rows.Close()
	var sd []models.ResponseSemesterDates
	for rows.Next() {
		var s models.ResponseSemesterDates
		err := rows.Scan(
			&s.ID,
			&s.SemesterID,
			&s.ScheduleType,
			&s.StartDate,
			&s.EndDate,
			&s.Notes,
			&s.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}
		sd = append(sd, s)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}
	return sd, nil
}

func (s *AuthService) GetProgramsByIds(c context.Context, locId []int64, org_id int64) ([]models.ResponseRequestProgramList, error) {
	var query string
	query += `
		SELECT p.id, p.program_name, p.created_at, p.timeframe_required
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
			&program.TimeFrameRequired,
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

func (s *AuthService) GetGroupAttendies(c context.Context, id *int64) ([]models.ResponseRequestStudentList, error) {
	var query string
	query += `
		SELECT s.id, s.first_name, s.last_name, s.middle_name, s.email, s.grade_level, s.active, s.semester_id, s.duration_required, s.timeframe, s.timeframe_start, s.timeframe_end
		FROM stu_tracker.Student_group_attendees sg 
		JOIN stu_tracker.Students s ON sg.student_id = s.id
		WHERE sg.student_group_id = $1;`

	rows, err := s.db.QueryContext(c, query, id)
	if err != nil {
		return nil, fmt.Errorf("error querying program: %w", err)
	}
	defer rows.Close()

	var sl []models.ResponseRequestStudentList
	for rows.Next() {
		var k models.ResponseRequestStudentList
		err := rows.Scan(
			&k.ID,
			&k.FirstName,
			&k.LastName,
			&k.MiddleName,
			&k.Email,
			&k.GradeLevel,
			&k.Active,
			&k.SemesterId,
			&k.DurationRequired,
			&k.Timeframe,
			&k.TimeframeStart,
			&k.TimeframeEnd,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}
		sl = append(sl, k)
	}

	// Check for any errors encountered during iteration
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}

	return sl, nil
}

func (s *AuthService) GetTeachers(c context.Context, location_id *int64) ([]models.ResponseTeachers, error) {
	var query string
	query += `
		SELECT id, name, room, grade_level, location_id, substitute
		FROM stu_tracker.Locations_teachers WHERE location_id = $1;`

	rows, err := s.db.QueryContext(c, query, location_id)
	if err != nil {
		return nil, fmt.Errorf("error querying program: %w", err)
	}
	defer rows.Close()

	var teachers []models.ResponseTeachers
	for rows.Next() {
		var teacher models.ResponseTeachers
		err := rows.Scan(
			&teacher.ID,
			&teacher.Name,
			&teacher.Room,
			&teacher.GradeLevel,
			&teacher.LocationID,
			&teacher.Substitute,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}
		teachers = append(teachers, teacher)
	}

	// Check for any errors encountered during iteration
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}

	return teachers, nil
}

func (s *AuthService) GetAdminLocations(c context.Context, admin_id *int64, orgid *int64) ([]models.ResponseAdminLocations, error) {
	query := `
		SELECT admin_id, location_id FROM stu_tracker.admin_location_access WHERE admin_id = $1 AND organization_id = $2;
	`
	rows, err := s.db.QueryContext(c, query, admin_id, orgid)
	if err != nil {
		return nil, err
	}
	var res []models.ResponseAdminLocations
	for rows.Next() {
		var m models.ResponseAdminLocations
		if err := rows.Scan(
			&m.AdminID,
			&m.LocationID,
		); err != nil {
			return nil, err
		}
		res = append(res, m)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}

	return res, nil
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
		ts.start_date, ts.end_date, ts.recurring, ts.notes, ts.created_at, ts.workweek, ts.start_time, ts.end_time, l.name
		FROM 
			stu_tracker.Tutor_schedules ts
		LEFT JOIN 
			stu_tracker.Programs p
		ON p.id = ts.program_id
		LEFT JOIN 
			stu_tracker.Locations l
		ON l.id = ts.location_id
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
			pq.Array(&schedule.WorkWeek),
			&schedule.StartTime,
			&schedule.EndTime,
			&schedule.LocationName,
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

func (s *AuthService) GetTutorSchedule(ctx context.Context, uid *int64, pid []int64) (*models.ResponseTutorSchedule, error) {
	if uid == nil || pid == nil {
		return nil, fmt.Errorf("params are null")
	}
	query := `
		SELECT DISTINCT
			gen_date::date AS date_value,
			ts.program_id,
			ts.location_id,
			ts.start_time, 
			ts.end_time,
			pg.program_name,
			ls.name
		FROM
			stu_tracker.Tutor_schedules ts
		JOIN LATERAL 
			generate_series(ts.start_date::date, COALESCE(ts.end_date, ts.start_date)::date, '1 day'::interval) AS t(gen_date)
		ON true
		LEFT JOIN stu_tracker.Locations ls
		ON ls.id = ts.location_id
		LEFT JOIN stu_tracker.Programs pg
		ON pg.id = ts.program_id
		WHERE ts.tutor_id = $1 AND ts.program_id = ANY($2) AND ts.schedule_type = 'inclusion' 
		AND (ts.workweek IS NULL OR ts.workweek = '{}'::text[] OR to_char(gen_date, 'DY') = ANY(ts.workweek));
	`
	rows, err := s.db.QueryContext(ctx, query, uid, pq.Array(pid))
	if err != nil {
		return nil, fmt.Errorf("error querying permissions: %w", err)
	}
	defer rows.Close()
	result := make(map[string][]models.TutorSchedule)
	for rows.Next() {
		var w models.TutorSchedule
		err := rows.Scan(
			&w.SessionDate,
			&w.ProgramID,
			&w.LocationID,
			&w.StartTime,
			&w.EndTime,
			&w.ProgramName,
			&w.LocationName,
		)
		key := w.SessionDate.Format("2006-01-02")
		result[key] = append(result[key], w)
		if err != nil {
			return nil, err
		}
	}

	r, err := s.SessionCompletedCheck(ctx, result, uid)
	if err != nil {
		return nil, err
	}

	return &models.ResponseTutorSchedule{
		Schedule: r,
	}, nil
}

func (s *AuthService) SessionCompletedCheck(c context.Context, result map[string][]models.TutorSchedule, uid *int64) (map[string][]models.TutorSchedule, error) {
	keys := make([]string, 0, len(result))
	for k := range result {
		keys = append(keys, k)
	}
	query := `
		SELECT
		ss.program_id,
		ss.location_id,
		ss.session_date
		FROM stu_tracker.Sessions ss
		WHERE ss.tutor_id = $1 AND ss.session_date::date = ANY($2::date[]);
	`
	rows, err := s.db.QueryContext(c, query, uid, pq.Array(keys))
	if err != nil {
		return nil, fmt.Errorf("unable to query for sessions given tutor_id and session dates: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var w models.RequestSessionVerify
		err := rows.Scan(
			&w.ProgramID,
			&w.LocationID,
			&w.SessionDate,
		)
		if err != nil {
			return nil, fmt.Errorf("unable to scan row %s", err)
		}
		key := w.SessionDate.Format("2006-01-02")
		values, ok := result[key]
		if !ok {
			continue
		}
		for k := range values {
			if *values[k].LocationID == *w.LocationID && *values[k].ProgramID == *w.ProgramID {
				values[k].SessionCompleted = true
			}
		}
		result[key] = values

	}
	return result, nil
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

func (s *AuthService) GetAnnouncementsAck(c context.Context, orgid *int64, aid *int64) ([]models.AnnouncementsListAck, error) {
	query := `SELECT ua.id, ua.tutor_id, ua.acknowledged, ua.acknowledged_at, t.first_name, t.last_name
	FROM stu_tracker.User_Acknowledgments ua
	LEFT JOIN stu_tracker.Tutors t
	ON t.id = ua.tutor_id
	WHERE ua.organization_id = $1 AND ua.announcement_id = $2;`
	rows, err := s.db.QueryContext(c, query, orgid, aid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var a []models.AnnouncementsListAck
	for rows.Next() {
		var m models.AnnouncementsListAck
		err := rows.Scan(
			&m.ID,
			&m.TutotID,
			&m.Acknowledged,
			&m.AcknowledgedAt,
			&m.FirstName,
			&m.LastName,
		)
		if err != nil {
			return nil, err
		}
		a = append(a, m)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}
	return a, nil
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

func (s *AuthService) GetLocationContactByID(c context.Context, orgid *int64, locid *int64) ([]models.ResponseLocationContact, error) {
	query := `SELECT id, first_name, last_name, email, phone, notes, description, location_id FROM stu_tracker.Location_contacts WHERE organization_id = $1 AND location_id = $2;`
	rows, err := s.db.QueryContext(c, query, orgid, locid)
	if err != nil {
		return nil, fmt.Errorf("error querying permissions: %w", err)
	}
	defer rows.Close()
	var loc []models.ResponseLocationContact
	for rows.Next() {
		var l models.ResponseLocationContact
		err := rows.Scan(
			&l.ID,
			&l.FirstName,
			&l.LastName,
			&l.Email,
			&l.Phone,
			&l.Notes,
			&l.Description,
			&l.LocationID,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}
		loc = append(loc, l)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}
	return loc, nil
}
