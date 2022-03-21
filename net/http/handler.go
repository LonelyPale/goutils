package http

import (
	"github.com/gin-gonic/gin"
)

var _ Filter = new(FilterFunc)

type Filter interface {
	Invoke(ctx *gin.Context, args []interface{}) ([]interface{}, error)
}

type FilterFunc func(ctx *gin.Context, args []interface{}) ([]interface{}, error)

func (f FilterFunc) Invoke(ctx *gin.Context, args []interface{}) ([]interface{}, error) {
	return f(ctx, args)
}

type Handler struct {
	requestFilters  []Filter
	responseFilters []Filter
}

func NewHandler(requestFilters []Filter, responseFilters []Filter) *Handler {
	return &Handler{
		requestFilters:  append([]Filter{}, requestFilters...),
		responseFilters: append([]Filter{}, responseFilters...),
	}
}

func (h *Handler) BindHandler(fn interface{}) *bindHandler {
	return NewBindHandler(fn).AddRequestFilter(h.requestFilters...).AddResponseFilter(h.responseFilters...)
}

func (h *Handler) Bind(fn interface{}) gin.HandlerFunc {
	return NewBindHandler(fn).AddRequestFilter(h.requestFilters...).AddResponseFilter(h.responseFilters...).Invoke
}
