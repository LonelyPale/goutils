package mqtt

type Config struct {
	Enable      bool   `value:"${mqtt.client.enable:=false}"`                    //是否启用 MQTT Client
	BrokerURI   string `value:"${mqtt.client.broker_uri:=tcp://127.0.0.1:1883}"` //broker URI
	ClientID    string `value:"${mqtt.client.client_id:=mqtt-client}"`           //客户端名称
	KeepAlive   int    `value:"${mqtt.client.keep_alive:=60}"`                   //ping请求发送间隔(单位秒)
	PingTimeout int    `value:"${mqtt.client.ping_timeout:=3}"`                  //ping请求超时时间(单位秒)
}

func DefaultConfig() *Config {
	return &Config{
		Enable:      true,
		BrokerURI:   "tcp://127.0.0.1:1883",
		ClientID:    "mqtt-client",
		KeepAlive:   60,
		PingTimeout: 3,
	}
}
