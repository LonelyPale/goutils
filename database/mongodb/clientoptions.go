// Created by LonelyPale at 2019-12-06

package mongodb

import (
	"github.com/LonelyPale/goutils/database/mongodb/config"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type ClientOptions struct {
	//custom options
	URI     string
	Timeout time.Duration

	//options.ClientOptions
	MaxPoolSize uint64
	MinPoolSize uint64

	mongoClientOptions *options.ClientOptions
}

func NewClientOptions(uri string) *ClientOptions {
	clientOptions := &ClientOptions{
		URI:                uri,
		MaxPoolSize:        10,               //最大连接池
		MinPoolSize:        3,                //最小连接池
		Timeout:            10 * time.Second, //单位秒
		mongoClientOptions: options.Client(),
	}
	clientOptions.Apply()
	return clientOptions
}

func NewClientOptionsFromConfig(conf *config.MongodbConfig) *ClientOptions {
	clientOptions := NewClientOptions(conf.URI)

	if conf.MaxPoolSize > 0 {
		clientOptions.MaxPoolSize = uint64(conf.MaxPoolSize)
	}
	if conf.MinPoolSize > 0 {
		clientOptions.MinPoolSize = uint64(conf.MinPoolSize)
	}
	if conf.Timeout > 0 {
		clientOptions.Timeout = time.Duration(conf.Timeout) * time.Second
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
