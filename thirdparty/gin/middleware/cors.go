package middleware

import (
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"

	"github.com/LonelyPale/goutils"
)

func Cors(origins ...string) gin.HandlerFunc {
	origin := goutils.If(len(origins) > 0, origins[0], "").(string)

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

		c.Writer.Header().Add("Access-Control-Allow-Origin", origin)
		c.Writer.Header().Add("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Add("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Writer.Header().Add("Access-Control-Allow-Headers", "Content-Type,Access-Control-Allow-Headers,Content-Length,Accept,Authorization,X-Requested-With")
		c.Writer.Header().Add("Access-Control-Allow-Headers", "token,uid,from,source-site,gate-token,x-forwarded-for")

		//放行所有OPTIONS方法
		if method == "OPTIONS" {
			c.Status(http.StatusOK)
			return
		}

		c.Next()
	}
}
