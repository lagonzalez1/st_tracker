package services

import (
	"database/sql"
	"fmt"
	"tracker/app/models"
)

func (s *AuthService) GetLocationsByID(id int64, role string) ([]models.ResponseRequestLocations, error) {
	var query string
	query += `
		SELECT loc.id, loc.name, loc.address, loc.city, loc.state, loc.zip_code, loc.created_at, loc.district_id
		FROM stu_tracker.Locations loc where loc.organization_id = $1;`
	rows, err := s.db.Query(query, id)
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
		locations = append(locations, location)
	}

	// Check for any errors encountered during iteration
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}

	return locations, nil
}

func (s *AuthService) GetSubjectById(id int64, role string) ([]models.SubjectList, error) {
	var query string
	query += `
		SELECT sb.id, sb.title, sb.description, sb.organization_id, sb.created_at
		FROM stu_tracker.Subjects sb WHERE sb.organization_id = $1;`
	rows, err := s.db.Query(query, id)
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

func (s *AuthService) GetSubjectByLocation(id int64, loc_id int64, role string) ([]models.SubjectList, error) {
	var query string
	query += `
		SELECT sb.id, sb.title, sb.description, sb.organization_id, sb.created_at
		FROM stu_tracker.Subjects sb where sb.organization_id = $1 AND sb.location_id = $2;`
	rows, err := s.db.Query(query, id)
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

func (s *AuthService) GetStudentsByID(id int64, role string, locationId int64) ([]models.ResponseRequestStudentList, error) {
	var query string
	query += `
			SELECT stu.id, stu.first_name, stu.last_name, stu.middle_name, 
			stu.location_id, stu.email, stu.grade_level, stu.active, stu.period, 
			stu.created_at, stu.semester_id
			FROM stu_tracker.Students stu
			INNER JOIN stu_tracker.Locations loc
			ON stu.location_id = loc.id
			WHERE loc.organization_id = $1 AND loc.id = $2;`

	rows, err := s.db.Query(query, id, locationId)
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

func (s *AuthService) GetAdminStaffById(id int64, role string) ([]models.ResponseRequestAdminList, error) {
	var query string
	query += `SELECT id,fullname, email, region, state 
			  FROM stu_tracker.Admin_staff
			  WHERE organization_id = $1;`

	rows, err := s.db.Query(query, id)
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

func (s *AuthService) GetDistrictsById(id int64, role string) ([]models.ResponseRequestDistrictList, error) {
	var query string
	query += `SELECT id, name, city ,state, region, created_at
			  FROM stu_tracker.district ds WHERE ds.organization_id = $1;`
	rows, err := s.db.Query(query, id)
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

func (s *AuthService) GetProgramsId(id int64, role string) ([]models.ResponseRequestProgramList, error) {
	var query string
	query += `SELECT id, program_name
			  FROM stu_tracker.programs pg
			  WHERE pg.organization_id = $1;`
	rows, err := s.db.Query(query, id)
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

func (s *AuthService) GetMaterialsById(id int64, role string) ([]models.ResponseRequestMaterialsList, error) {
	var query string
	query += `SELECT id, title, external_link, description, version, pre, mid, post, visable, created_at
			  FROM stu_tracker.materials mt
			  WHERE mt.organization_id = $1 AND mt.visable = FALSE;`
	rows, err := s.db.Query(query, id)
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
			&material.Visable,
			&material.CreatedAt,
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

func (s *AuthService) GetTutorsById(id int64, role string, locid int64) ([]models.ResponseRequestTutorsList, error) {
	var query string
	query += `SELECT id, first_name, last_name, email, created_at, location_id
			  FROM stu_tracker.Tutors tr
			  WHERE tr.organization_id = $1 AND tr.location_id = $2;`
	rows, err := s.db.Query(query, id, locid)
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

func (s *AuthService) GetSemestersById(id int64, role string) ([]models.ResponseRequestSemesterList, error) {
	var query string
	query += `SELECT id, year, title
			  FROM stu_tracker.Semester sm
			  WHERE sm.organization_id = $1;`
	rows, err := s.db.Query(query, id)
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

func (s *AuthService) GetAssessmentsById(id int64, role string) ([]models.ResponseAssessmentList, error) {
	var query string
	query += `SELECT id, title, description, letter, cycle, alpha_identifier, external_link, max_score, subject, material_id, created_at
			  FROM stu_tracker.Assessments aas WHERE aas.organization_id = $1;`
	rows, err := s.db.Query(query, id)
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
			&assessment.Subject,
			&assessment.MaterialID,
			&assessment.CreatedAt,
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

func (s *AuthService) GetProgramsByLocation(id int64, locId int64, role string) ([]models.ResponseRequestProgramList, error) {
	var query string
	query += `
		SELECT p.id, p.program_name, p.created_at
		FROM stu_tracker.Location_programs lp
		JOIN stu_tracker.Programs p ON lp.program_id = p.id
		WHERE lp.location_id = $1`

	rows, err := s.db.Query(query, locId)
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
