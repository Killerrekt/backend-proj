package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Port           string `mapstructure:"PORT"`
	DBHost         string `mapstructure:"POSTGRES_HOST"`
	DBUserName     string `mapstructure:"POSTGRES_USER"`
	DBUserPassword string `mapstructure:"POSTGRES_PASSWORD"`
	DBName         string `mapstructure:"POSTGRES_DB"`
	DBPort         string `mapstructure:"POSTGRES_PORT"`
	ClientOrigin   string `mapstructure:"CLIENT_ORIGIN"`
}

type RedisConfig struct {
	REDIS_HOST string `mapstructure:"REDIS_HOST"`
	REDIS_PORT string `mapstructure:"REDIS_PORT"`
	DB         int    `mapstructure:"REDIS_DB"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		return
	}

	return
}

func LoadRedisConfig() (redisConfig RedisConfig, err error) {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()
	redisConfig = RedisConfig{
		REDIS_HOST: viper.GetString("REDIS_HOST"),
		REDIS_PORT: viper.GetString("REDIS_PORT"),
		DB:         viper.GetInt("REDIS_DB"),
	}

	if redisConfig.REDIS_HOST == "" || redisConfig.REDIS_PORT == "" {
		err = fmt.Errorf("REDIS_HOST and REDIS_PORT must be set in the configuration")
	}

	return redisConfig, err
}
