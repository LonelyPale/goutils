package springweb

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/go-spring/spring-gin"
	"github.com/go-spring/spring-web"
)

var store sessions.Store

var SessionFilter = defaultSessionFilter

type WebSessionConfig struct {
	Enable    bool   `value:"${web.server.session.enable:=true}"`                              //是否启用 Session
	SecretKey string `value:"${web.server.session.secret_key:=github.com/LonelyPale/goutils}"` //加密密钥
	Name      string `value:"${web.server.session.name:=s}"`                                   //cookie 属性名
}

func defaultSessionFilter(config WebSessionConfig) SpringWeb.Filter {
	store = memstore.NewStore([]byte(config.SecretKey))
	return SpringGin.Filter(sessions.Sessions(config.Name, store))
}

func GetSession(ctx SpringWeb.WebContext) sessions.Session {
	return sessions.Default(SpringGin.GinContext(ctx))
}
