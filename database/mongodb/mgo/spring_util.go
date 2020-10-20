package mgo

import (
	"sync"

	"github.com/go-spring/spring-boot"

	"github.com/LonelyPale/goutils/database/mongodb"
	"github.com/LonelyPale/goutils/database/mongodb/config"
)

var (
	client *mongodb.Client
	first  sync.Once
)

func Client() *mongodb.Client {
	first.Do(func() {
		SpringBoot.GetBean(&client)
	})
	return client
}

func DB(db ...string) *mongodb.Database {
	if len(db) > 0 && len(db[0]) > 0 {
		return Client().Database(db[0])
	} else {
		defaultDBName := SpringBoot.GetStringProperty(config.KeyDefaultDBName)
		return Client().Database(defaultDBName)
	}
}
