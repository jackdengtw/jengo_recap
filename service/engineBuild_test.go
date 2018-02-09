package service

import (
	"strconv"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	_ "github.com/qetuantuan/jengo_recap/vo"
	"github.com/qetuantuan/jengo_recap/model"
	"github.com/qetuantuan/jengo_recap/queue"
)

type mockEngineBuildDao_testEngineBuild struct {
	B *model.Build
}

func (m *mockEngineBuildDao_testEngineBuild) UpdateBuildProperties(string, map[string]interface{}) error {
	return nil
}

func (m *mockEngineBuildDao_testEngineBuild) InsertBuild(build model.Build) (string, error) {
	m.B = &build
	m.B.Id = bson.NewObjectId().Hex()
	return m.B.Id, nil
}

func (m *mockEngineBuildDao_testEngineBuild) ListBuilds(map[string]interface{}, int, int) (b model.Builds, err error) {
	if m.B != nil {
		b = append(b, *m.B)
	} else {
		err = mgo.ErrNotFound
	}
	return
}

func (m *mockEngineBuildDao_testEngineBuild) GetBuild(string) (b model.Build, err error) {
	if m.B != nil {
		b = *m.B
	} else {
		err = mgo.ErrNotFound
	}
	return
}

/*
	ListBuilds(p *model.EngineListBuildsParams) (model.Builds, error)
	GetBuild(buildId string) (model.Build, error)
*/

var _ = Describe("Test Engine Build Service", func() {
	Describe("Create a build", func() {
		var (
			service *LocalEngineBuildService
			repoId1 int = 123
			buildId string
			err     error
		)
		BeforeEach(func() {
			service = &LocalEngineBuildService{
				Queue: queue.NewNativeTaskQueue(),
				Md:    &mockEngineBuildDao_testEngineBuild{},
			}

			buildId, err = service.CreateBuild(&model.EngineCreateBuildParams{
				// RepoId:  "repo1",
				UserId:  "user1",
				EventId: "event1",
				Repo: &model.PushEventRepository{
					ID: &repoId1,
				},
				Commits: []model.PushEventCommit{
					model.PushEventCommit{
					// PushEventCommit: vo.PushEventCommit{},
					},
				},
			})
		})
		It("Should return success", func() {
			Expect(err).NotTo(HaveOccurred())
		})
		It("Should be able to get the build", func() {
			build, err := service.GetBuild(buildId)
			Expect(err).NotTo(HaveOccurred())
			Expect(build.RepoId).To(Equal(strconv.Itoa(repoId1)))
			Expect(build.UserId).To(Equal("user1"))
			Expect(*build.EventId).To(Equal("event1"))
		})
		It("Should be able to list the builds", func() {
			builds, err := service.ListBuilds(&model.EngineListBuildsParams{
				UserId: "user1",
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(len(builds)).To(Equal(1))
			Expect(builds[0].RepoId).To(Equal(strconv.Itoa(repoId1)))
		})
	})
})
