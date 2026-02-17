package impl

import (
	"gopkg.in/ini.v1"
)

type AiConfig struct {
	Url         string  `ini:"api_url"`
	Key         string  `ini:"api_key"`
	Temperature float64 `ini:"temperature"`
}

type ModelConfig struct {
	Classifier string `ini:"classifier"`
	Writer     string `ini:"writer"`
}

type Config struct {
	AiConfig     AiConfig    `ini:"ai"`
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
