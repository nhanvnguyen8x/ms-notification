package schedulers

import (
	"errors"
	"github.com/go-co-op/gocron/v2"
	"github.com/sirupsen/logrus"
	"ms-notification/dtos"
	"strconv"
	"strings"
	"time"
)

type ScheduledManager struct {
	taskScheduler *TaskManager
	schedulerMap  map[string]gocron.Scheduler
}

func NewScheduledManager(taskScheduler *TaskManager) ScheduledManager {
	return ScheduledManager{
		taskScheduler: taskScheduler,
		schedulerMap:  make(map[string]gocron.Scheduler),
	}
}

func (sm *ScheduledManager) CreateSchedulers(schedulerTimings []*dtos.SchedulerTiming) error {
	for _, schedulerTiming := range schedulerTimings {
		var organizationID = schedulerTiming.OrganizationID
		location, err := time.LoadLocation(schedulerTiming.Timezone)
		if err != nil {
			logrus.Errorf("failed to load location %s", err.Error())
			return err
		}

		var scheduleLocation = gocron.WithLocation(location)
		scheduler, err := gocron.NewScheduler(scheduleLocation)
		if err != nil {
			logrus.Errorf("new scheduler err: %v", err)
			return err
		}

		sm.schedulerMap[organizationID] = scheduler
		logrus.Infof("CreateSchedulers for Organization: %s successfully", organizationID)
	}

	return nil
}

func (sm *ScheduledManager) ScheduleAtTimes(schedulerTimings []*dtos.SchedulerTiming) error {
	for _, schedulerTiming := range schedulerTimings {
		organizationID := schedulerTiming.OrganizationID
		var organizationScheduler = sm.schedulerMap[organizationID]

		timing, err := GetStartTime(schedulerTiming.StartTime)
		if err != nil {
			logrus.Errorf("Get start time %s has err: %v", schedulerTiming.StartTime, err)
			return err
		}

		schedulerJob, err := organizationScheduler.NewJob(
			gocron.DailyJob(1, gocron.NewAtTimes(gocron.NewAtTime(timing[0], timing[1], timing[2]))),
			gocron.NewTask(func() {
				if err := sm.taskScheduler.runDaily(schedulerTiming); err != nil {
					logrus.Errorf("sm.taskScheduler.runDaily() daily err: %s, organizationID: %s", err.Error(), organizationID)
					return
				}
			}),
		)

		if err != nil {
			logrus.Errorf("OrganizationScheduler %s new job err: %v", organizationID, err)
			return err
		}

		logrus.Infof("Scheduler daily job created for Organization: %s, JobID: %s", organizationID, schedulerJob.ID().String())
		organizationScheduler.Start()
	}

	return nil
}

func (sm *ScheduledManager) CreateScheduler(schedulerTiming *dtos.SchedulerTiming) error {
	location, err := time.LoadLocation(schedulerTiming.Timezone)
	if err != nil {
		logrus.Errorf("failed to load location %s", err.Error())
		return err
	}

	var organizationID = schedulerTiming.OrganizationID
	var scheduleLocation = gocron.WithLocation(location)
	scheduler, err := gocron.NewScheduler(scheduleLocation)
	if err != nil {
		logrus.Errorf("CreateScheduler for organization: %s err: %v", organizationID, err)
		return err
	}

	sm.schedulerMap[organizationID] = scheduler
	return nil
}

func (sm *ScheduledManager) ScheduleAtTime(schedulerTiming *dtos.SchedulerTiming) error {
	timing, err := GetStartTime(schedulerTiming.StartTime)
	if err != nil {
		logrus.Errorf("get start time err: %v", err)
		return err
	}

	var organizationID = schedulerTiming.OrganizationID
	schedulerJob, err := sm.schedulerMap[organizationID].NewJob(
		gocron.DailyJob(1, gocron.NewAtTimes(gocron.NewAtTime(timing[0], timing[1], timing[2]))),
		gocron.NewTask(func() {
			if err := sm.taskScheduler.runDaily(schedulerTiming); err != nil {
				logrus.Errorf("sm.taskScheduler.runDaily() daily err: %s", err.Error())
				return
			}
		}),
	)

	if err != nil {
		logrus.Errorf("scheduler new job err: %v", err)
		return err
	}

	logrus.Infof("scheduler daily job created: %s, %s", schedulerJob.ID(), schedulerJob.Name())
	sm.schedulerMap[organizationID].Start()
	return nil
}

func (sm *ScheduledManager) Validate() error {
	return nil
}

func (sm *ScheduledManager) ValidateTaskToday() (string, error) {
	return sm.taskScheduler.GetTasksWithinTodayContent(&dtos.SchedulerTiming{
		Timezone: "Asia/Ho_Chi_Minh",
	})
}

func (sm *ScheduledManager) ValidateTaskDueDay() (string, error) {
	return sm.taskScheduler.GetTasksWithDueDayContent(&dtos.SchedulerTiming{
		Timezone: "Asia/Ho_Chi_Minh",
	})
}

func (sm *ScheduledManager) StopOrganizationScheduler(organizationID string) error {
	if err := sm.schedulerMap[organizationID].StopJobs(); err != nil {
		logrus.Errorf("stop jobs err: %v", err)
		return err
	}

	return nil
}

func GetStartTime(startTime string) ([]uint, error) {
	var hours []uint
	parts := strings.Split(startTime, ":")
	if len(parts) != 3 {
		return nil, errors.New("invalid start time format")
	}

	for _, str := range parts {
		num, err := strconv.ParseUint(str, 10, 64)
		if err != nil {
			logrus.Errorf("Error converting string to int64: %v", err)
			return nil, err
		}
		hours = append(hours, uint(num))
	}

	return hours, nil
}
