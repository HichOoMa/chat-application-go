package config

import (
	"errors"

	"github.com/spf13/viper"
)

var AppConfig Config

type Config struct {
	Port         int16  `mapstructure:"PORT"`
	JwtIsuuer    string `mapstructure:"JWT_ISSUER"`
	JwtSecretKey string `mapstructure:"JWT_SECRET_KEY"`
	JwtExpired   int    `mapstructure:"JWT_EXPIRED"`
	MongoURI     string `mapstructure:"MONGO_URI"`
	DatabaseName string `mapstructure:"DATABASE_NAME"`
}

func InitializeAppConfig() error {
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	viper.AddConfigPath("internal/config")
	viper.AddConfigPath("/")
	viper.AllowEmptyEnv(true)
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		return errors.New("failed to load config file")
	}

	err = viper.Unmarshal(&AppConfig)
	if err != nil {
		return errors.New("failed to parse env to config struct")
	}

	if AppConfig.Port == 0 {
		return errors.New("required variabel environment is empty")
	}

	return nil
}
