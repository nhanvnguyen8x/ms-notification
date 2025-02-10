package schedulers

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"ms-notification/dtos"
	"ms-notification/infrastructures"
	"ms-notification/services"
	"strings"
	"time"
)

type TaskManager struct {
	taskService     services.TaskService
	userService     services.UserService
	configService   services.ConfigService
	sendpulseClient *infrastructures.SendPulseClient
}

func NewTaskManager(taskService services.TaskService,
	userService services.UserService,
	configService services.ConfigService,
	sendpulseClient *infrastructures.SendPulseClient) *TaskManager {
	return &TaskManager{
		taskService:     taskService,
		userService:     userService,
		configService:   configService,
		sendpulseClient: sendpulseClient,
	}
}

func (tm *TaskManager) runDaily(schedulerTiming *dtos.SchedulerTiming) error {
	var userTaskMap = make(map[string]*dtos.UserTask)

	users, err := tm.userService.GetAllUsers()
	if err != nil {
		logrus.Error("Get all users err: ", err)
		return err
	}

	userMap := make(map[string]*dtos.User)
	for _, user := range users {
		userMap[user.ID] = user
	}

	tasks, err := tm.taskService.GetAllTasks()
	if err != nil {
		logrus.Error("Get all tasks err: ", err)
		return err
	}

	assignedTodayTasks := tm.FilterAssignedTodayTask(tasks, schedulerTiming)
	for _, task := range assignedTodayTasks {
		logrus.Infof("Assign task to day: %s", task.TaskCode)
		var assignees = task.AssignedTo
		for _, assignedTo := range assignees {
			user, ok := userMap[assignedTo]
			if !ok {
				logrus.Error("Assign task to not exist user: ", assignedTo)
				continue
			}

			userTask, ok := userTaskMap[assignedTo]
			if !ok {
				userTaskMap[assignedTo] = &dtos.UserTask{
					User:               user,
					IncompleteTasks:    make([]*dtos.Task, 0),
					AssignedTodayTasks: make([]*dtos.Task, 0),
					DueDayTasks:        make([]*dtos.Task, 0),
				}
				userTaskMap[assignedTo].AssignedTodayTasks = append(userTaskMap[assignedTo].AssignedTodayTasks, task)
			} else {
				userTask.AssignedTodayTasks = append(userTask.AssignedTodayTasks, task)
			}
		}
	}

	var dueDayTasks = tm.FilterTaskWithDueDays(tasks, schedulerTiming)
	for _, task := range dueDayTasks {
		logrus.Infof("Due task to day: %s", task.TaskCode)
		var assignees = task.AssignedTo
		for _, assignedTo := range assignees {
			user, ok := userMap[assignedTo]
			if !ok {
				logrus.Error("Assign task to not exist user: ", assignedTo)
				continue
			}

			userTask, ok := userTaskMap[assignedTo]
			if !ok {
				userTaskMap[assignedTo] = &dtos.UserTask{
					User:               user,
					IncompleteTasks:    make([]*dtos.Task, 0),
					AssignedTodayTasks: make([]*dtos.Task, 0),
					DueDayTasks:        make([]*dtos.Task, 0),
					OutcomeTasks:       make([]*dtos.Task, 0),
				}
				userTaskMap[assignedTo].DueDayTasks = append(userTaskMap[assignedTo].DueDayTasks, task)
			} else {
				userTask.DueDayTasks = append(userTask.DueDayTasks, task)
			}
		}
	}

	// Handle Incomplete task
	var incompleteTasks = tm.FilterIncompleteTask(tasks)
	for _, task := range incompleteTasks {
		logrus.Infof("Incomplete task: %s", task.TaskCode)
		var assignees = task.AssignedTo
		for _, assignedTo := range assignees {
			user, ok := userMap[assignedTo]
			if !ok {
				logrus.Error("Assign task to not exist user: ", assignedTo)
				continue
			}

			userTask, ok := userTaskMap[assignedTo]
			if !ok {
				userTaskMap[assignedTo] = &dtos.UserTask{
					User:               user,
					IncompleteTasks:    make([]*dtos.Task, 0),
					AssignedTodayTasks: make([]*dtos.Task, 0),
					DueDayTasks:        make([]*dtos.Task, 0),
					OutcomeTasks:       make([]*dtos.Task, 0),
				}
				userTaskMap[assignedTo].IncompleteTasks = append(userTaskMap[assignedTo].IncompleteTasks, task)
			} else {
				userTask.IncompleteTasks = append(userTask.IncompleteTasks, task)
			}
		}
	}

	var outcomeTasks = tm.FilterDecisionOutcomeTask(tasks)
	for _, task := range outcomeTasks {
		logrus.Infof("Outcome task: %s", task.Outcome)
		var assignees = task.AssignedTo
		for _, assignedTo := range assignees {
			user, ok := userMap[assignedTo]
			if !ok {
				logrus.Error("Assign task to not exist user: ", assignedTo)
				continue
			}

			userTask, ok := userTaskMap[assignedTo]
			if !ok {
				userTaskMap[assignedTo] = &dtos.UserTask{
					User:               user,
					IncompleteTasks:    make([]*dtos.Task, 0),
					AssignedTodayTasks: make([]*dtos.Task, 0),
					DueDayTasks:        make([]*dtos.Task, 0),
					OutcomeTasks:       make([]*dtos.Task, 0),
				}
				userTaskMap[assignedTo].OutcomeTasks = append(userTaskMap[assignedTo].OutcomeTasks, task)
			} else {
				userTask.OutcomeTasks = append(userTask.OutcomeTasks, task)
			}
		}
	}
	for userID, userTask := range userTaskMap {
		logrus.Infof("Start Handle Email Notification for userID: %s - email: %s", userID, userTask.User.Email)
		var targetEmail = userTask.User.Email
		var targetName = userTask.User.FirstName
		var toAccount = &dtos.SendPulseEmailAccount{
			Name:  targetName,
			Email: targetEmail,
		}

		var incompleteTaskContent = tm.buildIncompleteTaskContent(userTask.IncompleteTasks)
		var todayTaskContent = tm.buildTodayTaskContent(userTask.AssignedTodayTasks)
		var dueDayTaskContent = tm.buildDueDayTaskContent(userTask.DueDayTasks)
		var outcomeTaskContent = tm.buildOutcomeTaskContent(userTask.OutcomeTasks)

		logrus.Infof("UserID: %s", userID)
		logrus.Infof("Target Email: %s", targetEmail)
		logrus.Infof("Incomplete Task Content: %s", incompleteTaskContent)
		logrus.Infof("To Day Task Content: %s", todayTaskContent)
		logrus.Infof("Due Day Task Content: %s", dueDayTaskContent)

		var content = fmt.Sprintf("%s\n%s\n%s\n%s", incompleteTaskContent, todayTaskContent, dueDayTaskContent, outcomeTaskContent)
		logrus.Infof("Content is: %s", content)

		if err := tm.HandleEmailNotification(schedulerTiming, content, toAccount); err != nil {
			logrus.Errorf("handle email notification err: %v", err)
		} else {
			logrus.Infof("Handle email notification successfully: %s", targetEmail)
		}
	}

	return nil
}

