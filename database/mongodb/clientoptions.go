// Created by LonelyPale at 2019-12-06

package mongodb

import (
	"github.com/LonelyPale/goutils/database/mongodb/config"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type ClientOptions struct {
	//custom options
	URI               string
	Timeout           time.Duration
	EnableTransaction bool

	//options.ClientOptions
	MaxPoolSize uint64
	MinPoolSize uint64

	mongoClientOptions *options.ClientOptions
}

func NewClientOptions(uri string) *ClientOptions {
	clientOptions := &ClientOptions{
		URI:                uri,
		Timeout:            10 * time.Second, //单位秒
		EnableTransaction:  false,            //默认不启用事务,启用事务需要打开副本集
		MaxPoolSize:        10,               //最大连接池
		MinPoolSize:        3,                //最小连接池
		mongoClientOptions: options.Client(),
	}
	clientOptions.Apply()
	return clientOptions
}

func NewClientOptionsFromConfig(conf *config.MongodbConfig) *ClientOptions {
	clientOptions := NewClientOptions(conf.URI)
	clientOptions.EnableTransaction = conf.EnableTransaction
	if conf.Timeout > 0 {
		clientOptions.Timeout = time.Duration(conf.Timeout) * time.Second
	}
	if conf.MaxPoolSize > 0 {
		clientOptions.MaxPoolSize = uint64(conf.MaxPoolSize)
	}
	if conf.MinPoolSize > 0 {
		clientOptions.MinPoolSize = uint64(conf.MinPoolSize)
	}

	clientOptions.Apply()
	return clientOptions
}

func NewClientOptionsFromFile(files ...string) *ClientOptions {
	conf := config.ReadConfigFile(files...)
	return NewClientOptionsFromConfig(conf.Mongodb)
}

func (c *ClientOptions) Apply() {
	c.mongoClientOptions.ApplyURI(c.URI)
	c.mongoClientOptions.SetMaxPoolSize(c.MaxPoolSize)
	c.mongoClientOptions.SetMinPoolSize(c.MinPoolSize)
}
