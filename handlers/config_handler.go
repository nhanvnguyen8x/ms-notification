package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"ms-notification/dtos"
	"ms-notification/schedulers"
	"ms-notification/services"
	"net/http"
)

type configHandler struct {
	scheduledManager schedulers.ScheduledManager
	taskService      services.TaskService
	configService    services.ConfigService
}

type ConfigHandler interface {
	UpdateScheduledTiming(ctx *gin.Context)
	ValidateTaskToday(ctx *gin.Context)
	ValidateTaskDueDay(ctx *gin.Context)
	ValidateNotification(ctx *gin.Context)
}

func NewConfigHandler(configService services.ConfigService, scheduledManager schedulers.ScheduledManager) ConfigHandler {
	return &configHandler{
		configService:    configService,
		scheduledManager: scheduledManager,
	}
}

// @BasePath /ms-notification

// UpdateScheduledTiming godoc
// @Summary Update Scheduled Timing config
// @Schemes
// @Description UpdateScheduledTiming
// @Tags UpdateScheduledTiming
// @Accept json
// @Produce json
// @Param   schedulerTiming body dtos.UpdateSchedulerTimingRequest true "Request body"
// @Success 200 {object} dtos.SchedulerTiming
// @Router /ms-notification/config/scheduler_timing [post]
func (ch *configHandler) UpdateScheduledTiming(ctx *gin.Context) {
	var request = &dtos.UpdateSchedulerTimingRequest{}
	if err := ctx.BindJSON(request); err != nil {
		logrus.Errorf("Failed to parse request body: %v", err)
		ctx.JSON(http.StatusBadRequest, "Failed to parse request body")
		return
	}

	var updatedSchedulerTiming = &dtos.SchedulerTiming{
		StartTime:      request.StartTime,
		Timezone:       request.Timezone,
		OrganizationID: request.OrganizationID,
	}

	_, err := ch.configService.UpdateScheduledTiming(request)
	if err != nil {
		logrus.Errorf("Failed to update timing:%v", err)
		ctx.JSON(http.StatusInternalServerError, "Failed to update timing")
		return
	}

	logrus.Errorf("Update scheduler timing successfully: %v", updatedSchedulerTiming)

	// Stop previous
	logrus.Infof("Staring - Stop previous scheduler")
	if err := ch.scheduledManager.StopOrganizationScheduler(updatedSchedulerTiming.OrganizationID); err != nil {
		logrus.Errorf("Failed to stop scheduler: %v", err)
		ctx.JSON(http.StatusInternalServerError, "Failed to stop scheduler")
		return
	}
	logrus.Infof("Ending - Stop previous scheduler successfully")

	// Create new scheduler
	logrus.Infof("Staring - Create a new scheduler")
	if err := ch.scheduledManager.CreateScheduler(updatedSchedulerTiming); err != nil {
		logrus.Errorf("Failed to create new scheduler: %s", err.Error())
		ctx.JSON(http.StatusInternalServerError, "Failed to create new scheduler")
		return
	}
	logrus.Infof("Ending - Create a new scheduler successfully")

	// Start new scheduler
	logrus.Infof("Staring - Start new scheduler with new timing: Startime: %s", request.StartTime)
	if err := ch.scheduledManager.ScheduleAtTime(updatedSchedulerTiming); err != nil {
		logrus.Errorf("Failed to schedule timing: %v", err)
		ctx.JSON(http.StatusInternalServerError, "Failed to stop scheduler")
		return
	}
	logrus.Infof("Ending - Start new scheduler with new timing successfully")

	ctx.JSON(http.StatusOK, updatedSchedulerTiming)
}

// @BasePath /ms-notification

// ValidateNotification godoc
// @Summary ValidateNotification config
// @Schemes
// @Description ValidateNotification
// @Tags ValidateNotification
// @Accept json
// @Produce json
// @Router /ms-notification/validate [post]
func (ch *configHandler) ValidateNotification(ctx *gin.Context) {
	if err := ch.scheduledManager.Validate(); err != nil {
		logrus.Errorf("Failed to validate scheduler: %v", err)
		ctx.JSON(http.StatusInternalServerError, "Failed to validate scheduler")
		return
	}

	ctx.JSON(http.StatusOK, "OK")
}

// ValidateTaskToday godoc
// @Summary ValidateTaskToday config
// @Schemes
// @Description ValidateTaskToday
// @Tags ValidateTaskToday
// @Accept json
// @Produce json
// @Router /ms-notification/validate_task_today [post]
func (ch *configHandler) ValidateTaskToday(ctx *gin.Context) {
	tasksWithinTodayContent, err := ch.scheduledManager.ValidateTaskToday()
	if err != nil {
		logrus.Errorf("Failed to validate scheduler: %v", err)
		ctx.JSON(http.StatusInternalServerError, "Failed to validate scheduler")
	}

	ctx.JSON(http.StatusOK, tasksWithinTodayContent)
}

// ValidateTaskDueDay godoc
// @Summary ValidateTaskDueDay config
// @Schemes
// @Description ValidateTaskDueDay
// @Tags ValidateTaskDueDay
// @Accept json
// @Produce json
// @Router /ms-notification/validate_task_due_day [post]
func (ch *configHandler) ValidateTaskDueDay(ctx *gin.Context) {
	tasksDueDayContent, err := ch.scheduledManager.ValidateTaskDueDay()
	if err != nil {
		logrus.Errorf("Failed to validate scheduler: %v", err)
		ctx.JSON(http.StatusInternalServerError, "Failed to validate scheduler")
	}

	ctx.JSON(http.StatusOK, tasksDueDayContent)
}
