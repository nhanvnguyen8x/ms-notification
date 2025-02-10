package dtos

import "time"

type SchedulerTiming struct {
	ID             int64     `json:"id"`
	StartTime      string    `json:"start_time"`
	Timezone       string    `json:"timezone"`
	OrganizationID string    `json:"org_id"`
	Purpose        string    `json:"purpose"`
	CreatedAt      time.Time `json:"created_at"`
}

type UpdateSchedulerTimingRequest struct {
	ID             int64  `json:"id"`
	StartTime      string `json:"start_time"`
	Timezone       string `json:"timezone"`
	OrganizationID string `json:"org_id"`
}
