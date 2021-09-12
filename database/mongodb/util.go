package mongodb

import (
	"github.com/go-spring/spring-boot"
	log "github.com/sirupsen/logrus"

	"github.com/lonelypale/goutils/database/mongodb/config"
)

const (
	logModule = "mongodb"
)

func GetClient(cfgs ...*config.MongodbConfig) *Client {
	var instance *Client
	var err error

	if len(cfgs) > 0 {
		opts := NewClientOptionsFromConfig(cfgs[0])
		instance, err = Connect(opts)
		if err != nil {
			log.WithField("module", logModule).Panic(err)
		}
	} else {
		if ok := SpringBoot.GetBean(&instance); !ok {
			log.WithField("module", logModule).Panic("Can't get spring mongodb instance")
		}
	}

	return instance
}
