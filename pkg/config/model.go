package config

type Configuration struct {
	MongoDB   MongoDBConfiguration
	RabbitMQ  RabbitMQConfiguration
	Openstack OpenstackConfiguration
	Vault     VaultConfiguration
	Vcenter   VcenterConfiguration
	Server    ServerConfiguration
	Auth      Auth
}

type MongoDBConfiguration struct {
	Url      string
	Database string
}

type RabbitMQConfiguration struct {
	Host     []string
	Port     string
	Username string
	Password string
}

type OpenstackConfiguration struct {
	AuthUrl            string
	ProjectId          string
	ProjectName        string
	UserDomainName     string
	Username           string
	Password           string
	RegionName         string
	Interface          string
	IdentityApiVersion string
}

type VaultConfiguration struct {
	Address string
	Secret  string
}

type VcenterConfiguration struct {
	Url      string
	Username string
}

type ServerConfiguration struct {
	Port string
}

type Auth struct {
	ValidationUrl string
}
