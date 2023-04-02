package sysinfo

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	"github.com/lonelypale/goutils"
	"github.com/lonelypale/goutils/errors"
	mygin "github.com/lonelypale/goutils/thirdparty/gin"
)

func RunServer(addr ...string) {
	gin.SetMode(gin.ReleaseMode)
	server := mygin.NewServer(resolveAddress(addr), NewRouter())
	server.Run()
}

type Router struct{}

func NewRouter() *Router {
	return &Router{}
}

func (r *Router) Setup(engine *gin.Engine) {
	api := engine.Group("/")
	api.GET("/sys-info", sysinfo)
}

func sysinfo(ctx *gin.Context) {
	info, err := New()
	if err != nil {
		handleError(ctx, err)
		return
	}

	result := goutils.NewSuccessMessage(info)
	ctx.JSON(http.StatusOK, result)
}

func handleError(ctx *gin.Context, err error) {
	if err == nil {
		err = errors.New("unknown error")
	}

	result := goutils.NewErrorMessage(err)
	ctx.JSON(http.StatusOK, result)
}

func resolveAddress(addr []string) string {
	switch len(addr) {
	case 0:
		if port := os.Getenv("PORT"); port != "" {
			log.Infof("Environment variable PORT=\"%s\"", port)
			return ":" + port
		}
		log.Info("Environment variable PORT is undefined. Using port :8080 by default")
		return ":8080"
	case 1:
		return addr[0]
	default:
		panic("too many parameters")
	}
}
