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
			ss.subject_id, 
			ss.notes, 
			ss.edited_at, 
			ss.created_at,
			pg.program_name AS program_name,
			sb.title AS subject_name
		FROM 
			stu_tracker.Sessions ss
		JOIN 
    		stu_tracker.Tutors t ON ss.tutor_id = t.id 
		JOIN 
			stu_tracker.Locations ll ON ll.id = ss.location_id
		JOIN 
			stu_tracker.Programs pg ON pg.id = ss.program_id 
		JOIN 
    		stu_tracker.Subjects sb ON ss.subject_id = sb.id `

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
	if ss.ProgramId != nil {
		conditions = append(conditions, fmt.Sprintf("ss.program_id = $%d", argIndex))
		args = append(args, ss.ProgramId)
		argIndex++

	}
	if ss.DateStart != "" {
		conditions = append(conditions, fmt.Sprintf("ss.start_time >= $%d", argIndex))
		args = append(args, ss.DateStart)
		argIndex++
	}
	if ss.DateEnd != "" {
		conditions = append(conditions, fmt.Sprintf("ss.start_time <= $%d", argIndex))
		args = append(args, ss.DateEnd)
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
