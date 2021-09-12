package middleware

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-gonic/gin"
)

var (
	SecretKey   = "github.com/lonelypale/goutils"
	SessionName = "s"
)

var store sessions.Store

func Session() gin.HandlerFunc {
	store = memstore.NewStore([]byte(SecretKey))
	return sessions.Sessions(SessionName, store)
}

func GetSession(c *gin.Context) sessions.Session {
	return sessions.Default(c)
}
