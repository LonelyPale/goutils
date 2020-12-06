package springweb

import (
	"github.com/gin-contrib/static"
	"github.com/go-spring/spring-gin"
	"github.com/go-spring/spring-web"
)

var StaticFilter = defaultStaticFilter

type WebStaticConfig struct {
	Enable              bool   `value:"${web.server.static.enable:=false}"`                //是否启用 Static 静态目录
	URLPrefix           string `value:"${web.server.static.url_prefix:=/static}"`          //url prefix
	LocalPath           string `value:"${web.server.static.local_path:=./static}"`         //目录路径
	AllowDirectoryIndex bool   `value:"${web.server.static.allow_directory_index:=false}"` //是否允许目录索引
}

func defaultStaticFilter(config WebStaticConfig) SpringWeb.Filter {
	return SpringGin.Filter(static.Serve(config.URLPrefix, static.LocalFile(config.LocalPath, config.AllowDirectoryIndex)))
}
