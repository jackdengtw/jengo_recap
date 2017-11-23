package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/golang/glog"
	"github.com/qetuantuan/jengo_recap/config"
	"github.com/qetuantuan/jengo_recap/dao"
	"github.com/qetuantuan/jengo_recap/handler"
	"github.com/qetuantuan/jengo_recap/service"
)

func main() {
	defer glog.Flush()

	mdHost := "127.0.0.1"

	// get config file
	flag.Parse()
	confFile := flag.String("config", "./config.json", "config file for jengo_recap application")

	// parse config file
	cfg := &config.ConfigEntity{}
	err := cfg.LoadFromFile(*confFile)
	if err != nil {
		glog.Infof("config file load error: err(%s)", err.Error())
		return
	}

	var server string
	var port int
	if len(cfg.Roles) != 1 {
		server = cfg.DefaultServer
		port = cfg.DefaultPort
	} else if cfg.Roles[0] == "user" {
		server = cfg.UserServer
		port = cfg.UserPort
	} else if cfg.Roles[0] == "gateway" {
		server = cfg.GatewayServer
		port = cfg.GatewayPort
	} else {
		glog.Fatal(fmt.Sprintf("Role %v not supported!", cfg.Roles[0]))
		os.Exit(-1)
	}

	md := &dao.MongoDao{Url: mdHost}
	if err := md.Init(); err != nil {
		glog.Fatal("init mongo failed! error:", err)
		return
	}
	service := &service.UserService{Md: md}
	// NewHandler implicitly register to handler.Handlers

	// User API
	_ = handler.NewCreateUserHandler(service)
	_ = handler.NewGetUserHandler(service)
	_ = handler.NewGetUserByLoginHandler(service)
	_ = handler.NewUpdateScmTokenHandler(service)

	// Gateway API
	_ = handler.NewGitHubHandler(cfg.EngineServer, cfg.EnginePort)

	router := NewRouter(handler.Handlers)

	glog.Fatal(http.ListenAndServe(
		fmt.Sprintf("%v:%v", server, port), router))
}
