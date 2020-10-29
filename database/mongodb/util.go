package mongodb

import (
	"sync"

	"github.com/go-spring/spring-boot"
	log "github.com/sirupsen/logrus"

	"github.com/LonelyPale/goutils/database/mongodb/config"
)

const (
	logModule = "mongodb"
)

var (
	singleton *Client
	once      sync.Once
)

func getClient(cfgs ...*config.MongodbConfig) *Client {
	var instance *Client
	if len(cfgs) > 0 {
		var err error
		opts := NewClientOptionsFromConfig(cfgs[0])
		instance, err = Connect(opts)
		if err != nil {
			log.WithField("module", logModule).Panic(err)
		}
	} else {
		if ok := SpringBoot.GetBean(&instance); !ok {
			log.WithField("module", logModule).Panic("Can't get mongodb instance")
		}
	}
	return instance
}

func GetInstance(cfgs ...*config.MongodbConfig) *Client {
	once.Do(func() {
		singleton = getClient(cfgs...)
	})
	return singleton
}

func DB(db ...string) *Database {
	return GetInstance().DB(db...)
}
