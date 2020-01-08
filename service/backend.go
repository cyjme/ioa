package service

type Backend struct {
	Id  string `mapstructure:"id"`
	Uri string `mapstructure:"uri"`
	Qps int    `mapstructure:"qps"`

	Weight int `mapstructure:"weight"`
}
