package repositories

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"ms-notification/dtos"
	"ms-notification/infrastructures"
)

type taskRepository struct {
	supabaseClient *infrastructures.SupabaseClient
}

type TaskRepository interface {
	GetAllTasks() ([]*dtos.Task, error)
	GetTasksWithinTheDay(startDate, endDate string) ([]*dtos.Task, error)
}

func NewTaskRepository(supabaseClient *infrastructures.SupabaseClient) TaskRepository {
	return &taskRepository{
		supabaseClient: supabaseClient,
	}
}

func (s *taskRepository) GetAllTasks() ([]*dtos.Task, error) {
	response, err := s.supabaseClient.SelectAll(TaskTableName)
	if err != nil {
		logrus.Errorf("Select all tasks error, %v", err)
		return nil, err
	}

	tasks := make([]*dtos.Task, 0)
	if err := json.Unmarshal(response, &tasks); err != nil {
		logrus.Errorf("Unmarshal tasks error, %v", err)
		return nil, err
	}

	return tasks, nil
}

func (s *taskRepository) GetTasksWithinTheDay(startDate, endDate string) ([]*dtos.Task, error) {
	response, err := s.supabaseClient.SelectTaskToday(TaskTableName, startDate, endDate)
	if err != nil {
		logrus.Errorf("Select all tasks error, %v", err)
		return nil, err
	}

	tasks := make([]*dtos.Task, 0)
	if err := json.Unmarshal(response, &tasks); err != nil {
		logrus.Errorf("Unmarshal tasks error, %v", err)
		return nil, err
	}

	return tasks, nil
}
