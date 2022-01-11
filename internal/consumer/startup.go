package consumer

import (
	"forklift/internal/consumer/functions"
	"forklift/internal/domain"
	"forklift/internal/domain/event"
	"forklift/pkg/config"
	"forklift/pkg/helper"
	"forklift/pkg/mongo"
	"forklift/pkg/rabbitmq"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"time"
)

func Init(cmd *cobra.Command, args []string) error {

	var configuration config.Configuration

	err := viper.Unmarshal(&configuration)
	if err != nil {
		panic("configuration is invalid!")
	}

	if err = helper.LoginOSP(&configuration); err != nil {
		logrus.Fatal(err)
	}

	var messageBus = rabbitmq.NewRabbitMqClient(
		configuration.RabbitMQ.Host,
		configuration.RabbitMQ.Username,
		configuration.RabbitMQ.Password,
		"",
		rabbitmq.RetryCount(0),
		rabbitmq.PrefetchCount(1))

	messageBus.AddPublisher("Forklift.Events:ExportVm", rabbitmq.Topic, event.ExportVmMessage{})
	messageBus.AddPublisher("Forklift.Events:ConvertVm", rabbitmq.Topic, event.ConvertVmMessage{})
	messageBus.AddPublisher("Forklift.Events:CreateImage", rabbitmq.Topic, event.CreateImageMessage{})
	messageBus.AddPublisher("Forklift.Events:CreateVolume", rabbitmq.Topic, event.CreateVolumeMessage{})
	messageBus.AddPublisher("Forklift.Events:CreateInstance", rabbitmq.Topic, event.CreateInstanceMessage{})
	messageBus.AddPublisher("Forklift.Events:ReserveFloatingIP", rabbitmq.Topic, event.ReserveFloatingIPMessage{})
	messageBus.AddPublisher("Forklift.Events:AssociateFloatingIP", rabbitmq.Topic, event.AssociateFloatingIPMessage{})
	messageBus.AddPublisher("Forklift.Events:CleanVm", rabbitmq.Topic, event.CleanVmMessage{})

	var mongoClient, mongoErr = mongo.NewClient(configuration.MongoDB.Url, 10*time.Second)
	if mongoErr != nil {
		return mongoErr
	}

	var taskRepository = domain.NewMigrationRepository(mongoClient, configuration.MongoDB.Database)
	var vcenterHelper = functions.NewVCenterHelper(taskRepository, configuration)
	var openstackHelper = functions.NewOpenstackHelper(taskRepository)
	var forkliftConsumer = NewForkliftConsumer(configuration, messageBus, taskRepository, vcenterHelper, openstackHelper)

	forkliftConsumer.Construct()

	return messageBus.RunConsumers()
}
