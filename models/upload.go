package models

type ResponseMultipleRegister struct {
	Status string `json:"status"`
	Count  int    `json:"count"`
}

type UploadRegister struct {
	OrganizationID *int64 `json:"organization_id"`
	ID             *int64 `json:"id"`
	Email          string `json:"email"`
	Role           string `json:"role"`
}

type ResponseMultipleRegisterUser struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Location  int64  `json:"location"`
}

type UploadRegisterStudents struct {
	ID             *int64 `json:"id"`
	Email          string `json:"email"`
	Role           string `json:"role"`
	OrganizationID *int64 `json:"organization_id"`
	SemesterID     *int64 `json:"semester_id"`
	LocationID     *int64 `json:"location_id"`
}

type UploadStudentRegister struct {
	ID         *int64 `json:"id"`
	FirstName  string `json:"first_name"`
	MiddleName string `json:"middle_name"`
	LastName   string `json:"last_name"`
	Email      string `json:"email"`
	Grade      *int64 `json:"grade"`
	Active     bool   `json:"active"`
}

type ResponseMultipleRegisterStudents struct {
	Status string                   `json:"status"`
	Count  int                      `json:"count"`
	List   []*UploadStudentRegister `json:"list"`
}
