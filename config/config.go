// Created by LonelyPale at 2019-07-27

package config

import (
	"errors"
	"github.com/BurntSushi/toml"
	"path/filepath"
)

// conf 必须是非空的结构体指针
func ReadConfigFile(file string, conf interface{}) error {
	if file == "" {
		return errors.New("ReadConfigFile: file path cannot be empty")
	}
	if conf == nil {
		return errors.New("ReadConfigFile: conf 必须是非空的结构体指针")
	}

	abspath, err := filepath.Abs(file)
	if err != nil {
		return err
	}

	if _, err := toml.DecodeFile(abspath, conf); err != nil {
		return err
	}

	return nil
}
