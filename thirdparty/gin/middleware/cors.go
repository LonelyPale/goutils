package middleware

import (
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
)

func Cors(origins ...string) gin.HandlerFunc {
	var origin string
	if len(origins) > 0 {
		origin = origins[0]
	}

	return func(c *gin.Context) {
		method := c.Request.Method

		if len(origin) == 0 {
			referer := c.GetHeader("Referer")
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
		if c.GetHeader("X-Real-Ip") == "" && c.GetHeader("X-Forwarded-For") == "" {
			c.Writer.Header().Add("Access-Control-Allow-Credentials", "true")
		}

		c.Writer.Header().Add("Access-Control-Allow-Origin", origin)
		c.Writer.Header().Add("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Writer.Header().Add("Access-Control-Allow-Headers", "Content-Type,Content-Length,Accept,Authorization,X-Requested-With")
		c.Writer.Header().Add("Access-Control-Expose-Headers", "Accept,Authorization,X-Requested-With")

		//放行所有OPTIONS方法
		if method == "OPTIONS" {
			c.Status(http.StatusOK)
			c.Abort()
			return
		}

		c.Next()
	}
}
