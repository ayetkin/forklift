package entity

import (
	"forklift/pkg/enums"
	"github.com/globalsign/mgo/bson"
	"time"
)

type Flavor struct {
	Name string `json:"name" bson:"Name,omitempty"`
	Id   string `json:"id,omitempty" bson:"Id,omitempty"`
}

type Key struct {
	Id          string  `json:"id,omitempty" bson:"Id,omitempty"`
	Name        string  `json:"name" bson:"Name,omitempty"`
	Fingerprint string  `json:"fingerprint,omitempty" bson:"Fingerprint,omitempty"`
	Size        float64 `json:"size,omitempty" bson:"Size,omitempty"`
}

type Project struct {
	Name string `json:"name" bson:"Name,omitempty"`
	Id   string `json:"id,omitempty" bson:"Id,omitempty"`
}

type Network struct {
	Name    string   `json:"name" bson:"Name,omitempty"`
	Id      string   `json:"id,omitempty" bson:"Id,omitempty"`
	Subnets []string `json:"subnets" bson:"Subnets,omitempty"`
}

type SecurityGroup struct {
	Name string `json:"name" bson:"Name,omitempty"`
	Id   string `json:"id,omitempty" bson:"Id,omitempty"`
}

type Image struct {
	Name string  `json:"name" bson:"Name,omitempty"`
	Id   string  `json:"id,omitempty" bson:"Id,omitempty"`
	Size float64 `json:"size" bson:"Size,omitempty"`
}

type Volume struct {
	Name     string  `json:"name" bson:"Name,omitempty"`
	Id       string  `json:"id,omitempty" bson:"Id,omitempty"`
	Bootable string  `json:"bootable,omitempty" bson:"Bootable,omitempty"`
	Size     float64 `json:"size" bson:"Size,omitempty"`
}

type FloatingIP struct {
	Name              string `json:"name" bson:"Name,omitempty"`
	Id                string `json:"id,omitempty" bson:"Id,omitempty"`
	FloatingIpAddress string `json:"floating_ip_address,omitempty" bson:"FloatingIpAddress,omitempty"`
}

type Instance struct {
	Name string `json:"name" bson:"Name,omitempty"`
	Id   string `json:"id,omitempty" bson:"Id,omitempty"`
}

type MigrationRequest struct {
	MessageId     string        `json:"messageId"`
	Dc            string        `json:"dc"`
	VmName        string        `json:"vmName"`
	InstanceName  string        `json:"instanceName"`
	Project       Project       `json:"project"`
	Flavor        Flavor        `json:"flavor"`
	PublicNetwork Network       `json:"publicNetwork"`
	Network       Network       `json:"network"`
	SecurityGroup SecurityGroup `json:"securityGroup"`
	Key           Key           `json:"key"`
}

type MigrationTask struct {
	MessageId     string         `json:"messageId" bson:"MessageId"`
	StartDate     time.Time      `json:"startDate" bson:"StartDate"`
	EndDate       time.Time      `json:"endDate" bson:"EndDate"`
	Status        enums.Status   `json:"status" bson:"Status"`
	Stage         enums.Stage    `json:"stage" bson:"Stage"`
	Message       string         `json:"message" bson:"Message"`
	Error         string         `json:"error" bson:"Error"`
	Dc            string         `json:"dc" bson:"Dc"`
	VmName        string         `json:"vmName" bson:"VmName"`
	InstanceName  string         `json:"instanceName" bson:"instanceName"`
	Instance      *Instance      `json:"instance" bson:"Instance"`
	Project       *Project       `json:"project" bson:"Project"`
	Flavor        *Flavor        `json:"flavor" bson:"Flavor"`
	Network       *Network       `json:"network" bson:"Network"`
	PublicNetwork *Network       `json:"publicNetwork" bson:"PublicNetwork"`
	SecurityGroup *SecurityGroup `json:"securityGroup" bson:"SecurityGroup"`
	Key           *Key           `json:"key" bson:"Key"`
	Image         *Image         `json:"image,omitempty" bson:"Image,omitempty"`
	Volume        *Volume        `json:"volume,omitempty" bson:"Volume,omitempty"`
	FloatingIP    *FloatingIP    `json:"floatingIP,omitempty" bson:"FloatingIP,omitempty"`
}

func NewMigrationTask(request MigrationRequest) *MigrationTask {
	return &MigrationTask{
		MessageId:     bson.NewObjectId().Hex(),
		StartDate:     time.Now().UTC().Add(3 * time.Hour),
		EndDate:       time.Time{},
		Status:        enums.Pending,
		Stage:         enums.PendingQueue,
		Message:       "",
		Error:         "",
		Dc:            request.Dc,
		VmName:        request.VmName,
		InstanceName:  request.InstanceName,
		Project:       &request.Project,
		Flavor:        &request.Flavor,
		Network:       &request.Network,
		PublicNetwork: &request.PublicNetwork,
		SecurityGroup: &request.SecurityGroup,
		Key:           &request.Key,
		Instance:      new(Instance),
		Image:         new(Image),
		Volume:        new(Volume),
		FloatingIP:    new(FloatingIP),
	}
}
