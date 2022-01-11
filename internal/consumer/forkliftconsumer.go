package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"forklift/internal/consumer/functions"
	"forklift/internal/domain"
	"forklift/internal/domain/entity"
	"forklift/internal/domain/event"
	"forklift/internal/domain/model"
	"forklift/pkg/config"
	"forklift/pkg/enums"
	"forklift/pkg/rabbitmq"
	log "github.com/sirupsen/logrus"
	"strings"
	"time"
)

type ForkliftConsumer struct {
	configuration   config.Configuration
	messageBus      rabbitmq.Client
	taskRepository  domain.MigrationTaskRepository
	vcenterHelper   *functions.VCenterHelper
	openstackHelper *functions.OpenstackHelper
}

func NewForkliftConsumer(configuration config.Configuration, messageBus rabbitmq.Client, taskRepository domain.MigrationTaskRepository, vcenterHelper *functions.VCenterHelper, openstackHelper *functions.OpenstackHelper) *ForkliftConsumer {
	return &ForkliftConsumer{
		configuration:   configuration,
		messageBus:      messageBus,
		taskRepository:  taskRepository,
		vcenterHelper:   vcenterHelper,
		openstackHelper: openstackHelper,
	}
}

func (c *ForkliftConsumer) Construct() {

	c.messageBus.AddConsumer("In.Forklift.Vcenter.ExportVm").
		SubscriberExchange("*", rabbitmq.Topic, "Forklift.Events:ExportVm").
		HandleConsumer(c.exportVm())

	c.messageBus.AddConsumer("In.Forklift.Vcenter.ConvertVm").
		SubscriberExchange("*", rabbitmq.Topic, "Forklift.Events:ConvertVm").
		HandleConsumer(c.convertVm())

	c.messageBus.AddConsumer("In.Forklift.Openstack.CreateImage").
		SubscriberExchange("*", rabbitmq.Topic, "Forklift.Events:CreateImage").
		HandleConsumer(c.createImage())

	c.messageBus.AddConsumer("In.Forklift.Openstack.CreateVolume").
		SubscriberExchange("*", rabbitmq.Topic, "Forklift.Events:CreateVolume").
		HandleConsumer(c.createVolume())

	c.messageBus.AddConsumer("In.Forklift.Openstack.CreateInstance").
		SubscriberExchange("*", rabbitmq.Topic, "Forklift.Events:CreateInstance").
		HandleConsumer(c.createInstance())

	c.messageBus.AddConsumer("In.Forklift.Openstack.ReserveFloatingIP").
		SubscriberExchange("*", rabbitmq.Topic, "Forklift.Events:ReserveFloatingIP").
		HandleConsumer(c.ReserveFloatingIP())

	c.messageBus.AddConsumer("In.Forklift.Openstack.AssociateFloatingIP").
		SubscriberExchange("*", rabbitmq.Topic, "Forklift.Events:AssociateFloatingIP").
		HandleConsumer(c.AssociateFloatingIP())

	c.messageBus.AddConsumer("In.Forklift.CleanVm").
		SubscriberExchange("*", rabbitmq.Topic, "Forklift.Events:CleanVm").
		HandleConsumer(c.cleanVm())
}

func (c *ForkliftConsumer) exportVm() func(message rabbitmq.Message) error {
	return func(message rabbitmq.Message) error {

		ctx := context.Background()

		var (
			eventMessage event.ExportVmMessage
			task         *entity.MigrationTask
			out          string
			err          error
		)

		if err = json.Unmarshal(message.Payload, &eventMessage); err != nil {
			return err
		}

		if task, err = c.taskRepository.GetMigrationByMessageId(eventMessage.MessageId); err != nil {
			return err
		}

		if task == nil {
			log.Info(ctx, fmt.Sprintf("%s task passed", eventMessage.MessageId))
			return nil
		}

		task.Stage = enums.ExportVm
		task.Status = enums.Running

		if err = c.UpdateTask(task); err != nil {
			return err
		}

		if out, err = c.vcenterHelper.ExportVm(*task); err != nil {
			log.WithField("exception.backtrace", err).Errorf("An error occurred while exporting vm from vcenter")
			if err = c.UpdateFailTask(ctx, task, out, err); err != nil {
				log.WithField("exception.backtrace", err).Errorf("An error occurred while updating failed task from db for exporting vm")
				return err
			}
			return err
		}

		task.Stage = enums.ConvertVm

		if err = c.UpdateTask(task); err != nil {
			return err
		}

		convertVmMessage := event.ConvertVmMessage{
			MessageId: task.MessageId,
		}

		if err = c.messageBus.Publish(ctx, "*", convertVmMessage); err != nil {
			log.WithField("exception.backtrace", err).Errorf("An error occurred while publishing StartTransferMessage ")
			return err
		}

		return nil
	}
}

