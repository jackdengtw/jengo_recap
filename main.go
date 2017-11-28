package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/golang/glog"
	"github.com/qetuantuan/jengo_recap/client"
	"github.com/qetuantuan/jengo_recap/config"
	"github.com/qetuantuan/jengo_recap/dao"
	"github.com/qetuantuan/jengo_recap/handler"
	"github.com/qetuantuan/jengo_recap/scm"
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
	} else if cfg.Roles[0] == "project" {
		server = cfg.ProjectServer
		port = cfg.ProjectPort
	} else {
		glog.Fatal(fmt.Sprintf("Role %v not supported!", cfg.Roles[0]))
		os.Exit(-1)
	}

	githubScm := scm.NewGithubScm(cfg.GatewayHookUrl)

	ud := &dao.UserMongoDao{Url: mdHost}
	if err := ud.Init(); err != nil {
		glog.Fatal("init user mongo failed! error:", err)
		os.Exit(-2)
	}
	uservice := &service.UserService{Md: ud, GithubScm: githubScm}
	// NewHandler implicitly register to handler.Handlers

	// User API
	_ = handler.NewCreateUserHandler(uservice)
	_ = handler.NewGetUserHandler(uservice)
	_ = handler.NewGetUserByLoginHandler(uservice)
	_ = handler.NewUpdateScmTokenHandler(uservice)

	// Gateway API
	_ = handler.NewGitHubHandler(cfg.EngineServer, cfg.EnginePort)

	// Project API
	md := &dao.ProjectMongoDao{Url: mdHost}
	err = md.Init()
	if err != nil {
		glog.Fatal("init Project Dao failed! error:", err)
		os.Exit(-2)
	}
	usClient := client.NewUserStoreClient(cfg.UserServer, cfg.UserPort)

	pService := &service.ProjectService{
		Md:        md,
		GithubScm: githubScm,
		UsClient:  usClient,
	}
	_ = handler.NewUpdateProjectHandler(pService)
	_ = handler.NewSwitchProjectHandler(pService)
	_ = handler.NewGetProjectHandler(pService)
	_ = handler.NewGetProjectsHandler(pService)

	runLogService := &service.RunLogService{Md: md}
	_ = handler.NewPutLogHandler(runLogService)
	_ = handler.NewGetLogHandler(runLogService)

	bService := &service.RunService{
		Md:        md,
		GithubScm: githubScm,
	}
	_ = handler.NewGetBuildsByFilterHandler(bService)
	_ = handler.NewGetBuildsByIdsHandler(bService)
	_ = handler.NewUpdateRunHandler(bService)
	_ = handler.NewPartialUpdateRunHandler(bService)
	_ = handler.NewInsertRunHandler(bService)

	router := NewRouter(handler.Handlers)

	glog.Fatal(http.ListenAndServe(
		fmt.Sprintf("%v:%v", server, port), router))
}
