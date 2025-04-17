package main

import "github.com/spf13/viper"

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

	viper.SetConfigFile(path)

	err := viper.ReadInConfig()
	if err != nil {
		return config, err
	}

	err = viper.Unmarshal(&config)
	return config, err
}
