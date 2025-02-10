package dtos

type User struct {
	ID                 string `json:"id"`
	Email              string `json:"email"`
	FirstName          string `json:"first_name"`
	LastName           string `json:"last_name"`
	ApplicantName      string `json:"applicant_name"`
	ApplicantHandPhone string `json:"applicant_hand_phone"`
}

type UserTask struct {
	User               *User
	IncompleteTasks    []*Task
	AssignedTodayTasks []*Task
	DueDayTasks        []*Task
	OutcomeTasks       []*Task
}
