package mgo

import (
	"sync"

	log "github.com/sirupsen/logrus"

	"github.com/LonelyPale/goutils/database/mongodb"
	"github.com/LonelyPale/goutils/database/mongodb/config"
)

const (
	logModule = "mongodb"
)

var (
	singleton *mongodb.Client
	once      sync.Once
)

func GetInstance(dbcfgs ...*config.MongodbConfig) *mongodb.Client {
	once.Do(func() {
		if len(dbcfgs) < 1 {
			log.WithField("module", logModule).Panic("first init MongodbConfig cannot be empty")
		}

		var err error
		opts := mongodb.NewClientOptionsFromConfig(dbcfgs[0])
		singleton, err = mongodb.Connect(opts)
		if err != nil {
			log.WithField("module", logModule).Panic(err)
		}
	})
	return singleton
}
