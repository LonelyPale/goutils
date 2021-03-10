package starter

import (
	"github.com/go-spring/spring-boot"
	"github.com/go-spring/spring-logger"

	"github.com/LonelyPale/goutils/database/mongodb"
	"github.com/LonelyPale/goutils/database/mongodb/config"
)

// go-spring spring-boot 启动器
func init() {
	SpringBoot.
		RegisterNameBeanFn("mongodb-client", func(cfg config.MongodbConfig) (*mongodb.Client, error) {
			client, err := mongodb.Connect(mongodb.NewClientOptionsFromConfig(&cfg))
			if err != nil {
				return nil, err
			}
			SpringLogger.Info("Mongodb client opened successfully.")
			return client, nil
		}).
		ConditionOnMissingBean((*mongodb.Client)(nil)).
		Destroy(func(client *mongodb.Client) {
			if err := mongodb.CloseClient(client); err != nil {
				SpringLogger.Error(err)
				return
			}
			SpringLogger.Info("Mongodb client closed successfully.")
		})
}
