package session

import (
	"github.com/gin-gonic/gin"
)

const (
	DefaultKey = "github.com/LonelyPale/goutils/session"
)

func Sessions(store Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		s := NewSession(store)
		c.Set(DefaultKey, s)
		c.Next()
	}
}

func Default(c *gin.Context) Session {
	return c.MustGet(DefaultKey).(Session)
}
