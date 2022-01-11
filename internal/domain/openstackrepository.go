package domain

import (
	"encoding/json"
	"forklift/internal/domain/entity"
	"forklift/pkg/helper"
	log "github.com/sirupsen/logrus"
)

type OpenstackRepository interface {
	GetFlavors() ([]entity.Flavor, error)
	GetProjects() ([]entity.Project, error)
	GetNetworks(projectName string) ([]entity.Network, error)
	GetSecurityGroups(projectName string) ([]entity.SecurityGroup, error)
	GetKeys(projectName string) ([]entity.Key, error)
}

type openstackRepository struct {
	logger log.Logger
}

func (r *openstackRepository) GetProjects() ([]entity.Project, error) {
	var (
		result []entity.Project
		out    []byte
		err    error
	)

	if out, err = runOpenstackCommand("openstack project list --format json"); err != nil {
		return nil, err
	}

	if err = json.Unmarshal(out, &result); err != nil {
		return nil, err
	}

	return result, err
}

func (r *openstackRepository) GetFlavors() ([]entity.Flavor, error) {
	var (
		result []entity.Flavor
		out    []byte
		err    error
	)

	if out, err = runOpenstackCommand("openstack flavor list --format json"); err != nil {
		return nil, err
	}

	if err = json.Unmarshal(out, &result); err != nil {
		return nil, err
	}

	return result, err
}

func (r *openstackRepository) GetNetworks(projectName string) ([]entity.Network, error) {
	var (
		result []entity.Network
		out    []byte
		err    error
	)

	if out, err = runOpenstackCommand("openstack network list --project " + projectName + " --format json"); err != nil {
		return nil, err
	}

	if err = json.Unmarshal(out, &result); err != nil {
		return nil, err
	}

	return result, err
}

func (r *openstackRepository) GetSecurityGroups(projectName string) ([]entity.SecurityGroup, error) {
	var (
		result []entity.SecurityGroup
		out    []byte
		err    error
	)

	if out, err = runOpenstackCommand("openstack security group list --project " + projectName + " --format json"); err != nil {
		return nil, err
	}

	if err = json.Unmarshal(out, &result); err != nil {
		return nil, err
	}

	return result, err
}

func (r *openstackRepository) GetKeys(projectName string) ([]entity.Key, error) {
	var (
		result []entity.Key
		out    []byte
		err    error
	)

	if out, err = runOpenstackCommand("openstack keypair list --project " + projectName + " --format json"); err != nil {
		return nil, err
	}

	if err = json.Unmarshal(out, &result); err != nil {
		return nil, err
	}

	return result, err
}

func runOpenstackCommand(cmd string) (out []byte, err error) {

	if out, err = helper.ExecuteCommand(cmd); err != nil {
		return nil, err
	}

	return out, nil
}

func NewOpenstackRepository() OpenstackRepository {
	return &openstackRepository{}
}
