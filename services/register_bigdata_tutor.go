package services

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"tracker/app/models"

	"github.com/xuri/excelize/v2"
	"golang.org/x/crypto/bcrypt"
)

func response_tutor_file(responseList []*models.ResponseMultipleRegisterUser) (*excelize.File, error) {
	f := excelize.NewFile()
	sheetName := "Multi_add_tutors"
	f.SetSheetName("Sheet1", sheetName)
	// Set headers
	headers := []string{"First Name", "Last Name", "Email", "Password", "Location"}
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
		f.SetCellValue(sheetName, fmt.Sprintf("E%d", rowNum), tutor.Location)
	}
	return f, nil
}

func (s *AuthService) GetPermissionsByRole(primaryRole string) ([]int, error) {
	if primaryRole == "" {
		return nil, fmt.Errorf("no role provided")
	}
	query := `SELECT id FROM stu_tracker.Permissions WHERE primary_role = $1;`
	rows, err := s.db.Query(query, primaryRole)
	if err != nil {
		return nil, fmt.Errorf("unable to query from Permissions: %w", err)
	}
	defer rows.Close()
	var ids []int
	for rows.Next() {
		var id int
		err := rows.Scan(&id)
		if err != nil {
			return nil, fmt.Errorf("unable to scan from rows: %v", err)
		}
		ids = append(ids, id)
	}
	return ids, nil
}

func buildPermissionQueryTutors(permissionIDs []int, tutorID int) (string, []interface{}) {
	if len(permissionIDs) == 0 {
		return "", nil // or consider returning an error
	}
	query := `INSERT INTO stu_tracker.Tutor_Permissions(tutor_id, permission_id) VALUES `
	var args []interface{}
	args = append(args, tutorID)

	var placeholders []string
	for i, permID := range permissionIDs {
		placeholders = append(placeholders, fmt.Sprintf("($1, $%d)", i+2))
		args = append(args, permID)
	}

	query += strings.Join(placeholders, ", ")
	return query, args
}

func (s *AuthService) RegisterMultipleTutors(c context.Context, rows [][]string, organizationID *int64) (*excelize.File, *models.ResponseMultipleRegister, error) {
	// Input validation
	if organizationID == nil {
		return nil, nil, fmt.Errorf("organization_id is null")
	}
	var responseList []*models.ResponseMultipleRegisterUser
	// Prepare SQL statement for inserting tutors
	stmt, err := s.db.PrepareContext(c, `
		INSERT INTO stu_tracker.Tutors 
		(first_name, last_name, email, password_hash, location_id, organization_id, active) 
		VALUES ($1, $2, $3, $4, $5, $6, $7) 
		RETURNING id;
	`)
	if err != nil {
		return nil, nil, err
	}
	defer stmt.Close()
	successCount := 0
	permissionIDs, err := s.GetPermissionsByRole("tutor")
	if err != nil {
		return nil, nil, fmt.Errorf("unable to get permissions by role")
	}
	fmt.Println("GET permissions by role: ", permissionIDs)
	// Loop through rows (skip header)
	for i, row := range rows {
		if i == 0 {
			continue
		}
		if len(row) > 5 {
			fmt.Printf("Skipping row %d: not enough columns\n", i+1)
			continue
		}
		firstName := row[0]
		lastName := row[1]
		email := row[2]
		boolStr := row[3]
		var locationID sql.NullInt64
		if len(row) == 5 {
			location := row[4]
			if strings.TrimSpace(location) != "" {
				parsedID, err := strconv.ParseInt(location, 10, 64)
				if err != nil {
					fmt.Printf("Row %d: invalid location id: %v\n", i+1, err)
					continue
				}
				locationID.Int64 = parsedID
				locationID.Valid = true
			} else {
				locationID.Valid = false
			}
		}

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
		rawPassword := fmt.Sprintf("!SSTAutopass%s%d!", firstName, i)
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(rawPassword), bcrypt.DefaultCost)
		if err != nil {
			fmt.Printf("Row %d: Failed to hash password: %v\n", i+1, err)
			continue
		}
		// Insert tutor and retrieve tutor ID
		var tutorID int64
		err = stmt.QueryRow(firstName, lastName, email, hashedPassword, locationID, organizationID, isActive).Scan(&tutorID)
		if err != nil {
			fmt.Printf("Row %d: Failed to insert tutor: %v\n", i+1, err)
			continue
		}
		permissionQuery, args := buildPermissionQueryTutors(permissionIDs, int(tutorID))
		_, err = s.db.QueryContext(c, permissionQuery, args...)
		if err != nil {
			fmt.Printf("Row %d: Failed to assign permissions: %v\n", i+1, err)
			continue
		}
		user.Password = rawPassword
		user.FirstName = firstName
		user.LastName = lastName
		user.Email = email
		user.Location = locationID.Int64
		responseList = append(responseList, &user)
		successCount++
	}

	file, err := response_tutor_file(responseList)
	if err != nil {
		fmt.Printf("unable to create reponse file: %v", err)
		return nil, nil, err
	}

	return file, &models.ResponseMultipleRegister{
		Status: "Upload success, affected count: ",
		Count:  successCount,
	}, nil
}
