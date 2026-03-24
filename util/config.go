package util

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	DBDriver        string        `mapstructure:"DB_DRIVER"`
	DBSource        string        `mapstructure:"DB_SOURCE"`
	ServerAddress   string        `mapstructure:"SERVER_ADDRESS"`
	AccessSecretKey string        `mapstructure:"ACCESS_SECRET_KEY"`
	AccessDuration  time.Duration `mapstructure:"ACCESS_DURATION"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	if err = viper.ReadInConfig(); err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
