// Created by LonelyPale at 2019-12-07

package config

import (
	"testing"
)

const configPath = "../database/mongodb/conf/mongodb.conf.toml"

type mongodbConfig struct {
	URI     string
	Timeout int
	DBName  string
}

type configFile struct {
	MongoDB *mongodbConfig
}

func TestReadConfig(t *testing.T) {
	conf := configFile{}
	err := ReadConfig(configPath, &conf)
	if err != nil {
		t.Error(err)
	}
	t.Log(conf.MongoDB)
}

func TestReadConfigMap(t *testing.T) {
	conf := make(map[string]interface{})
	err := ReadConfig(configPath, &conf)
	if err != nil {
		t.Error(err)
	}
	t.Log(conf)
	t.Log(conf["mongodb"])
}
