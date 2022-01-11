package api

import (
	"fmt"
	"forklift/internal/api/controller"
	"forklift/internal/api/docs"
	"forklift/internal/domain"
	"forklift/internal/domain/event"
	"forklift/pkg/config"
	"forklift/pkg/helper"
	"forklift/pkg/mongo"
	"forklift/pkg/rabbitmq"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	echoSwagger "github.com/swaggo/echo-swagger"
	"net/http"
	"time"
)

func Init(cmd *cobra.Command, args []string) error {
	docs.Init()

	var (
		configuration config.Configuration
		e             *echo.Echo
		err           error
	)

	err = viper.Unmarshal(&configuration)
	if err != nil {
		log.Fatal("Configuration is invalid!")
	}

	if err = helper.LoginOSP(&configuration); err != nil {
		log.Fatal(err)
	}

	e = echo.New()
	e.Debug = true
	e.HideBanner = true
	e.HidePort = true

	//e.Use(Process)

	e.GET("/api/healthz", func(c echo.Context) error {
		return c.String(http.StatusOK, "Service Up")
	})

	e.GET("/api/status", func(c echo.Context) error {
		return c.String(http.StatusOK, "Service Up")
	})

	e.GET("/api/swagger/*", echoSwagger.WrapHandler)

	var mongoClient, mongoErr = mongo.NewClient(configuration.MongoDB.Url, 10*time.Second)
	if mongoErr != nil {
		return mongoErr
	}

	var messageBus = rabbitmq.NewRabbitMqClient(
		configuration.RabbitMQ.Host,
		configuration.RabbitMQ.Username,
		configuration.RabbitMQ.Password,
		"",
		rabbitmq.RetryCount(0),
		rabbitmq.PrefetchCount(10),
	)

	messageBus.AddPublisher("Forklift.Events:ExportVm", rabbitmq.Topic, event.ExportVmMessage{})
	messageBus.AddPublisher("Forklift.Events:ConvertVm", rabbitmq.Topic, event.ConvertVmMessage{})
	messageBus.AddPublisher("Forklift.Events:CreateImage", rabbitmq.Topic, event.CreateImageMessage{})
	messageBus.AddPublisher("Forklift.Events:CreateVolume", rabbitmq.Topic, event.CreateVolumeMessage{})
	messageBus.AddPublisher("Forklift.Events:CreateInstance", rabbitmq.Topic, event.CreateInstanceMessage{})
	messageBus.AddPublisher("Forklift.Events:ReserveFloatingIP", rabbitmq.Topic, event.ReserveFloatingIPMessage{})
	messageBus.AddPublisher("Forklift.Events:AssociateFloatingIP", rabbitmq.Topic, event.AssociateFloatingIPMessage{})
	messageBus.AddPublisher("Forklift.Events:CleanVm", rabbitmq.Topic, event.CleanVmMessage{})

	var migrationTaskRepository = domain.NewMigrationRepository(mongoClient, configuration.MongoDB.Database)
	var vcenterRepository = domain.NewVcenterRepository(&configuration)
	var openstackRepository = domain.NewOpenstackRepository()

	controller.GetMigrationTasks(e, migrationTaskRepository)
	controller.GetPoweredOffVms(e, vcenterRepository)
	controller.GetDatacenters(e, vcenterRepository)
	controller.GetProjects(e, openstackRepository)
	controller.GetFlavors(e, openstackRepository)
	controller.GetNetworks(e, openstackRepository)
	controller.GetKeys(e, openstackRepository)
	controller.GetSecurityGroups(e, openstackRepository)
	controller.StartMigrationTask(e, migrationTaskRepository, messageBus)
	controller.RetryMigrationTask(e, migrationTaskRepository, messageBus)
	controller.DeleteMigrationTask(e, migrationTaskRepository, messageBus)

	log.Warningf("Started http web server [::]:%s", configuration.Server.Port)
	if err = e.Start(fmt.Sprintf(":%s", configuration.Server.Port)); err != nil {
		panic(err)
	}

	return nil
}