func (c *ForkliftConsumer) convertVm() func(message rabbitmq.Message) error {
	return func(message rabbitmq.Message) error {

		ctx := context.Background()

		var (
			eventMessage event.ConvertVmMessage
			task         *entity.MigrationTask
			out          string
			err          error
		)

		if err = json.Unmarshal(message.Payload, &eventMessage); err != nil {
			return err
		}

		if task, err = c.taskRepository.GetMigrationByMessageId(eventMessage.MessageId); err != nil {
			return err
		}

		if task == nil {
			log.Info(ctx, fmt.Sprintf("%s task passed", eventMessage.MessageId))
			return nil
		}

		if out, err = c.vcenterHelper.ConvertVM(*task); err != nil {
			log.WithField("exception.backtrace", err).Errorf("An error occurred while converting vm")
			if err = c.UpdateFailTask(ctx, task, out, err); err != nil {
				log.WithField("exception.backtrace", err).Errorf("An error occurred while converting vm fail")
				return err
			}
			return err
		}

		task.Stage = enums.CreateImage

		if err = c.UpdateTask(task); err != nil {
			return err
		}

		CreateImageMessage := event.CreateImageMessage{
			MessageId: task.MessageId,
		}

		if err = c.messageBus.Publish(ctx, "*", CreateImageMessage); err != nil {
			return err
		}

		return nil
	}
}

func (c *ForkliftConsumer) createImage() func(message rabbitmq.Message) error {
	return func(message rabbitmq.Message) error {

		ctx := context.Background()

		var (
			eventMessage        event.CreateImageMessage
			task                *entity.MigrationTask
			ResponseCreateImage *model.ResponseImage
			out                 string
			size                float64
			err                 error
		)

		if err = json.Unmarshal(message.Payload, &eventMessage); err != nil {
			return err
		}

		if task, err = c.taskRepository.GetMigrationByMessageId(eventMessage.MessageId); err != nil {
			return err
		}

		if task == nil {
			log.Info(ctx, fmt.Sprintf("%s task passed", eventMessage.MessageId))
			return nil
		}

		if out, ResponseCreateImage, err = c.openstackHelper.CreateImage(*task); err != nil {
			log.WithField("exception.backtrace", err).Errorf("An error occurred while creating image")
			if err = c.UpdateFailTask(ctx, task, out, err); err != nil {
				log.WithField("exception.backtrace", err).Errorf("An error occurred while creating image fail")
				return err
			}
			return err
		}

		task.Image.Name = ResponseCreateImage.Name
		task.Image.Id = ResponseCreateImage.Id
		task.Image.Size = ResponseCreateImage.Size
		task.Stage = enums.CreateVolume

		if task.Image.Size == 0 {
			if size, err = c.openstackHelper.ImageSize(task); err != nil {
				log.WithField("exception.backtrace", err).Errorf("An error occurred while getting image size")
				if err = c.UpdateFailTask(ctx, task, out, err); err != nil {
					log.WithField("exception.backtrace", err).Errorf("An error occurred while creating image fail")
					return err
				}
				return err
			}
			task.Image.Size = size
		}

		if err = c.UpdateTask(task); err != nil {
			return err
		}

		CreateVolumeMessage := event.CreateVolumeMessage{
			MessageId: task.MessageId,
		}

		if err = c.messageBus.Publish(ctx, "*", CreateVolumeMessage); err != nil {
			return err
		}

		return nil
	}
}

