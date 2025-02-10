package dtos

import "time"

type Task struct {
	TaskID            string    `json:"task_id"`
	TaskName          string    `json:"task_name"`
	TaskCode          string    `json:"task_code"`
	Remarks           string    `json:"remarks"`
	IsComplete        bool      `json:"is_complete"`
	CompletedDateTime time.Time `json:"completed_dt"`
	Action            string    `json:"action"`
	DueDate           time.Time `json:"due_date"`
	AssignedTo        []string  `json:"assigned_to"`
	CreatedAt         time.Time `json:"created_at"`
	CreatedBy         string    `json:"created_by"`
	UpdatedAt         time.Time `json:"updated_at"`
	DealID            string    `json:"deal_id"`
	OrganizationID    string    `json:"org_id"`
	Outcome           string    `json:"outcome"`
	DepTaskID         string    `json:"dep_task_id"`
}
