package service

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/qetuantuan/jengo_recap/vo"
	"github.com/qetuantuan/jengo_recap/model"
)

type mockLogDao_testBuildLog struct {
	L *model.BuildLog // keep one in mem is enough for now
}

func (l *mockLogDao_testBuildLog) AddLog(logs []byte) (id string, err error) {
	l.L = &model.BuildLog{
		Id:      bson.NewObjectId().Hex(),
		Content: logs,
	}
	id = l.L.Id
	return
}

func (l *mockLogDao_testBuildLog) GetLog(id string) (log model.BuildLog, err error) {
	if l.L != nil && l.L.Id == id {
		log = *l.L
	} else {
		err = mgo.ErrNotFound
	}
	return
}

var _ = Describe("Test Build Log Service", func() {
	Describe("Create a log", func() {
		var (
			service *LocalBuildLogService
			uri     string
			err     error

			testContent1 = "i am a log"
		)
		BeforeEach(func() {
			service = &LocalBuildLogService{
				Md: &mockLogDao_testBuildLog{},
			}

			uri, err = service.PutLog([]byte(testContent1))
		})
		It("Should return success", func() {
			Expect(err).NotTo(HaveOccurred())
		})
		It("Should be able to get the log", func() {
			content, err := service.GetLog(
				&vo.GetBuildLogParams{
					Offset: 0,
					Limit:  0,
					Uri:    uri,
				})
			Expect(err).NotTo(HaveOccurred())
			Expect(string(content)).To(Equal(testContent1))
		})
		It("Should get an error with non-existed uri", func() {
			_, err := service.GetLog(
				&vo.GetBuildLogParams{
					Offset: 0,
					Limit:  0,
					Uri:    "non-existed",
				})
			Expect(err).To(Equal(NotFoundError))
		})
		It("Should honor offset and limit", func() {
			content, err := service.GetLog(
				&vo.GetBuildLogParams{
					Offset: 2,
					Limit:  4,
					Uri:    uri,
				})
			Expect(err).NotTo(HaveOccurred())
			Expect(string(content)).To(Equal(testContent1[2:6]))
		})
	})
})
