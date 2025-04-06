package permissions

import (
	"database/sql"
	"fmt"
	"tracker/app/models"
)

const (
	Root  = "root"
	Admin = "admin"
	Tutor = "tutor"
)

//** Using redis here would be ideal **/

// ** On each new addition with permissions needed to be set **//
// ** Fetch all and sort **//
// ** Default permissions per role **//

func LoadPermissionsInit(db *sql.DB, role *string) ([]int, error) {
	var i []int
	if role == nil {
		return nil, fmt.Errorf("no role provided")
	}
	query := `SELECT id, name, description, role FROM stu_tracker.Permissions;`
	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error querying Session_students: %w", err)
	}
	defer rows.Close()
	var permissions []models.PermissionsModel
	for rows.Next() {
		var permission models.PermissionsModel
		err := rows.Scan(
			&permission.ID,
			&permission.Name,
			&permission.Description,
			&permission.Role,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}
		permissions = append(permissions, permission)
	}
	for i := 0; i < len(permissions); i++ {
		if *role == Admin {

		}
	}
	return i, nil
}
