package starter

import (
	"github.com/go-spring/spring-boot"
	"github.com/go-spring/spring-logger"

	"github.com/LonelyPale/goutils/database/redis"
)

// go-spring spring-boot 启动器
func init() {
	SpringBoot.
		RegisterNameBeanFn("redis-client", func(cfg redis.Config) (*redis.DB, error) {
			client, err := redis.NewRedisDB(&cfg)
			if err != nil {
				return nil, err
			}

			SpringLogger.Info("redis client opened successfully.")
			return client, nil
		}).
		ConditionOnMissingBean((*redis.DB)(nil)).
		Destroy(func(client *redis.DB) {
			if client == nil {
				SpringLogger.Error("redis client is nil")
				return
			}

			if err := client.Close(); err != nil {
				SpringLogger.Error(err)
				return
			}

			SpringLogger.Info("redis client closed successfully.")
		})
}
