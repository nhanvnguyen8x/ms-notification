package services

import (
	"ms-notification/dtos"
	"ms-notification/repositories"
)

type userService struct {
	userRepository repositories.UserRepository
}

type UserService interface {
	GetAllUsers() ([]*dtos.User, error)
}

func NewUserService(userRepository repositories.UserRepository) UserService {
	return &userService{
		userRepository: userRepository,
	}
}

func (s *userService) GetAllUsers() ([]*dtos.User, error) {
	return s.userRepository.GetAllUsers()
}
