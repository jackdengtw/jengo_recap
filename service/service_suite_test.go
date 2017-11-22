package service

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"io/ioutil"
	"os"
	"testing"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/dbtest"
)

func TestService(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Service Suite")
}

var (
	d       string
	server  *dbtest.DBServer
	session *mgo.Session
)

var _ = BeforeSuite(func() {
	d, server, session = SetupMongo()
})

var _ = AfterSuite(func() {
	TearDown(d, server, session)
})

func SetupMongo() (d string, server *dbtest.DBServer, session *mgo.Session) {
	d, _ = ioutil.TempDir(os.TempDir(), "mongotools-test")
	server = &dbtest.DBServer{}
	server.SetPath(d)
	// Note that the server will be started automagically
	session = server.Session()
	fmt.Println("In setup")
	return
}

func TearDown(d string, server *dbtest.DBServer, session *mgo.Session) {
	session.Close()
	server.Wipe()
	server.Stop()
	os.RemoveAll(d)
	fmt.Println("In teardown")
}
