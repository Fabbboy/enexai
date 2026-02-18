package impl

import (
	"github.com/go-playground/validator/v10"
	"gopkg.in/ini.v1"
)

type ApiConfig struct {
	Url         string  `ini:"url" validate:"required"`
	Key         string  `ini:"key"`
	Temperature float64 `ini:"temperature"`
}

type ModelConfig struct {
	Classifier string `ini:"classifier" validate:"required"`
	Writer     string `ini:"writer" validate:"required"`
}

type Config struct {
	ApiConfig    ApiConfig   `ini:"api"`
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

func (c *Config) Validate() error {
	validate := validator.New()
	return validate.Struct(c)
}
