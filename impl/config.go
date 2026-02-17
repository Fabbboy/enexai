package impl

import (
	"gopkg.in/ini.v1"
)

type AiConfig struct {
	Url         string  `ini:"api_url"`
	Key         string  `ini:"api_key"`
	Model       string  `ini:"model"`
	Temperature float64 `ini:"temperature"`
}

type Config struct {
	AiConfig AiConfig `ini:"ai"`
}

func LoadConfig(file string) (*Config, error) {
	config := &Config{}
	err := ini.MapTo(config, file)
	if err != nil {
		return nil, err
	}

	return config, nil
}
