// Created by LonelyPale at 2019-12-07

package config

import "testing"

const configPath = "./mongodb.conf.toml"

func TestReadConfig(t *testing.T) {
	conf := ReadConfig(configPath)
	t.Log(conf)
}
