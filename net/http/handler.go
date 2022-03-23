package http

import (
	"github.com/gin-gonic/gin"
)

type Filter func(ctx *gin.Context, args []interface{}) ([]interface{}, error)

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
