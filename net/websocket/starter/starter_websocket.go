package starter

import (
	"github.com/go-spring/spring-boot"
	"github.com/go-spring/spring-logger"

	"github.com/LonelyPale/goutils/net/websocket"
)

// go-spring spring-boot 启动器
func init() {
	SpringBoot.
		RegisterNameBeanFn("websocket-server", func(opts websocket.Options) *websocket.Server {
			if !opts.Enable {
				return nil
			}

			SpringLogger.Info("open websocket server.")
			server := websocket.NewServer(websocket.NewReader, websocket.NewWriter, &opts)
			go server.Run()
			return server
		}).
		ConditionOnMissingBean((*websocket.Server)(nil)).
		Destroy(func(server *websocket.Server) {
			if server == nil {
				return
			}

			SpringLogger.Info("close websocket server.")
			if err := server.Close(); err != nil {
				SpringLogger.Error(err)
			}
		})
}