func (tm *TaskManager) buildEmailContent(userTaskMap map[string]*dtos.UserTask) string {
	for userID, userTask := range userTaskMap {
		var targetEmail = userTask.User.Email
		var incompleteTaskContent = tm.buildIncompleteTaskContent(userTask.IncompleteTasks)
		var todayTaskContent = tm.buildTodayTaskContent(userTask.AssignedTodayTasks)
		var dueDayTaskContent = tm.buildTodayTaskContent(userTask.DueDayTasks)
		logrus.Infof("UserID: %s", userID)
		logrus.Infof("Target Email: %s", targetEmail)
		logrus.Infof("Incomplete Task Content: %s", incompleteTaskContent)
		logrus.Infof("To Day Task Content: %s", todayTaskContent)
		logrus.Infof("Due Day Task Content: %s", dueDayTaskContent)
	}

	return ""
}

func (tm *TaskManager) GetTasksWithinTodayContent(schedulerTiming *dtos.SchedulerTiming) (string, error) {
	// Handle new tasks assigned today
	location, err := time.LoadLocation(schedulerTiming.Timezone)
	if err != nil {
		logrus.Errorf("failed to load location %s", err.Error())
		return "", err
	}

	var today = time.Now().In(location)
	year, month, day := today.Date()
	var startTime = fmt.Sprintf("%d-%02d-%02d", year, month, day)

	var tomorrow = today.AddDate(0, 0, 1)
	year, month, day = tomorrow.Date()
	var endTime = fmt.Sprintf("%d-%02d-%02d", year, month, day)
	tasksWithinToday, err := tm.taskService.GetTasksWithinTheDay(startTime, endTime)
	if err != nil {
		logrus.Errorf("get tasks within today err: %v", err)
		return "", err
	}

	var tasksWithinTodayContent = tm.buildTodayTaskContent(tasksWithinToday)
	return tasksWithinTodayContent, nil
}

