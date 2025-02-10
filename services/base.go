package services

import "encoding/json"

var (
	ScheduledTimingTableName = "scheduler_timings"
	UserTableName            = "users"
	TaskTableName            = "deal_tasks"
)

func toString(data interface{}) string {
	byteData, err := json.Marshal(data)
	if err != nil {
		return ""
	}
	return string(byteData)
}
