// Created by LonelyPale at 2019-12-07

package config

import (
	"github.com/LonelyPale/goutils/config"
	"log"
	"os"
	"sync"
)

type mongodbConfig struct {
	URI         string
	Timeout     int
	MaxPoolSize int
	MinPoolSize int
	DBName      string
}

type configFile struct {
	MongoDB *mongodbConfig
}

var (
	conf configFile
	once sync.Once
)

const DefaultConfigFileName = "mongodb.conf.toml"

func ReadConfig(paths ...string) *mongodbConfig {
	once.Do(func() {
		var configPath string
		if len(paths) == 0 || paths[0] == "" {
			configPath = "." + string(os.PathSeparator) + DefaultConfigFileName
		} else {
			configPath = paths[0]
		}
		log.Printf("config file path: %s\n", configPath)

		err := config.ReadConfig(configPath, &conf)
		if err != nil {
			log.Fatal(err)
		}
	})

	return conf.MongoDB
}
