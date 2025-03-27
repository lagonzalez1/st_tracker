package services

import (
	"fmt"
	"strconv"
	"tracker/app/models"

	"github.com/xuri/excelize/v2"
	"golang.org/x/crypto/bcrypt"
)

func response_tutor_file(responseList []*models.ResponseMultipleRegisterUser) (*excelize.File, error) {
	f := excelize.NewFile()
	sheetName := "Multi_add_tutors"
	f.SetSheetName("Sheet1", sheetName)

	// Set headers
	headers := []string{"First Name", "Last Name", "Email", "Password"}
	for i, header := range headers {
		cell := fmt.Sprintf("%c1", 'A'+i)
		f.SetCellValue(sheetName, cell, header)
	}
	// Insert data
	for i, tutor := range responseList {
		rowNum := i + 2
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", rowNum), tutor.FirstName)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", rowNum), tutor.LastName)
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", rowNum), tutor.Email)
		f.SetCellValue(sheetName, fmt.Sprintf("D%d", rowNum), tutor.Password)
	}
	return f, nil
}

func (s *AuthService) RegisterMultipleTutors(rows [][]string, organizationID *int64) (*excelize.File, *models.ResponseMultipleRegister, error) {
	// Input validation
	if organizationID == nil {
		return nil, nil, fmt.Errorf("organization_id is null")
	}
	var responseList []*models.ResponseMultipleRegisterUser
	// Prepare SQL statement for inserting tutors
	stmt, err := s.db.Prepare(`
		INSERT INTO stu_tracker.Tutors 
		(first_name, last_name, email, password_hash, organization_id, location_id, active) 
		VALUES ($1, $2, $3, $4, $5, NULL, $6) 
		RETURNING id;
	`)
	if err != nil {
		return nil, nil, err
	}
	defer stmt.Close()
	successCount := 0
	// Loop through rows (skip header)
	for i, row := range rows {
		if i == 0 {
			continue
		}
		if len(row) != 4 {
			fmt.Printf("Skipping row %d: not enough columns\n", i+1)
			continue
		}
		firstName := row[0]
		lastName := row[1]
		email := row[2]
		boolStr := row[3]
		if email == "" || firstName == "" || lastName == "" {
			fmt.Printf("Skipping row %d: not enough columns\n", i+1)
			continue
		}
		// Convert boolean value
		isActive, err := strconv.ParseBool(boolStr)
		if err != nil {
			fmt.Printf("Row %d: Invalid boolean value: %s\n", i+1, boolStr)
			continue
		}
		var user models.ResponseMultipleRegisterUser
		// Generate hashed password
		rawPassword := fmt.Sprintf("SSTAutopass%s%d!", firstName, i)
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(rawPassword), bcrypt.DefaultCost)
		if err != nil {
			fmt.Printf("Row %d: Failed to hash password: %v\n", i+1, err)
			continue
		}
		// Insert tutor and retrieve tutor ID
		var tutorID int64
		err = stmt.QueryRow(firstName, lastName, email, hashedPassword, organizationID, isActive).Scan(&tutorID)
		if err != nil {
			fmt.Printf("Row %d: Failed to insert tutor: %v\n", i+1, err)
			continue
		}
		// Insert tutor permissions
		permissionQuery := `
			INSERT INTO stu_tracker.Tutor_Permissions (tutor_id, permission_id)
			VALUES ($1, 25), ($1, 24), ($1, 6), ($1, 39);
		`
		_, err = s.db.Exec(permissionQuery, tutorID)
		if err != nil {
			fmt.Printf("Row %d: Failed to assign permissions: %v\n", i+1, err)
			continue
		}
		user.Password = rawPassword
		user.FirstName = firstName
		user.LastName = lastName
		user.Email = email
		responseList = append(responseList, &user)
		successCount++
	}

	file, err := response_tutor_file(responseList)
	if err != nil {
		return nil, nil, err
	}

	return file, &models.ResponseMultipleRegister{
		Status: "Upload success, affected count: ",
		Count:  successCount,
	}, nil
}
