package repositories

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"ms-notification/dtos"
	"ms-notification/infrastructures"
)

type userRepository struct {
	supabaseClient *infrastructures.SupabaseClient
}

type UserRepository interface {
	GetAllUsers() ([]*dtos.User, error)
}

func NewUserRepository(supabaseClient *infrastructures.SupabaseClient) UserRepository {
	return &userRepository{
		supabaseClient: supabaseClient,
	}
}

func (s *userRepository) GetAllUsers() ([]*dtos.User, error) {
	response, err := s.supabaseClient.SelectAll(UserTableName)
	if err != nil {
		logrus.Errorf("Select All User Error: %v", err)
		return nil, err
	}

	users := make([]*dtos.User, 0)
	if err := json.Unmarshal(response, &users); err != nil {
		logrus.Errorf("Unmarshal User Error: %v", err)
		return nil, err
	}

	return users, nil
}
