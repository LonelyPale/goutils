package starter

import (
	"github.com/eclipse/paho.mqtt.golang"
	"github.com/go-spring/spring-boot"
	"github.com/go-spring/spring-logger"

	mqttClient "github.com/lonelypale/goutils/net/mqtt"
)

// go-spring spring-boot 启动器
func init() {
	SpringBoot.
		RegisterNameBeanFn("mqtt-client", func(cfg mqttClient.Config) (mqtt.Client, error) {
			client, err := mqttClient.Connect(&cfg)
			if err != nil {
				return nil, err
			}
			SpringLogger.Info("MQTT client connected successfully.")
			return client, nil
		}).
		ConditionOnMissingBean((*mqtt.Client)(nil)).
		Destroy(func(client mqtt.Client) {
			client.Disconnect(300)
			SpringLogger.Info("MQTT client closed successfully.")
		})
}
