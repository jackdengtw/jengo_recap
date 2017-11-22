package main

import (
	"log"
	"net/http"

	"github.com/golang/glog"
	"github.com/qetuantuan/jengo_recap/dao"
	"github.com/qetuantuan/jengo_recap/handler"
	"github.com/qetuantuan/jengo_recap/service"
)

func main() {

	mdHost := "127.0.0.1"
	md := &dao.MongoDao{Url: mdHost}
	if err := md.Init(); err != nil {
		glog.Fatal("init mongo failed! error:", err)
		return
	}
	service := &service.UserService{Md: md}
	// NewHandler implicitly register to handler.Handlers
	_ = handler.NewCreateUserHandler(service)
	_ = handler.NewGetUserHandler(service)
	_ = handler.NewGetUserByLoginHandler(service)
	_ = handler.NewUpdateScmTokenHandler(service)

	router := NewRouter(handler.Handlers)

	log.Fatal(http.ListenAndServe(":8088", router))
}
