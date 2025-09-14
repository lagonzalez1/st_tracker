package models

import "time"

type RegisterSchedule struct {
	TutorID      *int64     `json:"tutor_id"`
	ProgramID    *int64     `json:"program_id"`
	ScheduleType string     `json:"schedule_type"`
	StartDate    time.Time  `json:"start_date"`
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
	ScheduleType string     `json:"schedule_type"`
	Recurring    bool       `json:"recurring"`
	StartDate    time.Time  `json:"start_date"`
	EndDate      *time.Time `json:"end_date"`
	WorkWeek     []string   `json:"work_week"`
	Notes        string     `json:"notes"`
	CreatedAt    time.Time  `json:"created_at"`
}

type RemoveSchedule struct {
	ID *int64 `json:"id"`
}
