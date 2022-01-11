package model

import (
	"forklift/internal/domain/entity"
	"forklift/pkg/enums"
)

type MigrationRequest struct {
	Dc            string               `json:"dc"`
	VmName        string               `json:"vmName"`
	InstanceName  string               `json:"instanceName"`
	Project       entity.Project       `json:"project"`
	Flavor        entity.Flavor        `json:"flavor"`
	PublicNetwork entity.Network       `json:"publicNetwork"`
	Network       entity.Network       `json:"network"`
	SecurityGroup entity.SecurityGroup `json:"securityGroup"`
	Key           entity.Key           `json:"key"`
}

type RetryMigrationRequest struct {
	MessageId     string               `json:"messageId"`
	Dc            string               `json:"dc"`
	VmName        string               `json:"vmName"`
	InstanceName  string               `json:"instanceName"`
	Stage         enums.Stage          `json:"stage"`
	Project       entity.Project       `json:"project"`
	Flavor        entity.Flavor        `json:"flavor"`
	PublicNetwork entity.Network       `json:"publicNetwork"`
	Network       entity.Network       `json:"network"`
	SecurityGroup entity.SecurityGroup `json:"securityGroup"`
	Key           entity.Key           `json:"key"`
}

type VmListRequest struct {
	Dc string `json:"dc"`
}

type AddUserRequest struct {
	FullName  string `json:"fullName"`
	UserName  string `json:"userName"`
	UserEmail string `json:"userEmail"`
	UserId    string `json:"userId"`
	Groups    struct {
		Auth bool `json:"auth"`
	}
	Guid string `json:"guid"`
	Iat  int64  `json:"iat"`
	Exp  int64  `json:"exp"`
}
