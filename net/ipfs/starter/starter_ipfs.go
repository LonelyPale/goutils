package starter

import (
	"github.com/go-spring/spring-boot"
	"github.com/go-spring/spring-logger"

	"github.com/lonelypale/goutils/net/ipfs"
)

// go-spring spring-boot 启动器
func init() {
	SpringBoot.
		RegisterNameBeanFn("ipfs-client", func(cfg ipfs.Config) (*ipfs.Client, error) {
			client := ipfs.NewClient(&cfg)
			SpringLogger.Info("IPFS client connected successfully.")
			return client, nil
		}).
		ConditionOnMissingBean((*ipfs.Client)(nil)).
		Destroy(func(client *ipfs.Client) {
			client = nil
			SpringLogger.Info("IPFS client closed successfully.")
		})
}
