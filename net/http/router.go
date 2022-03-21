package http

import "github.com/gin-gonic/gin"

type Router struct {
	handler *Handler
}

func NewRouter(requestFilters []Filter, responseFilters []Filter) *Router {
	return &Router{
		handler: NewHandler(requestFilters, responseFilters),
	}
}

func (r *Router) Setup(engine *gin.Engine) {
	g := engine.Group("/")
	g.Any("ping", r.handler.Bind(ping))
}

func ping(_ *gin.Context) (string, error) {
	return "pong", nil
}
