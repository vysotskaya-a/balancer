package config

import (
	"github.com/spf13/viper"
	"log"
)

type Config struct {
	Port    string   `mapstructure:"port"`
	Targets []string `mapstructure:"targets"`
	Redis   struct {
		Addr string `mapstructure:"addr"`
		DB   int    `mapstructure:"db"`
	} `mapstructure:"redis"`
}

func LoadConfig() *Config {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./configs")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Ошибка загрузки конфига, %s", err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalf("Ошибка парсинга конфига, %s", err)
	}
	return &config
}
