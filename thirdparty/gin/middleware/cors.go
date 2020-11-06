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

		c.Writer.Header().Add("Access-Control-Allow-Origin", origin)
		c.Writer.Header().Add("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Add("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Writer.Header().Add("Access-Control-Allow-Headers", "Content-Type,Access-Control-Allow-Headers,Content-Length,Accept,Authorization,X-Requested-With")

		//放行所有OPTIONS方法
		if method == "OPTIONS" {
			c.Status(http.StatusOK)
			c.Abort()
			return
		}

		c.Next()
	}
}
