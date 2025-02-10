package services

import (
	"ms-notification/dtos"
	"ms-notification/repositories"
)

type configService struct {
	configRepository repositories.ConfigRepository
}

type ConfigService interface {
	GetScheduledTiming() ([]*dtos.SchedulerTiming, error)
	UpdateScheduledTiming(request *dtos.UpdateSchedulerTimingRequest) (*dtos.SchedulerTiming, error)
}

func NewConfigService(configRepository repositories.ConfigRepository) ConfigService {
	return &configService{
		configRepository: configRepository,
	}
}

func (s *configService) GetScheduledTiming() ([]*dtos.SchedulerTiming, error) {
	return s.configRepository.GetScheduledTiming()
}

func (s *configService) UpdateScheduledTiming(request *dtos.UpdateSchedulerTimingRequest) (*dtos.SchedulerTiming, error) {
	return s.configRepository.UpdateScheduledTiming(request)
}
