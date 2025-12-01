package models

import "time"

type RegisterSchedule struct {
	TutorID      *int64     `json:"tutor_id"`
	LocationID   *int64     `json:"location_id"`
	ProgramID    *int64     `json:"program_id"`
	ScheduleType string     `json:"schedule_type"`
	StartDate    time.Time  `json:"start_date"`
	StartTime    *time.Time `json:"start_time"`
	EndTime      *time.Time `json:"end_time"`
	WorkWeek     *[]string  `json:"work_week"`
	EndDate      *time.Time `json:"end_date"`
	Notes        string     `json:"notes"`
}

type RegisterScheduleResponse struct {
	ID     int64  `json:"id"`
	Status string `json:"status"`
}

type RegisterScheduleList struct {
	ID           *int64     `json:"id"`
	TutorID      *int64     `json:"tutor_id"`
	ProgramID    *int64     `json:"program_id"`
	ProgramName  string     `json:"program_name"`
	LocationName *string    `json:"location_name"`
	ScheduleType string     `json:"schedule_type"`
	Recurring    bool       `json:"recurring"`
	StartDate    time.Time  `json:"start_date"`
	EndDate      *time.Time `json:"end_date"`
	StartTime    *time.Time `json:"start_time"`
	EndTime      *time.Time `json:"end_time"`
	WorkWeek     []string   `json:"work_week"`
	Notes        string     `json:"notes"`
	CreatedAt    time.Time  `json:"created_at"`
}

type RemoveSchedule struct {
	ID *int64 `json:"id"`
}
