package functions

import (
	"encoding/json"
	"errors"
	"fmt"
	"forklift/internal/domain"
	"forklift/internal/domain/entity"
	"forklift/internal/domain/model"
	"forklift/pkg/helper"
	"strconv"
	"strings"
	"time"
)

type OpenstackHelper struct {
	commander      *Commander
	taskRepository domain.MigrationTaskRepository
}

func NewOpenstackHelper(taskRepository domain.MigrationTaskRepository) *OpenstackHelper {
	var openstackHelper = new(OpenstackHelper)
	openstackHelper.taskRepository = taskRepository
	openstackHelper.commander = NewCommander(taskRepository)
	return openstackHelper
}

func (h *OpenstackHelper) CreateImage(task entity.MigrationTask) (string, *model.ResponseImage, error) {
	var (
		result *model.ResponseImage
		os     string
		out    []byte
		cmd    string
		err    error
	)

	if os, err = h.DetectVmOS(task); err != nil {
		return "", nil, errors.New(fmt.Sprintf("Cant detect vm os! %v", err))
	}

	cmd = fmt.Sprintf("openstack image create --project %s --property hw_disk_bus='ide' --property hw_scsi_model=virtio-scsi --property hw_qemu_guest_agent=yes --property os_require_quiesce=yes --container-format bare --disk-format qcow2 --public --file vm/%s/%s.qcow2 %s --format json", task.Project.Name, task.VmName, task.VmName, task.InstanceName)

	if os == "Windows" {
		cmd = fmt.Sprintf("openstack image create --project %s  --property hw_disk_bus='ide' --property hw_scsi_model=virtio-scsi --property hw_qemu_guest_agent=yes --property os_require_quiesce=yes --property os_distro=windows --container-format bare --disk-format qcow2  --public --file vm/%s.qcow2 %s --format json", task.Project.Name, task.InstanceName, task.InstanceName)
	}

	if out, err = h.commander.executeCommand(task.MessageId, cmd); err != nil {
		return "", nil, err
	}

	if err = json.Unmarshal(out, &result); err != nil {
		return "", nil, err
	}

	return string(out), result, err
}

func (h *OpenstackHelper) ImageSize(task *entity.MigrationTask) (float64, error) {

	var (
		out  []byte
		cmd  string
		err  error
		size float64
	)

	time.Sleep(5 * time.Second)

	cmd = fmt.Sprintf("govc find /%s/vm/ -type m  -runtime.powerState poweredOff |grep %s |xargs -I $ govc ls -json $ |jq '.elements[0].Object.Summary.Storage.Uncommitted'", task.Dc, task.VmName)

	if out, err = helper.ExecuteCommand(cmd); err != nil {
		return size, err
	}

	if size, err = strconv.ParseFloat(string(out), 8); err != nil {
		return size, err
	}

	return size, err
}

func (h *OpenstackHelper) CreateVolume(task entity.MigrationTask) (string, *model.ResponseVolume, error) {

	var (
		ResponseCreateVolume *model.ResponseVolume
		out                  []byte
		cmd                  string
		err                  error
	)

	cmd = fmt.Sprintf("openstack volume create --image %s --size %s %s --format json",
		task.Image.Id, strconv.FormatFloat(task.Image.Size/float64(1073741824), 'f', 0, 64), task.InstanceName)

	if out, err = h.commander.executeCommand(task.MessageId, cmd); err != nil {
		return "", nil, err
	}

	if err = json.Unmarshal(out, &ResponseCreateVolume); err != nil {
		return "", nil, err
	}

	return string(out), ResponseCreateVolume, err
}

func (h *OpenstackHelper) VolumeStatus(task *entity.MigrationTask) (*model.ResponseVolume, error) {

	var (
		ResponseCreateVolume *model.ResponseVolume
		out                  []byte
		cmd                  string
		err                  error
	)

	cmd = fmt.Sprintf("openstack volume show %s --format json", task.Volume.Id)

	if out, err = helper.ExecuteCommand(cmd); err != nil {
		return nil, err
	}

	if err = json.Unmarshal(out, &ResponseCreateVolume); err != nil {
		return nil, err
	}

	return ResponseCreateVolume, err
}

func (h *OpenstackHelper) CreateInstance(task entity.MigrationTask) (string, *model.ResponseInstance, error) {

	var (
		result *model.ResponseInstance
		out    []byte
		cmd    string
		err    error
	)

	cmd = fmt.Sprintf("openstack server create --flavor %s --volume %s --network %s --security-group %s --key-name %s %s --format json",
		task.Flavor.Id, task.Volume.Id, task.Network.Id, task.SecurityGroup.Id, task.Key.Name, task.InstanceName)

	if out, err = h.commander.executeCommand(task.MessageId, cmd); err != nil {
		return "", nil, err
	}

	if err = json.Unmarshal(out, &result); err != nil {
		return "", nil, err
	}

	return string(out), result, err
}

func (h *OpenstackHelper) ReserveFloatingIP(task entity.MigrationTask) (string, *model.ResponseFloatingIP, error) {

	var (
		result *model.ResponseFloatingIP
		out    []byte
		cmd    string
		err    error
	)

	cmd = fmt.Sprintf("openstack floating ip create --subnet %s public --format json", task.PublicNetwork.Subnets[0])

	if out, err = h.commander.executeCommand(task.MessageId, cmd); err != nil {
		return "", nil, err
	}

	if err = json.Unmarshal(out, &result); err != nil {
		return "", nil, err
	}

	return string(out), result, err
}

func (h *OpenstackHelper) AssociateFloatingIP(task entity.MigrationTask) (string, error) {

	var (
		out []byte
		cmd string
		err error
	)

	if task.Instance == nil || task.FloatingIP == nil {
		return "", errors.New(fmt.Sprintf("Instance or Floating ip value nil!"))
	}

	cmd = fmt.Sprintf("openstack server add floating ip %s %s", task.Instance.Id, task.FloatingIP.FloatingIpAddress)

	if out, err = h.commander.executeCommand(task.MessageId, cmd); err != nil {
		return "", err
	}

	return string(out), err
}

func (h *OpenstackHelper) DetectVmOS(task entity.MigrationTask) (string, error) {

	var (
		out []byte
		cmd string
		err error
	)

	cmd = fmt.Sprintf("govc find /%s/ -type m -runtime.powerState poweredOff |grep %s |xargs -I $ govc vm.info $ |grep \"Guest name:\"", task.Dc, task.VmName)

	if out, err = helper.ExecuteCommand(cmd); err != nil {
		return "", err
	}

	if strings.Contains(string(out), "Windows") {
		return "Windows", nil
	}

	return "Linux", nil
}
