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
			SpringLogger.Info("open mongodb")
			return mongodb.Connect(mongodb.NewClientOptionsFromConfig(&cfg))
		}).
		ConditionOnMissingBean((*mongodb.Client)(nil)).
		Destroy(func(client *mongodb.Client) {
			SpringLogger.Info("close mongodb")
			if err := mongodb.CloseClient(client); err != nil {
				SpringLogger.Error(err)
			}
		})
}
