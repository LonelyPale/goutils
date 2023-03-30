package starter

import (
	SpringBoot "github.com/go-spring/spring-boot"
	SpringLogger "github.com/go-spring/spring-logger"

	"github.com/lonelypale/goutils/net/ipfs-cluster"
)

// go-spring spring-boot 启动器
func init() {
	SpringBoot.
		RegisterNameBeanFn("ipfs-cluster", func(cfg ipfscluster.Config) (*ipfscluster.Client, error) {
			client, err := ipfscluster.NewClient(&cfg)
			if err != nil {
				return nil, err
			}
			SpringLogger.Info("ipfs-cluster client connected successfully.")
			return client, nil
		}).
		ConditionOnMissingBean((*ipfscluster.Client)(nil)).
		Destroy(func(client *ipfscluster.Client) {
			client = nil
			SpringLogger.Info("ipfs-cluster client closed successfully.")
		})
}
