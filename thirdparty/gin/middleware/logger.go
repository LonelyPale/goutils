package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// 日志记录到文件
func LoggerToFile() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now() // 开始时间

		c.Next() // 处理请求

		endTime := time.Now()                 // 结束时间
		latencyTime := endTime.Sub(startTime) // 执行时间
		status := c.Writer.Status()           // 状态码
		clientIP := c.ClientIP()              // 请求IP
		reqMethod := c.Request.Method         // 请求方式
		reqUri := c.Request.RequestURI        // 请求路由

		// 日志格式
		if len(c.Errors) > 0 || status >= 400 {
			errStr := c.Errors.String() // 错误信息
			if errStr != "" {
				log.WithFields(log.Fields{
					"module":  "gin",
					"status":  status,
					"latency": latencyTime,
					"client":  clientIP,
					"method":  reqMethod,
					"uri":     reqUri,
					"error":   errStr,
				}).Error()
				return
			}

			log.WithFields(log.Fields{
				"module":  "gin",
				"status":  status,
				"latency": latencyTime,
				"client":  clientIP,
				"method":  reqMethod,
				"uri":     reqUri,
			}).Error()
			return
		}

		log.WithFields(log.Fields{
			"module":  "gin",
			"status":  status,
			"latency": latencyTime,
			"client":  clientIP,
			"method":  reqMethod,
			"uri":     reqUri,
		}).Info()
	}
}