func (tm *TaskManager) FilterTaskWithDueDays(allTasks []*dtos.Task, schedulerTiming *dtos.SchedulerTiming) []*dtos.Task {
	var dueDayTasks = make([]*dtos.Task, 0)
	location, err := time.LoadLocation(schedulerTiming.Timezone)
	if err != nil {
		logrus.Errorf("failed to load location %s", err.Error())
		return nil
	}

	var today = time.Now().In(location)
	for _, task := range allTasks {
		var taskDueDay = task.DueDate
		var threeDaysLater = today.Add(3 * 24 * time.Hour)
		if threeDaysLater.After(taskDueDay) {
			dueDayTasks = append(dueDayTasks, task)
		}
	}

	return dueDayTasks
}

func (tm *TaskManager) FilterIncompleteTask(allTasks []*dtos.Task) []*dtos.Task {
	var incompleteTasks = make([]*dtos.Task, 0)
	for _, task := range allTasks {
		if task.IsComplete == false {
			incompleteTasks = append(incompleteTasks, task)
		}
	}

	return incompleteTasks
}

func (tm *TaskManager) FilterDecisionOutcomeTask(allTasks []*dtos.Task) []*dtos.Task {
	return allTasks
}

func (tm *TaskManager) FilterAssignedTodayTask(allTasks []*dtos.Task, schedulerTiming *dtos.SchedulerTiming) []*dtos.Task {
	var assignedTodayTasks = make([]*dtos.Task, 0)
	location, err := time.LoadLocation(schedulerTiming.Timezone)
	if err != nil {
		logrus.Errorf("failed to load location %s", err.Error())
		return nil
	}

	var today = time.Now().In(location)
	for _, task := range allTasks {
		var createdAt = task.CreatedAt
		if createdAt.Year() == today.Year() && createdAt.Month() == today.Month() && createdAt.Day() == today.Day() {
			assignedTodayTasks = append(assignedTodayTasks, task)
		}
	}

	return assignedTodayTasks
}

func (tm *TaskManager) GetTasksWithDueDayContent(schedulerTiming *dtos.SchedulerTiming) (string, error) {
	// Handle Tasks with 3 day to due day
	var dueDayTasks = make([]*dtos.Task, 0)

	location, err := time.LoadLocation(schedulerTiming.Timezone)
	if err != nil {
		logrus.Errorf("failed to load location %s", err.Error())
		return "", err
	}

	var today = time.Now().In(location)
	taskWithDueDays, err := tm.taskService.GetAllTasks()
	if err != nil {
		logrus.Errorf("get GetAllTasks err: %v", err)
		return "", err
	}

	for _, task := range taskWithDueDays {
		var taskDueDay = task.DueDate
		var threeDaysLater = today.Add(3 * 24 * time.Hour)
		if threeDaysLater.After(taskDueDay) {
			dueDayTasks = append(dueDayTasks, task)
		}
	}

	var tasksWithinTodayContent = tm.buildDueDayTaskContent(dueDayTasks)
	return tasksWithinTodayContent, nil
}

