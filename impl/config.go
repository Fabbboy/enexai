package impl

import (
	"gopkg.in/ini.v1"
)

type ApiConfig struct {
	Url         string  `ini:"url"`
	Key         string  `ini:"key"`
	Temperature float64 `ini:"temperature"`
}

type ModelConfig struct {
	Classifier string `ini:"classifier"`
	Writer     string `ini:"writer"`
}

type Config struct {
	ApiConfig    ApiConfig   `ini:"ai"`
	ModelsConfig ModelConfig `ini:"models"`
}

func LoadConfig(file string) (*Config, error) {
	config := &Config{}
	err := ini.MapTo(config, file)
	if err != nil {
		return nil, err
	}

	return config, nil
}
