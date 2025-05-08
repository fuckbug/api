package main

import (
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Logger   loggerConf
	Port     int
	Postgres postgresConf
	Domain   string
}

type loggerConf struct {
	Level string
}

type postgresConf struct {
	Dsn string
}

func LoadConfig(path string) (Config, error) {
	config := Config{}

	if path != "" {
		viper.SetConfigFile(path)
		if err := viper.ReadInConfig(); err != nil {
			return config, err
		}
	}

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	err := viper.Unmarshal(&config)
	return config, err
}