func (c *ForkliftConsumer) createVolume() func(message rabbitmq.Message) error {
	return func(message rabbitmq.Message) error {

		ctx := context.Background()

		var (
			eventMessage         event.CreateVolumeMessage
			task                 *entity.MigrationTask
			ResponseCreateVolume *model.ResponseVolume
			out                  string
			err                  error
		)

		if err = json.Unmarshal(message.Payload, &eventMessage); err != nil {
			return err
		}

		if task, err = c.taskRepository.GetMigrationByMessageId(eventMessage.MessageId); err != nil {
			return err
		}

		if task == nil {
			log.Info(ctx, fmt.Sprintf("%s task passed", eventMessage.MessageId))
			return nil
		}

		if out, ResponseCreateVolume, err = c.openstackHelper.CreateVolume(*task); err != nil {
			log.WithField("exception.backtrace", err).Errorf("An error occurred while creating volume")
			if err = c.UpdateFailTask(ctx, task, out, err); err != nil {
				log.WithField("exception.backtrace", err).Errorf("An error occurred while creating volume fail")
				return err
			}
			return err
		}

		task.Volume.Name = ResponseCreateVolume.Name
		task.Volume.Id = ResponseCreateVolume.Id
		task.Volume.Size = ResponseCreateVolume.Size
		task.Volume.Bootable = ResponseCreateVolume.Bootable
		task.Stage = enums.CreateInstance

		if err = c.UpdateTask(task); err != nil {
			return err
		}

		CreateInstanceMessage := event.CreateInstanceMessage{
			MessageId: task.MessageId,
		}

		if err = c.messageBus.Publish(ctx, "*", CreateInstanceMessage); err != nil {
			return err
		}

		return nil
	}
}

func (c *ForkliftConsumer) createInstance() func(message rabbitmq.Message) error {
	return func(message rabbitmq.Message) error {

		var ctx = context.Background()

		var (
			eventMessage           event.CreateInstanceMessage
			task                   *entity.MigrationTask
			responseCreateInstance *model.ResponseInstance
			out                    string
			err                    error
		)

		if err = json.Unmarshal(message.Payload, &eventMessage); err != nil {
			return err
		}

		if task, err = c.taskRepository.GetMigrationByMessageId(eventMessage.MessageId); err != nil {
			return err
		}

		if task == nil {
			log.Info(ctx, fmt.Sprintf("%s task passed", eventMessage.MessageId))
			return nil
		}

		functions.NewWait(
			func(stop chan bool) {
				var (
					result               bool
					responseVolumeStatus *model.ResponseVolume
					taskFromDB           *entity.MigrationTask
				)

				if responseVolumeStatus, err = c.openstackHelper.VolumeStatus(task); err != nil {
					log.WithField("exception.backtrace", err).Errorf("An error occurred while getting volume status from OSP")
				}

				task.Volume.Bootable = responseVolumeStatus.Bootable

				if err = c.UpdateTask(task); err != nil {
					log.WithField("exception.backtrace", err).Errorf("An error occurred while updating task")
				}

				if taskFromDB, err = c.taskRepository.GetMigrationByMessageId(task.MessageId); err != nil {
					log.WithField("exception.backtrace", err).Errorf("An error occurred while getting message")
				}

				if taskFromDB.Volume.Bootable == "true" {
					result = true
				}

				if result {
					stop <- result
				}
			},
			"is volume bootable",
		)

		if out, responseCreateInstance, err = c.openstackHelper.CreateInstance(*task); err != nil {
			log.WithField("exception.backtrace", err).Errorf("An error occurred while creating instance")
			if err = c.UpdateFailTask(ctx, task, out, err); err != nil {
				log.WithField("exception.backtrace", err).Errorf("An error occurred while creating instance fail")
				return err
			}
			return err
		}
		task.Instance.Name = responseCreateInstance.Name
		task.Instance.Id = responseCreateInstance.Id
		task.Stage = enums.ReserveFloatingIP

		if err = c.UpdateTask(task); err != nil {
			return err
		}

		ReserveFloatingIPMessage := event.ReserveFloatingIPMessage{
			MessageId: task.MessageId,
		}

		if err = c.messageBus.Publish(ctx, "*", ReserveFloatingIPMessage); err != nil {
			return err
		}

		return nil
	}
}

