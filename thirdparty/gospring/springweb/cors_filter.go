package springweb

import (
	"net/http"
	"net/url"

	"github.com/go-spring/spring-web"
)

var CorsFilter = defaultCorsFilter

type WebCorsConfig struct {
	Enable bool   `value:"${web.server.cors.enable:=true}"` //是否启用 Cors
	Origin string `value:"${web.server.cors.origin:=}"`     //授权的跨域地址
}

func defaultCorsFilter(config WebCorsConfig) SpringWeb.Filter {
	return &corsFilter{config}
}

type corsFilter struct {
	WebCorsConfig
}

func (c *corsFilter) Invoke(ctx SpringWeb.WebContext, chain SpringWeb.FilterChain) {
	origin := c.Origin
	method := ctx.Request().Method

	if len(origin) == 0 {
		referer := ctx.GetHeader("Referer")
		if len(referer) > 0 {
			refurl, err := url.Parse(referer)
			if err != nil {
				origin = "*"
			} else {
				origin = refurl.Scheme + "://" + refurl.Host
			}
		} else {
			origin = "*"
		}
	}

	//fmt.Println()
	//fmt.Println("***** ***** ***** ***** *****")
	//for key, val := range c.Request.Header {
	//	fmt.Println(key, val)
	//}
	//fmt.Println("***** ***** ***** ***** *****")

	//过滤 nginx 等服务器转发请求时，也会配置 Access-Control-Allow-Credentials 的情况，
	//保证最终 Response 的 Headers 中只能有一个 Access-Control-Allow-Credentials 出现，
	//当 Headers 中有多个 Access-Control-Allow-Credentials 时，浏览器会报错。
	if ctx.GetHeader("X-Real-Ip") == "" && ctx.GetHeader("X-Forwarded-For") == "" {
		ctx.Header("Access-Control-Allow-Credentials", "true")
	}

	ctx.Header("Access-Control-Allow-Origin", origin)
	ctx.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	ctx.Header("Access-Control-Allow-Headers", "Content-Type,Content-Length,Accept,Authorization,X-Requested-With")
	ctx.Header("Access-Control-Expose-Headers", "Accept,Authorization,X-Requested-With")

	//放行所有OPTIONS方法
	if method == "OPTIONS" {
		ctx.Status(http.StatusOK)
		ctx.Abort()
		return
	}

	chain.Next(ctx)
}
