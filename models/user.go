package models

type User struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token        string                   `json:"token"`
	User         User                     `json:"user"`
	RefreshToken string                   `json:"refresh_token"`
	Permissions  LoginResponsePermissions `json:"permissions"`
}

type LoginResponsePermissions struct {
	DisableUpdate bool `json:"disable_update"`
	DisableCreate bool `json:"disable_create"`
	DisableDelete bool `json:"disable_delete"`
}

type RegisterResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

type RegisterRequestAdminRoot struct {
	Password     string `json:"password"`
	Email        string `json:"email"`
	Organization string `json:"organization"`
}

type RegisterResponseAdminRoot struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
}

type RegisterRequestAdminStaff struct {
}

type RegisterRequestDistrict struct {
}

type RegisterRequestProgramRoot struct {
}

type RegisterRequestTutors struct {
}

type RegisterRequestAssessments struct {
}

type RegisterRequestAssessmentsLog struct {
}

type RegisterRequestStudents struct {
}
type RegisterRequestSessions struct {
}
