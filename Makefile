program := jengo_recap
COMMIT := $(shell git log -1 --pretty=format:"%h")
SRC_PATH := ${GOPATH}/src/github.com/qetuantuan/$(program)
.PHONY: clean build test dist all

all:: test

check:: 
	[[ "${GOPATH}" != "" ]] || `echo "GOPATH not set!" && exit 1`

clean::
	rm -rf $(program)*

build:: check clean
	cd $(SRC_PATH)  && go build -o $(program)_$(COMMIT)

test:: build
	go test ./...

local_cd:: test
	ps -ef | grep $(program) | grep -v grep | awk '{print $$2}' | xargs sudo kill -9 > /dev/null 2>&1 || echo "no such process"
	nohup $(SRC_PATH)/$(program)_$(COMMIT) -logtostderr > $(SRC_PATH)/nohup.out 2>&1 &

dist:: check clean
	env GOOS=linux GOARCH=amd64 GO_GCFLAGS=-N go build -o $(program)_$(COMMIT)
	scp -C $(program)_* root@47.93.38.8:/root/workspace/$(program)/

init_env:: check
	mkdir -p `dirname ${SRC_PATH}`
	go get -v github.com/golang/glog
	go get -v github.com/gorilla/mux
	go get -v gopkg.in/mgo.v2/bson
	go get -v github.com/onsi/ginkgo
	go get -v github.com/onsi/gomega
