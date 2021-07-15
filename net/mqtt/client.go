package mqtt

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/eclipse/paho.mqtt.golang"
)

type Client = mqtt.Client
type Message = mqtt.Message

func NewClient(config *Config) mqtt.Client {
	//mqtt.DEBUG = log.New(os.Stdout, "", 0)
	mqtt.ERROR = log.New(os.Stdout, "", 0)

	opts := mqtt.NewClientOptions().
		AddBroker(config.BrokerURI).
		SetClientID(config.ClientID).
		SetKeepAlive(time.Duration(config.KeepAlive) * time.Second).
		SetPingTimeout(time.Duration(config.PingTimeout) * time.Second).
		SetAutoReconnect(config.AutoReconnect).
		SetMaxReconnectInterval(time.Duration(config.MaxReconnectInterval) * time.Second).
		SetUsername(config.Username).
		SetPassword(config.Password)
	opts.SetDefaultPublishHandler(defaultHandler) // 设置默认的消息处理回调函数

	return mqtt.NewClient(opts)
}

func Connect(config *Config) (mqtt.Client, error) {
	client := NewClient(config)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}
	return client, nil
}

// defaultHandler必须是安全的，可以被多个goroutine并发使用。
var defaultHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("TOPIC: %s\n", msg.Topic())
	fmt.Printf("MSG: %s\n", msg.Payload())
}
