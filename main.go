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
	"github.com/qetuantuan/jengo_recap/queue"
	"github.com/qetuantuan/jengo_recap/scm"
	"github.com/qetuantuan/jengo_recap/service"
	"github.com/qetuantuan/jengo_recap/worker"
)

func main() {
	defer glog.Flush()

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

	ud := &dao.UserMongoDao{Url: cfg.MongoHost}
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
	md := &dao.ProjectMongoDao{Url: cfg.MongoHost}
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

	// init queue
	q := queue.NewNativeTaskQueue()
	q.Start()

	// Engine API
	rd := &dao.RunDao{Url: cfg.MongoHost}
	err = rd.Init()
	if err != nil {
		glog.Fatal("init Engine Dao failed! error:", err)
		os.Exit(-2)
	}
	eService := &service.EngineRunService{
		client.NewProjectStoreClient(
			cfg.ProjectServer,
			cfg.ProjectPort,
		),
		rd,
		q,
	}

	_ = handler.NewCreateRunHandler(eService)
	_ = handler.NewDescribeRunsHandler(eService)
	_ = handler.NewDescribeRunHandler(eService)

	rService := &service.RunLogService{
		Md: md,
	}
	_ = handler.NewGetRunLogHandler(rService)

	router := NewRouter(handler.Handlers)

	// Workers

	// Start Start Worker
	go func() {
		w := worker.Start{
			Base: worker.Base{
				Queue:          q,
				ListeningTopic: queue.TopicStartGroup1Name,
				OutputTopic: []string{
					queue.TopicParseGroup1Name,
				},
			},
			ProjectClient: client.NewProjectStoreClient(
				cfg.ProjectServer,
				cfg.ProjectPort,
			),
			RunDao: rd,
		}
		w.Work()
	}()

	// Start Parse Worker
	go func() {
		w := worker.Parse{
			Base: worker.Base{
				Queue:          q,
				ListeningTopic: queue.TopicParseGroup1Name,
				OutputTopic: []string{
					queue.TopicScheduleGroup1Name,
				},
			},
			Pc: client.NewProjectStoreClient(
				cfg.ProjectServer,
				cfg.ProjectPort,
			),
			Uc: client.NewUserStoreClient(
				cfg.UserServer,
				cfg.UserPort,
			),
			Gc: scm.NewGithubScm(
				cfg.GatewayHookUrl,
			),
		}
		w.Work()
	}()

	// Start Finalize Worker
	go func() {
		w := worker.Finalize{
			Base: worker.Base{
				Queue:          q,
				ListeningTopic: queue.TopicFinalizeGroup1Name,
			},
			ProjectClient: client.NewProjectStoreClient(
				cfg.ProjectServer,
				cfg.ProjectPort,
			),
			RunDao: rd,
		}
		w.Work()
	}()

	// Start Schedule Worker

	// Start Execute Worker
	// Register Execute Worker to registry.WorkerRegistry

	glog.Fatal(http.ListenAndServe(
		fmt.Sprintf("%v:%v", server, port), router))
}
