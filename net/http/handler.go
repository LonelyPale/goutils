package http

import "github.com/gin-gonic/gin"

var _ Handler = new(HandlerFunc)

type Handler interface {
	Invoke(ctx *gin.Context, objs []interface{}) ([]interface{}, error)
}

type HandlerFunc func(ctx *gin.Context, objs []interface{}) ([]interface{}, error)

func (f HandlerFunc) Invoke(ctx *gin.Context, objs []interface{}) ([]interface{}, error) {
	return f(ctx, objs)
}

type handler struct {
	requestHandlers  []Handler
	responseHandlers []Handler
}

func NewHandler(requestHandlers []Handler, responseHandlers []Handler) *handler {
	return &handler{
		requestHandlers:  requestHandlers,
		responseHandlers: responseHandlers,
	}
}

func (h *handler) BIND(fn interface{}, requestHandlers []Handler, responseHandlers []Handler) gin.HandlerFunc {
	req := append([]Handler{}, h.requestHandlers...)
	req = append(req, requestHandlers...)
	resp := append([]Handler{}, h.requestHandlers...)
	resp = append(resp, responseHandlers...)
	return BIND(fn, req, resp)
}
