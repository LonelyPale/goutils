// Created by LonelyPale at 2019-12-06

package mongodb

import (
	"time"

	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/lonelypale/goutils/database/mongodb/config"
)

type ClientOptions struct {
	//custom options
	URI               string        //数据库连接字符串
	Timeout           time.Duration //单位秒
	EnableTransaction bool          //默认不启用事务,启用事务需要打开副本集
	DefaultDBName     string        //默认的数据库名称

	//options.ClientOptions
	MaxPoolSize uint64 //最大连接池
	MinPoolSize uint64 //最小连接池

	mongoClientOptions *options.ClientOptions
}

func NewClientOptions(uri string) *ClientOptions {
	clientOptions := &ClientOptions{
		URI:                uri,
		Timeout:            10 * time.Second,
		EnableTransaction:  false,
		DefaultDBName:      "test",
		MaxPoolSize:        10,
		MinPoolSize:        3,
		mongoClientOptions: options.Client(),
	}
	clientOptions.Apply()
	return clientOptions
}

func NewClientOptionsFromConfig(conf *config.MongodbConfig) *ClientOptions {
	clientOptions := NewClientOptions(conf.URI)
	clientOptions.EnableTransaction = conf.EnableTransaction
	clientOptions.DefaultDBName = conf.DefaultDBName
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
