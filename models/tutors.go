package models

import "time"

type TutorAccountability struct {
	SessionDate time.Time `json:"session_date"`
}

type ResponseTutorAccountability struct {
	SessionDates []TutorAccountability `json:"session_dates"`
	NonWorkdays  []time.Time           `json:"nonworkdays_dates"`
}

type ResponseTutorSchedule struct {
	Schedule map[string][]TutorSchedule `json:"schedule"`
}

type TutorSchedule struct {
	SessionDate      time.Time  `json:"session_date"`
	ProgramID        *int64     `json:"program_id"`
	LocationID       *int64     `json:"location_id"`
	StartTime        *time.Time `json:"start_time"`
	EndTime          *time.Time `json:"end_time"`
	LocationName     *string    `json:"location_name"`
	ProgramName      *string    `json:"program_name"`
	SessionCount     *int64     `json:"session_count"`
	SessionCompleted bool       `json:"session_completed"`
}

type RequestTutorAccountability struct {
	TutorID        *int64    `json:"tutor_id"`
	StartDate      time.Time `json:"start_time"`
	EndDate        time.Time `json:"end_date"`
	OrganizationID *int64    `json:"organization_id"`
}

type RequestSessionVerify struct {
	SessionDate time.Time `json:"session_date"`
	TutorID     *int64    `json:"tutor_id"`
	LocationID  *int64    `json:"location_id"`
	ProgramID   *int64    `json:"program_id"`
}

type EntitySchedule struct {
	ID              *int64     `json:"id"`
	TutorScheduleID *int64     `json:"tutor_schedule_id"`
	JobName         *string    `json:"job_name"`
	StartTime       time.Time  `json:"start_time"`
	EndTime         *time.Time `json:"end_time"`
	RecurrenceType  *string    `json:"recurrence_type"`
	SpecificDates   []string   `json:"specific_dates"`
	Frequency       []string   `json:"frequency"`
	StartDate       *time.Time `json:"start_date"`
	EndDate         *time.Time `json:"end_date"`
}
