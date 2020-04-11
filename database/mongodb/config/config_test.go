// Created by LonelyPale at 2019-12-07

package config

import "testing"

const configFile = "mongodb.conf.toml"

func TestReadConfigFile(t *testing.T) {
	conf := ReadConfigFile(configFile)
	t.Log(conf)
	t.Log(conf.Mongodb)
}