func (c ForkliftConsumer) ReserveFloatingIP() func(message rabbitmq.Message) error {
	return func(message rabbitmq.Message) error {

		var ctx = context.Background()

		var (
			eventMessage              event.ReserveFloatingIPMessage
			task                      *entity.MigrationTask
			responseReserveFloatingIP *model.ResponseFloatingIP
			out                       string
			err                       error
		)

		if err = json.Unmarshal(message.Payload, &eventMessage); err != nil {
			return err
		}

		if task, err = c.taskRepository.GetMigrationByMessageId(eventMessage.MessageId); err != nil {
			return err
		}

		if task == nil {
			log.Info(ctx, fmt.Sprintf("%s task passed", eventMessage.MessageId))
			return nil
		}

		if out, responseReserveFloatingIP, err = c.openstackHelper.ReserveFloatingIP(*task); err != nil {
			log.WithField("exception.backtrace", err).Errorf("An error occurred while associating floating ip")
			if err = c.UpdateFailTask(ctx, task, out, err); err != nil {
				log.WithField("exception.backtrace", err).Errorf("An error occurred while updating failed task from db for associating floating ip")
				return err
			}
			return err
		}

		task.FloatingIP.Name = responseReserveFloatingIP.Name
		task.FloatingIP.Id = responseReserveFloatingIP.Id
		task.FloatingIP.FloatingIpAddress = responseReserveFloatingIP.FloatingIpAddress

		task.Stage = enums.AssociateFloatingIP

		if err = c.UpdateTask(task); err != nil {
			return err
		}

		AssociateFloatingIPMessage := event.AssociateFloatingIPMessage{
			MessageId: task.MessageId,
		}

		if err = c.messageBus.Publish(ctx, "*", AssociateFloatingIPMessage); err != nil {
			return err
		}

		return nil
	}
}

func (c ForkliftConsumer) AssociateFloatingIP() func(message rabbitmq.Message) error {
	return func(message rabbitmq.Message) error {

		var ctx = context.Background()

		var (
			eventMessage event.AssociateFloatingIPMessage
			task         *entity.MigrationTask
			out          string
			err          error
		)

		if err = json.Unmarshal(message.Payload, &eventMessage); err != nil {
			return err
		}

		if task, err = c.taskRepository.GetMigrationByMessageId(eventMessage.MessageId); err != nil {
			return err
		}

		if task == nil {
			log.Info(ctx, fmt.Sprintf("%s task passed", eventMessage.MessageId))
			return nil
		}

		if out, err = c.openstackHelper.AssociateFloatingIP(*task); err != nil {
			log.WithField("exception.backtrace", err).Errorf("An error occurred while associating floating ip")
			if err = c.UpdateFailTask(ctx, task, out, err); err != nil {
				log.WithField("exception.backtrace", err).Errorf("An error occurred while updating failed task from db for associating floating ip")
				return err
			}
			return err
		}

		task.EndDate = time.Now().UTC().Add(3 * time.Hour)
		task.Stage = enums.Finished
		task.Status = enums.Completed
		task.Message = "Migration task successfully completed"

		if err = c.UpdateTask(task); err != nil {
			return err
		}

		if err = publishCleanVm(ctx, c.messageBus, task.VmName); err != nil {
			return err
		}

		return nil
	}
}

func (c *ForkliftConsumer) cleanVm() func(message rabbitmq.Message) error {
	return func(message rabbitmq.Message) error {

		var (
			eventMessage event.CleanVmMessage
			err          error
		)

		if err = json.Unmarshal(message.Payload, &eventMessage); err != nil {
			return err
		}

		if err = c.vcenterHelper.CleanVm(eventMessage.VmName); err != nil {
			return err
		}

		return nil
	}
}

func publishCleanVm(ctx context.Context, messageBus rabbitmq.Client, vmName string) (err error) {
	cleanVmMessage := event.CleanVmMessage{
		VmName: vmName,
	}

	if err = messageBus.Publish(ctx, "*", cleanVmMessage); err != nil {
		return err
	}

	return nil
}

func (c *ForkliftConsumer) UpdateFailTask(ctx context.Context, task *entity.MigrationTask, message string, errorParam error) error {

	var err error

	if strings.Contains(errorParam.Error(), "kill") {
		log.Warningf("%v: %s", errorParam, task.VmName)
		if err = publishCleanVm(ctx, c.messageBus, task.VmName); err != nil {
			return err
		}
		return nil
	}

	task.Error = errorParam.Error()
	task.Message = message
	task.Status = enums.Failed

	if err = c.UpdateTask(task); err != nil {
		return err
	}

	return err
}

func (c *ForkliftConsumer) UpdateTask(task *entity.MigrationTask) error {
	if err := c.taskRepository.UpdateMigrationTask(task); err != nil {
		return err
	}
	return nil
}
