package handler

import (
	"net/http"
)

var Handlers []HandlerReader = make([]HandlerReader, 0)

type HandlerReader interface {
	Name() string
	Method() string
	Pattern() string
	Query() string
	HandlerFunc() http.HandlerFunc
}

type BaseHandler struct {
	name        string
	method      string
	pattern     string
	query       string
	handlerFunc http.HandlerFunc
}

func (h *BaseHandler) Name() string {
	return h.name
}

func (h *BaseHandler) Method() string {
	return h.method
}

func (h *BaseHandler) Pattern() string {
	return h.pattern
}
func (h *BaseHandler) Query() string {
	return h.query
}

func (h *BaseHandler) HandlerFunc() http.HandlerFunc {
	return h.handlerFunc
}

// Register to a readonly array for later init
func Register(r HandlerReader) {
	Handlers = append(Handlers, r)
}
