package main

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"ms-notification/configs"
	_ "ms-notification/docs"
	"ms-notification/handlers"
	"ms-notification/infrastructures"
	"ms-notification/middlewares"
	"ms-notification/repositories"
	"ms-notification/schedulers"
	"ms-notification/services"
	"os"
	_ "time/tzdata"
)

func main() {
	// Read App Config
	appConfig := configs.NewApplicationConfig()

	// Init Logger
	logger := logrus.New()
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.InfoLevel)

	app := gin.Default()
	// Apply the recovery middleware
	app.Use(middlewares.Recovery())

	supabaseConfig := appConfig.Supabase
	supabaseClient := infrastructures.NewSupabaseClient(supabaseConfig.ApiUrl, supabaseConfig.ApiKey)

	sendpulseConfig := appConfig.SendPulse
	sendpulseClient := infrastructures.NewSendPulseClient(sendpulseConfig)

	userRepository := repositories.NewUserRepository(supabaseClient)
	taskRepository := repositories.NewTaskRepository(supabaseClient)
	configRepository := repositories.NewConfigRepository(supabaseClient)

	userService := services.NewUserService(userRepository)
	taskService := services.NewTaskService(taskRepository)
	configService := services.NewConfigService(configRepository)

	taskScheduler := schedulers.NewTaskManager(taskService, userService, configService, sendpulseClient)
	schedulerTimings, err := configService.GetScheduledTiming()
	if err != nil {
		logger.Errorf("Get all configs error: %s", err.Error())
		panic(err)
	}

	if len(schedulerTimings) == 0 {
		logger.Errorf("Get all configs empty")
		panic(err)
	}

	schedulerManager := schedulers.NewScheduledManager(taskScheduler)
	if err := schedulerManager.CreateSchedulers(schedulerTimings); err != nil {
		logrus.Errorf("failed to load location %s", err.Error())
		return
	}

	if err := schedulerManager.ScheduleAtTimes(schedulerTimings); err != nil {
		logrus.Errorf("ScheduleAtTimes err: %v", err)
		panic(err)
	} else {
		logrus.Infof("ScheduleAtTimes successfully")
	}

	configHandler := handlers.NewConfigHandler(configService, schedulerManager)
	webhookHandler := handlers.NewWebhookHandler(taskService)
	mainRouter := app.Group("/ms-notification")
	{
		mainRouter.GET("/webhook", webhookHandler.HandleTaskChange)
		mainRouter.POST("/config/scheduler_timing", configHandler.UpdateScheduledTiming)
		mainRouter.POST("/validate", configHandler.UpdateScheduledTiming)
		mainRouter.POST("/validate_task_today", configHandler.ValidateTaskToday)
		mainRouter.POST("/validate_task_due_day", configHandler.ValidateTaskDueDay)
	}

	app.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	if err := app.Run(appConfig.Server.Address); err != nil {
		logger.Errorf("Start HTTP Server error: %s", err.Error())
		panic(err)
	}
}
