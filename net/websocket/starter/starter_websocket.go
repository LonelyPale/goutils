package starter

import (
	"github.com/go-spring/spring-boot"
	"github.com/go-spring/spring-logger"

	"github.com/lonelypale/goutils/net/websocket"
)

// go-spring spring-boot 启动器
func init() {
	SpringBoot.RegisterNameBeanFn("websocket-container", func(config websocket.Config) *websocket.Container {
		return websocket.NewContainer(websocket.NewProcessor, &config)
	}).ConditionOnPropertyValue("web.server.websocket.enable", true)

	SpringBoot.RegisterNameBean("websocket-server", websocket.NewServer()).
		ConditionOnMissingBean((*websocket.Server)(nil)).
		ConditionOnPropertyValue("web.server.websocket.enable", true)

	SpringBoot.RegisterNameBean("websocket-server-starter", new(WebSocketStarter)).
		ConditionOnMissingBean((*WebSocketStarter)(nil)).
		ConditionOnPropertyValue("web.server.websocket.enable", true)
}

// WebSocketStarter WebSocket 服务器启动器
type WebSocketStarter struct {
	_ SpringBoot.ApplicationEvent `export:""`

	Server     *websocket.Server      `autowire:""`
	Containers []*websocket.Container `autowire:"[]?"`
}

func (starter *WebSocketStarter) OnStartApplication(ctx SpringBoot.ApplicationContext) {
	// 将收集到的 WS 容器赋值给 WS 服务器
	starter.Server.AddContainer(starter.Containers...)
	starter.Server.Start()
	SpringLogger.Info("open websocket server.")
}

func (starter *WebSocketStarter) OnStopApplication(ctx SpringBoot.ApplicationContext) {
	starter.Server.Stop()
	SpringLogger.Info("close websocket server.")
}
