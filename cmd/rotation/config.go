package main

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Logger     LoggerConf
	Storage    StorageConf
	GRPCServer GRPCServerConf
	MQ         MQConf
}

type LoggerConf struct {
	File  string `yaml:"logger.file"`
	Level string `yaml:"logger.level"`
}

type StorageConf struct {
	ConnStr string `yaml:"storage.connStr"`
}

type GRPCServerConf struct {
	HostPort string `yaml:"GRPCServer.hostPort"`
}

type MQConf struct {
	URI        string `yaml:"MQ.uri"`
	Exchange   string `yaml:"MQ.exchange"`
	Reliable   bool   `yaml:"MQ.reliable"`
	RoutingKey string `yaml:"MQ.routingKey"`
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
