// Created by LonelyPale at 2019-07-27

package config

import (
	"errors"
	"github.com/BurntSushi/toml"
	"path/filepath"
)

// conf 必须是非空的结构体指针
func ReadConfig(path string, conf interface{}) error {
	if path == "" {
		return errors.New("ReadConfig: path cannot be empty")
	}
	if conf == nil {
		return errors.New("ReadConfig: conf 必须是非空的结构体指针")
	}

	filePath, err := filepath.Abs(path)
	if err != nil {
		return err
	}

	if _, err := toml.DecodeFile(filePath, conf); err != nil {
		return err
	}

	return nil
}
