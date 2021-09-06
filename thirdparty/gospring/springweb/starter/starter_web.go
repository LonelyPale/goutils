package starter

import (
	"github.com/gin-gonic/gin"
	"github.com/go-spring/spring-boot"
	"github.com/go-spring/spring-gin"
	"github.com/go-spring/spring-web"
	"github.com/go-spring/starter-web"

	"github.com/LonelyPale/goutils"
	"github.com/LonelyPale/goutils/thirdparty/gospring/springweb"
)

type WebServerConfig struct {
	StarterWeb.WebServerConfig
	Cors    springweb.WebCorsConfig
	Session springweb.WebSessionConfig
	Static  springweb.WebStaticConfig
	Token   springweb.WebTokenConfig
}

func init() {
	SpringBoot.RegisterNameBeanFn("web-container", func(config WebServerConfig) SpringWeb.WebContainer {
		container := SpringGin.NewContainer(SpringWeb.ContainerConfig{
			Port: config.Port,
		})
		return ginHandler(container, config)
	}).ConditionOnOptionalPropertyValue("web.server.enable", true)

	SpringBoot.RegisterNameBeanFn("ssl-web-container", func(config WebServerConfig) SpringWeb.WebContainer {
		home := SpringBoot.GetStringProperty("app.home")
		key := goutils.Rootify(config.SSLKey, home)
		cert := goutils.Rootify(config.SSLCert, home)
		container := SpringGin.NewContainer(SpringWeb.ContainerConfig{
			EnableSSL: true,
			Port:      config.SSLPort,
			KeyFile:   key,
			CertFile:  cert,
		})
		return ginHandler(container, config)
	}).ConditionOnPropertyValue("web.server.ssl.enable", true)
}

func ginHandler(container *SpringGin.Container, config WebServerConfig) *SpringGin.Container {
	// 使用 gin 原生的中间件
	fLogger := SpringGin.Filter(gin.Logger())
	container.SetLoggerFilter(fLogger)

	// 使用 gin 原生的中间件
	fRecover := SpringGin.Filter(gin.Recovery())
	container.SetRecoveryFilter(fRecover)

	// gin cors 中间件
	if config.Cors.Enable {
		container.AddFilter(springweb.CorsFilter(config.Cors))
	}

	// gin session 中间件
	if config.Session.Enable {
		container.AddFilter(springweb.SessionFilter(config.Session))
	}

	// gin static 中间件
	if config.Static.Enable {
		container.AddFilter(springweb.StaticFilter(config.Static))
	}

	// gin token 中间件
	if config.Token.Enable {
		container.AddFilter(springweb.TokenFilter(config.Token))
	}

	return container
}
