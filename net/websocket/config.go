package websocket

import "time"

type Config struct {
	Enable            bool          `value:"${web.server.websocket.enable:=false}"`            //是否启用WebSocket
	Origin            bool          `value:"${web.server.websocket.origin:=true}"`             //是否启用跨域
	ReadDeadline      time.Duration `value:"${web.server.websocket.read_deadline:=0}"`         //消息单次读取超时时间，单位：秒
	WriteDeadline     time.Duration `value:"${web.server.websocket.write_deadline:=0}"`        //消息单次写入超时时间，单位：秒
	ReadBufferSize    int           `value:"${web.server.websocket.read_buffer_size:=20480}"`  //connect read buffer size: 20kb
	WriteBufferSize   int           `value:"${web.server.websocket.write_buffer_size:=20480}"` //connect write buffer size: 20kb
	MaxMessageSize    int64         `value:"${web.server.websocket.max_message_size:=65535}"`  //从消息管道读取消息的最大字节: 65535 byte
	ProcessorPoolSize int           `value:"${web.server.websocket.processor_pool_size:=10}"`  //处理器协程池大小
	InChanSize        int           `value:"${web.server.websocket.in_chan_size:=100}"`        //已读管道大小
	OutChanSize       int           `value:"${web.server.websocket.out_chan_size:=100}"`       //待写管道大小
}

func DefaultConfig() *Config {
	return &Config{
		Enable:            true,
		Origin:            true,
		ReadDeadline:      0,
		WriteDeadline:     0,
		ReadBufferSize:    20480,
		WriteBufferSize:   20480,
		MaxMessageSize:    65535,
		ProcessorPoolSize: 10,
		InChanSize:        100,
		OutChanSize:       100,
	}
}
