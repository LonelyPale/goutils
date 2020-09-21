// Created by LonelyPale at 2019-12-07

package config

import (
	"github.com/LonelyPale/goutils/config"
	"log"
	"os"
)

type MongodbConfig struct {
	URI               string
	Timeout           int //单位秒
	EnableTransaction bool
	MaxPoolSize       int
	MinPoolSize       int
}

func DefaultMongodbConfig() *MongodbConfig {
	return &MongodbConfig{
		URI:               "mongodb://user:password@ip1:port,ip2:port/?replicaSet=replicaSet",
		EnableTransaction: false,
		MaxPoolSize:       10,
		MinPoolSize:       3,
		Timeout:           10,
	}
}

type Config struct {
	Mongodb *MongodbConfig
}

var (
	CommonConfig = &Config{Mongodb: DefaultMongodbConfig()}
)

const DefaultConfigFile = "mongodb.conf.toml"

func ReadConfigFile(files ...string) *Config {
	var configFile string
	if len(files) == 0 || files[0] == "" {
		configFile = "." + string(os.PathSeparator) + DefaultConfigFile
	} else {
		configFile = files[0]
	}
	log.Printf("read config file: %s\n", configFile)

	err := config.ReadConfigFile(configFile, CommonConfig)
	if err != nil {
		log.Fatal(err)
	}

	return CommonConfig
}
