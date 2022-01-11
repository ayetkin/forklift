package controller

import (
	"context"
	"errors"
	"fmt"
	"forklift/internal/api/model"
	"forklift/internal/domain"
	"forklift/internal/domain/entity"
	"forklift/internal/domain/event"
	"forklift/pkg/enums"
	"forklift/pkg/rabbitmq"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

// GetMigrationTasks godoc
// @Produce  json
// @Success 200 {string} string
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Tags Migration
// @Router /api/migration/list [get]
func GetMigrationTasks(e *echo.Echo, taskRepository domain.MigrationTaskRepository) {
	e.GET("/api/migration/list", func(c echo.Context) error {
		tasks, err := taskRepository.GetMigrations()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, tasks)
	})
}

// StartMigrationTask godoc
// @Accept  json
// @Produce  json
// @Param request body entity.MigrationRequest true "request"
// @Success 200 {object} object
// @Failure 400 {string} string
// @Failure 401 {string} string
// @Failure 500 {string} string
// @Tags Migration
// @Router /api/migration/start [post]
func StartMigrationTask(e *echo.Echo, taskRepository domain.MigrationTaskRepository, messageBus rabbitmq.Client) {
	e.POST("/api/migration/start", func(c echo.Context) error {

		var (
			migrationTask *entity.MigrationTask
			request       = new(entity.MigrationRequest)
			ctx           = context.Background()
			err           error
		)

		if err = c.Bind(request); err != nil {
			log.Error(ctx, "request deserialize error.", err)
			return c.JSON(http.StatusBadRequest, fmt.Sprintf("request deserialize error: %s", err))
		}

		migrationTask = entity.NewMigrationTask(*request)

		if err = taskRepository.UpsertMigrationTask(migrationTask); err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}

		if err = publishEventMessage(ctx, messageBus, migrationTask.MessageId, migrationTask.Stage); err != nil {
			log.Error(ctx, "Event message publish error", err)
			return c.JSON(http.StatusBadRequest, err.Error())
		}

		return c.JSON(http.StatusOK, "ok")
	})
}

// RetryMigrationTask godoc
// @Accept  json
// @Produce  json
// @Param request body model.RetryMigrationRequest true "request"
// @Success 200 {object} object
// @Failure 400 {string} string
// @Failure 401 {string} string
// @Failure 500 {string} string
// @Tags Migration
// @Router /api/migration/retry [post]
func RetryMigrationTask(e *echo.Echo, migrationTaskRepository domain.MigrationTaskRepository, messageBus rabbitmq.Client) {
	e.POST("/api/migration/retry", func(c echo.Context) error {

		var (
			request             = new(model.RetryMigrationRequest)
			ctx                 = context.Background()
			migrationTaskFromDB *entity.MigrationTask
			err                 error
		)

		if err = c.Bind(request); err != nil {
			log.Errorf("Request deserialize error. %v", err)
			return c.JSON(http.StatusBadRequest, fmt.Sprintf("request deserialize error: %s", err))
		}

		if migrationTaskFromDB, err = migrationTaskRepository.GetMigrationByMessageId(request.MessageId); err != nil {
			if err.Error() != "not found" {
				return c.JSON(http.StatusInternalServerError, err.Error())
			}
		}

		if migrationTaskFromDB == nil {
			return c.JSON(http.StatusNotFound, fmt.Sprintf("%s MigrationTask not found", request.MessageId))
		}

		migrationTaskFromDB.Dc = request.Dc
		migrationTaskFromDB.VmName = request.VmName
		migrationTaskFromDB.InstanceName = request.InstanceName
		migrationTaskFromDB.Project = &request.Project
		migrationTaskFromDB.Flavor = &request.Flavor
		migrationTaskFromDB.PublicNetwork = &request.PublicNetwork
		migrationTaskFromDB.Network = &request.Network
		migrationTaskFromDB.SecurityGroup = &request.SecurityGroup
		migrationTaskFromDB.Key = &request.Key
		if request.Stage != "" {
			migrationTaskFromDB.Stage = request.Stage
		}

		migrationTaskFromDB.Status = enums.Running
		migrationTaskFromDB.Error = ""
		migrationTaskFromDB.Message = "Retrying migration"

		if err = migrationTaskRepository.UpsertMigrationTask(migrationTaskFromDB); err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}

		if err = publishEventMessage(ctx, messageBus, migrationTaskFromDB.MessageId, migrationTaskFromDB.Stage); err != nil {
			log.WithField("exception.backtrace", err).Errorf("An error occurred while publishing message")
			return c.JSON(http.StatusBadRequest, err.Error())
		}

		return c.JSON(http.StatusOK, "ok")
	})
}

// DeleteMigrationTask godoc
// @Accept  json
// @Param messageId path string true "Message ID"
// @Success 200
// @Failure 400 {object} string
// @Failure 404 {object} string
// @Tags Migration
// @Router /api/migration/delete/{messageID} [delete]
func DeleteMigrationTask(e *echo.Echo, taskRepository domain.MigrationTaskRepository, messageBus rabbitmq.Client) {
	e.DELETE("/api/migration/delete/:messageID", func(c echo.Context) error {

		var (
			task *entity.MigrationTask
			err  error
		)

		ctx := context.Background()

		messageID := strings.ToLower(c.Param("messageID"))

		if messageID == "" {
			return c.JSON(http.StatusBadRequest, "MessageID can not be empty!")
		}

		if task, err = taskRepository.GetMigrationByMessageId(messageID); err != nil {
			if err.Error() == "not found" {
				return c.JSON(http.StatusNotFound, "Migration task not found")
			} else {
				return c.JSON(http.StatusInternalServerError, err.Error())
			}
		}

		if err = taskRepository.UpsertMigrationDeletedTask(task); err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}

		if err = taskRepository.RemoveMigrationTaskByMessageId(messageID); err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}

		cleanVmMessage := event.CleanVmMessage{
			VmName: task.VmName,
		}

		if err = messageBus.Publish(ctx, "*", cleanVmMessage); err != nil {
			return err
		}

		return c.JSON(http.StatusOK, "Migration task successfully deleted")
	})
}

func publishEventMessage(ctx context.Context, messageBus rabbitmq.Client, messageId string, taskStage enums.Stage) error {

	var (
		message interface{}
		err     error
	)

	switch taskStage {
	case enums.PendingQueue:
		message = event.ExportVmMessage{MessageId: messageId}
	case enums.ExportVm:
		message = event.ExportVmMessage{MessageId: messageId}
	case enums.ConvertVm:
		message = event.ConvertVmMessage{MessageId: messageId}
	case enums.CreateImage:
		message = event.CreateImageMessage{MessageId: messageId}
	case enums.CreateVolume:
		message = event.CreateVolumeMessage{MessageId: messageId}
	case enums.CreateInstance:
		message = event.CreateInstanceMessage{MessageId: messageId}
	case enums.ReserveFloatingIP:
		message = event.ReserveFloatingIPMessage{MessageId: messageId}
	case enums.AssociateFloatingIP:
		message = event.AssociateFloatingIPMessage{MessageId: messageId}
	case enums.Finished:
		return errors.New(fmt.Sprintf("This task successfully complated! Retry aborted"))
	default:
		return errors.New(fmt.Sprintf("Invalid task stage for retrying migration"))
	}

	if err = messageBus.Publish(ctx, "*", message); err != nil {
		return err
	}

	return nil
}
