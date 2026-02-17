package impl

import (
	"gopkg.in/ini.v1"
)

type AiConfig struct {
	url         string  `ini:"api_url"`
	key         string  `ini:"api_key"`
	model       string  `ini:"model"`
	temperature float64 `ini:"temperature"`
}

type Config struct {
	aiConfig AiConfig `ini:"ai"`
}

func LoadConfig(file string) (*Config, error) {
	config := &Config{}
	err := ini.MapTo(config, file)
	if err != nil {
		return nil, err
	}

	return config, nil
}