func (tm *TaskManager) buildIncompleteTaskContent(incompleteTasks []*dtos.Task) string {
	if len(incompleteTasks) == 0 {
		return "Incomplete Tasks: NONE"
	}

	var taskNames []string
	for _, task := range incompleteTasks {
		taskNames = append(taskNames, fmt.Sprintf("Task Code: %s", task.TaskCode))
	}
	var taskName = strings.Join(taskNames, "\n\t")
	var incompleteTasksContent = fmt.Sprintf("Incomplete Tasks:\n %s", taskName)
	return incompleteTasksContent
}

func (tm *TaskManager) buildTodayTaskContent(todayTasks []*dtos.Task) string {
	if len(todayTasks) == 0 {
		return "Tasks assigned today: NONE"
	}

	var taskNames []string
	for _, task := range todayTasks {
		taskNames = append(taskNames, fmt.Sprintf("Task Code: %s - Created At: %s", task.TaskCode, task.CreatedAt.Format("2006-01-02 15:04:05")))
	}

	var taskName = strings.Join(taskNames, "\n\t")
	var incompleteTasksContent = fmt.Sprintf("Tasks assigned today:\n %s", taskName)
	return incompleteTasksContent
}

func (tm *TaskManager) buildDueDayTaskContent(dueDateTasks []*dtos.Task) string {
	if len(dueDateTasks) == 0 {
		return "Tasks with 3 to Due day:: NONE"
	}

	var taskNames []string
	for _, task := range dueDateTasks {
		taskNames = append(taskNames, fmt.Sprintf("Task Code: %s - Due day: %s", task.TaskCode, task.DueDate.Format("2006-01-02 15:04:05")))
	}

	var taskName = strings.Join(taskNames, "\n\t")
	var incompleteTasksContent = fmt.Sprintf("Tasks with 3 days to Due day:\n %s", taskName)
	return incompleteTasksContent
}

func (tm *TaskManager) buildOutcomeTaskContent(dueDateTasks []*dtos.Task) string {
	if len(dueDateTasks) == 0 {
		return "Decision outcome: NONE"
	}

	var taskNames []string
	for _, task := range dueDateTasks {
		taskNames = append(taskNames, fmt.Sprintf("Task Code: %s - Decision Outcome: %s", task.TaskCode, task.Outcome))
	}

	var taskName = strings.Join(taskNames, "\n\t")
	var incompleteTasksContent = fmt.Sprintf("Decision Outcomes:\n %s", taskName)
	return incompleteTasksContent
}

func (tm *TaskManager) HandleEmailNotification(schedulerTiming *dtos.SchedulerTiming, content string, toAccount *dtos.SendPulseEmailAccount) error {
	accessToken, err := tm.sendpulseClient.GetAccessToken()
	if err != nil {
		logrus.Errorf("get access token err: %v", err)
		return err
	}

	logrus.Infof("Get access token successfully")
	response, err := tm.sendpulseClient.SendTransactionalEmail(accessToken, tm.sendpulseClient.PrepareEmailHtml(schedulerTiming, content, toAccount))
	if err != nil {
		logrus.Errorf("send transactional email err: %v", err)
		return err
	}

	logrus.Infof("SendTransactionalEmail successfully: %s", string(response))
	return nil
}

func (tm *TaskManager) runWeekly() error {
	fmt.Println("Run Weekly task")

	return nil
}
