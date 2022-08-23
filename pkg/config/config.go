package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Port  string `mapstructure:"PORT"`
	DBURL string `mapstructure:"DB_URL"`
}

func LoadConfig() (*Config, error) {
	viper.AddConfigPath("./pkg/config/envs")
	viper.SetConfigName("local")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Panicf("Error reading config file: %s \n", err)
		return nil, err
	}

	var c Config
	if err := viper.Unmarshal(&c); err != nil {
		log.Panicf("Error unmarshalling config: %s \n", err)
		return nil, err
	}

	return &c, nil
}
