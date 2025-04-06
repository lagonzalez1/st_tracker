package models

type PermissionsModel struct {
	ID          *int64  `json:"id"`
	Name        *string `json:"name"`
	Description *string `json:"description"`
	Role        *string `json:"string"`
}
