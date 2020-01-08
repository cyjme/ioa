package config

import "errors"

type config interface {
	Add(key string, value string) bool
	Set(key string, value string) bool
	Delete(key string) bool
	Get(key string) string
	List() map[string]string
}

var Config config
var err error

func InitConfig(typ string, arg string) error {

	switch typ {
	case "yaml":
		Config, err = NewYamlConfig(arg)
		if err != nil {
			return err
		}
	case "json":
		Config, err = NewJsonConfig(arg)
		if err != nil {
			return err
		}
	default:
		return errors.New("init config error, you must provide a valid config type")
	}

	return nil
}
