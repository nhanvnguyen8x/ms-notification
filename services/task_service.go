package services

import (
	"ms-notification/dtos"
	"ms-notification/repositories"
)

type taskService struct {
	taskRepository repositories.TaskRepository
}

type TaskService interface {
	GetAllTasks() ([]*dtos.Task, error)
	GetTasksWithinTheDay(startDate, endDate string) ([]*dtos.Task, error)
}

func NewTaskService(taskRepository repositories.TaskRepository) TaskService {
	return &taskService{
		taskRepository: taskRepository,
	}
}

func (s *taskService) GetAllTasks() ([]*dtos.Task, error) {
	return s.taskRepository.GetAllTasks()
}

func (s *taskService) GetTasksWithinTheDay(startDate, endDate string) ([]*dtos.Task, error) {
	return s.taskRepository.GetTasksWithinTheDay(startDate, endDate)
}
