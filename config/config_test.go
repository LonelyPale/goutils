// Created by LonelyPale at 2019-12-07

package config

import (
	"testing"
)

const configFile = "test.conf.toml"

type ServerConfig struct {
	Name string
	Host string
	Port int
}

type Config struct {
	Server *ServerConfig
}

func TestReadConfigFile(t *testing.T) {
	conf := &Config{}
	err := ReadConfigFile(configFile, conf)
	if err != nil {
		t.Error(err)
	}
	t.Log(conf)
	t.Log(conf.Server)
}

func TestReadConfigFile2Map(t *testing.T) {
	conf := make(map[string]interface{})
	err := ReadConfigFile(configFile, &conf)
	if err != nil {
		t.Error(err)
	}
	t.Log(conf)
	t.Log(conf["Server"])
}
