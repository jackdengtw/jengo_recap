package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/qetuantuan/jengo_recap/handler"
)

func NewRouter(
	handlers []handler.HandlerReader) *mux.Router {

	// TODO: implement validate later
	router := mux.NewRouter().StrictSlash(false)
	for _, handler := range handlers {
		var h http.Handler
		h = handler.HandlerFunc()
		h = Logger(h, handler.Name())

		router.
			Methods(handler.Method()).
			Path(handler.Pattern()).
			Name(handler.Name()).
			Handler(h)
	}

	return router
}
