package repositories

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"ms-notification/dtos"
	"ms-notification/infrastructures"
)

type configRepository struct {
	supabaseClient *infrastructures.SupabaseClient
}

type ConfigRepository interface {
	GetScheduledTiming() ([]*dtos.SchedulerTiming, error)
	UpdateScheduledTiming(request *dtos.UpdateSchedulerTimingRequest) (*dtos.SchedulerTiming, error)
}

func NewConfigRepository(supabaseClient *infrastructures.SupabaseClient) ConfigRepository {
	return &configRepository{
		supabaseClient: supabaseClient,
	}
}

func (s *configRepository) GetScheduledTiming() ([]*dtos.SchedulerTiming, error) {
	response, err := s.supabaseClient.SelectAll(ScheduledTimingTableName)
	if err != nil {
		logrus.Errorf("Select Scheduler timing error, %v", err)
		return nil, err
	}

	logrus.Infof("Select Scheduler timing success, %s", string(response))

	configs := make([]*dtos.SchedulerTiming, 0)
	if err := json.Unmarshal(response, &configs); err != nil {
		logrus.Errorf("Unmarshal configs error, %v", err)
		return nil, err
	}

	return configs, nil
}

func (s *configRepository) UpdateScheduledTiming(request *dtos.UpdateSchedulerTimingRequest) (*dtos.SchedulerTiming, error) {
	response, err := s.supabaseClient.UpdateScheduleTiming(ScheduledTimingTableName, request)
	if err != nil {
		logrus.Errorf("Update Scheduled Timing config error, %v", err)
		return nil, err
	}

	logrus.Infof("Update Scheduled Timing success, %s", string(response))

	//var configs = &dtos.SchedulerTiming{}
	//if err := json.Unmarshal(response, &configs); err != nil {
	//	logrus.Errorf("Unmarshal configs error, %v", err)
	//	return nil, err
	//}

	return nil, nil
}
