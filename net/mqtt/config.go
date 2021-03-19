package mqtt

type Config struct {
	Enable               bool   `value:"${mqtt.client.enable:=false}"`                    //是否启用 MQTT Client
	BrokerURI            string `value:"${mqtt.client.broker_uri:=tcp://127.0.0.1:1883}"` //broker URI
	ClientID             string `value:"${mqtt.client.client_id:=mqtt-client}"`           //客户端名称
	KeepAlive            int    `value:"${mqtt.client.keep_alive:=60}"`                   //ping请求发送间隔(单位秒)
	PingTimeout          int    `value:"${mqtt.client.ping_timeout:=3}"`                  //ping请求超时时间(单位秒)
	AutoReconnect        bool   `value:"${mqtt.client.auto_reconnect:=true}"`             //是否自动重连
	MaxReconnectInterval int    `value:"${mqtt.client.max_reconnect_interval:=60}"`       //自动重连时间间隔(单位秒)
	Username             string `value:"${mqtt.client.username:=}"`                       //认证用户名
	Password             string `value:"${mqtt.client.password:=}"`                       //认证密码
}

func DefaultConfig() *Config {
	return &Config{
		Enable:               true,
		BrokerURI:            "tcp://127.0.0.1:1883",
		ClientID:             "mqtt-client",
		KeepAlive:            60,
		PingTimeout:          3,
		AutoReconnect:        true,
		MaxReconnectInterval: 60,
	}
}
