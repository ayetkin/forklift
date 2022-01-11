package model

type ResponseInstance struct {
	Name string `json:"name"`
	Id   string `json:"id,omitempty" bson:"Id,omitempty"`
}

type ResponseImage struct {
	Name string  `json:"name,omitempty" bson:"name,omitempty"`
	Id   string  `json:"id,omitempty" bson:"Id,omitempty"`
	Size float64 `json:"size,omitempty" bson:"Size,omitempty"`
}

type ResponseVolume struct {
	Name     string  `json:"name,omitempty" bson:"name,omitempty"`
	Id       string  `json:"id,omitempty" bson:"Id,omitempty"`
	Size     float64 `json:"size,omitempty" bson:"Size,omitempty"`
	Bootable string  `json:"bootable,omitempty" bson:"Bootable,omitempty"`
}

type ResponseFloatingIP struct {
	Name              string `json:"name,omitempty" bson:"Name,omitempty"`
	Id                string `json:"id,omitempty" bson:"Id,omitempty"`
	FloatingIpAddress string `json:"floating_ip_address" bson:"FloatingIpAddress"`
}
