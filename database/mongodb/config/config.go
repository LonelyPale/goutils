// Created by LonelyPale at 2019-12-07

package config

import (
	"log"
	"os"

	"github.com/LonelyPale/goutils/config"
)

//支持 go-spring 属性绑定
type MongodbConfig struct {
	URI               string `value:"${mongodb.uri:=mongodb://localhost}"`
	MinPoolSize       int    `value:"${mongodb.min_pool_size:=3}"`
	MaxPoolSize       int    `value:"${mongodb.max_pool_size:=10}"`
	Timeout           int    `value:"${mongodb.timeout:=10}"` //单位秒
	EnableTransaction bool   `value:"${mongodb.enable_transaction:=false}"`
}

func DefaultMongodbConfig() *MongodbConfig {
	return &MongodbConfig{
		URI:               "mongodb://user:password@ip1:port,ip2:port/?replicaSet=replicaSet",
		MinPoolSize:       3,
		MaxPoolSize:       10,
		Timeout:           10,
		EnableTransaction: false,
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
