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

type RegisterScheduleV2 struct {
	ID             *int64     `json:"id"`
	JobName        *string    `json:"job_name"`
	JobDescription *string    `json:"jon_description"`
	ProgramID      *int64     `json:"program_id"`
	TutorID        *int64     `json:"tutor_id"`
	SemesterID     *int64     `json:"semester_id"`
	SessionDate    time.Time  `json:"session_date"`
	StartTime      *time.Time `json:"start_time"`
	EndTime        *time.Time `json:"end_time"`
	LocationID     *int64     `json:"location_id"`
	ScheduleType   string     `json:"schedule_type"`
	Enabled        bool       `json:"enabled"`
	UpdatedAt      time.Time  `json:"updated_at"`
	CreatedAt      time.Time  `json:"created_at"`
}

type RegisterGlobalSchedule struct {
	ID                  int64          `json:"id,omitempty"`
	JobName             string         `json:"job_name"`
	JobDescription      *string        `json:"job_description,omitempty"`
	OrganizationID      int            `json:"organization_id"`
	TutorID             *int           `json:"tutor_id,omitempty"`
	LocationID          *int           `json:"location_id,omitempty"`
	GlobalRule          bool           `json:"global_rule"`
	CronJob             *string        `json:"cron_job,omitempty"`
	Metadata            map[string]any `json:"metadata,omitempty"`
	ProviderID          *int           `json:"provider_id,omitempty"`
	ProviderUID         *string        `json:"provider_uid,omitempty"`
	ProviderType        *string        `json:"provider_type,omitempty"`
	ProviderEmployeeID  *int64         `json:"provider_employee_id,omitempty"`
	ProviderEmployeeUID *string        `json:"provider_employee_uid,omitempty"`
	RecurrenceType      string         `json:"recurrence_type"` // weekly | date_range | specific_dates
	StartDate           time.Time      `json:"start_date"`
	EndDate             *time.Time     `json:"end_date,omitempty"`
	SpecificDates       []time.Time    `json:"specific_dates,omitempty"`
	Frequency           []string       `json:"frequency,omitempty"`
	StartTime           time.Time      `json:"start_time"`
	EndTime             time.Time      `json:"end_time"`
	ProgramID           *int           `json:"program_id,omitempty"`
	SemesterID          *int           `json:"semester_id,omitempty"`
	Enabled             bool           `json:"enabled"`
	Archive             bool           `json:"archive"`
	CreatedAt           time.Time      `json:"created_at,omitempty"`
	UpdatedAt           time.Time      `json:"updated_at,omitempty"`
}

type ResponseGlobalSchedule struct {
	ID                  int64          `json:"id,omitempty"`
	JobName             string         `json:"job_name"`
	JobDescription      *string        `json:"job_description,omitempty"`
	OrganizationID      *int64         `json:"organization_id"`
	TutorID             *int           `json:"tutor_id,omitempty"`
	LocationID          *int           `json:"location_id,omitempty"`
	GlobalRule          bool           `json:"global_rule"`
	CronJob             *string        `json:"cron_job,omitempty"`
	Metadata            map[string]any `json:"metadata,omitempty"`
	ProviderID          *int           `json:"provider_id,omitempty"`
	ProviderUID         *string        `json:"provider_uid,omitempty"`
	ProviderType        *string        `json:"provider_type,omitempty"`
	ProviderEmployeeID  *int64         `json:"provider_employee_id,omitempty"`
	ProviderEmployeeUID *string        `json:"provider_employee_uid,omitempty"`
	RecurrenceType      string         `json:"recurrence_type"`
	StartDate           time.Time      `json:"start_date"`
	EndDate             *time.Time     `json:"end_date,omitempty"`
	SpecificDates       []time.Time    `json:"specific_dates,omitempty"`
	Frequency           []string       `json:"frequency,omitempty"`
	StartTime           time.Time      `json:"start_time"`
	EndTime             time.Time      `json:"end_time"`
	ProgramID           *int           `json:"program_id,omitempty"`
	SemesterID          *int           `json:"semester_id,omitempty"`
	Enabled             bool           `json:"enabled"`
	Archive             bool           `json:"archive"`
	CreatedAt           time.Time      `json:"created_at,omitempty"`
	UpdatedAt           time.Time      `json:"updated_at,omitempty"`
}

type RegisterScheduleLink struct {
	TutorID    *int64 `json:"tutor_id"`
	ScheduleID *int64 `json:"schedule_id"`
	LocationID *int64 `json:"location_id"`
}
