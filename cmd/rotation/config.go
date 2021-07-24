package main

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Logger  LoggerConf
	Storage StorageConf
	// Slots   app.Slots `toml:"slots"` //удалить
}

type LoggerConf struct {
	File  string `toml:"logger.file"`
	Level string `toml:"logger.level"`
}

type StorageConf struct {
	ConnStr string `toml:"storage.connStr"`
}

// NewConfig make a config from configFilePath.
func NewConfig() Config {
	viper.SetConfigFile(configFilePath)
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("failed to read config: %v", err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalf("failed to unmarshal config: %v", err)
	}
	return config
}
