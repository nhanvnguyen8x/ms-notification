package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"ms-notification/services"
	"net/http"
)

type webhookHandler struct {
	taskService services.TaskService
}

type WebhookHandler interface {
	HandleTaskChange(ctx *gin.Context)
}

func NewWebhookHandler(service services.TaskService) WebhookHandler {
	return &webhookHandler{
		taskService: service,
	}
}

// @BasePath /ms-notification

// HandleTaskChange godoc
// @Summary Handle Task Change
// @Schemes
// @Description do ping
// @Tags example
// @Accept json
// @Produce json
// @Success 200 {string} Helloworld
// @Router /api/webhook [get]
func (h *webhookHandler) HandleTaskChange(ctx *gin.Context) {
	fmt.Println("handle task change")
	tasks, err := h.taskService.GetAllTasks()
	if err != nil {
		logrus.Errorf("get all tasks err: %v", err)
		ctx.JSON(http.StatusInternalServerError, nil)
		return
	}

	ctx.JSON(http.StatusOK, tasks)
}
