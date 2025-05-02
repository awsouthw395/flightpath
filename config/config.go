package config

import (
	"flightpath/server"
	"github.com/spf13/viper"
)

type Config struct {
	LogLevel string
	Server   server.Config
}

func (c *Config) Validate() error {
	err := c.Server.ValidateOrDefault()
	if err != nil {
		return err
	}
	if c.LogLevel == "" {
		c.LogLevel = "info"
	}

	return nil
}

func ReadConfig() (Config, error) {
	//FP = Flight Path
	viper.SetEnvPrefix("FP")

	viper.AutomaticEnv()

	config := Config{
		LogLevel: viper.GetString("LOG_LEVEL"),
		Server: server.Config{
			Port: viper.GetString("SERVER_PORT"),
		},
	}

	err := config.Validate()
	if err != nil {
		return Config{}, err
	}

	return config, nil
}
