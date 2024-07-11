package controllers

import (
	"encoding/json"
	"time"

	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/worker"
	"github.com/hibiken/asynq"
	"github.com/labstack/echo/v4"
	"github.com/rotisserie/eris"
)

// DispatchTask Dispatch task
// @Tags Admin-Task
// @Summary Dispatch task
// @Description Dispatch task
// @Accept  json
// @Produce  json
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/task/dispatch [post]
func DispatchTask(c echo.Context) error {
	var cc = c.(*models.CustomContext)

	var form models.DispatchTaskForm
	var err = cc.BindAndValidate(&form)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	payload, err := json.Marshal(form.Data)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	taskInfo, err := worker.GetInstance().Client.EnqueueContext(
		c.Request().Context(),
		asynq.NewTask(form.Name, payload),
		worker.QueueLow,
		asynq.MaxRetry(1),
		asynq.Retention(time.Hour*24),
	)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success(taskInfo)
}
