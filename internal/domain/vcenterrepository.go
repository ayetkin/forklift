package domain

import (
	"fmt"
	"forklift/pkg/config"
	"forklift/pkg/helper"
	log "github.com/sirupsen/logrus"
)

type VcenterRepository interface {
	GetAllPoweredOffVms(dc string) ([]string, error)
	GetAllDcs() ([]string, error)
}

type vcenterRepository struct {
	configuration *config.Configuration
}

func (r *vcenterRepository) GetAllPoweredOffVms(dc string) ([]string, error) {
	vcenterUrl, err := helper.VcenterUrl(r.configuration)
	if err != nil {
		log.WithField("error.backtrace", err).Errorf("An error occurred while generating vcneter url")
		return nil, err
	}
	cmd := fmt.Sprintf("govc find -u='%s' -k=true -json /%s/vm/ -type m -runtime.powerState poweredOff", vcenterUrl, dc)
	out, err := helper.ExecuteCommand(cmd)
	if err != nil {
		log.WithField("error.backtrace", err).Errorf("An error occurred while executing command")
		return nil, err
	}

	vmList, err := helper.JsonToArray(out)
	return vmList, err
}

func (r *vcenterRepository) GetAllDcs() ([]string, error) {
	vcenterUrl, err := helper.VcenterUrl(r.configuration)
	if err != nil {
		log.WithField("error.backtrace", err).Errorf("An error occurred while generating vcneter url")
		return nil, err
	}
	cmd := fmt.Sprintf("govc find -u='%s' -k=true -json -type d", vcenterUrl)
	out, err := helper.ExecuteCommand(cmd)
	if err != nil {
		log.WithField("error.backtrace", err).Errorf("An error occurred while executing command")
		return nil, err
	}

	dcList, err := helper.JsonToArray(out)
	return dcList, err
}

func NewVcenterRepository(configuration *config.Configuration) VcenterRepository {
	return &vcenterRepository{
		configuration: configuration,
	}
}
