package main

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"github.com/lonelypale/goutils/app"
	"github.com/lonelypale/goutils/net/http"
)

func main() {
	server := http.NewServer(http.ServerOptions{})
	req := append(make([]http.Filter, 0), reqFunc)
	resp := append(make([]http.Filter, 0), respFunc)
	server.AddRouter(http.NewRouter(req, resp))
	app.Boot(server)
}

func reqFunc(ctx *gin.Context, args []interface{}) ([]interface{}, error) {
	fmt.Println("reqFunc:", args)
	return args, nil
}

func respFunc(ctx *gin.Context, args []interface{}) ([]interface{}, error) {
	fmt.Println("respFunc:", args)
	args[0] = args[0].(string) + " !"
	return args, nil
}
