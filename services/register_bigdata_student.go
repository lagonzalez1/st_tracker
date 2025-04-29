package services

import (
	"context"
	"fmt"
	"strconv"
	"tracker/app/models"
)

func (s *AuthService) RegisterMultipleStudents(c context.Context, rows [][]string, organizationID *int64, semesterID int64, locationID int64) (*models.ResponseMultipleRegisterStudents, error) {
	// Input validation
	if organizationID == nil {
		return nil, fmt.Errorf("organization_id is null")
	}
	var responseList []*models.UploadStudentRegister
	// Prepare SQL statement for inserting tutors
	stmt, err := s.db.PrepareContext(c, `
		INSERT INTO stu_tracker.Students 
		(first_name, middle_name , last_name, grade_level, email, active, location_id, semester_id) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id;`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	successCount := 0
	// Loop through rows (skip header)
	for i, row := range rows {
		if i == 0 {
			continue
		}
		if len(row) != 6 {
			fmt.Printf("Skipping row %d: not enough columns\n", i+1)
			continue
		}
		firstName, middleName, lastName, grade, email, active := row[0], row[1], row[2], row[3], row[4], row[5]
		if firstName == "" || lastName == "" {
			fmt.Printf("Skipping row %d: not enough columns\n", i+1)
			continue
		}
		// Convert boolean value
		isActive, err := strconv.ParseBool(active)
		if err != nil {
			fmt.Printf("Row %d: Invalid boolean value: %s\n", i+1, active)
			continue
		}
		var user models.UploadStudentRegister
		err = stmt.QueryRow(firstName, middleName, lastName, grade, email, isActive, locationID, semesterID).Scan(&user.ID)
		if err != nil {
			fmt.Printf("Row %d: Failed to insert tutor: %v\n", i+1, err)
			continue
		}
		responseList = append(responseList, &user)
		successCount++
	}
	return &models.ResponseMultipleRegisterStudents{
		Status: "OK",
		Count:  successCount,
		List:   responseList,
	}, nil
}
