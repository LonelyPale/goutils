package springweb

import (
	"github.com/go-spring/spring-boot"
	"github.com/go-spring/spring-core"
	"github.com/go-spring/spring-gin"
	"github.com/go-spring/spring-web"

	"github.com/LonelyPale/goutils/database/redis"
	"github.com/LonelyPale/goutils/thirdparty/gin/middleware"
	"github.com/LonelyPale/goutils/token"
)

type WebTokenConfig struct {
	Enable    bool        `value:"${web.server.token.enable:=true}"`                                    //是否启用 Token
	SecretKey string      `value:"${web.server.token.secret_key:=github.com/LonelyPale/goutils/token}"` //签名密钥
	Expire    int         `value:"${web.server.token.expire:=86400}"`                                   //过期时间, 24小时, 单位: 秒
	Cache     redis.Cache `autowire:"?"`                                                                //是否启用 redis 缓存
}

func init() {
	// 这种方式可以避免使用 export 语法，就像 StringFilter 和 NumberFilter 那样。
	SpringBoot.RegisterFilter(SpringCore.ObjectBean(new(TokenFilter))).ConditionOnOptionalPropertyValue("web.server.token.enable", true)
}

var TokenType token.Token = (*token.StandardToken)(nil)

type TokenFilter struct {
	Config WebTokenConfig
}

func (t *TokenFilter) Invoke(ctx SpringWeb.WebContext, chain SpringWeb.FilterChain) {
	middleware.Token(TokenType, &token.Options{
		SecretKey: t.Config.SecretKey,
		Expire:    t.Config.Expire,
		Cache:     t.Config.Cache,
	})(SpringGin.GinContext(ctx))
}
