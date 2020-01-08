package service

type ServiceDefinition struct {
	Id  string `mapstructure:"id"`
	Uri string `mapstructure:"uri"`
	Qps int    `mapstructure:"qps"`

	Backends []Backend `mapstructure:"backends"`
}

type serviceDefinitionReader interface {
	getServiceDefinitions(config string) ([]ServiceDefinition, error)
}

type ServiceDefinitionWriter interface {
	AddServiceDefinition(config string, sd ServiceDefinition) error
	UpdateServiceDefinition(config string, sd ServiceDefinition) error
	DeleteServiceDefinition(config string, sid string) error
}