package springweb

import (
	"github.com/go-spring/spring-gin"
	"github.com/go-spring/spring-web"

	"github.com/LonelyPale/goutils/database/redis"
	"github.com/LonelyPale/goutils/thirdparty/gin/middleware"
	"github.com/LonelyPale/goutils/token"
)

var TokenFilter = defaultTokenFilter
var TokenType = (*token.StandardToken)(nil)

type WebTokenConfig struct {
	Enable    bool        `value:"${web.server.static.enable:=true}"`                                   //是否启用 Token
	SecretKey string      `value:"${web.server.token.secret_key:=github.com/LonelyPale/goutils/token}"` //签名密钥
	Expire    int         `value:"${web.server.token.expire:=86400}"`                                   //过期时间, 24小时, 单位: 秒
	Cache     redis.Cache `autowire:""`                                                                 //是否启用 redis 缓存
}

func defaultTokenFilter(config WebTokenConfig) SpringWeb.Filter {
	return SpringGin.Filter(middleware.Token(TokenType, &token.Options{
		SecretKey: config.SecretKey,
		Expire:    config.Expire,
		Cache:     config.Cache,
	}))
}
