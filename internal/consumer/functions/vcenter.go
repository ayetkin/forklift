package functions

import (
	"fmt"
	"forklift/internal/domain"
	"forklift/internal/domain/entity"
	"forklift/pkg/config"
	"forklift/pkg/helper"
	log "github.com/sirupsen/logrus"
	"net/url"
	"os"
)

type VCenterHelper struct {
	commander      *Commander
	taskRepository domain.MigrationTaskRepository
	configuration  config.Configuration
}

func NewVCenterHelper(taskRepository domain.MigrationTaskRepository, configuration config.Configuration) *VCenterHelper {
	var vcenterHelper = new(VCenterHelper)
	vcenterHelper.taskRepository = taskRepository
	vcenterHelper.commander = NewCommander(taskRepository)
	vcenterHelper.configuration = configuration
	return vcenterHelper
}

func (h *VCenterHelper) ExportVm(task entity.MigrationTask) (string, error) {

	var (
		out []byte
		cmd string
		u   *url.URL
		err error
	)

	if u, err = helper.VcenterUrl(&h.configuration); err != nil {
		return "", err
	}

	cmd = fmt.Sprintf("govc export.ovf -u='%s' -k=true -dc=%s -vm=%s vm/", u, task.Dc, task.VmName)

	if out, err = h.commander.executeCommand(task.MessageId, cmd); err != nil {
		return "", err
	}

	return string(out), nil
}

func (h *VCenterHelper) ConvertVM(task entity.MigrationTask) (string, error) {

	var (
		out []byte
		cmd string
		err error
	)

	cmd = fmt.Sprintf("qemu-img convert -f vmdk -O qcow2 vm/%s/%s-disk-0.vmdk vm/%s/%s.qcow2", task.VmName, task.VmName, task.VmName, task.VmName)

	if out, err = h.commander.executeCommand(task.MessageId, cmd); err != nil {
		return "", err
	}

	return string(out), nil
}

func (h *VCenterHelper) CleanVm(VmName string) error {

	if err := os.RemoveAll("vm/" + VmName + "/"); err != nil {
		return err
	}

	log.Warningf("Files Deleted: %s", VmName)

	return nil
}
