package models

import "time"

type TutorAccountability struct {
	SessionDate time.Time `json:"session_date"`
}

type ResponseTutorAccountability struct {
	SessionDates []TutorAccountability `json:"session_dates"`
	NonWorkdays  []time.Time           `json:"nonworkdays_dates"`
}

type RequestTutorAccountability struct {
	TutorID        *int64    `json:"tutor_id"`
	StartDate      time.Time `json:"start_time"`
	EndDate        time.Time `json:"end_date"`
	OrganizationID *int64    `json:"organization_id"`
}
